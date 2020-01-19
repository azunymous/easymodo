# easymodo
WIP generator for creating Kubernetes and Kustomization resources.

This can be used to quickly bootstrap a project onto kubernetes into multiple environments by 
generating Kubernetes resources and kustomizations. `easymodo` creates very simple deployment 
and service resources that are suitable for many applications and you can overlay these with 
configuration/secret mounting per environment.

- `create` commands are useful for bootstrapping and then changing to suit your needs.
- `modify` commands can be used predeploy such as for changing the version of the container.

You can use these commands both imperatively (for bootstrapping) or declarative (hooking into another build tool or applying the created resources immediately).


## Usage
`easymodo -h`

`easymodo create base [application name]`

`easymodo create overlay -s dev -c config.yaml="$(cat dev-config.yaml)"`

`easymodo modify image -s dev -i web-app:v1.2.3`

## Output

### Create
`create base` generates base deployment, service, ingress (if enabled) and kustomization resource YAMLs in the `./platform/base` directory along
with a dev overlay.

Use `easymodo base -h` to see all possible configuration flags such as setting the output directory or changing the image name.

`create overlay` defines a kustomization overlaying the base with a given namespace (via an argument or `-s`).

e.g To create an overlay for namespace `my-cool-app`:
```shell script
easymodo create overlay my-cool-app
```
Or for automatically adding a suffix to the app name:
```shell script
easymodo create overlay -s production
```
This will create the namespace `<app name>-production`

#### Configuration Generator generator
Configuration files can be provided with the `-c` or `--config` flag. easymodo expects key value pairs
of `<new config file name>="<configuration file>"`. 

For example, to read in the config file and generate a config map and deployment patch:
```shell script
easymodo create overlay app-dev -c config-yaml="$(cat ./dev-config.yaml)"
```

These will be mounted by default in `/config/` but the path can be overriden via the `-p` flag. 

#### Secret Generator generator
Similiar to configuration files, secrets generated via .env files can be mounted on the application container. This is with the
`-e` flag, taking the form of key value pairs. Unlike configmaps which are mounted as files, easymodo expects .env files containing secrets. These are exposed on the pod as environment variables.

### Modify
`modify image` generates an overlay with a different image.
For example:
```shell script
easymodo modify image app-dev -i gcr.io/dev/app:v1.2.3
```
With a suffix:
```shell script
easymodo modify image -s dev -i gcr.io/dev/app:v1.2.3
```

This can then be deployed inline:
```shell script
easymodo modify image app-dev -i gcr.io/my-project/dev/app:v1.2.3 | xargs kustomize build | kubectl apply -f - 
```