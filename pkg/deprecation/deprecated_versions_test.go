package deprecation

import (
	"reflect"
	"testing"
)

var versionsFile string = "../../config/versions.yaml"

func TestIsNewerOrEqualVersio(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"v1.7rc2", "1.22.0", false},
		{"1.0-7rc2", "v1.22.0-alpha.1", false},
		{"1.13.0", "1.14.0", false},
		{"v1.12.0", "1.14.0", false},
		{"v1.9.0-alpha.1", "1.10.0", false},
		{"v1.25.0-alpha.3", "1.23.0", true},
		{"v1.26.0-rc.0", "1.22.0", true},
		{"v1.22.0-alpha.1", "1.25.0", false},
		{"v1.22.0-alpha.1", "1.2.0-x.Y.0+metadata", true},
	}

	for _, c := range cases {
		var got bool
		got, _ = isNewerOrEqualVersion(c.v1, c.v2)
		if got != c.expected {
			t.Fatalf("Expected version: %v to be newer or equal than version: %v", c.v1, c.v2)

		}
	}

}

func TestIncrementSemVer(t *testing.T) {
	versions := []struct {
		version            string
		steps              int
		incrementedVersion string
	}{
		{"1.13.0", 1, "1.14.0"},
		{"v1.12.0", 2, "1.14.0"},
		{"v1.9.0-alpha.1", 1, "1.10.0"},
		{"v1.20.0-alpha.3", 3, "1.23.0"},
		{"v1.20.0-rc.0", 2, "1.22.0"},
		{"v1.22.0-alpha.1", 3, "1.25.0"},
	}

	for _, v := range versions {
		var got string
		got, _ = incrementSemVer(v.version, v.steps)
		if got != v.incrementedVersion {
			t.Fatalf("Version: %v is expected to be: %v after incrementing by %d steps, but got: %v", v.version, v.incrementedVersion, v.steps, got)

		}
	}
}

func TestIsDeprecatedVersion(t *testing.T) {
	apis := []struct {
		kind       string
		apiVersion string
		k8sVersion string
		status     bool
	}{
		{"Deployment", "extensions/v1beta1", "v1.12.0", true},
		{"Deployment", "extensions/v1beta1", "v1.8.0", false},
		{"Deployment", "apps/v1beta2", "v1.10.0", true},
		{"StatefulSet", "apps/v1beta1", "v1.10.0", true},
		{"NetworkPolicy", "extensions/v1beta1", "v1.10.0", true},
		{"Ingress", "extensions/v1beta1", "v1.13.0", false},
		{"Ingress", "extensions/v1beta1", "v1.14.0", true},
		{"Ingress", "networking.k8s.io/v1beta1", "v1.19.0", true},
		{"ReplicaSet", "extensions/v1beta1", "v1.17.0", false},
		{"PriorityClass", "scheduling.k8s.io/v1beta1", "v1.17.0", true},
	}

	for _, api := range apis {
		var got bool
		got = isDeprecatedVersion(api.kind, api.apiVersion, api.k8sVersion, versionsFile)
		if got != api.status {
			t.Fatalf("The API Version: %v deprecation status is: %v. Expected: %v", api.apiVersion, api.status, got)

		}
	}

}

func TestIsRemovedVersion(t *testing.T) {
	apis := []struct {
		kind       string
		apiVersion string
		k8sVersion string
		status     bool
	}{
		{"Deployment", "extensions/v1beta1", "v1.12.0", false},
		{"Deployment", "apps/v1beta2", "v1.16.0", true},
		{"StatefulSet", "apps/v1beta1", "v1.16.0", true},
		{"NetworkPolicy", "extensions/v1beta1", "v1.16.0", true},
		{"NetworkPolicy", "extensions/v1beta1", "v1.15.0", false},
		{"Ingress", "extensions/v1beta1", "v1.21.0", false},
		{"Ingress", "extensions/v1beta1", "v1.22.0", true},
		{"Ingress", "networking.k8s.io/v1beta1", "v1.20.0", false},
		{"ReplicaSet", "extensions/v1beta1", "v1.17.0", true},
		{"PriorityClass", "scheduling.k8s.io/v1beta1", "v1.17.0", true},
		{"PodDisruptionBudgetList", "policy/v1beta1", "v1.17.0", false},
	}

	for _, api := range apis {
		var got bool
		got = isRemovedVersion(api.kind, api.apiVersion, api.k8sVersion, versionsFile)
		if got != api.status {
			t.Fatalf("The API Version: %v removal status is: %v. Expected: %v", api.apiVersion, api.status, got)

		}
	}

}

func TestGetDeprecatedKindInfo(t *testing.T) {
	versions := []struct {
		kind       string
		apiVersion string
		expected   map[string]string
	}{
		{
			"Deployment",
			"extensions/v1beta1",
			map[string]string{
				"replacementApi":      "apps/v1",
				"removedInVersion":    "v1.16.0",
				"deprecatedInVersion": "v1.9.0",
			},
		},
		{
			"StatefulSet",
			"apps/v1beta1",
			map[string]string{
				"replacementApi":      "apps/v1",
				"removedInVersion":    "v1.16.0",
				"deprecatedInVersion": "v1.9.0",
			},
		},
		{
			"PodDisruptionBudget",
			"policy/v1beta1",
			map[string]string{
				"replacementApi":      "n/a",
				"removedInVersion":    "n/a",
				"deprecatedInVersion": "v1.22.0",
			},
		},
	}

	for _, v := range versions {
		var got map[string]string
		got = getDeprecatedKindInfo(v.kind, v.apiVersion, versionsFile)
		if !reflect.DeepEqual(got, v.expected) {
			t.Fatalf("The API Version info: %v doesn't match the expected result, \nExpected: %v. ", got, v.expected)

		}
	}
}

func TestCheckDeprecations(t *testing.T) {
	versions := []struct {
		kind       string
		apiVersion string
		k8sVersion string
		expected   map[string]string
	}{
		{
			"Deployment",
			"extensions/v1beta1",
			"v1.16.0",
			map[string]string{
				"apiVersion":               "extensions/v1beta1",
				"kind":                     "Deployment",
				"deprecated":               "true",
				"removed":                  "true",
				"replacementApi":           "apps/v1",
				"removedInVersion":         "v1.16.0",
				"deprecatedInVersion":      "v1.9.0",
				"removedInNextRelease":     "true",
				"removedInNextTwoReleases": "true",
			},
		},

		{
			"Ingress",
			"extensions/v1beta1",
			"v1.19.0",
			map[string]string{
				"apiVersion":               "extensions/v1beta1",
				"kind":                     "Ingress",
				"deprecated":               "true",
				"removed":                  "false",
				"replacementApi":           "networking.k8s.io/v1",
				"removedInVersion":         "v1.22.0",
				"deprecatedInVersion":      "v1.14.0",
				"removedInNextRelease":     "false",
				"removedInNextTwoReleases": "false",
			},
		},
	}

	for _, v := range versions {
		var got map[string]string
		got = CheckDeprecations(v.kind, v.apiVersion, v.k8sVersion, versionsFile)
		if !reflect.DeepEqual(got, v.expected) {
			t.Fatalf("CheckDeprecation \nExpected: %v doesn't match the expected result, \nExpected: %v. ", got, v.expected)

		}
	}
}
