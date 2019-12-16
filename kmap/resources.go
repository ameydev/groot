package kmap

import (
	v1 "k8s.io/api/core/v1"

	appsv1 "k8s.io/api/apps/v1"
)

type Resource struct {
	name           string
	status         string
	kind           string
	isStandAlone   bool
	info           map[string]string
	children       []*Resource
	parent         *Resource
	Labels         map[string]string
	selector       map[string]string
	ownerReference []string
	spec           v1.PodSpec
	hasParent      bool
}

type ResourcePool struct {
	resources []*Resource
}

func (r *ResourcePool) addToResourcePool(resource *Resource) {
	r.resources = append(r.resources, resource)
}

func convertResource(resource interface{}, rPool *ResourcePool) *ResourcePool {
	switch resource {

	case resource.(*v1.PodList):
		pods := resource.(*v1.PodList)
		rPool = addPodDetails(pods, rPool)
	case resource.(*v1.ComponentStatusList):
		componentStatuses := resource.(*v1.ComponentStatusList)
		rPool = addComponentStatuses(componentStatuses, rPool)
	case resource.(*v1.ConfigMapList):
		cms := resource.(*v1.ConfigMapList)
		rPool = addConfigMaps(cms, rPool)
	case resource.(*v1.EndpointsList):
		endPoints := resource.(*v1.EndpointsList)
		rPool = addEndpoints(endPoints, rPool)
	case resource.(*v1.EventList):
		events := resource.(*v1.EventList)
		rPool = addEvents(events, rPool)
	case resource.(*v1.LimitRangeList):
		limitRanges := resource.(*v1.LimitRangeList)
		rPool = addLimitRanges(limitRanges, rPool)
	case resource.(*v1.NamespaceList):
		namespaces := resource.(*v1.NamespaceList)
		rPool = addNamespaces(namespaces, rPool)
	case resource.(*v1.PersistentVolumeList):
		pvs := resource.(*v1.PersistentVolumeList)
		rPool = addPVs(pvs, rPool)
	case resource.(*v1.PersistentVolumeClaimList):
		pvcs := resource.(*v1.PersistentVolumeClaimList)
		rPool = addPVCs(pvcs, rPool)
	case resource.(*v1.PodTemplateList):
		podTemplates := resource.(*v1.PodTemplateList)
		rPool = addPodTemplates(podTemplates, rPool)
	case resource.(*v1.ResourceQuotaList):
		resQuotas := resource.(*v1.ResourceQuotaList)
		rPool = addResourceQuotas(resQuotas, rPool)
	case resource.(*v1.SecretList):
		secrets := resource.(*v1.SecretList)
		rPool = addSecrets(secrets, rPool)
	case resource.(*v1.ServiceList):
		services := resource.(*v1.ServiceList)
		rPool = addServices(services, rPool)
	case resource.(*v1.ServiceAccountList):
		serviceAccs := resource.(*v1.ServiceAccountList)
		rPool = addServiceAccounts(serviceAccs, rPool)

		// these will be from the appsV1
	case resource.(*appsv1.DaemonSetList):
		daemonsets := resource.(*appsv1.DaemonSetList)
		rPool = addDaemonSets(daemonsets, rPool)
	case resource.(*appsv1.DeploymentList):
		deployments := resource.(*appsv1.DeploymentList)
		rPool = addDeployments(deployments, rPool)
	case resource.(*appsv1.ReplicaSetList):
		rsets := resource.(*appsv1.ReplicaSetList)
		rPool = addReplicaSets(rsets, rPool)
	case resource.(*appsv1.StatefulSetList):
		ssets := resource.(*appsv1.StatefulSetList)
		rPool = addStateFulSets(ssets, rPool)
	}
	return rPool
}

