package cmd

/*
Variable for the package containing shared flag information.
Features self described functions for getting these values and functions with the 'Flag' suffix for
getting the pointer to these values.
*/
var global = Flags{
	verify: false,
}

// ResetOptionalFlags is for resetting the flags. This is for testing purposes, when a command will
// called multiple times from a test.
func ResetOptionalFlags() {
	global.directory = "platform"
	global.context = ""

	global.configFiles = map[string]string{}
	global.secretEnvs = map[string]string{}
	global.namespaceResource = false
	global.suffix = ""
	global.verify = false
	global.ingress = ""
	global.replicas = 1
	global.image = ""
	global.kustomizations = []string{}
	global.limits = map[string]string{}
	global.requests = map[string]string{}
	global.output = ""
}

type Flags struct {
	configFiles       map[string]string
	secretEnvs        map[string]string
	directory         string
	context           string
	suffix            string
	namespace         string
	namespaceResource bool
	verify            bool
	ingress           string
	replicas          int
	image             string
	kustomizations    []string
	limits            map[string]string
	requests          map[string]string
	output            string
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

func Context() string {
	return global.context
}

func ContextFlag() *string {
	return &global.context
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
func Replicas() int {
	return global.replicas
}

func ReplicasFlag() *int {
	return &global.replicas
}

func Image() string {
	return global.image
}

func ImageFlag() *string {
	return &global.image
}
func Kustomizations() []string {
	return global.kustomizations
}

func KustomizationsFlag() *[]string {
	return &global.kustomizations
}

func Output() string {
	return global.output
}

func OutputFlag() *string {
	return &global.output
}

func Verify() bool {
	return global.verify
}

func VerifyFlag() *bool {
	return &global.verify
}

func Limits() map[string]string {
	return global.limits
}

func LimitsFlag() *map[string]string {
	return &global.limits
}

func Requests() map[string]string {
	return global.requests
}

func RequestsFlag() *map[string]string {
	return &global.requests
}
