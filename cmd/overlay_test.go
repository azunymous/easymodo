package cmd

import (
	"bytes"
	"github.com/azunymous/easymodo/cmd/fs"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func setUpOverlayCommand() (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	cmd, buf, err := setUpCommand()
	_ = os.Chdir("testdata")
	base := afero.NewOsFs()
	roBase := afero.NewReadOnlyFs(base)
	ufs := afero.NewCopyOnWriteFs(roBase, afero.NewMemMapFs())
	fs.SetFsTo(ufs)

	return cmd, buf, err
}

func TestCreatesOverlayDir(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"overlay",
		"app-dev",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.True(t, stat.IsDir())
}

func TestCreatesOverlayKustomizationFile(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"overlay",
		"app-dev",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev", "kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed", "app-dev", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
}

func TestCreatesOverlayDeploymentMergeKustomization(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	configuration := `example:
  - a
  - b
  - c
configuration: test`

	cmd.SetArgs([]string{
		"overlay",
		"app-dev",
		"-c configuration.yaml=" + configuration,
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev", "kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-config", "app-dev", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
}
