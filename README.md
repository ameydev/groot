# Ktree A.K.A. Groot
kubernetes helper tool for finding out status of a kubernetes resource with their respective linked resources in a cluster.

# Tree Structured Map of K8S resources:
With just one `groot` command we can find the current relational mapping of the `k8s` resources deployed in a particular namespace.


# Troubleshoot the namespace:
Tired of firing ` kubectl get/describe pod/deploy/svc` commands just to know status of your deployments? why not just `groot` it! :) To see `status` of your `Services`, `deployments` and `pods`(and rest of the resources as well, Soon.)
We can even check if some `Configmap` or `Secrets` are not mapped with any `Pods` or not.


It will print the following information about the resource queried:
  - Status of the resource
  - Mappings/relations with other resources.
eg:

```
$ groot --namespace sock-shop

Namespace  -  sock-shop  Status:  
	|-------Service  -  carts  Status:  
		|-------Endpoint  -  carts  Status:  
		|-------Deployment  -  carts  Status:  1/1
			|-------Replicaset  -  carts-54d97c56b6  Status:  
				|-------Pod  -  carts-54d97c56b6-ngfhj  Status:  Running
					|-------Secret  -  default-token-qcjmp  Status:  
	|-------Service  -  carts-db  Status:  
		|-------Endpoint  -  carts-db  Status:  
		|-------Deployment  -  carts-db  Status:  2/2
			|-------Replicaset  -  carts-db-5678cc578f  Status:  
				|-------Pod  -  carts-db-5678cc578f-6jp5w  Status:  Running
					|-------Secret  -  default-token-qcjmp  Status:  
				|-------Pod  -  carts-db-5678cc578f-g7t5r  Status:  Running
					|-------Secret  -  default-token-qcjmp  Status:  
```

# Credits to ksearch:
The `groot` uses `ksearch` to get the list of k8s resources. `ksearch` is a kubectl plugin that will help us list all (literally all) the resources in a namespace and the resources can be searched using names as well.
Know more about `ksearch` (https://github.com/infracloudio/ksearch) 


# Installation 

1. Download the CLI
```
curl -LO https://github.com/ameydev/groot/releases/download/v0.0.1/groot

```
2. Make the binary executable
```
chmod +x ./groot

```
3. Move the binary in to your `PATH`.
```
sudo mv ./groot /usr/local/bin/groot

```
4. Test to ensure the version
```
groot --version

```


# Contributions:

Contributions and suggestions are always welcome. Feel free to fork the repo or file issue.
