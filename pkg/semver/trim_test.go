package semver

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrim(t *testing.T) {
	type Test struct {
		Name     string
		Major    string
		Minor    string
		Patch    string
		Prebuild string
		Prefix   string
		Suffix   string
	}

	var tests = []Test{
		{Name: "Default", Major: "0", Minor: "0", Patch: "0"},
		{Name: "Patch", Major: "0", Minor: "0", Patch: "1"},
		{Name: "Minor", Major: "0", Minor: "2", Patch: "0"},
		{Name: "Major", Major: "3", Minor: "0", Patch: "0"},
		{Name: "DiscardPrefix", Major: "1", Minor: "0", Patch: "0", Prefix: "v"},
		{Name: "DiscardSuffix", Major: "1", Minor: "0", Patch: "0", Suffix: "a"},
		{Name: "DiscardSuffixAlt", Major: "1", Minor: "0", Patch: "0", Suffix: "-alt"},
		{Name: "DiscardPrebuild", Major: "2", Minor: "0", Patch: "0", Prebuild: "-pre+001"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var version = fmt.Sprintf(`%s%s.%s.%s%s%s`, test.Prefix, test.Major, test.Minor, test.Patch,
				test.Suffix, test.Prebuild)

			var want = strings.ReplaceAll(version, test.Prefix, "")
			want = strings.ReplaceAll(want, test.Suffix, "")
			want = strings.ReplaceAll(want, test.Prebuild, "")

			var got, err = Trim(test.Prefix, test.Suffix, version)

			assert.Equal(t, want, got, `want: "%s", got: "%s"`, want, got)

			if test.Prefix != "" {
				assert.False(t, strings.HasPrefix(got, test.Prefix))
			}

			if test.Suffix != "" {
				assert.False(t, strings.HasSuffix(got, test.Suffix))
			}

			if test.Prebuild != "" {
				assert.False(t, strings.HasSuffix(got, test.Prebuild))
			}

			assert.NoError(t, err)
		})
	}

	type ErrorTest struct {
		Name    string
		Version string
	}

	var errorTests = []ErrorTest{
		{Name: "ReturnErrorOnInvalidVersion", Version: "invalid"},
	}

	for _, test := range errorTests {
		t.Run(test.Name, func(t *testing.T) {
			var _, got = Trim("v", "a", test.Version)
			assert.Error(t, got)
		})
	}
}
