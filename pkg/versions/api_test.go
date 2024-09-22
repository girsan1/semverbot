package versions

import (
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/restechnica/semverbot/internal/fakes"
	"github.com/restechnica/semverbot/internal/mocks"
	"github.com/restechnica/semverbot/pkg/cli"
	"github.com/restechnica/semverbot/pkg/git"
	"github.com/restechnica/semverbot/pkg/modes"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func TestAPI_GetVersion(t *testing.T) {
	type Test struct {
		Name    string
		Prefix  string
		Suffix  string
		Version string
	}

	var tests = []Test{
		{Name: "ReturnVersion", Prefix: "v", Suffix: "", Version: "0.0.0"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return(test.Version, nil)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix, GitAPI: gitAPI}

			var got, err = versionAPI.GetVersion()

			assert.NoError(t, err)
			assert.Equal(t, test.Version, got, `want: "%s, got: "%s"`, test.Version, got)
		})
	}

	type GitErrorTest struct {
		Error  error
		Name   string
		Prefix string
		Suffix string
	}

	var gitErrorTests = []GitErrorTest{
		{Name: "ReturnErrorOnGitError", Prefix: "v", Suffix: "", Error: fmt.Errorf("some-error")},
	}

	for _, test := range gitErrorTests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return("", test.Error)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix, GitAPI: gitAPI}

			var _, got = versionAPI.GetVersion()

			assert.Error(t, got)
			assert.Equal(t, test.Error, got, `want: "%s, got: "%s"`, test.Error, got)
		})
	}

	type SemverErrorTest struct {
		Error    error
		Name     string
		Prefix   string
		Suffix   string
		Versions string
	}

	var semverErrorTests = []SemverErrorTest{
		{Name: "ReturnErrorOnInvalidVersions", Versions: "invalid1 invalid2", Error: fmt.Errorf("could not find a valid semver version")},
		{Name: "ReturnErrorOnNoVersions", Versions: "", Error: fmt.Errorf("could not find a valid semver version")},
	}

	for _, test := range semverErrorTests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return(test.Versions, nil)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix, GitAPI: gitAPI}

			var _, got = versionAPI.GetVersion()

			assert.Error(t, got)
			assert.Equal(t, test.Error, got, `want: "%s, got: "%s"`, test.Error, got)
		})
	}
}

func TestAPI_GetVersionOrDefault(t *testing.T) {
	type Test struct {
		Name    string
		Prefix  string
		Suffix  string
		Version string
	}

	var tests = []Test{
		{Name: "ReturnVersionWithoutError", Prefix: "v", Suffix: "", Version: "0.0.0"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return(test.Version, nil)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix, GitAPI: gitAPI}

			var got, err = versionAPI.GetVersion()

			assert.NoError(t, err)
			assert.Equal(t, test.Version, got, `want: "%s, got: "%s"`, test.Version, got)
		})
	}

	type ErrorTest struct {
		Error  error
		Prefix string
		Suffix string
		Name   string
	}

	var errorTests = []ErrorTest{
		{Name: "ReturnDefaultVersionOnGitApiError", Prefix: "v", Suffix: "", Error: fmt.Errorf("some-error")},
	}

	for _, test := range errorTests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return("", test.Error)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix, GitAPI: gitAPI}

			var got = versionAPI.GetVersionOrDefault(cli.DefaultVersion)

			assert.Equal(t, cli.DefaultVersion, got, `want: "%s, got: "%s"`, cli.DefaultVersion, got)
		})
	}
}

func TestAPI_PredictVersion(t *testing.T) {
	type Test struct {
		Mode    modes.Mode
		Name    string
		Prefix  string
		Suffix  string
		Version string
		Want    string
	}

	var tests = []Test{
		{Name: "ReturnPatchPrediction", Prefix: "v", Suffix: "", Mode: modes.NewPatchMode(), Version: "0.0.0", Want: "0.0.1"},
		{Name: "ReturnMinorPrediction", Prefix: "v", Suffix: "", Mode: modes.NewMinorMode(), Version: "0.0.0", Want: "0.1.0"},
		{Name: "ReturnMajorPrediction", Prefix: "v", Suffix: "", Mode: modes.NewMajorMode(), Version: "0.0.0", Want: "1.0.0"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return(test.Version, nil)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix, GitAPI: gitAPI}

			var got, err = versionAPI.PredictVersion(test.Version, test.Mode)

			assert.NoError(t, err)
			assert.Equal(t, test.Want, got, `want: "%s, got: "%s"`, test.Want, got)
		})
	}

	type ErrorTest struct {
		Error   error
		Name    string
		Prefix  string
		Suffix  string
		Version string
	}

	var errorTests = []ErrorTest{
		{Name: "ReturnErrorOnModeIncrementError", Prefix: "v", Error: fmt.Errorf("some-error"), Version: "invalid"},
	}

	for _, test := range errorTests {
		t.Run(test.Name, func(t *testing.T) {
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix}

			var mode = mocks.NewMockMode()
			mode.On("Increment", mock.Anything, mock.Anything, mock.Anything).Return(test.Version, test.Error)

			var _, got = versionAPI.PredictVersion("0.0.0", mode)

			assert.Error(t, got)
			assert.Equal(t, test.Error, got, `want: "%s, got: "%s"`, test.Error, got)
		})
	}
}

