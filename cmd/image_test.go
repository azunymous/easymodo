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
	"strings"
	"testing"
)

func setUpImageCommand() (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	cmd, buf, err := setUpCommand()
	_ = os.Chdir("testdata")
	_ = os.Chdir("image")
	base := afero.NewOsFs()
	roBase := afero.NewReadOnlyFs(base)
	ufs := afero.NewCopyOnWriteFs(roBase, afero.NewMemMapFs())
	fs.SetFsTo(ufs)

	return cmd, buf, err
}

func TestCreatesImageOverlayDir(t *testing.T) {
	cmd, buf, err := setUpImageCommand()
	cmd.SetArgs([]string{
		"modify",
		"image",
		"app-dev",
		"-i app:v1.0.0",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev-v1.0.0")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.True(t, stat.IsDir())
	cleanup()
}

func TestCreatesImageOverlayKustomization(t *testing.T) {
	cmd, buf, err := setUpImageCommand()
	cmd.SetArgs([]string{
		"modify",
		"image",
		"app-dev",
		"-i app:v1.0.0",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev-v1.0.0", "kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("different-image", "app-dev-v1.0.0", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesImageOverlayKustomizationForSuffix(t *testing.T) {
	cmd, buf, err := setUpImageCommand()
	cmd.SetArgs([]string{
		"modify",
		"image",
		"-s", "dev",
		"-i app:v1.0.0",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "dev-v1.0.0", "kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("different-image-suffix", "dev-v1.0.0", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesImageOverlayPatch(t *testing.T) {
	cmd, buf, err := setUpImageCommand()
	cmd.SetArgs([]string{
		"modify",
		"image",
		"app-dev",
		"-i app:v1.0.0",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev-v1.0.0", "deployment-image-patch.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("different-image", "app-dev-v1.0.0", "deployment-image-patch.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestOutputsImageOverlayDirectory(t *testing.T) {
	cmd, _, _ := setUpImageCommand()
	cmd.SetArgs([]string{
		"modify",
		"image",
		"app-dev",
		"-i app:v1.0.0",
	})

	f, _ := ioutil.TempFile(os.TempDir(), "*")
	w = f
	_ = cmd.Execute()
	out, _ := ioutil.ReadFile(f.Name())
	assert.Equal(t, path.Join(wd, "testdata", "image", "platform", "app-dev-v1.0.0"), strings.TrimSpace(string(out)))
	cleanup()
}

func TestCreatesImageOverlayDirForUnknownVersion(t *testing.T) {
	cmd, buf, err := setUpImageCommand()
	cmd.SetArgs([]string{
		"modify",
		"image",
		"app-dev",
		"-i appWithNoVersion.com/my/cool/app",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev-UNKNOWN")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.True(t, stat.IsDir())
	cleanup()
}

func TestCreatesImageOverlayDirForAlphabeticalVersion(t *testing.T) {
	cmd, buf, err := setUpImageCommand()
	cmd.SetArgs([]string{
		"modify",
		"image",
		"app-dev",
		"-i appWithVersion.com/app:latest",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev-latest")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.True(t, stat.IsDir())
	cleanup()
}

func TestCreatesImageOverlayDirWithContext(t *testing.T) {
	cmd, buf, err := setUpImageCommand()
	cmd.SetArgs([]string{
		"modify",
		"image",
		"-s", "dev",
		"-i app:v1.0.0",
		"-d", "platform-with-context",
		"--context", "usa",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join("platform-with-context", "usa", "dev-v1.0.0")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.True(t, stat.IsDir())
	cleanup()
}

func TestCreatesImageOverlayKustomizationWithContext(t *testing.T) {
	cmd, buf, err := setUpImageCommand()
	cmd.SetArgs([]string{
		"modify",
		"image",
		"-s", "dev",
		"-i app:v1.0.0",
		"-d", "platform-with-context",
		"--context", "usa",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join("platform-with-context", "usa", "dev-v1.0.0", "kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("different-image-with-context", "usa", "dev-v1.0.0", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}