func addPodDetails(pods *v1.PodList, rPool *ResourcePool) *ResourcePool {
	if len(pods.Items) > 0 {

		for _, pod := range pods.Items {
			var podResource Resource
			podResource.name = pod.Name
			podResource.kind = "Pod"
			podResource.status = string(pod.Status.Phase)
			podResource.Labels = pod.Labels
			podResource.spec = pod.Spec
			rPool.addToResourcePool(&podResource)
		}

	}
	return rPool
}
func addPodTemplates(podTemplates *v1.PodTemplateList, rPool *ResourcePool) *ResourcePool {
	if len(podTemplates.Items) > 0 {
		for _, podTemplate := range podTemplates.Items {
			var podTemplateResource Resource
			podTemplateResource.name = podTemplate.Name
			podTemplateResource.kind = "PodTemplate"
			podTemplateResource.Labels = podTemplate.Labels
			rPool.addToResourcePool(&podTemplateResource)

		}

	}
	return rPool
}
func addComponentStatuses(componentStatuses *v1.ComponentStatusList, rPool *ResourcePool) *ResourcePool {
	if len(componentStatuses.Items) > 0 {
		for _, componentStatus := range componentStatuses.Items {
			var podTemplateResource Resource
			podTemplateResource.name = componentStatus.Name
			podTemplateResource.kind = "ComponentStatus"
			podTemplateResource.Labels = componentStatus.Labels
			rPool.addToResourcePool(&podTemplateResource)

		}

	}
	return rPool
}
func addConfigMaps(cms *v1.ConfigMapList, rPool *ResourcePool) *ResourcePool {
	if len(cms.Items) > 0 {
		for _, configMap := range cms.Items {
			var blankResource Resource
			blankResource.name = configMap.Name
			blankResource.kind = "ConfigMap"
			blankResource.Labels = configMap.Labels
			rPool.addToResourcePool(&blankResource)
		}
	}
	return rPool
}
func addEndpoints(endPoints *v1.EndpointsList, rPool *ResourcePool) *ResourcePool {
	if len(endPoints.Items) > 0 {

		for _, endpoint := range endPoints.Items {
			var blankResource Resource
			blankResource.name = endpoint.Name
			blankResource.kind = "Endpoint"
			blankResource.Labels = endpoint.Labels
			rPool.addToResourcePool(&blankResource)
		}

	}
	return rPool
}
func addEvents(events *v1.EventList, rPool *ResourcePool) *ResourcePool {
	if len(events.Items) > 0 {

		for _, event := range events.Items {
			var blankResource Resource
			blankResource.name = event.Name
			blankResource.kind = "Event"
			var data = map[string]string{}
			data["Type"] = event.Type
			data["ObjectKind"] = event.InvolvedObject.Kind
			data["ObjectName"] = event.InvolvedObject.Name
			data["Message"] = event.Message
			data["Reason"] = event.Reason
			blankResource.info = data
			rPool.addToResourcePool(&blankResource)
		}

	}

	return rPool
}
func addLimitRanges(limitRanges *v1.LimitRangeList, rPool *ResourcePool) *ResourcePool {
	if len(limitRanges.Items) > 0 {

		for _, limitRange := range limitRanges.Items {
			var blankResource Resource
			blankResource.name = limitRange.Name
			blankResource.kind = "LimitRange"
			blankResource.Labels = limitRange.Labels
			rPool.addToResourcePool(&blankResource)
		}

	}
	return rPool
}
func addNamespaces(namespaces *v1.NamespaceList, rPool *ResourcePool) *ResourcePool {
	if len(namespaces.Items) > 0 {

		for _, namespace := range namespaces.Items {
			var blankResource Resource
			blankResource.name = namespace.Name
			blankResource.kind = "Namespace"
			blankResource.Labels = namespace.Labels
			rPool.addToResourcePool(&blankResource)
		}

	}
	return rPool
}
func addPVs(pvs *v1.PersistentVolumeList, rPool *ResourcePool) *ResourcePool {
	if len(pvs.Items) > 0 {

		for _, pv := range pvs.Items {
			var blankResource Resource
			blankResource.name = pv.Name
			blankResource.kind = "PersistentVolume"
			blankResource.status = string(pv.Status.Phase)
			blankResource.Labels = pv.Labels
			rPool.addToResourcePool(&blankResource)
		}

	}
	return rPool
}
func addPVCs(pvcs *v1.PersistentVolumeClaimList, rPool *ResourcePool) *ResourcePool {
	if len(pvcs.Items) > 0 {

		for _, pvc := range pvcs.Items {
			var blankResource Resource
			blankResource.name = pvc.Name
			blankResource.kind = "PersistentVolumeClaim"
			blankResource.status = string(pvc.Status.Phase)
			blankResource.Labels = pvc.Labels
			var data = map[string]string{}
			data["Volume"] = pvc.Spec.VolumeName
			blankResource.info = data
			rPool.addToResourcePool(&blankResource)
		}

	}
	return rPool
}

