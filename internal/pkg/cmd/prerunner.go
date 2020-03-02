package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/confluentinc/ccloud-sdk-go"
	"github.com/confluentinc/mds-sdk-go"
	"github.com/jonboulle/clockwork"
	"github.com/spf13/cobra"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/confluentinc/cli/internal/pkg/analytics"
	v2 "github.com/confluentinc/cli/internal/pkg/config/v2"
	v3 "github.com/confluentinc/cli/internal/pkg/config/v3"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/cli/internal/pkg/log"
	"github.com/confluentinc/cli/internal/pkg/update"
	"github.com/confluentinc/cli/internal/pkg/version"
)

// PreRun is a helper class for automatically setting up Cobra PersistentPreRun commands
type PreRunner interface {
	Anonymous(command *CLICommand) func(cmd *cobra.Command, args []string) error
	Authenticated(command *AuthenticatedCLICommand) func(cmd *cobra.Command, args []string) error
	AuthenticatedWithMDS(command *AuthenticatedCLICommand) func(cmd *cobra.Command, args []string) error
	HasAPIKey(command *HasAPIKeyCLICommand) func(cmd *cobra.Command, args []string) error
}

// PreRun is the standard PreRunner implementation
type PreRun struct {
	UpdateClient update.Client
	CLIName      string
	Logger       *log.Logger
	Clock        clockwork.Clock
	Analytics    analytics.Client
	FlagResolver FlagResolver
	Version      *version.Version
}

type CLICommand struct {
	*cobra.Command
	Config    *DynamicConfig
	Version   *version.Version
	prerunner PreRunner
}

type AuthenticatedCLICommand struct {
	*CLICommand
	Client    *ccloud.Client
	MDSClient *mds.APIClient
	Context   *DynamicContext
	State     *v2.ContextState
}

type HasAPIKeyCLICommand struct {
	*CLICommand
	Context *DynamicContext
}

func (a *AuthenticatedCLICommand) AuthToken() string {
	return a.State.AuthToken
}
func (a *AuthenticatedCLICommand) EnvironmentId() string {
	return a.State.Auth.Account.Id
}

func NewAuthenticatedCLICommand(command *cobra.Command, cfg *v3.Config, prerunner PreRunner) *AuthenticatedCLICommand {
	cmd := &AuthenticatedCLICommand{
		CLICommand: NewCLICommand(command, cfg, prerunner),
		Context:    nil,
		State:      nil,
	}
	command.PersistentPreRunE = prerunner.Authenticated(cmd)
	cmd.Command = command
	return cmd
}

func NewAuthenticatedWithMDSCLICommand(command *cobra.Command, cfg *v3.Config, prerunner PreRunner) *AuthenticatedCLICommand {
	cmd := &AuthenticatedCLICommand{
		CLICommand: NewCLICommand(command, cfg, prerunner),
		Context:    nil,
		State:      nil,
	}
	command.PersistentPreRunE = prerunner.AuthenticatedWithMDS(cmd)
	cmd.Command = command
	return cmd
}

func NewHasAPIKeyCLICommand(command *cobra.Command, cfg *v3.Config, prerunner PreRunner) *HasAPIKeyCLICommand {
	cmd := &HasAPIKeyCLICommand{
		CLICommand: NewCLICommand(command, cfg, prerunner),
		Context:    nil,
	}
	command.PersistentPreRunE = prerunner.HasAPIKey(cmd)
	cmd.Command = command
	return cmd
}

func NewAnonymousCLICommand(command *cobra.Command, cfg *v3.Config, prerunner PreRunner) *CLICommand {
	cmd := NewCLICommand(command, cfg, prerunner)
	command.PersistentPreRunE = prerunner.Anonymous(cmd)
	cmd.Command = command
	return cmd
}

func NewCLICommand(command *cobra.Command, cfg *v3.Config, prerunner PreRunner) *CLICommand {
	return &CLICommand{
		Config:    NewDynamicConfig(cfg, nil, nil),
		Command:   command,
		prerunner: prerunner,
	}
}

func (a *AuthenticatedCLICommand) AddCommand(command *cobra.Command) {
	command.PersistentPreRunE = a.PersistentPreRunE
	a.Command.AddCommand(command)
}

