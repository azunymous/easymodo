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

func TestCreatesGroupKustomizationWithAbsPathFromDifferentRelativeOutputDirectory(t *testing.T) {
	cmd, buf, err := setUpGroupCommand()
	cmd.SetArgs([]string{
		"group",
		"-k", path.Join(wd, "testdata", "group", "platform", "dev"),
		"-o", "../",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join("../kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("group-dev-different-output", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesGroupKustomizationWithAbsPathAndAbsoluteOutputDirectory(t *testing.T) {
	cmd, buf, err := setUpGroupCommand()
	cmd.SetArgs([]string{
		"group",
		"-k", path.Join(wd, "testdata", "group", "platform", "dev"),
		"-o", wd,
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(wd, "kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("group-dev-abs-output", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesGroupKustomizationWithRelativePathAndRelativeOutputDirectory(t *testing.T) {
	cmd, buf, err := setUpGroupCommand()
	cmd.SetArgs([]string{
		"group",
		"-k", "platform/dev",
		"-o", "../release",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join("../release/kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("group-dev-relative-output", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesGroupKustomizationWithMultipleFiles(t *testing.T) {
	cmd, buf, err := setUpGroupCommand()
	cmd.SetArgs([]string{
		"group",
		"-k platform/dev",
		"-k", path.Join(wd, "testdata", "group", "platform", "app-dev"),
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join("kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("group-dev-multiple", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesGroupKustomizationPassesWithVerifcationForMultipleFiles(t *testing.T) {
	cmd, buf, err := setUpGroupCommand()
	cmd.SetArgs([]string{
		"group",
		"-k", "platform/dev",
		"-k", path.Join(wd, "testdata", "group", "platform", "app-dev"),
		"-v",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join("kustomization.yaml")
	stat, fErr := fs.Get().Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("group-dev-multiple", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(fs.Get(), p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
	cleanup()
}

func TestCreatesGroupKustomizationFailsWithVerification(t *testing.T) {
	cmd, buf, err := setUpGroupCommand()
	cmd.SetArgs([]string{
		"group",
		"-k", "platform/does-not-exist",
		"-v",
	})

	assert.Panics(t, func() { _ = cmd.Execute() }, "The command did not panic")
	println(buf.String())
	println(err.String())

	p := path.Join("kustomization.yaml")
	_, fErr := fs.Get().Stat(p)
	assert.NotNil(t, fErr)
	cleanup()
}

func TestCreatesGroupKustomizationFailsWithFile(t *testing.T) {
	cmd, buf, err := setUpGroupCommand()
	cmd.SetArgs([]string{
		"group",
		"-k", "platform/dev/kustomization.yaml",
		"-v",
	})

	assert.Panics(t, func() { _ = cmd.Execute() }, "The command did not panic")
	println(buf.String())
	println(err.String())

	p := path.Join("kustomization.yaml")
	_, fErr := fs.Get().Stat(p)
	assert.NotNil(t, fErr)
	cleanup()
}
