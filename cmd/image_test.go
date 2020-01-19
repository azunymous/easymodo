package cmd

import (
	"bytes"
	"github.com/azunymous/easymodo/fs"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func setUpVersionCommand() (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	cmd, buf, err := setUpCommand()
	wd, _ := os.Getwd()
	println(wd)
	_ = os.Chdir("testdata")
	_ = os.Chdir("image")
	base := afero.NewOsFs()
	roBase := afero.NewReadOnlyFs(base)
	ufs := afero.NewCopyOnWriteFs(roBase, afero.NewMemMapFs())
	fs.SetFsTo(ufs)

	return cmd, buf, err
}

func TestCreatesVersionDir(t *testing.T) {
	cmd, buf, err := setUpVersionCommand()
	cmd.SetArgs([]string{
		"modify",
		"image",
		"app-dev",
		"-v v1.0.0",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev-temp")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.True(t, stat.IsDir())
	cleanup()
}
