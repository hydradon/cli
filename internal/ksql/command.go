package ksql

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	ksqlv2 "github.com/confluentinc/ccloud-sdk-go-v2/ksql/v2"

	pauth "github.com/confluentinc/cli/v3/pkg/auth"
	"github.com/confluentinc/cli/v3/pkg/ccloudv2"
	pcmd "github.com/confluentinc/cli/v3/pkg/cmd"
	"github.com/confluentinc/cli/v3/pkg/config"
)

type ksqlCommand struct {
	*pcmd.AuthenticatedCLICommand
}

type ksqlCluster struct {
	Id                    string `human:"ID" serialized:"id"`
	Name                  string `human:"Name" serialized:"name"`
	OutputTopicPrefix     string `human:"Topic Prefix" serialized:"topic_prefix"`
	KafkaClusterId        string `human:"Kafka" serialized:"kafka"`
	Storage               int32  `human:"Storage" serialized:"storage"`
	Endpoint              string `human:"Endpoint" serialized:"endpoint"`
	Status                string `human:"Status" serialized:"status"`
	DetailedProcessingLog bool   `human:"Detailed Processing Log" serialized:"detailed_processing_log"`
}

func New(cfg *config.Config, prerunner pcmd.PreRunner) *cobra.Command {
	cmd := &cobra.Command{
		Use:         "ksql",
		Short:       "Manage ksqlDB.",
		Annotations: map[string]string{pcmd.RunRequirement: pcmd.RequireNonAPIKeyCloudLoginOrOnPremLogin},
	}

	cmd.AddCommand(newClusterCommand(cfg, prerunner))

	return cmd
}

func (c *ksqlCommand) formatClusterForDisplayAndList(cluster *ksqlv2.KsqldbcmV2Cluster) *ksqlCluster {
	detailedProcessingLog := true
	if cluster.Spec.HasUseDetailedProcessingLog() {
		detailedProcessingLog = cluster.Spec.GetUseDetailedProcessingLog()
	}

	return &ksqlCluster{
		Id:                    cluster.GetId(),
		Name:                  cluster.Spec.GetDisplayName(),
		OutputTopicPrefix:     cluster.Status.GetTopicPrefix(),
		KafkaClusterId:        cluster.Spec.KafkaCluster.GetId(),
		Storage:               cluster.Status.GetStorage(),
		Endpoint:              cluster.Status.GetHttpEndpoint(),
		Status:                c.getClusterStatus(cluster),
		DetailedProcessingLog: detailedProcessingLog,
	}
}

func (c *ksqlCommand) getClusterStatus(cluster *ksqlv2.KsqldbcmV2Cluster) string {
	status := cluster.Status.GetPhase()
	if cluster.Status.GetIsPaused() {
		status = "PAUSED"
	} else if status == "PROVISIONED" {
		provisioningFailed, err := c.checkProvisioningFailed(cluster.Status.GetHttpEndpoint())
		if err != nil {
			status = "UNKNOWN"
		} else if provisioningFailed {
			status = "PROVISIONING FAILED"
		}
	}
	return status
}

func (c *ksqlCommand) checkProvisioningFailed(endpoint string) (bool, error) {
	state, err := c.Context.AuthenticatedState()
	if err != nil {
		return false, err
	}

	dataplaneToken, err := pauth.GetDataplaneToken(state, c.Context.GetPlatformServer())
	if err != nil {
		return false, err
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: dataplaneToken})

	slingClient := sling.New().Client(oauth2.NewClient(context.Background(), ts)).Base(endpoint)
	var failure map[string]any
	response, err := slingClient.New().Get("/info").Receive(nil, &failure)
	if err != nil || response == nil {
		return false, err
	}

	if response.StatusCode == http.StatusServiceUnavailable {
		errorCode, ok := failure["error_code"].(float64)
		if !ok {
			return false, fmt.Errorf("failed to cast 'error_code' to float64")
		}
		// If we have a 50321 we know that ACLs are misconfigured
		if int(errorCode) == 50321 {
			return true, nil
		}
	}
	return false, nil
}

func (c *ksqlCommand) validArgs(cmd *cobra.Command, args []string) []string {
	if len(args) > 0 {
		return nil
	}

	return c.validArgsMultiple(cmd, args)
}

func (c *ksqlCommand) validArgsMultiple(cmd *cobra.Command, args []string) []string {
	if err := c.PersistentPreRunE(cmd, args); err != nil {
		return nil
	}

	environmentId, err := c.Context.EnvironmentId()
	if err != nil {
		return nil
	}

	return autocompleteClusters(environmentId, c.V2Client)
}

func autocompleteClusters(environment string, client *ccloudv2.Client) []string {
	clusters, err := client.ListKsqlClusters(environment)
	if err != nil {
		return nil
	}

	suggestions := make([]string, len(clusters))
	for i, cluster := range clusters {
		suggestions[i] = fmt.Sprintf("%s\t%s", cluster.GetId(), cluster.Spec.GetDisplayName())
	}
	return suggestions
}
