package kmap

import (
	"fmt"

	"github.com/infracloudio/ksearch/pkg/util"

	"k8s.io/client-go/kubernetes"
)

var rPool *ResourcePool
var groot *Resource
var indentationCount int = 0
var indentation string

// FindThemAll will fetch namespaced resources lists from ksearch
func FindThemAll(clientset *kubernetes.Clientset, namespace *string) error {
	var kinds = "Services,Endpoints,Deployments,ReplicaSets,StatefulSets,PodTemplates,Pods,ConfigMaps,PersistentVolumeClaims,PersistentVolumes,Secrets,ServiceAccounts,DaemonSets"

	getter := make(chan interface{})

	go util.Getter(*namespace, clientset, kinds, getter)
	resourceNamespace := &Resource{name: *namespace, kind: "Namespace"}
	rPool = &ResourcePool{resources: []*Resource{resourceNamespace}}
	for {
		resource, ok := <-getter
		if !ok {
			break
		}

		rPool = convertResource(resource, rPool)

	}
	groot = resourceNamespace
	groot = mapThemAll(groot, rPool)
	printTree(groot)
	return nil
}

func mapThemAll(groot *Resource, rPool *ResourcePool) *Resource {
	resources := rPool.resources
	for _, resource := range resources {

		if resource.kind == "Service" && resource.parent == nil {
			serviceGroot, uPool := findServiceChildren(resource, rPool)
			groot.children = append(groot.children, serviceGroot)
			resource.parent = groot
			rPool = uPool

		}
		if resource.kind == "Deployment" && resource.parent == nil {
			deploymentGroot, uPool := findDeployChildren(resource, rPool)
			groot.children = append(groot.children, deploymentGroot)
			resource.parent = groot
			rPool = uPool

		}
		if resource.kind == "Replicaset" && resource.parent == nil {
			rsGroot, uPool := findRSChildren(resource, rPool)
			groot.children = append(groot.children, rsGroot)
			resource.parent = groot
			rPool = uPool

		}
		if resource.kind == "StatefulSet" && resource.parent == nil {
			stsgroot, uPool := findStatefulSetChildren(resource, rPool)
			groot.children = append(groot.children, stsgroot)
			resource.parent = groot
			rPool = uPool

		}

		if resource.kind == "Pod" && resource.parent == nil {
			podGroot, uPool := findPodChildren(resource, rPool)
			groot.children = append(groot.children, podGroot)
			resource.parent = groot
			rPool = uPool

		}
		if resource.kind == "PersistentVolumeClaim" && resource.parent == nil {
			pvcGroot, uPool := findPVCChildren(resource, rPool)
			groot.children = append(groot.children, pvcGroot)
			resource.parent = groot
			rPool = uPool

		}
		if resource.kind != "Namespace" && resource.parent == nil {
			groot.children = append(groot.children, resource)
			resource.parent = groot

		} else {
			continue
		}

	}
	return groot
}

func findNSChildren(groot *Resource, rPool *ResourcePool) *Resource {
	for _, resource := range rPool.resources {
		if resource.parent == nil && resource.kind != "Namespace" {

			groot.children = append(groot.children, resource)
		}
	}
	return groot
}

func findPodChildren(pod *Resource, rPool *ResourcePool) (*Resource, *ResourcePool) {
	// Pod children are supposed to be configMaps, secrets and volumeMounts

	for _, resource := range rPool.resources {

		if resource.kind == "ConfigMap" || resource.kind == "Secret" || resource.kind == "PersistentVolumeClaim" {
			for _, volume := range pod.spec.Volumes {
				configMap := volume.VolumeSource.ConfigMap
				if configMap != nil && !resource.hasParent {
					if volume.VolumeSource.ConfigMap.LocalObjectReference.Name == resource.name {
						pod.children = append(pod.children, resource)
						resource.parent = pod
						resource.hasParent = true
						break
					}
				}
				pvc := volume.VolumeSource.PersistentVolumeClaim
				if pvc != nil && !resource.hasParent {
					if volume.VolumeSource.PersistentVolumeClaim.ClaimName == resource.name {
						resource, uPool := findPVCChildren(resource, rPool)
						pod.children = append(pod.children, resource)
						resource.parent = pod
						rPool = uPool
						resource.hasParent = true
						break

					}
				}
				secret := volume.VolumeSource.Secret
				if secret != nil && !resource.hasParent {
					if volume.VolumeSource.Secret.SecretName == resource.name {
						pod.children = append(pod.children, resource)
						resource.parent = pod
						resource.hasParent = true
						break
					}
				}
			}

		}
	}

	return pod, rPool
}

