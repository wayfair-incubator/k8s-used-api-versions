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

// Version describes the deprecated API version for different kinds
type Version struct {
	// APIVersion is the name of the API version used by specific kind.
	APIVersion string `json:"version" yaml:"version"`
	// Kind is the Object type such as "Deployment" or "Ingress"
	Kind string `json:"kind" yaml:"kind"`
	// Kubernetes version in which the API version is deprecated in
	DeprecatedInVersion string `json:"deprecatedInVersion" yaml:"deprecatedInVersion"`
	// Kubernetes version in which the API version is removed in
	RemovedInVersion string `json:"removedInVersion" yaml:"removedInVersion"`
	// ReplacementAPI is the new supported API version
	ReplacementAPI string `json:"replacementApi" yaml:"replacementApi"`
}
