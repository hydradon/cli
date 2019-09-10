package iam

import (
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/config"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/go-printer"

	"context"

	mds "github.com/confluentinc/mds-sdk-go"
)

var (
	resourcePatternListFields = []string{"ResourceType", "Name", "PatternType"}
	resourcePatternListLabels = []string{"Role", "ResourceType", "Name", "PatternType"}
)

type rolebindingOptions struct {
	role             string
	resource         string
	prefix           bool
	principal        string
	scopeClusters    mds.ScopeClusters
	resourcesRequest mds.ResourcesRequest
}

type rolebindingCommand struct {
	*cobra.Command
	config *config.Config
	ch     *pcmd.ConfigHelper
	client *mds.APIClient
	ctx    context.Context
}

// NewRolebindingCommand returns the sub-command object for interacting with RBAC rolebindings.
func NewRolebindingCommand(config *config.Config, ch *pcmd.ConfigHelper, client *mds.APIClient) *cobra.Command {
	cmd := &rolebindingCommand{
		Command: &cobra.Command{
			Use:   "rolebinding",
			Short: "Manage RBAC and IAM role bindings.",
			Long:  "Manage Role Based Access (RBAC) and Identity and Access Management (IAM) role bindings.",
		},
		config: config,
		ch:     ch,
		client: client,
		ctx:    context.WithValue(context.Background(), mds.ContextAccessToken, config.AuthToken),
	}

	cmd.init()
	return cmd.Command
}

func (c *rolebindingCommand) init() {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List role bindings.",
		Long:  "List the role bindings for a particular principal and scope.",
		RunE:  c.list,
		Args:  cobra.NoArgs,
	}
	listCmd.Flags().String("principal", "", "Principal whose rolebindings should be listed.")
	listCmd.Flags().String("role", "", "List rolebindings under a specific role given to a principal.")
	listCmd.Flags().String("kafka-cluster-id", "", "Kafka cluster ID for scope of rolebinding listings.")
	listCmd.Flags().String("schema-registry-cluster-id", "", "Schema Registry cluster ID for scope of rolebinding listings.")
	listCmd.Flags().String("ksql-cluster-id", "", "KSQL cluster ID for scope of rolebinding listings.")
	listCmd.Flags().String("connect-cluster-id", "", "Kafka Connect cluster ID for scope of rolebinding listings.")
	listCmd.Flags().SortFlags = false
	check(listCmd.MarkFlagRequired("principal"))
	c.AddCommand(listCmd)

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a role binding.",
		RunE:  c.create,
		Args:  cobra.NoArgs,
	}
	createCmd.Flags().String("role", "", "Role name of the new role binding.")
	createCmd.Flags().String("resource", "", "Qualified resource name for the role binding.")
	createCmd.Flags().Bool("prefix", false, "Whether the provided resource name is treated as a prefix pattern.")
	createCmd.Flags().String("principal", "", "Qualified principal name for the role binding.")
	createCmd.Flags().String("kafka-cluster-id", "", "Kafka cluster ID for the role binding.")
	createCmd.Flags().String("schema-registry-cluster-id", "", "Schema Registry cluster ID for the role binding.")
	createCmd.Flags().String("ksql-cluster-id", "", "KSQL cluster ID for the role binding.")
	createCmd.Flags().String("connect-cluster-id", "", "Kafka Connect cluster ID for the role binding.")
	createCmd.Flags().SortFlags = false
	check(createCmd.MarkFlagRequired("role"))
	check(createCmd.MarkFlagRequired("principal"))
	c.AddCommand(createCmd)

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an existing role binding.",
		RunE:  c.delete,
		Args:  cobra.NoArgs,
	}
	deleteCmd.Flags().String("role", "", "Role name of the existing role binding.")
	deleteCmd.Flags().String("resource", "", "Qualified resource name associated with the role binding.")
	deleteCmd.Flags().Bool("prefix", false, "Whether the provided resource name is treated as a prefix pattern.")
	deleteCmd.Flags().String("principal", "", "Qualified principal name associated with the role binding.")
	deleteCmd.Flags().String("kafka-cluster-id", "", "Kafka cluster ID for the role binding.")
	deleteCmd.Flags().String("schema-registry-cluster-id", "", "Schema Registry cluster ID for the role binding.")
	deleteCmd.Flags().String("ksql-cluster-id", "", "KSQL cluster ID for the role binding.")
	deleteCmd.Flags().String("connect-cluster-id", "", "Kafka Connect cluster ID for the role binding.")
	deleteCmd.Flags().SortFlags = false
	check(createCmd.MarkFlagRequired("role"))
	check(deleteCmd.MarkFlagRequired("principal"))
	c.AddCommand(deleteCmd)
}

