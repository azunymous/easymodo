# TODO
- [ ] Document exported functions and packages
- [ ] Create development kustomization folder and resources
- [ ] Generate kustomization resource templates and Go code
- [ ] Generate files in tmp directory before copying to actual directory
- [ ] Idempotence
- [ ] Fix tests to not rely on resetting global flags

# Plan
#### Declarative:
- *init* defines the base kustomization of the application
- *overlay* defines an overlay kustomization
