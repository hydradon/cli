package flink

import (
	"github.com/spf13/cobra"

	flinkv2 "github.com/confluentinc/ccloud-sdk-go-v2/flink/v2"

	pcmd "github.com/confluentinc/cli/v3/pkg/cmd"
	"github.com/confluentinc/cli/v3/pkg/errors"
	"github.com/confluentinc/cli/v3/pkg/examples"
	"github.com/confluentinc/cli/v3/pkg/output"
)

func (c *command) newComputePoolUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update [id]",
		Short:             "Update a Flink compute pool.",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: pcmd.NewValidArgsFunction(c.validComputePoolArgs),
		RunE:              c.computePoolUpdate,
		Example: examples.BuildExampleString(
			examples.Example{
				Text: `Update name and CFU count of a Flink compute pool.`,
				Code: `confluent flink compute-pool update my-compute-pool --name "new name" --cfu 2`,
			},
		),
	}

	cmd.Flags().String("name", "", "Name of the compute pool.")
	cmd.Flags().Int32("cfu", 0, "Number of Confluent Flink Units (CFU).")
	pcmd.AddEnvironmentFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddOutputFlag(cmd)

	return cmd
}

func (c *command) computePoolUpdate(cmd *cobra.Command, args []string) error {
	if err := errors.CheckNoUpdate(cmd.Flags(), "name", "cfu"); err != nil {
		return err
	}

	id := c.Context.GetCurrentFlinkComputePool()
	if len(args) > 0 {
		id = args[0]
	}
	if id == "" {
		return errors.NewErrorWithSuggestions("no compute pool selected", "Select a compute pool with `confluent flink compute-pool use` or as an argument.")
	}

	environmentId, err := c.Context.EnvironmentId()
	if err != nil {
		return err
	}

	computePool, err := c.V2Client.DescribeFlinkComputePool(id, environmentId)
	if err != nil {
		return err
	}

	environment, err := c.V2Client.GetOrgEnvironment(environmentId)
	if err != nil {
		return err
	}

	update := flinkv2.FcpmV2ComputePoolUpdate{
		Id: flinkv2.PtrString(id),
		Spec: &flinkv2.FcpmV2ComputePoolSpecUpdate{
			MaxCfu: flinkv2.PtrInt32(computePool.Spec.GetMaxCfu()),
			Environment: &flinkv2.GlobalObjectReference{
				Id:           environmentId,
				Related:      environment.Metadata.GetSelf(),
				ResourceName: environment.Metadata.GetResourceName(),
			},
		},
	}

	cfu, err := cmd.Flags().GetInt32("cfu")
	if err != nil {
		return err
	}
	if cfu != 0 {
		update.Spec.MaxCfu = flinkv2.PtrInt32(cfu)
	}

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	if name != "" {
		update.Spec.DisplayName = flinkv2.PtrString(name)
	}

	updatedComputePool, err := c.V2Client.UpdateFlinkComputePool(id, update)
	if err != nil {
		return err
	}

	table := output.NewTable(cmd)
	table.Add(&computePoolOut{
		IsCurrent: computePool.GetId() == c.Context.GetCurrentFlinkComputePool(),
		Id:        computePool.GetId(),
		Name:      updatedComputePool.Spec.GetDisplayName(),
		Cfu:       updatedComputePool.Spec.GetMaxCfu(),
		Region:    computePool.Spec.GetRegion(),
		Status:    computePool.Status.GetPhase(),
	})
	return table.Print()
}