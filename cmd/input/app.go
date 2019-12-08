package input

type Application struct {
	Name          string
	Stateful      bool
	ImageName     string
	ContainerName string
	ContainerPort int
	Protocol      string
}
