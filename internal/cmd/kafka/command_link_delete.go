package kafka

import (
	"github.com/spf13/cobra"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/cli/internal/pkg/resource"
	"github.com/confluentinc/cli/internal/pkg/utils"
)

func (c *linkCommand) newDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <link>",
		Short: "Delete a previously created cluster link.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.delete,
	}

	pcmd.AddClusterFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddEnvironmentFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddContextFlag(cmd, c.CLICommand)

	return cmd
}

func (c *linkCommand) delete(cmd *cobra.Command, args []string) error {
	linkName := args[0]

	kafkaREST, err := c.GetKafkaREST()
	if kafkaREST == nil {
		if err != nil {
			return err
		}
		return errors.New(errors.RestProxyNotAvailableMsg)
	}

	clusterId, err := getKafkaClusterLkcId(c.AuthenticatedStateFlagCommand)
	if err != nil {
		return err
	}

	if httpResp, err := kafkaREST.CloudClient.DeleteKafkaLink(clusterId, linkName); err != nil {
		return kafkaRestError(kafkaREST.CloudClient.GetUrl(), err, httpResp)
	}

	utils.Printf(cmd, errors.DeletedResourceMsg, resource.ClusterLink, linkName)
	return nil
}
