package cmd

import (
	"bytes"
	"github.com/azunymous/easymodo/fs"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func setUpVerifyCommand() (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	cmd, buf, err := setUpCommand()
	_ = os.Chdir("testdata")
	_ = os.Chdir("verify")
	base := afero.NewOsFs()
	roBase := afero.NewReadOnlyFs(base)
	ufs := afero.NewCopyOnWriteFs(roBase, afero.NewMemMapFs())
	fs.SetFsTo(ufs)

	return cmd, buf, err
}

func TestVerifiesKustomization(t *testing.T) {
	cmd, _, _ := setUpVerifyCommand()
	cmd.SetArgs([]string{
		"verify",
	})

	assert.NotPanics(t, func() { _ = cmd.Execute() })
	cleanup()
}
