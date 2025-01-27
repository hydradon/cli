package connect

import (
	"fmt"

	"github.com/spf13/cobra"

	connectv1 "github.com/confluentinc/ccloud-sdk-go-v2/connect/v1"

	pcmd "github.com/confluentinc/cli/v3/pkg/cmd"
	"github.com/confluentinc/cli/v3/pkg/errors"
	"github.com/confluentinc/cli/v3/pkg/examples"
	"github.com/confluentinc/cli/v3/pkg/output"
)

func (c *clusterCommand) newResumeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "resume <id-1> [id-2] ... [id-N]",
		Short:             "Resume connectors.",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: pcmd.NewValidArgsFunction(c.validArgsMultiple),
		RunE:              c.resume,
		Annotations:       map[string]string{pcmd.RunRequirement: pcmd.RequireNonAPIKeyCloudLogin},
		Example: examples.BuildExampleString(
			examples.Example{
				Text: `Resume connectors "lcc-000001" and "lcc-000002":`,
				Code: "confluent connect cluster resume lcc-000001 lcc-000002",
			},
		),
	}

	pcmd.AddClusterFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddContextFlag(cmd, c.CLICommand)
	pcmd.AddEnvironmentFlag(cmd, c.AuthenticatedCLICommand)

	return cmd
}

func (c *clusterCommand) resume(_ *cobra.Command, args []string) error {
	kafkaCluster, err := c.Context.GetKafkaClusterForCommand(c.V2Client)
	if err != nil {
		return err
	}

	environmentId, err := c.Context.EnvironmentId()
	if err != nil {
		return err
	}

	connectorsByName, err := c.V2Client.ListConnectorsWithExpansions(environmentId, kafkaCluster.ID, "id,info")
	if err != nil {
		return err
	}

	connectorsById := make(map[string]connectv1.ConnectV1ConnectorExpansion)
	for _, connector := range connectorsByName {
		connectorsById[connector.Id.GetId()] = connector
	}

	for _, id := range args {
		connector, ok := connectorsById[id]
		if !ok {
			return fmt.Errorf(errors.UnknownConnectorIdErrorMsg, id)
		}

		if err := c.V2Client.ResumeConnector(connector.Info.GetName(), environmentId, kafkaCluster.ID); err != nil {
			return err
		}

		output.Printf(c.Config.EnableColor, "Resumed connector \"%s\".\n", id)
	}

	return nil
}
