package pipeline

import (
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/spf13/cobra"

	streamdesignerv1 "github.com/confluentinc/ccloud-sdk-go-v2/stream-designer/v1"

	pcmd "github.com/confluentinc/cli/v3/pkg/cmd"
	"github.com/confluentinc/cli/v3/pkg/errors"
	"github.com/confluentinc/cli/v3/pkg/examples"
)

func (c *command) newCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new pipeline.",
		Args:  cobra.NoArgs,
		RunE:  c.create,
		Example: examples.BuildExampleString(
			examples.Example{
				Text: `Create a Stream Designer pipeline with name "my-pipeline".`,
				Code: `confluent pipeline create --name my-pipeline --description "example pipeline"`,
			},
			examples.Example{
				Text: "Create a Stream Designer pipeline with a KSQL cluster and a Schema Registry cluster.",
				Code: `confluent pipeline create --name my-pipeline --description "example pipeline" --ksql-cluster lksqlc-12345 --use-schema-registry`,
			},
		),
	}

	cmd.Flags().String("name", "", "Name of the pipeline.")
	cmd.Flags().String("description", "", "Description of the pipeline.")
	pcmd.AddKsqlClusterFlag(cmd, c.AuthenticatedCLICommand)
	cmd.Flags().Bool("use-schema-registry", false, "Add the Schema Registry cluster to the pipeline.")
	cmd.Flags().String("sql-file", "", "Path to a KSQL file containing the pipeline's source code.")
	cmd.Flags().StringArray("secret", []string{}, "A named secret that can be referenced in pipeline source code, for example, \"secret_name=secret_content\".\n"+
		"This flag can be supplied multiple times. The secret mapping must have the format <secret-name>=<secret-value>,\n"+
		"where <secret-name> consists of 1-128 lowercase, uppercase, numeric or underscore characters but may not begin with a digit.\n"+
		"The <secret-value> can be of any format but may not be empty.")
	pcmd.AddOutputFlag(cmd)
	pcmd.AddClusterFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddEnvironmentFlag(cmd, c.AuthenticatedCLICommand)

	cobra.CheckErr(cmd.MarkFlagFilename("sql-file", "sql"))
	cobra.CheckErr(cmd.MarkFlagRequired("name"))

	return cmd
}

func (c *command) create(cmd *cobra.Command, _ []string) error {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	description, err := cmd.Flags().GetString("description")
	if err != nil {
		return err
	}

	ksqlCluster, err := cmd.Flags().GetString("ksql-cluster")
	if err != nil {
		return err
	}

	sqlFile, err := cmd.Flags().GetString("sql-file")
	if err != nil {
		return err
	}

	secrets, err := cmd.Flags().GetStringArray("secret")
	if err != nil {
		return err
	}

	kafkaCluster, err := c.Context.GetKafkaClusterForCommand(c.V2Client)
	if err != nil {
		return err
	}

	environmentId, err := c.Context.EnvironmentId()
	if err != nil {
		return err
	}

	// validate ksql id
	if ksqlCluster != "" {
		if _, err := c.V2Client.DescribeKsqlCluster(ksqlCluster, environmentId); err != nil {
			return err
		}
	}

	// read pipeline source code file if provided
	sourceCode := ""
	if sqlFile != "" {
		fileContent, err := os.ReadFile(sqlFile)
		if err != nil {
			return err
		}
		sourceCode = string(fileContent)
	}

	// parse and construct secret mappings
	secretMappings, err := createSecretMappings(secrets, secretMappingWithoutEmptyValue)
	if err != nil {
		return err
	}

	// required fields
	createPipeline := streamdesignerv1.SdV1Pipeline{
		Spec: &streamdesignerv1.SdV1PipelineSpec{
			DisplayName:  streamdesignerv1.PtrString(name),
			Description:  streamdesignerv1.PtrString(description),
			SourceCode:   &streamdesignerv1.SdV1SourceCodeObject{Sql: sourceCode},
			Secrets:      &secretMappings,
			Environment:  &streamdesignerv1.ObjectReference{Id: environmentId},
			KafkaCluster: &streamdesignerv1.ObjectReference{Id: kafkaCluster.ID},
		},
	}

	// add KSQL cluster if present
	if ksqlCluster != "" {
		createPipeline.Spec.KsqlCluster = &streamdesignerv1.ObjectReference{Id: ksqlCluster}
	}

	// check if we need to configure Schema Registry cluster for this pipeline
	useSchemaRegistry, err := cmd.Flags().GetBool("use-schema-registry")
	if err != nil {
		return err
	}

	if useSchemaRegistry {
		clusters, err := c.V2Client.GetSchemaRegistryClustersByEnvironment(environmentId)
		if err != nil {
			return err
		}
		if len(clusters) == 0 {
			return errors.NewSRNotEnabledError()
		}
		createPipeline.Spec.StreamGovernanceCluster = &streamdesignerv1.ObjectReference{Id: clusters[0].GetId()}
	}

	pipeline, err := c.V2Client.CreatePipeline(createPipeline)
	if err != nil {
		return err
	}

	return printTable(cmd, pipeline)
}

func createSecretMappings(secrets []string, regex string) (map[string]string, error) {
	secretMappings := make(map[string]string)

	// The name of a secret may consist of lowercase letters, uppercase letters, digits,
	// and the '_' (underscore) and may not begin with a digit.
	pattern := regexp.MustCompile(regex)

	for _, secret := range secrets {
		if !pattern.MatchString(secret) {
			return nil, fmt.Errorf(`invalid secret pattern "%s"`, secret)
		}

		matches := pattern.FindStringSubmatch(secret)
		name, value := matches[1], matches[2]
		secretMappings[name] = value
	}
	return secretMappings, nil
}

func getOrderedSecretNames(secrets *map[string]string) []string {
	if secrets == nil {
		return []string{}
	}

	names := make([]string, 0, len(*secrets))
	for n := range *secrets {
		names = append(names, n)
	}

	sort.Strings(names)
	return names
}
