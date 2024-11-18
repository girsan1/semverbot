package semver

import (
	"fmt"
	"regexp"
	"strings"

	blangsemver "github.com/blang/semver/v4"
)

// Find finds the biggest valid semver version in a slice of strings.
// The initial order of the versions does not matter.
// Returns the biggest valid semver version if found, otherwise an error stating no valid semver version has been found.
func Find(prefix string, suffix string, versions []string) (found string, err error) {
	var parsedVersions blangsemver.Versions
	var parsedVersion blangsemver.Version

	var validVersion = regexp.MustCompile(`^\d+\.\d+\.\d+$`) //to allow simple version like 0.23.234
	var filteredVersions []string
	for _, version := range versions {
		if prefix == "v" && validVersion.MatchString(version) || strings.Contains(version, prefix) && strings.Contains(version, suffix) {
			filteredVersions = append(filteredVersions, version)
		}
	}

	for _, version := range filteredVersions {
		if parsedVersion, err = Parse(prefix, suffix, version); err != nil {
			continue
		}

		parsedVersions = append(parsedVersions, parsedVersion)
	}

	if len(parsedVersions) == 0 {
		return found, fmt.Errorf("could not find a valid semver version")
	}

	blangsemver.Sort(parsedVersions)

	var targetVersion = parsedVersions[len(parsedVersions)-1]

	// necessary because blangsemver's Version.String() strips any prefix
	for _, version := range filteredVersions {
		if strings.Contains(version, targetVersion.String()) {
			found = version
			break
		}
	}

	return found, nil
}
