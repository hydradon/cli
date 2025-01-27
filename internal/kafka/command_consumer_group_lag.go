package kafka

import (
	"github.com/spf13/cobra"

	"github.com/confluentinc/cli/v3/pkg/config"
)

type lagOut struct {
	ClusterId       string `human:"Cluster" serialized:"cluster_id"`
	ConsumerGroupId string `human:"Consumer Group" serialized:"consumer_group_id"`
	Lag             int64  `human:"Lag" serialized:"lag"`
	LogEndOffset    int64  `human:"Log End Offset" serialized:"log_end_offset"`
	CurrentOffset   int64  `human:"Current Offset" serialized:"current_offset"`
	ConsumerId      string `human:"Consumer" serialized:"consumer_id"`
	InstanceId      string `human:"Instance" serialized:"instance_id"`
	ClientId        string `human:"Client" serialized:"client_id"`
	Topic           string `human:"Topic" serialized:"topic"`
	PartitionId     int32  `human:"Partition" serialized:"partition_id"`
}

func (c *consumerCommand) newLagCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lag",
		Short: "View consumer group lag.",
	}

	if cfg.IsCloudLogin() {
		cmd.AddCommand(c.newLagDescribeCommand())
		cmd.AddCommand(c.newLagListCommand())
		cmd.AddCommand(c.newLagSummarizeCommand())
	} else {
		cmd.AddCommand(c.newLagDescribeCommandOnPrem())
		cmd.AddCommand(c.newLagListCommandOnPrem())
		cmd.AddCommand(c.newLagSummarizeCommandOnPrem())
	}

	return cmd
}
