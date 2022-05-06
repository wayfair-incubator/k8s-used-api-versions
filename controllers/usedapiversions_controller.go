/*
Copyright 2021.

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

package controllers

import (
	"context"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/runtime"
	discovery "k8s.io/client-go/discovery"
	restclient "k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	apiversionv1beta1 "github.com/wayfair-incubator/k8s-used-api-versions/api/v1beta1"
	"github.com/wayfair-incubator/k8s-used-api-versions/pkg/deprecation"
)

var (
	usedApiVersionsInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wf_operator_used_api_versions",
			Help: "The status of the used API versions",
		},
		[]string{"name",
			"used_api_versions_namespace",
			"kind",
			"api_version",
			"deprecated",
			"removed",
			"replacement_api",
			"removed_in_version",
			"deprecated_in_version",
			"removed_in_next_release",
			"removed_in_next_2_releases"},
	)
)

// UsedApiVersionsReconciler reconciles a UsedApiVersions object
type UsedApiVersionsReconciler struct {
	client.Client
	Log          logr.Logger
	Scheme       *runtime.Scheme
	ClientConfig *restclient.Config
	VersionsFile string
}

// NewUsedApiVersionsReconciler creates a new UsedApiVersionsReconciler.
func NewUsedApiVersionsReconciler(
	cli client.Client,
	log logr.Logger,
	scheme *runtime.Scheme,
) *UsedApiVersionsReconciler {
	metrics.Registry.MustRegister(
		usedApiVersionsInfo,
	)

	return &UsedApiVersionsReconciler{
		Client: cli,
		Log:    log,
		Scheme: scheme,
	}
}

//+kubebuilder:rbac:groups=api-version.wayfair.com,resources=usedapiversions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api-version.wayfair.com,resources=usedapiversions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api-version.wayfair.com,resources=usedapiversions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the UsedApiVersions object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *UsedApiVersionsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	r.updateUsedApiVersionsMetrics(ctx, log)
	var usedApiVersions apiversionv1beta1.UsedApiVersions
	if err := r.Get(ctx, req.NamespacedName, &usedApiVersions); err != nil {
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	k8sVersion, _ := r.getKubernetesVersion(log)
	var usedAPIStatus []apiversionv1beta1.APIVersionStatus
	for _, apiVersionMeta := range usedApiVersions.Spec.UsedApiVersions {
		var usedAPI apiversionv1beta1.APIVersionStatus
		usedAPI = getUsedAPIVersionsStatus(apiVersionMeta.Kind, apiVersionMeta.APIVersion, k8sVersion, r.VersionsFile)
		usedAPIStatus = append(usedAPIStatus, usedAPI)
	}

	updateFinalStatus(usedAPIStatus, &usedApiVersions)
	usedApiVersions.Status.ApiVersionsStatus = usedAPIStatus

	if err := r.Status().Update(ctx, &usedApiVersions); err != nil {
		log.Error(err, "unable to update usedApiVersions Status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// updateFinalStatus updates the finalStatus struct fields based on the deprecation status
func updateFinalStatus(usedAPIStatus []apiversionv1beta1.APIVersionStatus, usedApiVersions *apiversionv1beta1.UsedApiVersions) {

	// Reset values to zero
	usedApiVersions.Status.FinalStatus.Deprecated = 0
	usedApiVersions.Status.FinalStatus.Removed = 0
	usedApiVersions.Status.FinalStatus.RemovedInNextRelease = 0
	usedApiVersions.Status.FinalStatus.RemovedInNextTwoReleases = 0

	for _, s := range usedAPIStatus {
		if s.Deprecated == true {
			usedApiVersions.Status.FinalStatus.Deprecated += 1
		}
		if s.Removed == true {
			usedApiVersions.Status.FinalStatus.Removed += 1
		}
		if s.RemovedInNextRelease == true {
			usedApiVersions.Status.FinalStatus.RemovedInNextRelease += 1
		}
		if s.RemovedInNextTwoReleases == true {
			usedApiVersions.Status.FinalStatus.RemovedInNextTwoReleases += 1
		}
	}
}

// getUsedAPIVersionsStatus returns the overall deprecation status.
func getUsedAPIVersionsStatus(kind, apiVersion, k8sVersion, VersionsFile string) (apiVersionStatus apiversionv1beta1.APIVersionStatus) {
	deprecations := deprecation.CheckDeprecations(kind, apiVersion, k8sVersion, VersionsFile)
	apiVersionStatus.APIVersion = apiVersion
	apiVersionStatus.Kind = kind
	apiVersionStatus.Deprecated, _ = strconv.ParseBool(deprecations["deprecated"])
	apiVersionStatus.Removed, _ = strconv.ParseBool(deprecations["removed"])
	apiVersionStatus.DeprecatedInVersion = deprecations["deprecatedInVersion"]
	apiVersionStatus.RemovedInVersion = deprecations["removedInVersion"]
	apiVersionStatus.ReplacementAPI = deprecations["replacementApi"]
	apiVersionStatus.RemovedInNextRelease, _ = strconv.ParseBool(deprecations["removedInNextRelease"])
	apiVersionStatus.RemovedInNextTwoReleases, _ = strconv.ParseBool(deprecations["removedInNextRelease"])

	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *UsedApiVersionsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiversionv1beta1.UsedApiVersions{}).
		Complete(r)
}

// updateUsedApiVersionsMetrics updates and export metrics for all the UsedApiVersions kinds.
func (r *UsedApiVersionsReconciler) updateUsedApiVersionsMetrics(ctx context.Context, log logr.Logger) {
	var usedApiVersionsList apiversionv1beta1.UsedApiVersionsList
	err := r.Client.List(ctx, &usedApiVersionsList)
	if err != nil {
		log.Error(err, "error in collecting used apiVersions metrics")
		return
	}
	k8sVersion, _ := r.getKubernetesVersion(log)

	usedApiVersionsInfo.Reset()
	for _, u := range usedApiVersionsList.Items {
		for _, apiVersionMeta := range u.Spec.UsedApiVersions {
			deprecations := deprecation.CheckDeprecations(apiVersionMeta.Kind, apiVersionMeta.APIVersion, k8sVersion, "config/versions.yaml")
			usedApiVersionsInfo.With(prometheus.Labels{
				"name":                        u.Name,
				"used_api_versions_namespace": u.Namespace,
				"kind":                        apiVersionMeta.Kind,
				"api_version":                 apiVersionMeta.APIVersion,
				"deprecated":                  deprecations["deprecated"],
				"removed":                     deprecations["removed"],
				"replacement_api":             deprecations["replacementApi"],
				"removed_in_version":          deprecations["removedInVersion"],
				"deprecated_in_version":       deprecations["deprecatedInVersion"],
				"removed_in_next_release":     deprecations["removedInNextRelease"],
				"removed_in_next_2_releases":  deprecations["removedInNextTwoReleases"],
			}).Set(1)
		}
	}
	log.Info("Updated used apiVersions metrics.")
}

func (r *UsedApiVersionsReconciler) getKubernetesVersion(log logr.Logger) (string, error) {
	discoveryClient, _ := discovery.NewDiscoveryClientForConfig(r.ClientConfig)
	kubeVersion, err := discoveryClient.ServerVersion()

	if err != nil {
		log.Error(err, "failed to get Kubernetes server version.")
		return "", err
	}
	return kubeVersion.String(), nil
}

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(usedApiVersionsInfo)
}
