package repositories

import (
	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"context"
	"fmt"
	"strings"

	apierrors "code.cloudfoundry.org/korifi/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

//+kubebuilder:rbac:groups=korifi.cloudfoundry.org,resources=cfapps;cfbuilds;cfpackages;cfprocesses;cfspaces;cftasks,verbs=list
//+kubebuilder:rbac:groups=korifi.cloudfoundry.org,resources=cfdomains;cfroutes,verbs=list
//+kubebuilder:rbac:groups=korifi.cloudfoundry.org,resources=cfservicebindings;cfserviceinstances,verbs=list

type ListOptionBuilder interface {
	build(value string) metav1.ListOptions
}

type CFNameFieldSelector struct {
	ListOptionBuilder
}

func (f CFNameFieldSelector) build(value string) metav1.ListOptions {
	return metav1.ListOptions{
		FieldSelector: "metadata.name=" + value,
	}
}

type CFLabelFieldSelector struct {
	key string
	ListOptionBuilder
}

func (f CFLabelFieldSelector) build(value string) metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector: f.key + "=" + value,
	}
}

var (
	CFAppsGVR = schema.GroupVersionResource{
		Group:    "korifi.cloudfoundry.org",
		Version:  "v1alpha1",
		Resource: "cfapps",
	}

	CFBuildsGVR = schema.GroupVersionResource{
		Group:    "korifi.cloudfoundry.org",
		Version:  "v1alpha1",
		Resource: "cfbuilds",
	}

	CFDomainsGVR = schema.GroupVersionResource{
		Group:    "korifi.cloudfoundry.org",
		Version:  "v1alpha1",
		Resource: "cfdomains",
	}

	CFDropletsGVR = schema.GroupVersionResource{
		Group:    "korifi.cloudfoundry.org",
		Version:  "v1alpha1",
		Resource: "cfbuilds",
	}

	CFPackagesGVR = schema.GroupVersionResource{
		Group:    "korifi.cloudfoundry.org",
		Version:  "v1alpha1",
		Resource: "cfpackages",
	}

	CFProcessesGVR = schema.GroupVersionResource{
		Group:    "korifi.cloudfoundry.org",
		Version:  "v1alpha1",
		Resource: "cfprocesses",
	}

	CFRoutesGVR = schema.GroupVersionResource{
		Group:    "korifi.cloudfoundry.org",
		Version:  "v1alpha1",
		Resource: "cfroutes",
	}

	CFServiceBindingsGVR = schema.GroupVersionResource{
		Group:    "korifi.cloudfoundry.org",
		Version:  "v1alpha1",
		Resource: "cfservicebindings",
	}

	CFServiceInstancesGVR = schema.GroupVersionResource{
		Group:    "korifi.cloudfoundry.org",
		Version:  "v1alpha1",
		Resource: "cfserviceinstances",
	}

	CFOrgsGVR = schema.GroupVersionResource{
		Group:    "korifi.cloudfoundry.org",
		Version:  "v1alpha1",
		Resource: "cforgs",
	}

	CFSpacesGVR = schema.GroupVersionResource{
		Group:    "korifi.cloudfoundry.org",
		Version:  "v1alpha1",
		Resource: "cfspaces",
	}

	CFTasksGVR = schema.GroupVersionResource{
		Group:    "korifi.cloudfoundry.org",
		Version:  "v1alpha1",
		Resource: "cftasks",
	}

	ResourceMap = map[string]schema.GroupVersionResource{
		AppResourceType:             CFAppsGVR,
		BuildResourceType:           CFBuildsGVR,
		DropletResourceType:         CFDropletsGVR,
		DomainResourceType:          CFDomainsGVR,
		PackageResourceType:         CFPackagesGVR,
		ProcessResourceType:         CFProcessesGVR,
		RouteResourceType:           CFRoutesGVR,
		ServiceBindingResourceType:  CFServiceBindingsGVR,
		ServiceInstanceResourceType: CFServiceInstancesGVR,
		OrgResourceType:             CFOrgsGVR,
		SpaceResourceType:           CFSpacesGVR,
		TaskResourceType:            CFTasksGVR,
	}

	SelectorMap = map[string]ListOptionBuilder{
		AppResourceType:             CFNameFieldSelector{},
		BuildResourceType:           CFNameFieldSelector{},
		DropletResourceType:         CFNameFieldSelector{},
		DomainResourceType:          CFNameFieldSelector{},
		PackageResourceType:         CFNameFieldSelector{},
		ProcessResourceType:         CFNameFieldSelector{},
		RouteResourceType:           CFNameFieldSelector{},
		ServiceBindingResourceType:  CFNameFieldSelector{},
		ServiceInstanceResourceType: CFNameFieldSelector{},
		SpaceResourceType:           CFLabelFieldSelector{key: korifiv1alpha1.CFSpaceGUIDLabelKey},
		OrgResourceType:             CFLabelFieldSelector{key: korifiv1alpha1.CFOrgGUIDLabelKey},
		TaskResourceType:            CFNameFieldSelector{},
	}
)

