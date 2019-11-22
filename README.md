# groot
kubernetes helper utility cli tool for finding out the k8s resources in the cluster and show them in mapped with each other

It should give the overview of all(or some in case we input the list of it having CRDs) the resources, showing
  - Status of the resource
  - mappings/relations with each other. etc
eg. 

```
$ ktree -n dev

service: frontend (status:ok)
- deployment (status:ok)
    - 3 pods (status:ok)
       --  configmap-1 (status:ok)
       --  volume 1 (status:ok)

service: backend (status:some error with volume-1)
- deployment 
    - 5 pods (status:ok)
       --  configmap-1 (status:ok)
       -- volume -1 ( status: some error )
       --> label matched with DB
```
