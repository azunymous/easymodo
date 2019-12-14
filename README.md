# easymodo
WIP generator for creating kustomization files

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