type NamespaceRetriever struct {
	client dynamic.Interface
}

func NewNamespaceRetriever(client dynamic.Interface) NamespaceRetriever {
	return NamespaceRetriever{
		client: client,
	}
}

func (nr NamespaceRetriever) NameFor(ctx context.Context, resourceGUID, resourceType string) (string, error) {
	gvr, ok := ResourceMap[resourceType]
	if !ok {
		return "", fmt.Errorf("ResourceMap: resource type %q unknown", resourceType)
	}

	builder, ok := SelectorMap[resourceType]

	if !ok {
		return "", fmt.Errorf("SelectorMap: resource type %q unknown", resourceType)
	}

	opts := builder.build(resourceGUID)

	list, err := nr.client.Resource(gvr).List(ctx, opts)
	if err != nil {
		return "", fmt.Errorf("failed to list %v: %w", resourceType, apierrors.FromK8sError(err, resourceType))
	}

	if len(list.Items) == 0 {
		return "", apierrors.NewNotFoundError(fmt.Errorf("resource %q not found", resourceGUID), resourceType)
	}

	if len(list.Items) > 1 {
		return "", fmt.Errorf("get-%s duplicate records exist", strings.ToLower(resourceType))
	}

	metadata := list.Items[0].Object["metadata"].(map[string]interface{})

	ns := metadata["name"].(string)

	if ns == "" {
		return "", fmt.Errorf("get-%s: resource is not namespace-scoped", strings.ToLower(resourceType))
	}

	return ns, nil
}

func (nr NamespaceRetriever) NamespaceFor(ctx context.Context, resourceGUID, resourceType string) (string, error) {
	gvr, ok := ResourceMap[resourceType]
	if !ok {
		return "", fmt.Errorf("ResourceMap: resource type %q unknown", resourceType)
	}

	builder, ok := SelectorMap[resourceType]

	if !ok {
		return "", fmt.Errorf("SelectorMap: resource type %q unknown", resourceType)
	}

	opts := builder.build(resourceGUID)

	list, err := nr.client.Resource(gvr).List(ctx, opts)
	if err != nil {
		return "", fmt.Errorf("failed to list %v: %w", resourceType, apierrors.FromK8sError(err, resourceType))
	}

	if len(list.Items) == 0 {
		return "", apierrors.NewNotFoundError(fmt.Errorf("resource %q not found", resourceGUID), resourceType)
	}

	if len(list.Items) > 1 {
		return "", fmt.Errorf("get-%s duplicate records exist", strings.ToLower(resourceType))
	}

	metadata := list.Items[0].Object["metadata"].(map[string]interface{})

	ns := metadata["namespace"].(string)

	if ns == "" {
		return "", fmt.Errorf("get-%s: resource is not namespace-scoped", strings.ToLower(resourceType))
	}

	return ns, nil
}
