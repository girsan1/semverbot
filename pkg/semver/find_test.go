package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	type Test struct {
		Name      string
		Prefix    string
		Suffix    string
		Versions  []string
		WantIndex int
	}

	var tests = []Test{
		{Name: "FindVersionIfValid", Prefix: "v", Suffix: "", Versions: []string{"v1.0.1", "v0.1.1", "v0.1.0"}, WantIndex: 0},
		{Name: "FindVersionWithCustomPrefixIfValid", Prefix: "test-", Suffix: "", Versions: []string{"test-1.0.1", "test-0.1.1", "test-0.1.0"}, WantIndex: 0},
		{Name: "FindVersionWithCustomSuffixIfValid", Prefix: "v", Suffix: "a", Versions: []string{"v1.0.1a", "v0.1.1a", "v0.1.0a"}, WantIndex: 0},
		{Name: "SkipVersionIfInvalid", Prefix: "v", Suffix: "", Versions: []string{"invalid1", "invalid2", "v0.1.0"}, WantIndex: 2},
		{Name: "FindVersionWhenDifferentOrder", Prefix: "v", Suffix: "", Versions: []string{"v1.3.1", "v0.2.0", "v2.3.0"}, WantIndex: 2},
		{Name: "FindVersionWhenMultiplePrefixes", Prefix: "v", Suffix: "", Versions: []string{"v1.3.1", "v0.2.0", "2.3.0"}, WantIndex: 2},
		{Name: "FindVersionWhenMultipleSuffixes", Prefix: "v", Suffix: "n", Versions: []string{"v1.3.1a", "v0.2.0-alt", "v2.3.0n"}, WantIndex: 2},
		{Name: "FindVersionWhenMultiplePrefixesWithSameVersion", Prefix: "c", Suffix: "", Versions: []string{"v1.3.1", "v0.2.0", "c1.3.1"}, WantIndex: 2},
		{Name: "FindVersionWhenMultipleSuffixesWithSameVersion", Prefix: "", Suffix: "-s", Versions: []string{"1.3.1-n", "v0.2.0-s", "1.3.1-s"}, WantIndex: 2},
		{Name: "FindVersionWhenMultiplePrefixesWithSameVersion2", Prefix: "egress-admarketplace-", Suffix: "", Versions: []string{
			//"egress-admarketplace-0.0.0",
			//"egress-admarketplace-0.1.0",
			"egress-claravine-0.0.0",
			"egress-claravine-0.0.1",
			"egress-claravine-0.0.10",
			"egress-claravine-0.0.11",
			"egress-claravine-0.0.12",
			"egress-claravine-0.0.13",
			"egress-claravine-0.0.14",
			"egress-claravine-0.0.15",
			"egress-claravine-0.0.16",
			"egress-claravine-0.0.17",
			"egress-claravine-0.0.18",
			"egress-claravine-0.0.19",
			"egress-claravine-0.0.2",
			"egress-claravine-0.0.20",
			"egress-claravine-0.0.21",
			"egress-claravine-0.0.22",
			"egress-claravine-0.0.23",
			"egress-claravine-0.0.24",
			"egress-claravine-0.0.3",
			"egress-claravine-0.0.4",
			"egress-claravine-0.0.5",
			"egress-claravine-0.0.6",
			"egress-claravine-0.0.7",
			"egress-claravine-0.0.8",
			"egress-claravine-0.0.9",
			"egress-cp-code-kpi-0.0.0",
			"egress-cp-code-kpi-0.0.1",
			"egress-cp-code-kpi-0.0.2",
			"egress-cp-code-kpi-0.0.3",
			"egress-cp-code-kpi-0.0.4",
			"egress-cp-code-kpi-0.0.5",
			"egress-cp-code-kpi-0.0.6",
			"egress-store-conversion-0.1.0",
			"egress-store-conversion-0.1.1",
			"egress-store-conversion-0.2.0",
			"egress-store-conversion-0.2.1",
			"egress-store-conversion-0.3.0",
			"egress-store-conversion-0.3.1",
			"egress-store-conversion-0.3.2",
			"egress-store-conversion-0.3.3",
			"egress-store-conversion-0.4.0",
			"egress-store-conversion-0.4.1",
			"egress-store-conversion-0.4.2",
			"egress-store-conversion-0.4.3",
			"egress-store-conversion-0.5.0",
			"egress-store-conversion-0.5.1",
			"egress-store-conversion-0.5.2",
			"egress-store-conversion-0.5.3",
			"egress-store-conversion-0.6.0",
			"egress-store-conversion-0.7.0",
			"egress-store-conversion-0.8.0",
			"egress-admarketplace-0.0.0",
			"egress-admarketplace-0.1.0"}, WantIndex: 52},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var versions = test.Versions
			var want = versions[test.WantIndex]
			var got, err = Find(test.Prefix, test.Suffix, versions)

			assert.Equal(t, want, got, `want: "%s", got: "%s"`, want, got)
			assert.NoError(t, err)
		})
	}

	type ErrorTest struct {
		Name     string
		Versions []string
	}

	var errorTests = []ErrorTest{
		{Name: "ReturnErrorOnInvalidVersions", Versions: []string{"invalid", "semver", "versions"}},
		{Name: "ReturnErrorOnNoVersions", Versions: []string{}},
	}

	for _, test := range errorTests {
		t.Run(test.Name, func(t *testing.T) {
			var _, got = Find("v", "a", test.Versions)
			assert.Error(t, got)
		})
	}
}
