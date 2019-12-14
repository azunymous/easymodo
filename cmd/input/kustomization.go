package input

type Kustomization struct {
	Res       []string
	Patches   []string
	Config    map[string][]string
	Namespace string
}

func NewKustomization(res []string, namespace string) *Kustomization {
	return &Kustomization{Res: res, Namespace: namespace}
}

func (k *Kustomization) AddResource(fileName string) {
	k.Res = append(k.Res, fileName)
}

func (k *Kustomization) AddPatch(patchFileName string) {
	k.Patches = append(k.Patches, patchFileName)
}

func (k *Kustomization) AddConfig(name, patchFileName string) {
	k.Config[name] = append(k.Config[name], patchFileName)
}
