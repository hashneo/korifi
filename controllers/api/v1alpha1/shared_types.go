package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	CFLabelKeyPrefix     = "cloudfoundry.org/"
	KorifiLabelKeyPrefix = "korifi." + CFLabelKeyPrefix

	CFAppGUIDLabelKey       = KorifiLabelKeyPrefix + "app-guid"
	CFAppRevisionKey        = KorifiLabelKeyPrefix + "app-rev"
	CFAppRevisionKeyDefault = "0"
	CFPackageGUIDLabelKey   = KorifiLabelKeyPrefix + "package-guid"
	CFBuildGUIDLabelKey     = KorifiLabelKeyPrefix + "build-guid"
	CFProcessGUIDLabelKey   = KorifiLabelKeyPrefix + "process-guid"
	CFProcessTypeLabelKey   = KorifiLabelKeyPrefix + "process-type"
	CFDomainGUIDLabelKey    = KorifiLabelKeyPrefix + "domain-guid"
	CFRouteGUIDLabelKey     = KorifiLabelKeyPrefix + "route-guid"
	CFTaskGUIDLabelKey      = KorifiLabelKeyPrefix + "task-guid"
	CFOrgGUIDLabelKey       = KorifiLabelKeyPrefix + "org-guid"
	CFSpaceGUIDLabelKey     = KorifiLabelKeyPrefix + "space-guid"
	CFUserGUIDLabelKey      = KorifiLabelKeyPrefix + "user-guid"

	CFDefaultDomainLabelKey = KorifiLabelKeyPrefix + "default-domain"

	CFBindingTypeLabelKey = KorifiLabelKeyPrefix + "binding-type"

	StagingConditionType   = "Staging"
	ReadyConditionType     = "Ready"
	SucceededConditionType = "Succeeded"

	PropagateRoleBindingAnnotation    = CFLabelKeyPrefix + "propagate-cf-role"
	PropagateServiceAccountAnnotation = CFLabelKeyPrefix + "propagate-service-account"
	PropagatedFromLabel               = CFLabelKeyPrefix + "propagated-from"
)

type Lifecycle struct {
	// The CF Lifecycle type.
	// Only "buildpack" is currently allowed
	Type LifecycleType `json:"type"`
	// Data used to specify details for the Lifecycle
	Data LifecycleData `json:"data"`
}

// LifecycleType inform the platform of how to build droplets and run apps
// allow only values "buildpack"
// +kubebuilder:validation:Enum=buildpack
type LifecycleType string

// LifecycleData is shared by CFApp and CFBuild
type LifecycleData struct {
	// Buildpacks to include in auto-detection when building the app image.
	// If no values are specified, then all available buildpacks will be used for auto-detection
	Buildpacks []string `json:"buildpacks,omitempty"`

	// Stack to use when building the app image
	Stack string `json:"stack"`
}

// Registry is used by CFPackage and CFBuild/Droplet to identify Registry and secrets to access the image provided
type Registry struct {
	// The location of the source image
	Image string `json:"image"`
	// A list of secrets required to pull the image from its repository
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

// RequiredLocalObjectReference is a reference to an object in the same namespace.
// Unlike k8s.io/api/core/v1/LocalObjectReference, name is required.
type RequiredLocalObjectReference struct {
	Name string `json:"name"`
}
