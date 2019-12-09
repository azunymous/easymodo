package cmd

import (
	"bytes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"path/filepath"
	"testing"
)

const platformDirDefault = "platform"
const baseDirDefault = "base"

func TestCreatesPlatformDir(t *testing.T) {
	appFs = afero.NewMemMapFs()
	cmd := rootCmd
	buf := new(bytes.Buffer)
	err := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(err)
	cmd.SetArgs([]string{
		"init",
		"app",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	stat, fErr := appFs.Stat(platformDirDefault)
	assert.Nil(t, fErr)
	assert.True(t, stat.IsDir())
}

func TestCreatesBaseDir(t *testing.T) {
	appFs = afero.NewMemMapFs()
	cmd := rootCmd
	buf := new(bytes.Buffer)
	err := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(err)
	cmd.SetArgs([]string{
		"init",
		"app",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	stat, fErr := appFs.Stat(path.Join(platformDirDefault, baseDirDefault))
	assert.Nil(t, fErr)
	assert.True(t, stat.IsDir())
}

func TestCreatesKustomizationFile(t *testing.T) {
	appFs = afero.NewMemMapFs()
	cmd := rootCmd
	buf := new(bytes.Buffer)
	err := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(err)
	cmd.SetArgs([]string{
		"init",
		"app",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, baseDirDefault, "kustomization.yaml")
	stat, fErr := appFs.Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("testdata", "basic", "kustomization.yaml"))
	actual, fErr := afero.ReadFile(appFs, p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
}

func TestCreatesDeploymentFile(t *testing.T) {
	appFs = afero.NewMemMapFs()
	cmd := rootCmd
	buf := new(bytes.Buffer)
	err := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(err)
	cmd.SetArgs([]string{
		"init",
		"app",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, baseDirDefault, "deployment.yaml")
	stat, fErr := appFs.Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("testdata", "basic", "deployment.yaml"))
	actual, fErr := afero.ReadFile(appFs, p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
}

func TestCreatesServiceFile(t *testing.T) {
	appFs = afero.NewMemMapFs()
	cmd := rootCmd
	buf := new(bytes.Buffer)
	err := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(err)
	cmd.SetArgs([]string{
		"init",
		"app",
	})

	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	p := path.Join(platformDirDefault, baseDirDefault, "service.yaml")
	stat, fErr := appFs.Stat(p)
	assert.Nil(t, fErr)
	assert.False(t, stat.IsDir())

	expect, _ := ioutil.ReadFile(filepath.Join("testdata", "basic", "service.yaml"))
	actual, fErr := afero.ReadFile(appFs, p)
	if fErr != nil {
		t.Fatal(fErr)
	}
	assert.YAMLEq(t, string(expect), string(actual))
}

func TestCreatesFlagPlatformDir(t *testing.T) {
	appFs = afero.NewMemMapFs()
	cmd := rootCmd
	buf := new(bytes.Buffer)
	err := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(err)
	cmd.SetArgs([]string{
		"init",
		"app",
		"-d", "flagPlatform",
	})
	_ = cmd.Execute()
	println(buf.String())
	println(err.String())

	stat, fErr := appFs.Stat("flagPlatform")
	assert.Nil(t, fErr)
	assert.True(t, stat.IsDir())
}
