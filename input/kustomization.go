package input

type Kustomization struct {
	Res       []string
	Patches   []string
	Config    map[string][]string
	Secrets   map[string][]string
	Namespace string
}

func NewKustomization(res []string, namespace string) *Kustomization {
	return &Kustomization{Res: res, Namespace: namespace}
}

func (k *Kustomization) AddResource(fileName string) {
	k.Res = append(k.Res, fileName)
}

func (k *Kustomization) AddPatch(patchFilename string) {
	k.Patches = append(k.Patches, patchFilename)
}

func (k *Kustomization) AddConfig(name, configFilename string) {
	k.Config[name] = append(k.Config[name], configFilename)
}

func (k *Kustomization) AddSecret(name, secretFilename string) {
	k.Secrets[name] = append(k.Secrets[name], secretFilename)
}
