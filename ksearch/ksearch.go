package ksearch

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

type Resource struct {
	name     string
	status   string
	kind     string
	children []*Resource
}

func SearchResources(clientset *kubernetes.Clientset, namespace *string) error {

	pods, err := clientset.CoreV1().Pods(*namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting the pod from clientset", err)
		return err
	}

	printPodDetails(pods)

	componentStatuses, err := clientset.CoreV1().ComponentStatuses().List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting components statusses from clientset", err)
		return err
	}
	printComponentStatuses(componentStatuses)

	configmaps, err := clientset.CoreV1().ConfigMaps(*namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting configmaps from cleintset", err)
		return err
	}
	printConfigMaps(configmaps)

	endPoints, err := clientset.CoreV1().Endpoints(*namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting endpoints from clientset", err)
		return err
	}
	printEndpoints(endPoints)

	limitRanges, err := clientset.CoreV1().LimitRanges(*namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting limitrange from clientset", err)
		return err
	}
	printLimitRanges(limitRanges)

	pvs, err := clientset.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting PVs throuhg CLIentset", err)
		return err
	}
	printPVs(pvs)

	pvcs, err := clientset.CoreV1().PersistentVolumeClaims(*namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting PVCs throuhg CLIentset", err)
		return err
	}
	printPVCs(pvcs)

	podTemplates, err := clientset.CoreV1().PodTemplates(*namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting podTemplates throuhg CLIentset", err)
		return err
	}
	printPodTemplates(podTemplates)

	resQuotas, err := clientset.CoreV1().ResourceQuotas(*namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting resourceQuots throuhg CLIentset", err)
		return err
	}
	printResourceQuotas(resQuotas)

	secrets, err := clientset.CoreV1().Secrets(*namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting secrets throuhg CLIentset", err)
		return err
	}
	printSecrets(secrets)

	services, err := clientset.CoreV1().Services(*namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting services  throuhg CLIentset", err)
		return err
	}
	printServices(services)

	serviceAccs, err := clientset.CoreV1().ServiceAccounts(*namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting serviceacc throuhg CLIentset", err)
		return err
	}
	printServiceAccounts(serviceAccs)

	daemonsets, err := clientset.AppsV1().DaemonSets(*namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting ds from clientset", err)
		return err
	}
	printDaemonSets(daemonsets)

	deployments, err := clientset.AppsV1().Deployments(*namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting deployment from clientset", err)
		return err
	}
	printDeployments(deployments)
	replicasets, err := clientset.AppsV1().ReplicaSets(*namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting replicasets from clientset ", err)
		return err
	}
	printReplicaSets(replicasets)
	statefulSets, err := clientset.AppsV1().StatefulSets(*namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Info("There was an error getting statefulset from clientset", err)
		return err
	}
	printStateFulSets(statefulSets)

	return nil
}

var indentationCount int = 1
var indentation string

func getIndentation(ind int) string {

	indentation = ""
	for i := 0; i < ind; i++ {
		indentation += "\t"
	}
	return indentation

}
func printPodDetails(pods *v1.PodList) {
	fmt.Println("\nPods\n----\n")
	var tempInd string = getIndentation(indentationCount)

	for _, pod := range pods.Items {
		fmt.Println(tempInd+"- Pods - "+pod.Name, "", pod.Status.Phase, "")
	}
}

func printComponentStatuses(componentStatuses *v1.ComponentStatusList) {
	fmt.Println("\nComponentStatuses\n-------------\n")
	var tempInd string = getIndentation(indentationCount)

	for _, componentStatus := range componentStatuses.Items {
		fmt.Println(tempInd+"- ComponentStatuses - "+componentStatus.Name, "", componentStatus.Conditions[0].Type, "")
	}
}

func printConfigMaps(cms *v1.ConfigMapList) {
	fmt.Println("\nConfigMaps\n--------------\n")
	var tempInd string = getIndentation(indentationCount)
	for _, configMap := range cms.Items {
		fmt.Println(tempInd + "- ConfigMaps - " + configMap.Name)
	}
}

func printEndpoints(endPoints *v1.EndpointsList) {
	fmt.Println("\nEndpoints\n--------------\n")
	var tempInd string = getIndentation(indentationCount)
	for _, endpoint := range endPoints.Items {
		fmt.Println(tempInd + "- Endpoints - " + endpoint.Name)
	}

}

func printLimitRanges(limitRanges *v1.LimitRangeList) {
	fmt.Println("\nLimitRanges\n--------------\n")
	var tempInd string = getIndentation(indentationCount)

	for _, limitRange := range limitRanges.Items {
		fmt.Println(tempInd+"- LimitRanges - ", limitRange.Name)
	}

}

