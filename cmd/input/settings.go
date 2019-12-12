package input

var global = Flags{
	directory: "platform",
}

type Flags struct {
	directory string
}

func Directory() string {
	return global.directory
}

func DirectoryFlag() *string {
	return &global.directory
}
