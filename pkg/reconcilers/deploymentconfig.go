package reconcilers

import (
	"fmt"
	"reflect"

	"github.com/3scale/3scale-operator/pkg/common"
	"github.com/3scale/3scale-operator/pkg/helper"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	appsv1 "github.com/openshift/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func DeploymentConfigResourcesMutator(existingObj, desiredObj common.KubernetesObject) (bool, error) {
	existing, ok := existingObj.(*appsv1.DeploymentConfig)
	if !ok {
		return false, fmt.Errorf("%T is not a *appsv1.DeploymentConfig", existingObj)
	}
	desired, ok := desiredObj.(*appsv1.DeploymentConfig)
	if !ok {
		return false, fmt.Errorf("%T is not a *appsv1.DeploymentConfig", desiredObj)
	}

	return DeploymentConfigContainerResourcesReconciler(desired, existing), nil
}

func DeploymentConfigResourcesAndAffinityAndTolerationsMutator(existingObj, desiredObj common.KubernetesObject) (bool, error) {
	existing, ok := existingObj.(*appsv1.DeploymentConfig)
	if !ok {
		return false, fmt.Errorf("%T is not a *appsv1.DeploymentConfig", existingObj)
	}
	desired, ok := desiredObj.(*appsv1.DeploymentConfig)
	if !ok {
		return false, fmt.Errorf("%T is not a *appsv1.DeploymentConfig", desiredObj)
	}

	update := false

	tmpUpdate := DeploymentConfigContainerResourcesReconciler(desired, existing)
	update = update || tmpUpdate

	tmpUpdate = DeploymentConfigAffinityReconciler(desired, existing)
	update = update || tmpUpdate

	tmpUpdate = DeploymentConfigTolerationsReconciler(desired, existing)
	update = update || tmpUpdate

	return update, nil
}

func GenericDeploymentConfigMutator(existingObj, desiredObj common.KubernetesObject) (bool, error) {
	existing, ok := existingObj.(*appsv1.DeploymentConfig)
	if !ok {
		return false, fmt.Errorf("%T is not a *appsv1.DeploymentConfig", existingObj)
	}
	desired, ok := desiredObj.(*appsv1.DeploymentConfig)
	if !ok {
		return false, fmt.Errorf("%T is not a *appsv1.DeploymentConfig", desiredObj)
	}

	update := false

	tmpUpdate := DeploymentConfigReplicasReconciler(desired, existing)
	update = update || tmpUpdate

	tmpUpdate = DeploymentConfigContainerResourcesReconciler(desired, existing)
	update = update || tmpUpdate

	tmpUpdate = DeploymentConfigAffinityReconciler(desired, existing)
	update = update || tmpUpdate

	tmpUpdate = DeploymentConfigTolerationsReconciler(desired, existing)
	update = update || tmpUpdate

	return update, nil
}

func DeploymentConfigReplicasReconciler(desired, existing *appsv1.DeploymentConfig) bool {
	update := false

	if desired.Spec.Replicas != existing.Spec.Replicas {
		existing.Spec.Replicas = desired.Spec.Replicas
		update = true
	}

	return update
}

func DeploymentConfigAffinityReconciler(desired, existing *appsv1.DeploymentConfig) bool {
	updated := false

	if !reflect.DeepEqual(existing.Spec.Template.Spec.Affinity, desired.Spec.Template.Spec.Affinity) {
		diff := cmp.Diff(existing.Spec.Template.Spec.Affinity, desired.Spec.Template.Spec.Affinity)
		log.Info(fmt.Sprintf("%s spec.template.spec.Affinity has changed: %s", common.ObjectInfo(desired), diff))
		existing.Spec.Template.Spec.Affinity = desired.Spec.Template.Spec.Affinity
		updated = true
	}

	return updated
}

func DeploymentConfigTolerationsReconciler(desired, existing *appsv1.DeploymentConfig) bool {
	updated := false

	if !reflect.DeepEqual(existing.Spec.Template.Spec.Tolerations, desired.Spec.Template.Spec.Tolerations) {
		diff := cmp.Diff(existing.Spec.Template.Spec.Tolerations, desired.Spec.Template.Spec.Tolerations)
		log.Info(fmt.Sprintf("%s spec.template.spec.Tolerations has changed: %s", common.ObjectInfo(desired), diff))
		existing.Spec.Template.Spec.Tolerations = desired.Spec.Template.Spec.Tolerations
		updated = true
	}

	return updated
}

func DeploymentConfigContainerResourcesReconciler(desired, existing *appsv1.DeploymentConfig) bool {
	desiredName := common.ObjectInfo(desired)
	update := false

	//
	// Check container resource requirements
	//
	if len(desired.Spec.Template.Spec.Containers) != 1 {
		panic(fmt.Sprintf("%s desired spec.template.spec.containers length changed to '%d', should be 1", desiredName, len(desired.Spec.Template.Spec.Containers)))
	}

	if len(existing.Spec.Template.Spec.Containers) != 1 {
		log.Info(fmt.Sprintf("%s spec.template.spec.containers length changed to '%d', recreating dc", desiredName, len(existing.Spec.Template.Spec.Containers)))
		existing.Spec.Template.Spec.Containers = desired.Spec.Template.Spec.Containers
		update = true
	}

	if !helper.CmpResources(&existing.Spec.Template.Spec.Containers[0].Resources, &desired.Spec.Template.Spec.Containers[0].Resources) {
		diff := cmp.Diff(existing.Spec.Template.Spec.Containers[0].Resources, desired.Spec.Template.Spec.Containers[0].Resources, cmpopts.IgnoreUnexported(resource.Quantity{}))
		log.Info(fmt.Sprintf("%s spec.template.spec.containers[0].resources have changed: %s", desiredName, diff))
		existing.Spec.Template.Spec.Containers[0].Resources = desired.Spec.Template.Spec.Containers[0].Resources
		update = true
	}

	return update
}
