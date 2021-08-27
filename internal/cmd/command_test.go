package cmd

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/config"
	v2 "github.com/confluentinc/cli/internal/pkg/config/v2"
	v3 "github.com/confluentinc/cli/internal/pkg/config/v3"
	"github.com/confluentinc/cli/internal/pkg/log"
	"github.com/confluentinc/cli/internal/pkg/utils"
	pversion "github.com/confluentinc/cli/internal/pkg/version"
)

var (
	mockBaseConfig = &config.BaseConfig{Params: &config.Params{Logger: log.New()}}
	mockVersion    = new(pversion.Version)
)

func TestHelp_NoContext(t *testing.T) {
	cfg := &v3.Config{BaseConfig: mockBaseConfig}

	out, err := runWithConfig(cfg)
	require.NoError(t, err)

	commands := []string{
		"completion", "config", "help", "kafka", "local", "login", "logout", "signup", "update", "version",
	}
	if runtime.GOOS == "windows" {
		commands = utils.Remove(commands, "local")
	}

	for _, command := range commands {
		require.Contains(t, out, command)
	}
}

func TestHelp_Cloud(t *testing.T) {
	cfg := &v3.Config{
		BaseConfig:     mockBaseConfig,
		Contexts:       map[string]*v3.Context{"cloud": {PlatformName: "confluent.cloud"}},
		CurrentContext: "cloud",
	}

	out, err := runWithConfig(cfg)
	require.NoError(t, err)

	commands := []string{
		"admin", "api-key", "audit-log", "cloud-signup", "completion", "config", "connect", "environment", "help",
		"iam", "init", "kafka", "ksql", "login", "logout", "price", "prompt", "schema-registry", "service-account",
		"shell", "update", "version",
	}

	for _, command := range commands {
		require.Contains(t, out, command)
	}
}

func TestHelp_CloudWithAPIKey(t *testing.T) {
	cfg := &v3.Config{
		BaseConfig: mockBaseConfig,
		Contexts: map[string]*v3.Context{
			"cloud-with-api-key": {
				PlatformName: "confluent.cloud",
				Credential:   &v2.Credential{CredentialType: v2.APIKey},
			},
		},
		CurrentContext: "cloud-with-api-key",
	}

	out, err := runWithConfig(cfg)
	require.NoError(t, err)

	commands := []string{
		"admin", "audit-log", "completion", "config", "help", "init", "kafka", "login", "logout", "signup", "update",
		"version",
	}

	for _, command := range commands {
		require.Contains(t, out, command)
	}
}

func TestHelp_OnPrem(t *testing.T) {
	cfg := &v3.Config{
		BaseConfig:     mockBaseConfig,
		Contexts:       map[string]*v3.Context{"on-prem": {PlatformName: "https://example.com"}},
		CurrentContext: "on-prem",
	}

	out, err := runWithConfig(cfg)
	require.NoError(t, err)

	commands := []string{
		"audit-log", "cluster", "completion", "config", "connect", "help", "iam", "kafka", "ksql", "local", "login",
		"logout", "schema-registry", "secret", "signup", "update", "version",
	}
	if runtime.GOOS == "windows" {
		commands = utils.Remove(commands, "local")
	}

	for _, command := range commands {
		require.Contains(t, out, command)
	}
}

func runWithConfig(cfg *v3.Config) (string, error) {
	cli := NewConfluentCommand(cfg, true, mockVersion)
	return pcmd.ExecuteCommand(cli.Command, "help")
}
