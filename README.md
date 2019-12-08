# easymodo
WIP generator for creating initial kustomization files

## Usage

`easymodo init [application name]`

`easymodo -h`

## Output

Generates base deployment, service, ingress (if enabled) and kustomization resource YAMLs in the `./platform/base` directory along
with a dev overlay.

Use `easymodo init -h` to see all possible configuration flags such as setting the output directory or changing the image name.

