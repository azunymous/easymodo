package input

var global = Flags{
	force: false,
}

type Flags struct {
	configFiles       map[string]string
	directory         string
	suffix            string
	namespace         string
	namespaceResource bool
	force             bool
}

func ConfigFiles() map[string]string {
	return global.configFiles
}

func ConfigFilesFlag() *map[string]string {
	return &global.configFiles
}

func Directory() string {
	return global.directory
}

func DirectoryFlag() *string {
	return &global.directory
}

func Suffix() string {
	return global.suffix
}

func SuffixFlag() *string {
	return &global.suffix
}

func Namespace() string {
	return global.namespace
}

func NamespaceFlag() *string {
	return &global.namespace
}

func NamespaceResource() bool {
	return global.namespaceResource
}

func NamespaceResourceFlag() *bool {
	return &global.namespaceResource
}

func Force() bool {
	return global.force
}

func ForceFlag() *bool {
	return &global.force
}
