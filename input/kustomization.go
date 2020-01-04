package input

// Kustomization defines the struct for what is required for a kustomization
type Kustomization struct {
	Res       []string
	Patches   []string
	Config    map[string][]string
	Secrets   map[string][]string
	Namespace string
}

// NewKustomization creates a new kustomization.
func NewKustomization(res []string, namespace string) *Kustomization {
	return &Kustomization{Res: res, Namespace: namespace}
}

// AddResource adds a new kubernetes resource or base to the kustomization
func (k *Kustomization) AddResource(fileName string) {
	k.Res = append(k.Res, fileName)
}

// AddPatch adds a new mergePatch to the kustomization
func (k *Kustomization) AddPatch(patchFilename string) {
	k.Patches = append(k.Patches, patchFilename)
}

// AddConfig adds a file to a config map generator in the kustomization
func (k *Kustomization) AddConfig(name, configFilename string) {
	k.Config[name] = append(k.Config[name], configFilename)
}

// AddSecret adds an env file to a secret generator in the kustomization
func (k *Kustomization) AddSecret(name, secretFilename string) {
	k.Secrets[name] = append(k.Secrets[name], secretFilename)
}
