# TODO

- [ ] Document exported functions and packages
- [ ] Create development kustomization folder and resources
- [ ] Generate kustomization resource templates and Go code
- [ ] Generate files in tmp directory before copying to actual directory
- [ ] Idempotence
- [ ] Refactor tests to have less repeated code

# Plan
#### Declarative:
- *init* defines the base kustomization of the application
- *overlay* defines an overlay kustomization

#### Imperative
 - *add* creates additional namespace overlays.
