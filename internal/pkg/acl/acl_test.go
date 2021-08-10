package acl

import (
	"fmt"
	errMsgs "github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/kafka-rest-sdk-go/kafkarestv3"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseAclRequest(t *testing.T) {
	var suite = []struct {
		args        []string
		expectedAcl AclRequestDataWithError
	}{
		{
			args: []string{"--operation", "READ", "--principal", "User:Alice", "--cluster-scope", "--host", "127.0.0.1", "--allow"},
			expectedAcl: AclRequestDataWithError{
				ResourceType: kafkarestv3.ACLRESOURCETYPE_CLUSTER,
				ResourceName: "kafka-cluster",
				PatternType:  kafkarestv3.ACLPATTERNTYPE_LITERAL,
				Principal:    "User:Alice",
				Host:         "127.0.0.1",
				Operation:    kafkarestv3.ACLOPERATION_READ,
				Permission:   kafkarestv3.ACLPERMISSION_ALLOW,
				Errors:       nil,
			},
		},
		{
			args: []string{"--operation", "delete", "--principal", "Group:Admin", "--topic", "Test", "--prefix", "--deny"},
			expectedAcl: AclRequestDataWithError{
				ResourceType: kafkarestv3.ACLRESOURCETYPE_TOPIC,
				ResourceName: "Test",
				PatternType:  kafkarestv3.ACLPATTERNTYPE_PREFIXED,
				Principal:    "Group:Admin",
				Host:         "*",
				Operation:    kafkarestv3.ACLOPERATION_DELETE,
				Permission:   kafkarestv3.ACLPERMISSION_DENY,
				Errors:       nil,
			},
		},
		{
			args: []string{"--operation", "fake", "--principal", "User:Alice", "--cluster-scope", "--transactional-id", "123"},
			expectedAcl: AclRequestDataWithError{
				Errors:       multierror.Append(errors.New("Invalid operation value: FAKE"), fmt.Errorf("exactly one of %v must be set",
					convertToFlags(kafkarestv3.ACLRESOURCETYPE_TOPIC, kafkarestv3.ACLRESOURCETYPE_GROUP,
						kafkarestv3.ACLRESOURCETYPE_CLUSTER, kafkarestv3.ACLRESOURCETYPE_TRANSACTIONAL_ID))),
			},
		},
		{
			args: []string{"--operation", "READ", "--principal", "User:Alice", "--transactional-id", "123", "--allow", "--deny"},
			expectedAcl: AclRequestDataWithError{
				Errors:       multierror.Append(errors.Errorf(errMsgs.OnlySetAllowOrDenyErrorMsg)),
			},
		},
	}
	req := require.New(t)
	for _, s := range suite {
		cmd := &cobra.Command{}
		cmd.Flags().AddFlagSet(AclFlags())
		_ = cmd.ParseFlags(s.args)
		acl := ParseAclRequest(cmd)
		if s.expectedAcl.Errors != nil {
			req.NotNil(acl.Errors)
			req.Equal(s.expectedAcl.Errors.Error(), acl.Errors.Error())
		} else {
			req.Nil(acl.Errors)
			req.Equal(s.expectedAcl, *acl)
		}
	}
}

func TestValidateCreateDeleteAclRequestData(t *testing.T) {
	var suite = []struct {
		initialAcl  AclRequestDataWithError
		expectedAcl AclRequestDataWithError
	}{
		{
			initialAcl: AclRequestDataWithError{
				ResourceType: kafkarestv3.ACLRESOURCETYPE_CLUSTER,
				Permission:   kafkarestv3.ACLPERMISSION_ALLOW,
			},
			expectedAcl: AclRequestDataWithError{
				PatternType: kafkarestv3.ACLPATTERNTYPE_LITERAL,
				ResourceType: kafkarestv3.ACLRESOURCETYPE_CLUSTER,
				Permission:   kafkarestv3.ACLPERMISSION_ALLOW,
			},
		},
		{
			initialAcl: AclRequestDataWithError{Host: "*"},
			expectedAcl: AclRequestDataWithError{Errors: multierror.Append(errors.Errorf(errMsgs.MustSetAllowOrDenyErrorMsg), errors.Errorf(errMsgs.MustSetResourceTypeErrorMsg,
				convertToFlags(kafkarestv3.ACLRESOURCETYPE_TOPIC, kafkarestv3.ACLRESOURCETYPE_GROUP,
					kafkarestv3.ACLRESOURCETYPE_CLUSTER, kafkarestv3.ACLRESOURCETYPE_TRANSACTIONAL_ID)))},
		},
	}
	req := require.New(t)
	for _, s := range suite {
		validatedAcl := ValidateCreateDeleteAclRequestData(&s.initialAcl)
		if s.expectedAcl.Errors != nil {
			req.NotNil(validatedAcl.Errors)
			req.Equal(s.expectedAcl.Errors.Error(), validatedAcl.Errors.Error())
		} else {
			req.Nil(validatedAcl.Errors)
			req.Equal(s.expectedAcl, *validatedAcl)
		}
	}
}