func (h *HasAPIKeyCLICommand) AddCommand(command *cobra.Command) {
	command.PersistentPreRunE = h.PersistentPreRunE
	h.Command.AddCommand(command)
}

// Anonymous provides PreRun operations for commands that may be run without a logged-in user
func (r *PreRun) Anonymous(command *CLICommand) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		command.Version = r.Version
		command.Config.Resolver = r.FlagResolver
		if err := log.SetLoggingVerbosity(cmd, r.Logger); err != nil {
			return errors.HandleCommon(err, cmd)
		}
		if err := r.notifyIfUpdateAvailable(cmd, r.CLIName, command.Version.Version); err != nil {
			return errors.HandleCommon(err, cmd)
		}
		ctx, err := command.Config.Context(cmd)
		if err != nil {
			return err
		}
		err = r.validateToken(cmd, ctx)
		switch err.(type) {
		case *ccloud.ExpiredTokenError:
			err := ctx.DeleteUserAuth()
			if err != nil {
				return err
			}
			ErrPrintln(cmd, "Your token has expired. You are now logged out.")
			analyticsError := r.Analytics.SessionTimedOut()
			if analyticsError != nil {
				r.Logger.Debug(analyticsError.Error())
			}
		}
		r.Analytics.TrackCommand(cmd, args)
		return nil
	}
}

// Authenticated provides PreRun operations for commands that require a logged-in Confluent Cloud user.
func (r *PreRun) Authenticated(command *AuthenticatedCLICommand) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		err := r.Anonymous(command.CLICommand)(cmd, args)
		if err != nil {
			return errors.HandleCommon(err, cmd)
		}
		err = r.setClients(command)
		if err != nil {
			return errors.HandleCommon(err, cmd)
		}
		ctx, err := command.Config.Context(cmd)
		if err != nil {
			return errors.HandleCommon(err, cmd)
		}
		if ctx == nil {
			return errors.HandleCommon(errors.ErrNoContext, cmd)
		}
		command.Context = ctx
		command.State, err = ctx.AuthenticatedState(cmd)
		if err != nil {
			return errors.HandleCommon(err, cmd)
		}
		return r.validateToken(cmd, ctx)
	}
}

// Authenticated provides PreRun operations for commands that require a logged-in Confluent Cloud user.
func (r *PreRun) AuthenticatedWithMDS(command *AuthenticatedCLICommand) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		err := r.Anonymous(command.CLICommand)(cmd, args)
		if err != nil {
			return errors.HandleCommon(err, cmd)
		}
		err = r.setClients(command)
		if err != nil {
			return errors.HandleCommon(err, cmd)
		}
		ctx, err := command.Config.Context(cmd)
		if err != nil {
			return errors.HandleCommon(err, cmd)
		}
		if ctx == nil {
			return errors.HandleCommon(errors.ErrNoContext, cmd)
		}
		if !ctx.HasMDSLogin() {
			return errors.HandleCommon(errors.ErrNotLoggedIn, cmd)
		}
		command.Context = ctx
		command.State = ctx.State
		return r.validateToken(cmd, ctx)
	}
}

// HasAPIKey provides PreRun operations for commands that require an API key.
func (r *PreRun) HasAPIKey(command *HasAPIKeyCLICommand) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		err := r.Anonymous(command.CLICommand)(cmd, args)
		if err != nil {
			return errors.HandleCommon(err, cmd)
		}
		ctx, err := command.Config.Context(cmd)
		if err != nil {
			return errors.HandleCommon(err, cmd)
		}
		if ctx == nil {
			return errors.HandleCommon(errors.ErrNoContext, cmd)
		}
		command.Context = ctx
		hasAPIKey, err := ctx.HasAPIKey(cmd, ctx.KafkaClusterContext.GetActiveKafkaClusterId())
		if err != nil {
			return errors.HandleCommon(err, cmd)
		}
		if !hasAPIKey {
			err = &errors.UnspecifiedAPIKeyError{ClusterID: ctx.KafkaClusterContext.GetActiveKafkaClusterId()}
			return errors.HandleCommon(err, cmd)
		}
		return nil
	}
}