func (c *rolebindingCommand) validatePrincipalFormat(principal string) error {
	if len(strings.Split(principal, ":")) == 1 {
		return errors.New("Principal must be specified in this format: <Principal Type>:<Principal Name>")
	}

	return nil
}

func (c *rolebindingCommand) parseAndValidateResourcePattern(typename string, prefix bool) (mds.ResourcePattern, error) {
	var result mds.ResourcePattern
	if prefix {
		result.PatternType = "PREFIXED"
	} else {
		result.PatternType = "LITERAL"
	}

	parts := strings.Split(typename, ":")
	if len(parts) != 2 {
		return result, errors.New("Resource must be specified in this format: <Resource Type>:<Resource Name>")
	}
	result.ResourceType = parts[0]
	result.Name = parts[1]

	return result, nil
}

func (c *rolebindingCommand) validateRoleAndResourceType(roleName string, resourceType string) error {
	role, _, err := c.client.RoleDefinitionsApi.RoleDetail(c.ctx, roleName)
	if err != nil {
		return errors.Wrapf(err, "Failed to look up role %s. Was an invalid role name specified?", roleName)
	}

	allResourceTypes := []string{}
	found := false
	for _, operation := range role.AccessPolicy.AllowedOperations {
		allResourceTypes = append(allResourceTypes, operation.ResourceType)
		if operation.ResourceType == resourceType {
			found = true
			break
		}
	}

	if !found {
		return errors.New("Invalid resource type " + resourceType + " specified. It must be one of " + strings.Join(allResourceTypes, ", "))
	}

	return nil
}

func (c *rolebindingCommand) parseAndValidateScope(cmd *cobra.Command) (*mds.ScopeClusters, error) {
	scope := &mds.ScopeClusters{}

	nonKafkaScopesSet := 0

	cmd.Flags().Visit(func(flag *pflag.Flag) {
		switch flag.Name {
		case "kafka-cluster-id":
			scope.KafkaCluster = flag.Value.String()
		case "schema-registry-cluster-id":
			scope.SchemaRegistryCluster = flag.Value.String()
			nonKafkaScopesSet++
		case "ksql-cluster-id":
			scope.KsqlCluster = flag.Value.String()
			nonKafkaScopesSet++
		case "connect-cluster-id":
			scope.ConnectCluster = flag.Value.String()
			nonKafkaScopesSet++
		}
	})

	if scope.KafkaCluster == "" && nonKafkaScopesSet > 0 {
		return nil, errors.HandleCommon(errors.New("Must also specify a --kafka-cluster-id to uniquely identify the scope."), cmd)
	}

	if scope.KafkaCluster == "" && nonKafkaScopesSet == 0 {
		return nil, errors.HandleCommon(errors.New("Must specify at least one cluster ID flag to indicate role binding scope."), cmd)
	}

	if nonKafkaScopesSet > 1 {
		return nil, errors.HandleCommon(errors.New("Cannot specify more than one non-Kafka cluster ID for a scope."), cmd)
	}

	return scope, nil
}

func (c *rolebindingCommand) list(cmd *cobra.Command, args []string) error {
	role := "*"
	if cmd.Flags().Changed("role") {
		r, err := cmd.Flags().GetString("role")
		if err != nil {
			return errors.HandleCommon(err, cmd)
		}
		role = r
	}

	principal, err := cmd.Flags().GetString("principal")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	err = c.validatePrincipalFormat(principal)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	scopeClusters, err := c.parseAndValidateScope(cmd)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	var roleNamesWithMultiplicity []string
	var resourcePatterns []mds.ResourcePattern
	roleNames := []string{role}
	if role == "*" {
		roleNames, _, err = c.client.UserAndRoleMgmtApi.ScopedPrincipalRolenames(c.ctx, principal, mds.Scope{Clusters: *scopeClusters})
		if err != nil {
			return errors.HandleCommon(err, cmd)
		}
	}

	for _, r := range roleNames {
		// This only gets resource-scoped bindings...
		rps, _, err := c.client.UserAndRoleMgmtApi.GetRoleResourcesForPrincipal(c.ctx, principal, r, mds.Scope{Clusters: *scopeClusters})
		if err != nil {
			return errors.HandleCommon(err, cmd)
		}
		// ...so manually append cluster-scoped bindings when needed
		if len(rps) == 0 && isClusterScopedRole(r) {
			rps = append(rps, mds.ResourcePattern{
				ResourceType: "Cluster",
				Name:         "",
				PatternType:  "",
			})
		}
		resourcePatterns = append(resourcePatterns, rps...)
		for range rps {
			roleNamesWithMultiplicity = append(roleNamesWithMultiplicity, r)
		}
	}

	var data [][]string
	for i, pattern := range resourcePatterns {
		data = append(data, append([]string{roleNamesWithMultiplicity[i]}, printer.ToRow(&pattern, resourcePatternListFields)...))
	}
	printer.RenderCollectionTable(data, resourcePatternListLabels)

	return nil
}

