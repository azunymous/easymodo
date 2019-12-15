# easymodo
WIP generator for creating Kubernetes and Kustomization resources.

This can be used to quickly bootstrap a project onto kubernetes into multiple environments by generating Kubernetes resources and kustomizations. `easymodo` creates very simple deployment and service resources that are suitable for many applications and you can overlay these with configuration/secret mounting per environment.

You can use these commands both imperatively (for bootstrapping) or declaratively (hooking into another build tool or applying the created resources immediately).

## Usage
`easymodo -h`

`easymodo init [application name]`

`easymodo overlay -s dev -c config.yaml="$(cat dev-config.yaml)"`

## Output

`init` generates base deployment, service, ingress (if enabled) and kustomization resource YAMLs in the `./platform/base` directory along
with a dev overlay.

Use `easymodo init -h` to see all possible configuration flags such as setting the output directory or changing the image name.

`overlay` defines a kustomization overlaying the base with a given namespace (via an argument or `-s`).

e.g To create an overlay for namespace `my-cool-app`:
```shell script
easymodo overlay my-cool-app
```
Or for automatically adding a suffix to the app name:
```shell script
easymodo overlay -s production
```
This will create the namespace `<app name>-production`

### Configuration Generator generator
Configuration files can be provided with the `-c` or `--config` flag. easymodo expects key value pairs
of `<new config file name>="<configuration file>"`. 

For example, to read in the config file and generate a config map and deployment patch:
```shell script
easymodo overlay app-dev -c config-yaml="$(cat ./dev-config.yaml)"
```

These will be mounted by default in `/config/` but the path can be overriden via the `-p` flag. 

### Secret Generator generator
Similiar to configuration files, secrets generated via .env files can be mounted on the application container. This is with the
`-e` flag, taking the form of key value pairs. Unlike configmaps which are mounted as files, easymodo expects .env files containing secrets. These are exposed on the pod as environment variables.
