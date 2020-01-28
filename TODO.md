# TODO
- [ ] Verify directory structure and all environments are buildable
- [ ] Use image kustomize feature instead of deployment patch for image command via flag
- [ ] Document exported functions and packages for resource templates
- [ ] Generate kustomization resource templates and Go code
- [ ] Generate files in tmp directory before copying to actual directory
- [ ] Idempotence
- [ ] Replace log.Fatalf scenarios with alternative exits and add unit tests 

---
# Cluster based directory structure
- [ ] Implement cluster based directory structure

Currently, easymodo assumes and creates a flat structure inside the platform directory.
```
-platform/
 ├── base/
 │   ├── deployment.yaml
 │   ├── service.yaml
 │   └── kustomization.yaml
 ├── dev/
 │   ├── ingress.yaml
 │   └── kustomization.yaml
 └── prod/
     ├── ingress.yaml
     ├── config.yaml
     ├── secrets.env
     ├── deployment-config-patch.yaml
     ├── deployment-secrets-patch.yaml
     └── kustomization.yaml
```

This is ideal for a single cluster, where you only need one context. It can also work for multiple 
clusters where your Kubernetes resources are identical across clusters per namespace.

For multiple clusters where environments can vary, there will be overlapping kubernetes resources in
some cases and vast differences in others, making the folder structure more complex.

## Potential layout
```
-platform/
 ├── base/
 │   ├── deployment.yaml
 │   ├── service.yaml
 │   └── kustomization.yaml
 ├── us-west-1/
 │   ├── dev/
 │   └── prod/
 └── us-east-1/
     ├── dev/
     ├── test/
     ├── stage/
     └── prod/
```