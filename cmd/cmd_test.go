package cmd

import (
	tp "github.com/dsrvlabs/vatz/manager/types"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestInitCmd(t *testing.T) {
	tests := []struct {
		Desc       string
		Args       []string
		ExpectFile string
	}{
		// TODO: Testing default config could corrupt real environment.
		// So, remove it intensionally.
		//{
		//	Desc:       "Init with default",
		//	Args:       []string{"init"},
		//	ExpectFile: "default.yaml",
		//},
		{
			Desc:       "Init with selected filename",
			Args:       []string{"init", "--output", "hello.yaml", "--home", "./"},
			ExpectFile: "hello.yaml",
		},
	}

	for _, test := range tests {
		root := cobra.Command{}
		root.AddCommand(createInitCommand(tp.TEST))
		root.SetArgs(test.Args)

		err := root.Execute()
		defer os.Remove(test.ExpectFile)

		assert.Nil(t, err)

		_, err = os.ReadFile(test.ExpectFile)
		assert.Nil(t, err)
	}
}