func addResourceQuotas(resQuotas *v1.ResourceQuotaList, rPool *ResourcePool) *ResourcePool {
	if len(resQuotas.Items) > 0 {

		for _, resQ := range resQuotas.Items {
			var blankResource Resource
			blankResource.name = resQ.Name
			blankResource.kind = "ResourceQuota"
			blankResource.Labels = resQ.Labels
			rPool.addToResourcePool(&blankResource)
		}
	}
	return rPool
}
func addSecrets(secrets *v1.SecretList, rPool *ResourcePool) *ResourcePool {
	if len(secrets.Items) > 0 {

		for _, secret := range secrets.Items {
			var blankResource Resource
			blankResource.name = secret.Name
			blankResource.kind = "Secret"
			blankResource.Labels = secret.Labels
			rPool.addToResourcePool(&blankResource)
		}

	}
	return rPool
}
func addServices(services *v1.ServiceList, rPool *ResourcePool) *ResourcePool {
	if len(services.Items) > 0 {

		for _, service := range services.Items {
			var blankResource Resource
			blankResource.name = service.Name
			blankResource.kind = "Service"
			blankResource.Labels = service.Labels
			blankResource.selector = service.Spec.Selector
			rPool.addToResourcePool(&blankResource)
		}
	}
	return rPool
}
func addServiceAccounts(serviceAccs *v1.ServiceAccountList, rPool *ResourcePool) *ResourcePool {
	if len(serviceAccs.Items) > 0 {

		for _, serviceAcc := range serviceAccs.Items {
			var blankResource Resource
			blankResource.name = serviceAcc.Name
			blankResource.kind = "ServiceAccount"
			blankResource.Labels = serviceAcc.Labels
			rPool.addToResourcePool(&blankResource)
		}

	}
	return rPool
}
func addDaemonSets(daemonsets *appsv1.DaemonSetList, rPool *ResourcePool) *ResourcePool {
	if len(daemonsets.Items) > 0 {

		for _, ds := range daemonsets.Items {
			var blankResource Resource
			blankResource.name = ds.Name
			blankResource.kind = "DaemonSet"
			blankResource.status = string(ds.Status.NumberReady) + "/" + string(ds.Status.DesiredNumberScheduled)
			blankResource.Labels = ds.Labels
			blankResource.selector = ds.Spec.Selector.MatchLabels
			rPool.addToResourcePool(&blankResource)
		}

	}
	return rPool
}
func addDeployments(deployments *appsv1.DeploymentList, rPool *ResourcePool) *ResourcePool {
	if len(deployments.Items) > 0 {
		for _, deployment := range deployments.Items {
			var blankResource Resource
			blankResource.name = deployment.Name
			blankResource.kind = "Deployment"
			var readyReplica int32 = deployment.Status.ReadyReplicas
			var availableReplica int32 = deployment.Status.AvailableReplicas
			blankResource.status = Int32toString(readyReplica) + "/" + Int32toString(availableReplica)
			blankResource.Labels = deployment.Labels
			blankResource.selector = deployment.Spec.Selector.MatchLabels
			rPool.addToResourcePool(&blankResource)
		}
	}
	return rPool
}
func addReplicaSets(rsets *appsv1.ReplicaSetList, rPool *ResourcePool) *ResourcePool {
	if len(rsets.Items) > 0 {
		for _, rs := range rsets.Items {
			var blankResource Resource
			blankResource.name = rs.Name
			blankResource.kind = "Replicaset"
			blankResource.Labels = rs.Labels
			blankResource.selector = rs.Spec.Selector.MatchLabels
			rPool.addToResourcePool(&blankResource)
		}
	}
	return rPool
}
func addStateFulSets(ssets *appsv1.StatefulSetList, rPool *ResourcePool) *ResourcePool {
	if len(ssets.Items) > 0 {
		for _, sset := range ssets.Items {
			var blankResource Resource
			blankResource.name = sset.Name
			blankResource.kind = "StatefulSet"
			blankResource.status = Int32toString(sset.Status.ReadyReplicas) + "/" + Int32toString(sset.Status.CurrentReplicas)
			blankResource.Labels = sset.Labels
			blankResource.selector = sset.Spec.Selector.MatchLabels
			rPool.addToResourcePool(&blankResource)
		}
	}
	return rPool
}

func Int32toString(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}
