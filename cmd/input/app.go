package input

type Application struct {
	Name          string
	Stateful      bool
	Image         string
	ContainerName string
	ContainerPort int
	Protocol      string
	Host          string
}
