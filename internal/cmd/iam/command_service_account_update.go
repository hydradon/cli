package iam

import (
	"context"

	orgv1 "github.com/confluentinc/cc-structs/kafka/org/v1"
	"github.com/spf13/cobra"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/cli/internal/pkg/examples"
	"github.com/confluentinc/cli/internal/pkg/resource"
	"github.com/confluentinc/cli/internal/pkg/utils"
)

func (c *serviceAccountCommand) newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update <id>",
		Short:             "Update a service account.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: pcmd.NewValidArgsFunction(c.validArgs),
		RunE:              pcmd.NewCLIRunE(c.update),
		Example: examples.BuildExampleString(
			examples.Example{
				Text: `Update the description of service account "sa-123456".`,
				Code: `confluent iam service-account update sa-123456 --description "Update demo service account information."`,
			},
		),
	}

	cmd.Flags().String("description", "", "Description of the service account.")
	_ = cmd.MarkFlagRequired("description")

	return cmd
}

func (c *serviceAccountCommand) update(cmd *cobra.Command, args []string) error {
	description, err := cmd.Flags().GetString("description")
	if err != nil {
		return err
	}

	if err := requireLen(description, descriptionLength, "description"); err != nil {
		return err
	}

	if resource.LookupType(args[0]) != resource.ServiceAccount {
		return errors.New(errors.BadServiceAccountIDErrorMsg)
	}
	serviceAccountId := args[0]

	user := &orgv1.User{
		ResourceId:         serviceAccountId,
		ServiceDescription: description,
	}

	if err := c.Client.User.UpdateServiceAccount(context.Background(), user); err != nil {
		return err
	}

	utils.ErrPrintf(cmd, errors.UpdateSuccessMsg, "description", "service account", serviceAccountId, description)
	return nil
}
