package cmd

import (
	"github.com/spf13/cobra"

	dynamicconfig "github.com/confluentinc/cli/v3/pkg/dynamic-config"
	"github.com/confluentinc/cli/v3/pkg/version"
)

type CLICommand struct {
	*cobra.Command
	Config  *dynamicconfig.DynamicConfig
	Version *version.Version
}

func NewCLICommand(cmd *cobra.Command) *CLICommand {
	return &CLICommand{
		Command: cmd,
		Config:  &dynamicconfig.DynamicConfig{},
	}
}

func NewAnonymousCLICommand(cmd *cobra.Command, prerunner PreRunner) *CLICommand {
	c := NewCLICommand(cmd)
	cmd.PersistentPreRunE = chain(prerunner.Anonymous(c, false), prerunner.ParseFlagsIntoContext(c))
	return c
}
