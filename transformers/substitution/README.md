# Substitution Transformation

Adds support for [bash string replacement functions](https://github.com/drone/envsubst)

This is a simple kustomize plugin which allows you to perform environment substitution on your kustomize files. This mght be interesting if you want to add variables to exisitng components, without applying 1000s of patches.

The plugin is packaged as docker container (Containerized KRM) [Read More](https://kubectl.docs.kubernetes.io/guides/extending_kustomize/containerized_krm_functions/).

You can disable environment substitution for a single file by giving it an annotation (FluxCD native):

```
apiVersion: operatingsystemmanager.k8c.io/v1alpha1
kind: OperatingSystemProfile
metadata:
  name: osp-flatcar
  namespace: kube-system
  annotations:
    kustomize.toolkit.fluxcd.io/substitute: "disabled"
....   
```

If you want to skip a single variable, you can use `$${variable}`, this will print `${variable}`

To run it locally, you can use the following command:

```
kustomize build --enable-alpha-plugins $MYAPP
```

## Example

[See a simple exmaple here](../../exmaples/01-cluster-component/)

All your variables you want to substitude are declared under the `.values` section. As if now `map[string]string` is supported.

```
apiVersion: transformers.subst.github.com/v1
# That's relevant
kind: Substitution
metadata:
  name: transformer
  annotations:
    # That's relevant
    config.kubernetes.io/function: |
      container:
        image: ghcr.io/buttahtoast/transformers:latest  
values:
  cluster: "cluster-name"
  cluster_id: "cluster-id"
```

