package kmap

import (
	"fmt"

	"github.com/infracloudio/ksearch/pkg/util"

	// "github.com/ameydev/groot/printer"

	"k8s.io/client-go/kubernetes"
)

var rPool *ResourcePool
var groot *Resource
var indentationCount int = 1
var indentation string

// var kinds = "Pods,ComponentStatuses,ConfigMaps,Endpoints,LimitRanges,Namespaces,PersistentVolumes,PersistentVolumeClaims,PodTemplates,ResourceQuotas,Secrets,Services,ServiceAccounts,DaemonSets,Deployments,ReplicaSets,StatefulSets"

func FindThemAll(clientset *kubernetes.Clientset, namespace *string) error {
	// var kinds = "Pods,ComponentStatuses,ConfigMaps,Endpoints,LimitRanges,Namespaces,PersistentVolumes,PersistentVolumeClaims,PodTemplates,ResourceQuotas,Secrets,Services,ServiceAccounts,DaemonSets,Deployments,ReplicaSets,StatefulSets"
	var kinds = "Services,Endpoints,Deployments,ReplicaSets,StatefulSets,PodTemplates,Pods,ConfigMaps,PersistentVolumeClaims,PersistentVolumes,Secrets,ServiceAccounts,DaemonSets"

	getter := make(chan interface{})

	go util.Getter(*namespace, clientset, kinds, getter)
	resourceNamespace := &Resource{name: *namespace, kind: "Namespace"}
	rPool = &ResourcePool{resources: []*Resource{resourceNamespace}}
	for {
		resource, ok := <-getter
		if !ok {
			// fmt.Print("domee", resource)
			// return nil
			break
		}

		rPool = convertResource(resource, rPool)
		// fmt.Print("Resource", resource)

	}
	groot = resourceNamespace
	groot = mapThemAll(groot, rPool)
	printTree(groot)
	return nil
}

func mapThemAll(groot *Resource, rPool *ResourcePool) *Resource {
	resources := rPool.resources
	// fmt.Println("In Mapp")
	for _, resource := range resources {
		// resource.parent = groot
		// fmt.Println(resource.kind, " - ", resource.name, " Status: ", resource.status)

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
			// fmt.Println("The reource was namespace")
			continue
		}

	}
	return groot
}

// func getServiceGroot(serviceResource *Resource, rPool *ResourcePool) {

// 	return resource, &rPool
// }

func findNSChildren(groot *Resource, rPool *ResourcePool) *Resource {
	for _, resource := range rPool.resources {
		if resource.parent == nil && resource.kind != "Namespace" {

			groot.children = append(groot.children, resource)
		}
	}
	return groot
}

func getIndentation(ind int) string {
	indentation = ""
	for i := 0; i < ind; i++ {
		indentation += "\t"
	}
	// indentationCount += 1
	return indentation
}

func printTree(groot *Resource) {

	fmt.Println(getIndentation(indentationCount)+groot.kind, " - ", groot.name, " Status: ", groot.status)

	if len(groot.children) > 0 {
		indentationCount += 1
		for _, child := range groot.children {
			printTree(child)
		}
		indentationCount -= 1
	}

}
func findPodChildren(pod *Resource, rPool *ResourcePool) (*Resource, *ResourcePool) {
	// Pod children are supposed to be configMaps, secrets and volumeMounts

	for _, resource := range rPool.resources {

		if resource.kind == "ConfigMap" || resource.kind == "Secret" || resource.kind == "PersistentVolumeClaim" {
			for _, volume := range pod.spec.Volumes {
				configMap := volume.VolumeSource.ConfigMap
				if configMap != nil {
					if volume.VolumeSource.ConfigMap.LocalObjectReference.Name == resource.name {
						pod.children = append(pod.children, resource)
						resource.parent = pod
						break
					}
				}
				pvc := volume.VolumeSource.PersistentVolumeClaim
				if pvc != nil {
					if volume.VolumeSource.PersistentVolumeClaim.ClaimName == resource.name {
						// pod.children = append(pod.children, resource)
						// resource.parent = pod
						resource, uPool := findPVCChildren(resource, rPool)
						pod.children = append(pod.children, resource)
						resource.parent = pod
						rPool = uPool
						break

					}
				}
				secret := volume.VolumeSource.Secret
				if secret != nil {
					if volume.VolumeSource.Secret.SecretName == resource.name {
						pod.children = append(pod.children, resource)
						resource.parent = pod
						break
					}
				}
			}

		}
	}
	// groot.children = append(groot.children, pod)
	return pod, rPool
}

func findDeployChildren(deploy *Resource, rPool *ResourcePool) (*Resource, *ResourcePool) {

	for _, resource := range rPool.resources {

		if resource.kind == "Replicaset" {
			// fmt.Faddf(w, "%v\t%v\t%v\t%v\t%v\n", deployment.Name, deployment.Status.ReadyReplicas, "", deployment.Status.AvailableReplicas, "")
			selector := deploy.selector
			if selector != nil {
				// filterPodsWithLabel(pods, selector)
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

					}
				}

			}
		}
	}
	return deploy, rPool
}

func findServiceChildren(service *Resource, rPool *ResourcePool) (*Resource, *ResourcePool) {
	// var uPool *ResourcePool
	for _, resource := range rPool.resources {

		if resource.kind == "Deployment" || resource.kind == "Endpoint" {
			selector := service.selector
			if selector != nil {
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
					}
				}

			}
		}
	}
	return service, rPool
}

func findRSChildren(rs *Resource, rPool *ResourcePool) (*Resource, *ResourcePool) {

	for _, resource := range rPool.resources {

		if resource.kind == "Pod" {
			// fmt.Faddf(w, "%v\t%v\t%v\t%v\t%v\n", deployment.Name, deployment.Status.ReadyReplicas, "", deployment.Status.AvailableReplicas, "")
			selector := rs.selector
			if selector != nil {
				// filterPodsWithLabel(pods, selector)
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

		if resource.kind == "PersistentVolume" {

			if pvc.info["Volume"] == resource.name {
				// resource, uPool := findPodChildren(resource, rPool)
				pvc.children = append(pvc.children, resource)
				resource.parent = pvc
				break
				// rPool = uPool

			}
		}
	}

	return pvc, rPool
}