func printNamespaces(namespaces *v1.NamespaceList) {
	fmt.Println("\nNamespaces\n--------------\n")
	var tempInd string = getIndentation(indentationCount)

	for _, namespace := range namespaces.Items {
		fmt.Println(tempInd+"- NameSpaces -", namespace.Name)
	}

}

func printPVs(pvs *v1.PersistentVolumeList) {
	var tempInd string = getIndentation(indentationCount)

	fmt.Println("\nPersistentVolumes\n--------------\n")
	for _, pv := range pvs.Items {
		fmt.Println(tempInd+"- Persistent Volumes - ", pv.Name, ": Status - ", pv.Status,
			"\n\t\t PVClaim reference Name - ", pv.Spec.ClaimRef.Name)
	}

}

func printPVCs(pvcs *v1.PersistentVolumeClaimList) {
	fmt.Println("\nPersistentVolumeClaims\n--------------\n")
	var tempInd string = getIndentation(indentationCount)
	for _, pvc := range pvcs.Items {
		fmt.Println(tempInd+"- Persistent VolumeClaims - ", pvc.Name, ": Status - ", pvc.Status)

	}

}

func printPodTemplates(podTemplates *v1.PodTemplateList) {
	fmt.Println("\nPodTemplates\n--------------\n")
	var tempInd string = getIndentation(indentationCount)

	for _, podTemplate := range podTemplates.Items {
		fmt.Println(tempInd+"- PodTemplates - ", podTemplate.Name)
	}

}

func printResourceQuotas(resQuotas *v1.ResourceQuotaList) {
	fmt.Println("\nResourceQuotas\n--------------\n")
	var tempInd string = getIndentation(indentationCount)

	for _, resQ := range resQuotas.Items {
		fmt.Println(tempInd+"- ResourseQuotas - ", resQ.Name, " Creation Time -", resQ.CreationTimestamp)
	}

}

func printSecrets(secrets *v1.SecretList) {
	fmt.Println("\nSecrets\n--------------\n")
	var tempInd string = getIndentation(indentationCount)

	for _, secret := range secrets.Items {
		fmt.Println(tempInd+"- Secrets -", secret.Name, " - secret type: ", secret.Type)
	}

}

func printServices(services *v1.ServiceList) {
	fmt.Println("\nServices\n--------------\n")
	var tempInd string = getIndentation(indentationCount)

	for _, service := range services.Items {
		fmt.Println(tempInd+"- Services -", service.Name, ": Status -")
	}

}

func printServiceAccounts(serviceAccs *v1.ServiceAccountList) {
	fmt.Println("\nServiceAccounts\n--------------\n")
	var tempInd string = getIndentation(indentationCount)

	for _, serviceAcc := range serviceAccs.Items {
		fmt.Println(tempInd+"- ServiceAccounts -", serviceAcc.Name)
	}

}

func printDaemonSets(daemonsets *appsv1.DaemonSetList) {
	fmt.Println("\nDaemonSets\n--------------\n")
	var tempInd string = getIndentation(indentationCount)

	for _, ds := range daemonsets.Items {
		fmt.Println("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", ds.Namespace, ds.Name, ds.Status.DesiredNumberScheduled, ds.Status.CurrentNumberScheduled, ds.Status.NumberReady, "", ds.Status.NumberAvailable, ds.Spec.Template.Spec.NodeSelector, "")
		fmt.Println(tempInd+"- DaemonSets -", ds.Name, ": Status(DesiredNumberScheduled/CurrentNumberScheduled - NumberReady/NumberAvailable) - (%v/%v - %v/%v)", ds.Status.DesiredNumberScheduled, ds.Status.CurrentNumberScheduled, ds.Status.NumberReady, ds.Status.NumberAvailable)
	}

}

func printDeployments(deployments *appsv1.DeploymentList) {
	fmt.Println("\nDeployments\n--------------\n")
	var tempInd string = getIndentation(indentationCount)
	for _, deployment := range deployments.Items {
		fmt.Println(tempInd+"- Deployment - "+deployment.Name, " ", deployment.Status.ReadyReplicas, "/", deployment.Status.AvailableReplicas)

	}

}

func printReplicaSets(rsets *appsv1.ReplicaSetList) {
	fmt.Println("\nReplicaSets\n--------------\n")
	var tempInd string = getIndentation(indentationCount)
	for _, rs := range rsets.Items {
		fmt.Println(tempInd+"- ReplicaSets -", rs.Name, "", "", "", "")
	}

}

func printStateFulSets(ssets *appsv1.StatefulSetList) {
	fmt.Println("\nStatefulSets\n--------------\n")
	var tempInd string = getIndentation(indentationCount)
	for _, sset := range ssets.Items {
		fmt.Println("%v\t%v\t%v\n", sset.Name, sset.Status.ReadyReplicas, "")
		fmt.Println(tempInd+"- ReplicaSets -", sset.Name, ": Status - ", sset.Status.ReadyReplicas)
	}

}
