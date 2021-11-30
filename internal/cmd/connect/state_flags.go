package connect

import (
	"github.com/spf13/pflag"

	"github.com/confluentinc/cli/internal/pkg/cmd"
)

var clusterSubcommandFlags = map[string]*pflag.FlagSet{
	"list": cmd.ContextSet(),
}

var subcommandFlags = map[string]*pflag.FlagSet{
	"describe": cmd.ClusterEnvironmentContextSet(),
	"list":     cmd.ClusterEnvironmentContextSet(),
	"create":   cmd.ClusterEnvironmentContextSet(),
	"delete":   cmd.ClusterEnvironmentContextSet(),
	"update":   cmd.ClusterEnvironmentContextSet(),
	"pause":    cmd.ClusterEnvironmentContextSet(),
	"resume":   cmd.ClusterEnvironmentContextSet(),
}
