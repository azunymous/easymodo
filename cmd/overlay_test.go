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

func TestCreatesOverlayDirWithSuffix(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"overlay",
		"-s", "suffix",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "suffix")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.True(t, stat.IsDir())
}

func TestCreatesOverlayNamespaceResourceFile(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"overlay",
		"app-dev",
		"-r",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev", "namespace.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed", "app-dev", "namespace.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
}

func TestCreatesOverlayIngressResourceFile(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"overlay",
		"app-dev",
		"--ingress=example.com",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev", "ingress.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-ingress", "app-dev", "ingress.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
}

func TestCreatesOverlayKustomizationWithIngress(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"overlay",
		"app-dev",
		"--ingress=example.com",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev", "kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-ingress", "app-dev", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
}

func TestCreatesOverlayDeploymentMergeWithConfigKustomization(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	configuration := `example:
  - a
  - b
  - c
configuration: test`

	cmd.SetArgs([]string{
		"overlay",
		"app-dev",
		"-c", "configuration.yaml=" + configuration,
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

func TestCreatesProvidedConfigFile(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	configuration := `example:
  - a
  - b
  - c
configuration: test`

	cmd.SetArgs([]string{
		"overlay",
		"app-dev",
		"-c",
		"configuration.yaml=" + configuration,
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev", "configuration.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-config", "app-dev", "configuration.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
}

func TestCreateOverlayDeploymentConfigMergePatch(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	configuration := `example:
  - a
  - b
  - c
configuration: test`

	cmd.SetArgs([]string{
		"overlay",
		"app-dev",
		"-c",
		"configuration.yaml=" + configuration,
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev", "deployment-config-patch.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-config", "app-dev", "deployment-config-patch.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
}

func TestCreateOverlayDeploymentConfigMergePatchWithConfigPath(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	configuration := `example:
  - a
  - b
  - c
configuration: test`

	cmd.SetArgs([]string{
		"overlay",
		"app-dev",
		"-c",
		"configuration.yaml=" + configuration,
		"-p", "/other/",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev", "deployment-config-patch.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("misc", "patch-with-different-config-path.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
}

func TestCreatesOverlayDeploymentMergeWithSecretEnvKustomization(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	env := `ENVIRONMENT=DEVELOPMENT`
	cmd.SetArgs([]string{
		"overlay",
		"app-dev",
		"-e", "dev.env=" + env,
	})

	println(cmd.Args)
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev", "kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-secret-env", "app-dev", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
}

func TestCreatesProvidedSecretEnv(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	env := `ENVIRONMENT=DEVELOPMENT`
	cmd.SetArgs([]string{
		"overlay",
		"app-dev",
		"-e", "dev.env=" + env,
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev", "dev.env")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-secret-env", "app-dev", "dev.env"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.Equal(t, string(expect), string(actual))
}

func TestCreatesOverlayDeploymentMergeWithSecretEnvPatch(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	env := `ENVIRONMENT=DEVELOPMENT`
	cmd.SetArgs([]string{
		"overlay",
		"app-dev",
		"-e", "dev.env=" + env,
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev", "deployment-secret-patch.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-secret-env", "app-dev", "deployment-secret-patch.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.Equal(t, string(expect), string(actual))
}