// notifyIfUpdateAvailable prints a message if an update is available
func (r *PreRun) notifyIfUpdateAvailable(cmd *cobra.Command, name string, currentVersion string) error {
	updateAvailable, latestVersion, err := r.UpdateClient.CheckForUpdates(name, currentVersion, false)
	if err != nil {
		// This is a convenience helper to check-for-updates before arbitrary commands. Since the CLI supports running
		// in internet-less environments (e.g., local or on-prem deploys), swallow the error and log a warning.
		r.Logger.Warn(err)
		return nil
	}
	if updateAvailable {
		msg := "Updates are available for %s from (current: %s, latest: %s). To install them, please run:\n$ %s update\n\n"
		if !strings.HasPrefix(latestVersion, "v") {
			latestVersion = "v" + latestVersion
		}
		ErrPrintf(cmd, msg, name, currentVersion, latestVersion, name)
	}
	return nil
}

func (r *PreRun) setClients(cliCmd *AuthenticatedCLICommand) error {
	ctx, err := cliCmd.Config.Context(cliCmd.Command)
	if err != nil {
		return err
	}
	ccloudClient, err := r.createCCloudClient(ctx, cliCmd.Command, cliCmd.Version)
	if err != nil {
		return err
	}
	cliCmd.Client = ccloudClient
	cliCmd.MDSClient = r.createMDSClient(ctx, cliCmd.Version)
	cliCmd.Config.Client = ccloudClient
	return nil
}

func (r *PreRun) createCCloudClient(ctx *DynamicContext, cmd *cobra.Command, ver *version.Version) (*ccloud.Client, error) {
	var baseURL string
	var authToken string
	var logger *log.Logger
	var userAgent string
	if ctx != nil {
		baseURL = ctx.Platform.Server
		state, err := ctx.AuthenticatedState(cmd)
		if err != nil && err != errors.ErrNotLoggedIn {
			return nil, err
		}
		if err == nil {
			authToken = state.AuthToken
		}
		logger = ctx.Logger
		userAgent = ver.UserAgent
	}
	return ccloud.NewClientWithJWT(context.Background(), authToken, &ccloud.Params{
		BaseURL: baseURL, Logger: logger, UserAgent: userAgent,
	}), nil
}

func (r *PreRun) createMDSClient(ctx *DynamicContext, ver *version.Version) *mds.APIClient {
	mdsConfig := mds.NewConfiguration()
	if ctx == nil {
		return mds.NewAPIClient(mdsConfig)
	}
	mdsConfig.BasePath = ctx.Platform.Server
	mdsConfig.UserAgent = ver.UserAgent
	if ctx.Platform.CaCertPath == "" {
		return mds.NewAPIClient(mdsConfig)
	}
	caCertPath := ctx.Platform.CaCertPath
	// Try to load certs. On failure, warn, but don't error out because this may be an auth command, so there may
	// be a --ca-cert-path flag on the cmd line that'll fix whatever issue there is with the cert file in the config
	caCertFile, err := os.Open(caCertPath)
	if err == nil {
		defer caCertFile.Close()
		mdsConfig.HTTPClient, err = SelfSignedCertClient(caCertFile, r.Logger)
		if err != nil {
			r.Logger.Warnf("Unable to load certificate from %s. %s. Resulting SSL errors will be fixed by logging in with the --ca-cert-path flag.", caCertPath, err.Error())
			mdsConfig.HTTPClient = DefaultClient()
		}
	} else {
		r.Logger.Warnf("Unable to load certificate from %s. %s. Resulting SSL errors will be fixed by logging in with the --ca-cert-path flag.", caCertPath, err.Error())
		mdsConfig.HTTPClient = DefaultClient()

	}
	return mds.NewAPIClient(mdsConfig)
}

func (r *PreRun) validateToken(cmd *cobra.Command, ctx *DynamicContext) error {
	// validate token (not expired)
	var authToken string
	if ctx != nil {
		authToken = ctx.State.AuthToken
	}
	var claims map[string]interface{}
	token, err := jwt.ParseSigned(authToken)
	if err != nil {
		return errors.HandleCommon(new(ccloud.InvalidTokenError), cmd)
	}
	if err := token.UnsafeClaimsWithoutVerification(&claims); err != nil {
		return errors.HandleCommon(err, cmd)
	}
	if exp, ok := claims["exp"].(float64); ok {
		if float64(r.Clock.Now().Unix()) > exp {
			return errors.HandleCommon(new(ccloud.ExpiredTokenError), cmd)
		}
	}
	return nil
}
