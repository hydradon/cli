package kafka

import (
	"fmt"

	"github.com/spf13/cobra"

	cmkv2 "github.com/confluentinc/ccloud-sdk-go-v2/cmk/v2"

	pcmd "github.com/confluentinc/cli/v3/pkg/cmd"
	"github.com/confluentinc/cli/v3/pkg/config"
	"github.com/confluentinc/cli/v3/pkg/errors"
	"github.com/confluentinc/cli/v3/pkg/examples"
	"github.com/confluentinc/cli/v3/pkg/form"
	"github.com/confluentinc/cli/v3/pkg/output"
)

func (c *clusterCommand) newUpdateCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update <id>",
		Short:             "Update a Kafka cluster.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: pcmd.NewValidArgsFunction(c.validArgs),
		RunE:              c.update,
		Annotations:       map[string]string{pcmd.RunRequirement: pcmd.RequireNonAPIKeyCloudLogin},
		Example: examples.BuildExampleString(
			examples.Example{
				Text: "Update the name and CKU count of a Kafka cluster:",
				Code: `confluent kafka cluster update lkc-123456 --name "New Cluster Name" --cku 3`,
			},
		),
	}

	cmd.Flags().String("name", "", "Name of the Kafka cluster.")
	cmd.Flags().Uint32("cku", 0, `Number of Confluent Kafka Units. For Kafka clusters of type "dedicated" only. When shrinking a cluster, you must reduce capacity one CKU at a time.`)
	pcmd.AddContextFlag(cmd, c.CLICommand)
	if cfg.IsCloudLogin() {
		pcmd.AddEnvironmentFlag(cmd, c.AuthenticatedCLICommand)
	}
	pcmd.AddOutputFlag(cmd)

	return cmd
}

func (c *clusterCommand) update(cmd *cobra.Command, args []string) error {
	if err := errors.CheckNoUpdate(cmd.Flags(), "cku", "name"); err != nil {
		return err
	}

	environmentId, err := c.Context.EnvironmentId()
	if err != nil {
		return err
	}

	id := args[0]
	currentCluster, _, err := c.V2Client.DescribeKafkaCluster(id, environmentId)
	if err != nil {
		return errors.NewErrorWithSuggestions(
			fmt.Sprintf(errors.KafkaClusterNotFoundErrorMsg, id),
			errors.ChooseRightEnvironmentSuggestions,
		)
	}

	update := cmkv2.CmkV2ClusterUpdate{
		Id:   cmkv2.PtrString(id),
		Spec: &cmkv2.CmkV2ClusterSpecUpdate{Environment: &cmkv2.EnvScopedObjectReference{Id: environmentId}},
	}

	if cmd.Flags().Changed("name") {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		if name == "" {
			return errors.New("`--name` flag value must not be empty")
		}
		update.Spec.SetDisplayName(name)
	}

	if cmd.Flags().Changed("cku") {
		cku, err := cmd.Flags().GetUint32("cku")
		if err != nil {
			return err
		}
		updatedCku, err := c.validateResize(int32(cku), &currentCluster)
		if err != nil {
			return err
		}
		update.Spec.Config = &cmkv2.CmkV2ClusterSpecUpdateConfigOneOf{CmkV2Dedicated: &cmkv2.CmkV2Dedicated{Kind: "Dedicated", Cku: updatedCku}}
	}

	updatedCluster, err := c.V2Client.UpdateKafkaCluster(id, update)
	if err != nil {
		return errors.NewWrapErrorWithSuggestions(err, "failed to update Kafka cluster", "A cluster can't be updated while still provisioning. If you just created this cluster, retry in a few minutes.")
	}

	ctx := c.Context.Config.Context()
	c.Context.Config.SetOverwrittenCurrentKafkaCluster(ctx.KafkaClusterContext.GetActiveKafkaClusterId())
	ctx.KafkaClusterContext.SetActiveKafkaCluster(id)

	return c.outputKafkaClusterDescription(cmd, &updatedCluster, true)
}

func (c *clusterCommand) validateResize(cku int32, currentCluster *cmkv2.CmkV2Cluster) (int32, error) {
	// Ensure the cluster is a Dedicated Cluster
	if currentCluster.GetSpec().Config.CmkV2Dedicated == nil {
		return 0, errors.New("failed to update Kafka cluster: cluster resize is only supported on dedicated clusters")
	}
	// Durability Checks
	if currentCluster.Spec.GetAvailability() == highAvailability && cku <= 1 {
		return 0, errors.New("`--cku` value must be greater than 1 for high durability")
	}
	if cku == 0 {
		return 0, errors.New(errors.CkuMoreThanZeroErrorMsg)
	}
	// Cluster can't be resized while it's provisioning or being expanded already.
	// Name _can_ be changed during these times, though.
	if err := isClusterResizeInProgress(currentCluster); err != nil {
		return 0, err
	}
	// If shrink
	if cku < currentCluster.GetSpec().Config.CmkV2Dedicated.Cku {
		promptMessage := ""
		// metrics api auth via jwt
		if err := c.validateKafkaClusterMetrics(currentCluster, true); err != nil {
			promptMessage += fmt.Sprintf("\n%v\n", err)
		}
		if err := c.validateKafkaClusterMetrics(currentCluster, false); err != nil {
			promptMessage += fmt.Sprintf("\n%v\n", err)
		}
		if promptMessage != "" {
			if ok, err := confirmShrink(promptMessage); !ok || err != nil {
				return 0, err
			}
		}
	}
	return cku, nil
}

func (c *clusterCommand) validateKafkaClusterMetrics(currentCluster *cmkv2.CmkV2Cluster, isLatestMetric bool) error {
	window := "3 day"
	if isLatestMetric {
		window = "15 min"
	}

	if err := c.validateClusterLoad(*currentCluster.Id, isLatestMetric); err != nil {
		return errors.Errorf("Looking at metrics in the last %s window:\n%v", window, err)
	}

	return nil
}

func confirmShrink(promptMessage string) (bool, error) {
	f := form.New(form.Field{ID: "proceed", Prompt: fmt.Sprintf("Validated cluster metrics and found that: %s\nDo you want to proceed with shrinking your kafka cluster?", promptMessage), IsYesOrNo: true})
	if err := f.Prompt(form.NewPrompt()); err != nil {
		return false, errors.New("cluster resize error: failed to read your confirmation")
	}
	if !f.Responses["proceed"].(bool) {
		output.Println(false, "Not proceeding with kafka cluster shrink")
		return false, nil
	}
	return true, nil
}
