package cmd

import (
	"bytes"
	"github.com/azunymous/easymodo/fs"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func setUpGroupCommand() (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	cmd, buf, err := setUpCommand()
	_ = os.Chdir("testdata")
	_ = os.Chdir("group")
	base := afero.NewOsFs()
	roBase := afero.NewReadOnlyFs(base)
	ufs := afero.NewCopyOnWriteFs(roBase, afero.NewMemMapFs())
	fs.SetFsTo(ufs)

	return cmd, buf, err
}

func TestCreatesGroupKustomization(t *testing.T) {
	cmd, buf, err := setUpGroupCommand()
	cmd.SetArgs([]string{
		"group",
		"-k platform/dev",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join("kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("group-dev", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesGroupKustomizationWithAbsPath(t *testing.T) {
	cmd, buf, err := setUpGroupCommand()
	cmd.SetArgs([]string{
		"group",
		"-k", path.Join(wd, "testdata", "group", "platform", "dev"),
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join("kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("group-dev", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}
