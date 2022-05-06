[![Release](https://img.shields.io/github/v/release/wayfair-incubator/oss-template?display_name=tag)](CHANGELOG.md)
[![Lint](https://github.com/wayfair-incubator/oss-template/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/wayfair-incubator/oss-template/actions/workflows/lint.yml)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.0-4baaaa.svg)](CODE_OF_CONDUCT.md)
[![Maintainer](https://img.shields.io/badge/Maintainer-Wayfair-7F187F)](https://wayfair.github.io)

# k8s-used-api-versions

The `k8s-used-api-versions` Operator collects the API versions used by different components such as operators, controllers, etc., and export these API versions in a Prometheus metrics format. These metrics include the status of the used API versions such as whether the API version is deprecated, removed, etc. It also provides the same result via `kubectl`

## Motivation

According to the Kubernetes deprecation policy, Kubernetes might deprecate or remove API versions when launching a new release. The deprecated API versions should be upgraded to the new supported API versions before they became "removed", otherwise the component will be in a broken state.

When maintaining a lot of custom components such as Controllers, Operators, etc., you might forget to check and upgrade the hardcoded API versions. This is why this operator is developed, it allows you to expose the used (hardcoded) API versions via custom resource **UsedApiVersions** and the operator will check the status of these API versions and export them via Prometheus metrics format. You can integrate it with Prometheus and Alertmanager to receive notifications before upgrading the cluster and break these components.

## Setup

The Operator can be run as a deployment in the cluster. See [deployment.yaml](config/manager/manager.yaml) for an example.

## Usage

After deploying the operator, you just need to create the custom resource with the used API versions.

  - Example of the custom resource: [UsedApiVersions](./config/samples/api-version_v1beta1_usedapiversions.yaml)

  - Example of the exported metrics:

```sh
wf_operator_used_api_versions{api_version="apps/v1beta2",deprecated="false",deprecated_in_version="n/a",kind="ReplicaSet",name="ingress-operator",removed="true",removed_in_next_2_releases="true",removed_in_next_release="true",removed_in_version="v1.16.0",replacement_api="apps/v1"} 1
wf_operator_used_api_versions{api_version="extensions/v1beta1",deprecated="true",deprecated_in_version="v1.14.0",kind="Ingress",name="ingress-operator",removed="false",removed_in_next_2_releases="false",removed_in_next_release="false",removed_in_version="v1.22.0",replacement_api="networking.k8s.io/v1"} 1
```

The operator will update the status of the custom resource, so you can get the same result via `kubectl`

```sh
$ kubectl get UsedApiVersions -n ingress ingress-operator -oyaml

status:
  apiVersionsStatus:
  - apiVersion: extensions/v1beta1
    deprecated: true
    deprecatedInVersion: v1.14.0
    kind: Ingress
    removed: false
    removedInNextRelease: false
    removedInNextTwoReleases: false
    removedInVersion: v1.22.0
    replacementApi: networking.k8s.io/v1
```

Also, you can get a quick overview of all the deployed components

```sh
$ kubectl get UsedApiVersions
NAME              KIND              AGE    DEPRECATED   REMOVED
example-operator  UsedApiVersions   27h    3            1
ns-controller     UsedApiVersions   27h    1            0
ingress-operator  UsedApiVersions   137m   1            0
```

## Configuration

These command line arguments are available

``--metrics-bind-address``
    The address the metric endpoint binds to (Default: `:8080`)

``--health-probe-bind-address``
    The address the probe endpoint binds to (Default: `:8081`)

``--leader-elect``
    Enable leader election for controller manager (Default: `false`).
    Enabling this will ensure there is only one active controller manager

``--versions-file``
    The versions file used to check deprecations (Default: `config/versions.yaml`)

## Development

This Operator was developed using Kubebuilder, so it's highly recommended not to update the CRD manually. You can use the kubebuilder markers to do the changes, then run

```sh
make install 
make run
```

## Roadmap

See the [open issues](https://github.com/org_name/repo_name/issues) for a list of proposed features (and known issues).

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**. For detailed contributing guidelines, please see [CONTRIBUTING.md](CONTRIBUTING.md)

## License

Distributed under the `<License name>` License. See `LICENSE` for more information.

## Contact

Your Name - [@twitter_handle](https://twitter.com/twitter_handle) - email

Project Link: [https://github.com/org_name/repo_name](https://github.com/org_name/repo_name)

## Acknowledgements

This template was adapted from
[https://github.com/othneildrew/Best-README-Template](https://github.com/othneildrew/Best-README-Template).