func TestAPI_PushVersion(t *testing.T) {
	type Test struct {
		Mode    modes.Mode
		Name    string
		Prefix  string
		Suffix  string
		Version string
		Want    string
	}

	var tests = []Test{
		{Name: "PushWithPrefix", Mode: modes.NewPatchMode(), Prefix: "v", Version: "0.0.1", Want: "v0.0.1"},
		{Name: "PushWithoutPrefix", Mode: modes.NewPatchMode(), Prefix: "", Version: "0.0.1", Want: "0.0.1"},
		{Name: "PushWithSuffix", Mode: modes.NewPatchMode(), Prefix: "v", Suffix: "a", Version: "0.0.1", Want: "v0.0.1a"},
		{Name: "PushWithSuffixAlt", Mode: modes.NewPatchMode(), Prefix: "v", Suffix: "-alt", Version: "0.0.1", Want: "v0.0.1-alt"},
		{Name: "PushWithoutSuffix", Mode: modes.NewPatchMode(), Prefix: "", Suffix: "", Version: "0.0.1", Want: "0.0.1"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var gitAPI = fakes.NewFakeGitAPI()
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix, GitAPI: gitAPI}

			var err = versionAPI.PushVersion(test.Version)

			var pushedTags = versionAPI.GitAPI.(*fakes.FakeGitAPI).PushedTags
			var got = pushedTags[len(pushedTags)-1]

			assert.NoError(t, err)
			assert.Equal(t, test.Want, got, `want: "%s, got: "%s"`, test.Want, got)
		})
	}

	type ErrorTest struct {
		Error   error
		Name    string
		Version string
	}

	var errorTests = []ErrorTest{
		{Name: "ReturnErrorOnGitApiError", Error: fmt.Errorf("some-error"), Version: "invalid"},
	}

	for _, test := range errorTests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Run", mock.Anything, mock.Anything).Return(test.Error)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: "v", Suffix: "", GitAPI: gitAPI}

			var got = versionAPI.PushVersion("0.0.1")

			assert.Error(t, got)
			assert.Equal(t, test.Error, got, `want: "%s, got: "%s"`, test.Error, got)
		})
	}
}

func TestAPI_ReleaseVersion(t *testing.T) {
	type Test struct {
		Mode    modes.Mode
		Name    string
		Prefix  string
		Suffix  string
		Version string
		Want    string
	}

	var tests = []Test{
		{Name: "ReleaseWithPrefix", Mode: modes.NewPatchMode(), Prefix: "v", Version: "0.0.1", Want: "v0.0.1"},
		{Name: "ReleaseWithoutPrefix", Mode: modes.NewPatchMode(), Prefix: "", Version: "0.0.1", Want: "0.0.1"},
		{Name: "ReleaseWithSuffix", Mode: modes.NewPatchMode(), Prefix: "v", Suffix: "a", Version: "0.0.1", Want: "v0.0.1a"},
		{Name: "ReleaseWithSuffixAlt", Mode: modes.NewPatchMode(), Prefix: "v", Suffix: "-alt", Version: "0.0.1", Want: "v0.0.1-alt"},
		{Name: "ReleaseWithoutSuffix", Mode: modes.NewPatchMode(), Prefix: "", Suffix: "", Version: "0.0.1", Want: "0.0.1"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var gitAPI = fakes.NewFakeGitAPI()
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix, GitAPI: gitAPI}

			var err = versionAPI.ReleaseVersion(test.Version)

			var localTags = versionAPI.GitAPI.(*fakes.FakeGitAPI).LocalTags
			var got = localTags[len(localTags)-1]

			assert.NoError(t, err)
			assert.Equal(t, test.Want, got, `want: "%s, got: "%s"`, test.Want, got)
		})
	}

	type ErrorTest struct {
		Error   error
		Name    string
		Version string
	}

	var errorTests = []ErrorTest{
		{Name: "ReturnErrorOnGitApiError", Error: fmt.Errorf("some-error"), Version: "invalid"},
	}

	for _, test := range errorTests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Run", mock.Anything, mock.Anything).Return(test.Error)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: "v", Suffix: "", GitAPI: gitAPI}

			var got = versionAPI.ReleaseVersion("0.0.1")

			assert.Error(t, got)
			assert.Equal(t, test.Error, got, `want: "%s, got: "%s"`, test.Error, got)
		})
	}
}

func TestAPI_UpdateVersion(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		var gitAPI = fakes.NewFakeGitAPI()
		var versionAPI = API{GitAPI: gitAPI}

		var err = versionAPI.UpdateVersion()

		assert.NoError(t, err)
	})

	type ErrorTest struct {
		Error   error
		Name    string
		Version string
	}

	var errorTests = []ErrorTest{
		{Name: "ReturnErrorOnGitApiError", Error: fmt.Errorf("some-error"), Version: "invalid"},
	}

	for _, test := range errorTests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return("", test.Error)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{GitAPI: gitAPI}

			var got = versionAPI.UpdateVersion()

			assert.Error(t, got)
			assert.Equal(t, test.Error, got, `want: "%s, got: "%s"`, test.Error, got)
		})
	}
}

func TestNewAPI(t *testing.T) {
	t.Run("ValidateState", func(t *testing.T) {
		var api = NewAPI("v", "")
		assert.NotNil(t, api.GitAPI)
	})
}
