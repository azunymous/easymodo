package cmd

/*
Variable for the package containing shared flag information.
Features self described functions for getting these values and functions with the 'Flag' suffix for
getting the pointer to these values.
*/
var global = Flags{
	force: false,
}

// ResetOptionalFlags is for resetting the flags. This is for testing purposes, when a command will
// called multiple times from a test.
func ResetOptionalFlags() {
	global.configFiles = map[string]string{}
	global.secretEnvs = map[string]string{}
	global.namespaceResource = false
	global.suffix = ""
	global.ingress = ""
}

type Flags struct {
	configFiles       map[string]string
	secretEnvs        map[string]string
	directory         string
	suffix            string
	namespace         string
	namespaceResource bool
	force             bool
	ingress           string
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

func Ingress() string {
	return global.ingress
}

func IngressFlag() *string {
	return &global.ingress
}

func Force() bool {
	return global.force
}

func ForceFlag() *bool {
	return &global.force
}