func findDeployChildren(deploy *Resource, rPool *ResourcePool) (*Resource, *ResourcePool) {

	for _, resource := range rPool.resources {
		// var isMapped bool

		if resource.kind == "Replicaset" || resource.kind == "PodTemplate" {
			// fmt.Faddf(w, "%v\t%v\t%v\t%v\t%v\n", deployment.Name, deployment.Status.ReadyReplicas, "", deployment.Status.AvailableReplicas, "")
			selector := deploy.selector
			if selector != nil && !resource.hasParent {
				// filterPodsWithLabel(pods, selector)
				if resource.Labels != nil {
					for key, val := range resource.Labels {
						result, ok := selector[key]
						if !ok {
							continue
						}
						if result == val {
							resource, uPool := findRSChildren(resource, rPool)
							deploy.children = append(deploy.children, resource)
							resource.parent = deploy
							rPool = uPool
							resource.hasParent = true
							break

						}
					}
				}
				if !resource.hasParent {
					for key, val := range resource.selector {
						result, ok := selector[key]
						if !ok {
							continue
						}
						if result == val {
							resource, uPool := findRSChildren(resource, rPool)
							deploy.children = append(deploy.children, resource)
							resource.parent = deploy
							rPool = uPool
							break
						}
					}
				}

			}

		}
	}

	return deploy, rPool
}

func findStatefulSetChildren(StatefulSet *Resource, rPool *ResourcePool) (*Resource, *ResourcePool) {

	for _, resource := range rPool.resources {

		if resource.kind == "Pod" || resource.kind == "PodTemplate" {
			selector := StatefulSet.selector
			if selector != nil && !resource.hasParent {
				if resource.Labels != nil {
					for key, val := range resource.Labels {
						result, ok := selector[key]
						if !ok {
							continue
						}
						if result == val {
							resource, uPool := findPodChildren(resource, rPool)
							StatefulSet.children = append(StatefulSet.children, resource)
							resource.parent = StatefulSet
							rPool = uPool
							resource.hasParent = true
							break

						}
					}
				}
				if !resource.hasParent {
					for key, val := range resource.selector {
						result, ok := selector[key]
						if !ok {
							continue
						}
						if result == val {
							resource, uPool := findRSChildren(resource, rPool)
							StatefulSet.children = append(StatefulSet.children, resource)
							resource.parent = StatefulSet
							rPool = uPool
							break
						}
					}
				}

			}

		}
	}

	return StatefulSet, rPool
}

func findServiceChildren(service *Resource, rPool *ResourcePool) (*Resource, *ResourcePool) {
	// var uPool *ResourcePool
	for _, resource := range rPool.resources {

		if resource.kind == "Deployment" || resource.kind == "Endpoint" {
			selector := service.selector
			if selector != nil && !resource.hasParent {
				if resource.Labels != nil {
					for key, val := range resource.Labels {
						result, ok := selector[key]
						if !ok {
							continue
						}
						if result == val {
							resource, uPool := findDeployChildren(resource, rPool)
							service.children = append(service.children, resource)
							resource.parent = service
							rPool = uPool
							resource.hasParent = true
							break
						}
					}
				}
				if resource.selector != nil && !resource.hasParent {
					for key, val := range resource.selector {
						result, ok := selector[key]
						if !ok {
							continue
						}
						if result == val {
							resource, uPool := findDeployChildren(resource, rPool)
							service.children = append(service.children, resource)
							resource.parent = service
							rPool = uPool
							resource.hasParent = true
							break
						}
					}
				}

			}
		}
	}
	return service, rPool
}

func findRSChildren(rs *Resource, rPool *ResourcePool) (*Resource, *ResourcePool) {

	for _, resource := range rPool.resources {

		if resource.kind == "Pod" && !resource.hasParent {
			selector := rs.selector
			if selector != nil {
				for key, val := range resource.Labels {
					result, ok := selector[key]
					if !ok {
						continue
					}
					if result == val {
						resource, uPool := findPodChildren(resource, rPool)
						rs.children = append(rs.children, resource)
						resource.parent = rs
						rPool = uPool
						resource.hasParent = true
						break

					}
				}

			}
		}
	}
	return rs, rPool
}

func findPVCChildren(pvc *Resource, rPool *ResourcePool) (*Resource, *ResourcePool) {

	for _, resource := range rPool.resources {

		if resource.kind == "PersistentVolume" && !resource.hasParent {

			if pvc.info["Volume"] == resource.name {
				// resource, uPool := findPodChildren(resource, rPool)
				pvc.children = append(pvc.children, resource)
				resource.parent = pvc
				resource.hasParent = true
				break
				// rPool = uPool

			}
		}
	}

	return pvc, rPool
}

func getTreeOutine() string {
	return "|-------"
}
func getIndentation(ind int) string {
	indentation = ""
	for i := 0; i < ind; i++ {
		indentation += "\t"
	}
	if indentationCount >= 1 {
		indentation += getTreeOutine()
	}

	return indentation
}

func printTree(groot *Resource) {

	if groot.kind == "Event" {
		if groot.info["Type"] != "Normal" {
			fmt.Println(getIndentation(indentationCount)+groot.kind, " - Type - ", groot.info["Type"], "\t Reason - ", groot.info["Reason"], "\t ObjectKind/Name - (", groot.info["ObjectKind"], "/", groot.info["ObjectName"], ")")
			fmt.Println(getIndentation(indentationCount)+"*** Message *** ", groot.info["Message"])
		}

	} else {
		fmt.Println(getIndentation(indentationCount)+groot.kind, " - ", groot.name, " Status: ", groot.status)
	}

	if len(groot.children) > 0 {
		indentationCount += 1
		for _, child := range groot.children {
			printTree(child)
		}
		indentationCount -= 1
	}

}
