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

func setUpOverlayCommand() (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	cmd, buf, err := setUpCommand()
	_ = os.Chdir("testdata")
	_ = os.Chdir("overlay")
	base := afero.NewOsFs()
	roBase := afero.NewReadOnlyFs(base)
	ufs := afero.NewCopyOnWriteFs(roBase, afero.NewMemMapFs())
	fs.SetFsTo(ufs)

	return cmd, buf, err
}

func TestCreatesOverlayDir(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"create",
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
	cleanup()
}

func TestCreatesOverlayKustomizationFile(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"create",
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
	cleanup()
}

func TestCreatesOverlayDirWithSuffix(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"create",
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
	cleanup()
}

func TestCreatesOverlayNamespaceResourceFile(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"-n",
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
	cleanup()
}

func TestCreatesOverlayIngressResourceFile(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"create",
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
	cleanup()
}

func TestCreatesOverlayKustomizationWithIngress(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"create",
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
	cleanup()
}

func TestCreatesOverlayDeploymentMergeWithConfigKustomization(t *testing.T) {
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
	cleanup()
}

func TestCreatesProvidedConfigFile(t *testing.T) {
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
	cleanup()
}

func TestCreateOverlayDeploymentConfigMergePatch(t *testing.T) {
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
	cleanup()
}

func TestCreateOverlayDeploymentConfigMergePatchWithConfigPath(t *testing.T) {
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
	cleanup()
}

func TestCreatesOverlayDeploymentMergeWithSecretEnvKustomization(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	env := `ENVIRONMENT=DEVELOPMENT`
	cmd.SetArgs([]string{
		"create",
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
	cleanup()
}

func TestCreatesProvidedSecretEnv(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	env := `ENVIRONMENT=DEVELOPMENT`
	cmd.SetArgs([]string{
		"create",
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
	cleanup()
}

func TestCreatesOverlayDeploymentMergeWithSecretEnvPatch(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	env := `ENVIRONMENT=DEVELOPMENT`
	cmd.SetArgs([]string{
		"create",
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
	cleanup()
}

func TestCreateOverlayDeploymentReplicaMergePatch(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"-r", "2",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev", "deployment-replica-patch.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-replicas", "app-dev", "deployment-replica-patch.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesOverlayIngressResourceFileWithDifferentServicePort(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()
	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"--ingress=example.com",
		"-d=platform-with-different-service-port",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join("platform-with-different-service-port", "app-dev", "ingress.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-ingress-with-different-service-port", "app-dev", "ingress.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreateOverlayDeploymentResourceRequestPatch(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"--requests", "cpu=100m,memory=250Mi",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev", "deployment-requests-patch.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-requests", "app-dev", "deployment-requests-patch.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreateOverlayDeploymentResourceLimitsPatch(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"--limits", "cpu=100m,memory=250Mi",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, "app-dev", "deployment-limits-patch.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-limits", "app-dev", "deployment-limits-patch.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreateOverlayDeploymentMemoryResourceLimitsPatch(t *testing.T) {
	cmd, buf, err := setUpOverlayCommand()

	cmd.SetArgs([]string{
		"create",
		"overlay",
		"app-dev",
		"--requests", "memory=500Mi",
		"--limits", "memory=1Gi",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	pR := path.Join(platformDirDefault, "app-dev", "deployment-requests-patch.yaml")
	stat, fErr := fs.Get().Stat(pR)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	pL := path.Join(platformDirDefault, "app-dev", "deployment-limits-patch.yaml")
	stat, fErr = fs.Get().Stat(pL)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("overlayed-with-memory-resources", "app-dev", "deployment-requests-patch.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), pR)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))

	expect, _ = ioutil.ReadFile(filepath.Join("overlayed-with-memory-resources", "app-dev", "deployment-limits-patch.yaml"))
	actual, fErr = afero.ReadFile(fs.Get(), pL)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}
