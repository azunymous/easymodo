package cmd

var global = Flags{
	force: false,
}

func ResetOptionalFlags() {
	global.configFiles = map[string]string{}
	global.secretEnvs = map[string]string{}
	global.suffix = ""
}

type Flags struct {
	configFiles       map[string]string
	secretEnvs        map[string]string
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

func SecretEnvs() map[string]string {
	return global.secretEnvs
}

func SecretEnvsFlag() *map[string]string {
	return &global.secretEnvs
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
