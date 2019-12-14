# groot
kubernetes helper tool for finding out status of a kubernetes resource with their respective linked resources in a cluster.

It will print the following information about the resource queried:
  - Status of the resource
  - Mappings/relations with other resources.
eg:

```
$ groot -n dev

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
