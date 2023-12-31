package spare

import (
	"errors"
	"os"

	"github.com/charmbracelet/log"
	"github.com/nao1215/gorky/file"
	"github.com/nao1215/rainbow/cmd/subcmd"
	"github.com/nao1215/spare/config"
	"github.com/spf13/cobra"
)

// newInitCmd return init sub command.
func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "init",
		Short:   "Generate .spare.yml at current directory",
		Example: "   spare init",
		RunE: func(cmd *cobra.Command, args []string) error {
			return subcmd.Run(cmd, args, &initCmd{})
		},
	}
}

type initCmd struct{}

// Parse parses the arguments and flags.
func (i *initCmd) Parse(_ *cobra.Command, _ []string) error {
	return nil
}

// Do generate .spare.yml at current directory.
// If .spare.yml already exists, return error.
func (i *initCmd) Do() error {
	if file.IsFile(config.ConfigFilePath) {
		return config.ErrConfigFileAlreadyExists
	}

	file, err := os.Create(config.ConfigFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	if err := config.NewConfig().Write(file); err != nil {
		return err
	}
	log.Info("[ CREATE ]", "config file name", config.ConfigFilePath)
	log.Info("[  INFO  ] If you need to change the setting values, please refer to the documentation")
	// TODO: add link to documentation
	return nil
}
