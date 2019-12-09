package kmap

import (
	"fmt"

	"github.com/infracloudio/ksearch/pkg/util"

	// "github.com/ameydev/groot/ksearch"
	"k8s.io/client-go/kubernetes"
)

var rPool *ResourcePool

// var kinds = "Pods,ComponentStatuses,ConfigMaps,Endpoints,LimitRanges,Namespaces,PersistentVolumes,PersistentVolumeClaims,PodTemplates,ResourceQuotas,Secrets,Services,ServiceAccounts,DaemonSets,Deployments,ReplicaSets,StatefulSets"

func FindThemAll(clientset *kubernetes.Clientset, namespace *string) error {
	var kinds = "Pods,ComponentStatuses,ConfigMaps,Endpoints,LimitRanges,Namespaces,PersistentVolumes,PersistentVolumeClaims,PodTemplates,ResourceQuotas,Secrets,Services,ServiceAccounts,DaemonSets,Deployments,ReplicaSets,StatefulSets"

	getter := make(chan interface{})

	go util.Getter(*namespace, clientset, kinds, getter)
	resourceNamespace := &Resource{name: *namespace}
	rPool = &ResourcePool{resources: []*Resource{resourceNamespace}}
	for {
		// time.Sleep(time.Second)
		resource, ok := <-getter

		if ok {
			// fmt.Println("Some issue with", resource)
			// fmt.Println(resource)
			rPool = convertResource(resource, rPool)
		} else {
			// fmt.Println("some err ", resource)
			break
		}

		// else {
		// 	rPool = convertResource(resource, rPool)
		// }

		// fmt.Print(rPool)
		// rPool =addToResourcePool( &rPool)
		// printers.Printer(resource, *resName)
	}

	mapThemAll(rPool)
	return nil
}

func mapThemAll(rPool *ResourcePool) {
	resources := rPool.resources
	// fmt.Println("resources", resources)

	// fmt.Println("resource 0", *resources[0])
	for index, rs := range resources {
		fmt.Println(index, rs.kind, " - ", rs.name, " Status: ", rs.status)

	}
}
