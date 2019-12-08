package input

type Kustomization struct {
	Res       []string
	Namespace string
}

func NewKustomization(res []string, namespace string) *Kustomization {
	return &Kustomization{Res: res, Namespace: namespace}
}