func (c *rolebindingCommand) parseCommon(cmd *cobra.Command) (*rolebindingOptions, error) {
	role, err := cmd.Flags().GetString("role")
	if err != nil {
		return nil, errors.HandleCommon(err, cmd)
	}

	resource, err := cmd.Flags().GetString("resource")
	if err != nil {
		return nil, errors.HandleCommon(err, cmd)
	}

	prefix := cmd.Flags().Changed("prefix")

	principal, err := cmd.Flags().GetString("principal")
	if err != nil {
		return nil, errors.HandleCommon(err, cmd)
	}
	err = c.validatePrincipalFormat(principal)
	if err != nil {
		return nil, errors.HandleCommon(err, cmd)
	}

	scopeClusters, err := c.parseAndValidateScope(cmd)
	if err != nil {
		return nil, errors.HandleCommon(err, cmd)
	}

	resourcesRequest := mds.ResourcesRequest{}
	if resource != "" {
		parsedResourcePattern, err := c.parseAndValidateResourcePattern(resource, prefix)
		if err != nil {
			return nil, errors.HandleCommon(err, cmd)
		}
		err = c.validateRoleAndResourceType(role, parsedResourcePattern.ResourceType)
		if err != nil {
			return nil, errors.HandleCommon(err, cmd)
		}
		resourcePatterns := []mds.ResourcePattern{
			parsedResourcePattern,
		}
		resourcesRequest = mds.ResourcesRequest{
			Scope:            mds.Scope{Clusters: *scopeClusters},
			ResourcePatterns: resourcePatterns,
		}
	}

	return &rolebindingOptions{
			role,
			resource,
			prefix,
			principal,
			*scopeClusters,
			resourcesRequest,
		},
		nil
}

func (c *rolebindingCommand) create(cmd *cobra.Command, args []string) error {
	options, err := c.parseCommon(cmd)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	var resp *http.Response
	if options.resource != "" {
		resp, err = c.client.UserAndRoleMgmtApi.AddRoleResourcesForPrincipal(c.ctx, options.principal, options.role, options.resourcesRequest)
	} else {
		resp, err = c.client.UserAndRoleMgmtApi.AddRoleForPrincipal(c.ctx, options.principal, options.role, mds.Scope{Clusters: options.scopeClusters})
	}

	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return errors.HandleCommon(errors.Wrapf(err, "No error, but received HTTP status code %d.  Please file a support ticket with details", resp.StatusCode), cmd)
	}

	return nil
}

func (c *rolebindingCommand) delete(cmd *cobra.Command, args []string) error {
	options, err := c.parseCommon(cmd)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	var resp *http.Response
	if options.resource != "" {
		resp, err = c.client.UserAndRoleMgmtApi.RemoveRoleResourcesForPrincipal(c.ctx, options.principal, options.role, options.resourcesRequest)
	} else {
		resp, err = c.client.UserAndRoleMgmtApi.DeleteRoleForPrincipal(c.ctx, options.principal, options.role, mds.Scope{Clusters: options.scopeClusters})
	}

	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return errors.HandleCommon(errors.Wrapf(err, "No error, but received HTTP status code %d.  Please file a support ticket with details", resp.StatusCode), cmd)
	}

	return nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// TODO please move this to a backend route
func isClusterScopedRole(role string) bool {
	clusterScopedRoles := []string{
		"SystemAdmin",
		"ClusterAdmin",
		"SecurityAdmin",
		"UserAdmin",
		"Operator",
	}
	for _, r := range clusterScopedRoles {
		if r == role {
			return true
		}
	}
	return false
}
