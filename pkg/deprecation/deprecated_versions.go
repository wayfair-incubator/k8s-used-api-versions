/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package deprecation

import (
	"io/ioutil"
	"strconv"

	semver "github.com/hashicorp/go-version"
	"sigs.k8s.io/yaml"
)

// getDeprecatedVersions gets the deprecated apiVersions from the versions file
func getDeprecatedVersions(versionsFile string) (*Versions, error) {
	v, err := ioutil.ReadFile(versionsFile)
	if err != nil {
		return nil, err
	}

	deprecatedVersions := new(Versions)
	err = yaml.Unmarshal(v, deprecatedVersions)
	return deprecatedVersions, err
}

// isNewerOrEqualVersion compares two semVersions and checks if the first version
// is equal or greater than the second version.
// The first version is the kubernetes cluster version
// The second version is the version from versions.yaml file
func isNewerOrEqualVersion(k8sVersion, yamlFileVersion string) (bool, error) {
	v1, err := semver.NewVersion(k8sVersion)
	if err != nil {
		return false, err
	}
	v2, err := semver.NewVersion(yamlFileVersion)
	if err != nil {
		return false, err
	}
	if v1.GreaterThanOrEqual(v2) {
		return true, nil
	}
	return false, nil
}

// incrementSemVer increments the provided semVersion by specific value
func incrementSemVer(version string, steps int) (string, error) {
	v1, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	versionParts := v1.Segments()
	major, minor, patch := versionParts[0], versionParts[1], versionParts[2]
	minor = minor + steps
	incrementedSemVer := strconv.Itoa(major) + "." + strconv.Itoa(minor) + "." + strconv.Itoa(patch)

	return incrementedSemVer, nil

}

// isDeprecatedVersion checks if the provided apiVersion of specific kind is deprecated
// based on the current k8s version and the deprecation file "versions.yaml"
func isDeprecatedVersion(kind, apiVersion, k8sVersion, versionsFile string) bool {
	// var v *Versions
	v, _ := getDeprecatedVersions(versionsFile)
	// fmt.Println(v.DeprecatedVersions)
	for _, dep := range v.DeprecatedVersions {
		if kind == dep.Kind && apiVersion == dep.APIVersion {
			if dep.DeprecatedInVersion == "" {
				return false
			}
			n, _ := isNewerOrEqualVersion(k8sVersion, dep.DeprecatedInVersion)
			if n {
				return true
			}
		}
	}
	return false
}

// isRemovedVersion checks if the provided apiVersion of specific kind is removed
// based on the current k8s version and the deprecation file "versions.yaml"
func isRemovedVersion(kind, apiVersion, k8sVersion, versionsFile string) bool {
	// var v *Versions
	v, _ := getDeprecatedVersions(versionsFile)
	for _, dep := range v.DeprecatedVersions {

		if kind == dep.Kind && apiVersion == dep.APIVersion {
			if dep.RemovedInVersion == "" {
				return false
			}
			n, _ := isNewerOrEqualVersion(k8sVersion, dep.RemovedInVersion)
			if n {
				return true
			}
		}
	}
	return false
}

// getDeprecatedKindInfo gets some information about the deprecated or removed apiVersion of specific kind such as:
// replacementApi: The new apiVersion that should be used instead of the current deprecated or removed apiVersion
// deprecated_in_version: The apiVersion was deprecated in which k8s version.
// removed_in_version: The apiVersion was removed in which k8s version
func getDeprecatedKindInfo(kind, apiVersion, versionsFile string) map[string]string {
	// var v *Versions
	var replacementApi, removedInVersion, deprecatedInVersion string
	result := make(map[string]string)
	v, _ := getDeprecatedVersions(versionsFile)
	for _, dep := range v.DeprecatedVersions {
		if kind == dep.Kind && apiVersion == dep.APIVersion {
			if dep.ReplacementAPI == "" {
				replacementApi = "n/a"
			} else {
				replacementApi = dep.ReplacementAPI
			}

			if dep.RemovedInVersion == "" {
				removedInVersion = "n/a"
			} else {
				removedInVersion = dep.RemovedInVersion
			}

			if dep.DeprecatedInVersion == "" {
				deprecatedInVersion = "n/a"
			} else {
				deprecatedInVersion = dep.DeprecatedInVersion
			}
		}
	}
	result["replacementApi"] = replacementApi
	result["removedInVersion"] = removedInVersion
	result["deprecatedInVersion"] = deprecatedInVersion

	return result
}

// CheckDeprecations is the main function used to check the overall deprecation status.
func CheckDeprecations(kind, apiVersion, k8sVersion, versionsFile string) map[string]string {
	result := make(map[string]string)
	deprecated := isDeprecatedVersion(kind, apiVersion, k8sVersion, versionsFile)
	removed := isRemovedVersion(kind, apiVersion, k8sVersion, versionsFile)
	deprecatedKindInfo := getDeprecatedKindInfo(kind, apiVersion, versionsFile)
	nextVersion, _ := incrementSemVer(k8sVersion, 1)
	nextTwoVersion, _ := incrementSemVer(k8sVersion, 2)
	removedInNextRelease := isRemovedVersion(kind, apiVersion, nextVersion, versionsFile)
	removedInNextTwoReleases := isRemovedVersion(kind, apiVersion, nextTwoVersion, versionsFile)
	result["deprecated"] = strconv.FormatBool(deprecated)
	result["removed"] = strconv.FormatBool(removed)
	result["replacementApi"] = deprecatedKindInfo["replacementApi"]
	result["removedInVersion"] = deprecatedKindInfo["removedInVersion"]
	result["deprecatedInVersion"] = deprecatedKindInfo["deprecatedInVersion"]
	result["kind"] = kind
	result["apiVersion"] = apiVersion
	result["removedInNextRelease"] = strconv.FormatBool(removedInNextRelease)
	result["removedInNextTwoReleases"] = strconv.FormatBool(removedInNextTwoReleases)

	return result
}
