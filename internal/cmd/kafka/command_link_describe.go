package kafka

import (
	"github.com/spf13/cobra"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/cli/internal/pkg/output"
)

var (
	describeLinkConfigFields           = []string{"ConfigName", "ConfigValue", "ReadOnly", "Sensitive", "Source", "Synonyms"}
	structuredDescribeLinkConfigFields = camelToSnake(describeLinkConfigFields)
	humanDescribeLinkConfigFields      = camelToSpaced(describeLinkConfigFields)
)

type LinkConfigWriter struct {
	ConfigName  string
	ConfigValue string
	ReadOnly    bool
	Sensitive   bool
	Source      string
	Synonyms    []string
}

func (c *linkCommand) newDescribeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe <link>",
		Short: "Describe a previously created cluster link.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.describe,
	}

	pcmd.AddClusterFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddEnvironmentFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddContextFlag(cmd, c.CLICommand)
	pcmd.AddOutputFlag(cmd)

	return cmd
}

func (c *linkCommand) describe(cmd *cobra.Command, args []string) error {
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

	listLinkConfigsRespData, httpResp, err := kafkaREST.CloudClient.ListKafkaLinkConfigs(clusterId, linkName)
	if err != nil {
		return kafkaRestError(kafkaREST.CloudClient.GetUrl(), err, httpResp)
	}

	outputWriter, err := output.NewListOutputWriter(cmd, describeLinkConfigFields, humanDescribeLinkConfigFields, structuredDescribeLinkConfigFields)
	if err != nil {
		return err
	}

	if len(listLinkConfigsRespData.Data) < 1 {
		return outputWriter.Out()
	}

	outputWriter.AddElement(&LinkConfigWriter{
		ConfigName:  "dest.cluster.id",
		ConfigValue: listLinkConfigsRespData.Data[0].ClusterId,
		ReadOnly:    true,
		Sensitive:   true,
	})

	for _, config := range listLinkConfigsRespData.Data {
		outputWriter.AddElement(&LinkConfigWriter{
			ConfigName:  config.Name,
			ConfigValue: config.Value,
			ReadOnly:    config.ReadOnly,
			Sensitive:   config.Sensitive,
			Source:      config.Source,
			Synonyms:    config.Synonyms,
		})
	}

	return outputWriter.Out()
}
