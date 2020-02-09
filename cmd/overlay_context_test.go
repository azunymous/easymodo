package cmd

import (
	"github.com/azunymous/easymodo/fs"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"path/filepath"
	"testing"
)

func TestCreatesOverlayDirWithContextWithContext(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"--context", "usa",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "usa", "app-dev")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.True(t, stat.IsDir())
	cleanup()
}

func TestCreatesOverlayKustomizationFileWithContext(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"--context", "usa",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "usa", "app-dev", "kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-context", "usa", "app-dev", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesOverlayDirWithSuffixWithContext(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"create",
		"overlay",
		"-s", "suffix",
		"--context", "usa",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "usa", "suffix")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.True(t, stat.IsDir())
	cleanup()
}

func TestCreatesOverlayNamespaceResourceFileWithContext(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"-n",
		"--context", "usa",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "usa", "app-dev", "namespace.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-context", "usa", "app-dev", "namespace.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesOverlayIngressResourceFileWithContext(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"--ingress=example.com",
		"--context", "usa",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "usa", "app-dev", "ingress.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-ingress-with-context", "usa", "app-dev", "ingress.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesOverlayKustomizationWithIngressWithContext(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"--ingress=example.com",
		"--context", "usa",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "usa", "app-dev", "kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-ingress-with-context", "usa", "app-dev", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesOverlayDeploymentMergeWithConfigKustomizationWithContext(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	configuration := `example:
  - a
  - b
  - c
configuration: test`

	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"-c", "configuration.yaml=" + configuration,
		"--context", "usa",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "usa", "app-dev", "kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-config-with-context", "usa", "app-dev", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesProvidedConfigFileWithContext(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	configuration := `example:
  - a
  - b
  - c
configuration: test`

	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"-c",
		"configuration.yaml=" + configuration,
		"--context", "usa",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "usa", "app-dev", "configuration.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-config-with-context", "usa", "app-dev", "configuration.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreateOverlayDeploymentConfigMergePatchWithContext(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	configuration := `example:
  - a
  - b
  - c
configuration: test`

	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"-c",
		"configuration.yaml=" + configuration,
		"--context", "usa",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "usa", "app-dev", "deployment-config-patch.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-config-with-context", "usa", "app-dev", "deployment-config-patch.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesOverlayDeploymentMergeWithSecretEnvKustomizationWithContext(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	env := `ENVIRONMENT=DEVELOPMENT`
	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"-e", "dev.env=" + env,
		"--context", "usa",
	})

	println(cmd.Args)
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "usa", "app-dev", "kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-secret-env-with-context", "usa", "app-dev", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesProvidedSecretEnvWithContext(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	env := `ENVIRONMENT=DEVELOPMENT`
	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"-e", "dev.env=" + env,
		"--context", "usa",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "usa", "app-dev", "dev.env")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-secret-env-with-context", "usa", "app-dev", "dev.env"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.Equal(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesOverlayDeploymentMergeWithSecretEnvPatchWithContext(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	env := `ENVIRONMENT=DEVELOPMENT`
	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"-e", "dev.env=" + env,
		"--context", "usa",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "usa", "app-dev", "deployment-secret-patch.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-secret-env-with-context", "usa", "app-dev", "deployment-secret-patch.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.Equal(t, string(expect), string(actual))
	cleanup()
}

func TestCreateOverlayDeploymentReplicaMergePatchWithContext(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"-r", "2",
		"--context", "usa",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "usa", "app-dev", "deployment-replica-patch.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-replicas-with-context", "usa", "app-dev", "deployment-replica-patch.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}
