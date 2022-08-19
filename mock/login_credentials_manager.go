// Code generated by mocker. DO NOT EDIT.
// github.com/travisjeffery/mocker
// Source: login_credentials_manager.go

package mock

import (
	sync "sync"

	github_com_confluentinc_ccloud_sdk_go_v1 "github.com/confluentinc/ccloud-sdk-go-v1"
	github_com_confluentinc_cli_internal_pkg_auth "github.com/confluentinc/cli/internal/pkg/auth"
	github_com_confluentinc_cli_internal_pkg_config_v1 "github.com/confluentinc/cli/internal/pkg/config/v1"
	github_com_confluentinc_cli_internal_pkg_netrc "github.com/confluentinc/cli/internal/pkg/netrc"
	github_com_spf13_cobra "github.com/spf13/cobra"
)

// MockLoginCredentialsManager is a mock of LoginCredentialsManager interface
type MockLoginCredentialsManager struct {
	lockGetCloudCredentialsFromEnvVar sync.Mutex
	GetCloudCredentialsFromEnvVarFunc func(orgResourceId string) func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error)

	lockGetOnPremCredentialsFromEnvVar sync.Mutex
	GetOnPremCredentialsFromEnvVarFunc func() func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error)

	lockGetCredentialsFromConfig sync.Mutex
	GetCredentialsFromConfigFunc func(cfg *github_com_confluentinc_cli_internal_pkg_config_v1.Config) func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error)

	lockGetCredentialsFromNetrc sync.Mutex
	GetCredentialsFromNetrcFunc func(filterParams github_com_confluentinc_cli_internal_pkg_netrc.NetrcMachineParams) func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error)

	lockGetCloudCredentialsFromPrompt sync.Mutex
	GetCloudCredentialsFromPromptFunc func(cmd *github_com_spf13_cobra.Command, orgResourceId string) func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error)

	lockGetOnPremCredentialsFromPrompt sync.Mutex
	GetOnPremCredentialsFromPromptFunc func(cmd *github_com_spf13_cobra.Command) func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error)

	lockGetPrerunCredentialsFromConfig sync.Mutex
	GetPrerunCredentialsFromConfigFunc func(cfg *github_com_confluentinc_cli_internal_pkg_config_v1.Config) func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error)

	lockGetOnPremPrerunCredentialsFromEnvVar sync.Mutex
	GetOnPremPrerunCredentialsFromEnvVarFunc func() func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error)

	lockGetOnPremPrerunCredentialsFromNetrc sync.Mutex
	GetOnPremPrerunCredentialsFromNetrcFunc func(arg0 *github_com_spf13_cobra.Command, arg1 github_com_confluentinc_cli_internal_pkg_netrc.NetrcMachineParams) func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error)

	lockSetCloudClient sync.Mutex
	SetCloudClientFunc func(client *github_com_confluentinc_ccloud_sdk_go_v1.Client)

	calls struct {
		GetCloudCredentialsFromEnvVar []struct {
			OrgResourceId string
		}
		GetOnPremCredentialsFromEnvVar []struct {
		}
		GetCredentialsFromConfig []struct {
			Cfg *github_com_confluentinc_cli_internal_pkg_config_v1.Config
		}
		GetCredentialsFromNetrc []struct {
			FilterParams github_com_confluentinc_cli_internal_pkg_netrc.NetrcMachineParams
		}
		GetCloudCredentialsFromPrompt []struct {
			Cmd           *github_com_spf13_cobra.Command
			OrgResourceId string
		}
		GetOnPremCredentialsFromPrompt []struct {
			Cmd *github_com_spf13_cobra.Command
		}
		GetPrerunCredentialsFromConfig []struct {
			Cfg *github_com_confluentinc_cli_internal_pkg_config_v1.Config
		}
		GetOnPremPrerunCredentialsFromEnvVar []struct {
		}
		GetOnPremPrerunCredentialsFromNetrc []struct {
			Arg0 *github_com_spf13_cobra.Command
			Arg1 github_com_confluentinc_cli_internal_pkg_netrc.NetrcMachineParams
		}
		SetCloudClient []struct {
			Client *github_com_confluentinc_ccloud_sdk_go_v1.Client
		}
	}
}

// GetCloudCredentialsFromEnvVar mocks base method by wrapping the associated func.
func (m *MockLoginCredentialsManager) GetCloudCredentialsFromEnvVar(orgResourceId string) func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error) {
	m.lockGetCloudCredentialsFromEnvVar.Lock()
	defer m.lockGetCloudCredentialsFromEnvVar.Unlock()

	if m.GetCloudCredentialsFromEnvVarFunc == nil {
		panic("mocker: MockLoginCredentialsManager.GetCloudCredentialsFromEnvVarFunc is nil but MockLoginCredentialsManager.GetCloudCredentialsFromEnvVar was called.")
	}

	call := struct {
		OrgResourceId string
	}{
		OrgResourceId: orgResourceId,
	}

	m.calls.GetCloudCredentialsFromEnvVar = append(m.calls.GetCloudCredentialsFromEnvVar, call)

	return m.GetCloudCredentialsFromEnvVarFunc(orgResourceId)
}

// GetCloudCredentialsFromEnvVarCalled returns true if GetCloudCredentialsFromEnvVar was called at least once.
func (m *MockLoginCredentialsManager) GetCloudCredentialsFromEnvVarCalled() bool {
	m.lockGetCloudCredentialsFromEnvVar.Lock()
	defer m.lockGetCloudCredentialsFromEnvVar.Unlock()

	return len(m.calls.GetCloudCredentialsFromEnvVar) > 0
}

// GetCloudCredentialsFromEnvVarCalls returns the calls made to GetCloudCredentialsFromEnvVar.
func (m *MockLoginCredentialsManager) GetCloudCredentialsFromEnvVarCalls() []struct {
	OrgResourceId string
} {
	m.lockGetCloudCredentialsFromEnvVar.Lock()
	defer m.lockGetCloudCredentialsFromEnvVar.Unlock()

	return m.calls.GetCloudCredentialsFromEnvVar
}

// GetOnPremCredentialsFromEnvVar mocks base method by wrapping the associated func.
func (m *MockLoginCredentialsManager) GetOnPremCredentialsFromEnvVar() func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error) {
	m.lockGetOnPremCredentialsFromEnvVar.Lock()
	defer m.lockGetOnPremCredentialsFromEnvVar.Unlock()

	if m.GetOnPremCredentialsFromEnvVarFunc == nil {
		panic("mocker: MockLoginCredentialsManager.GetOnPremCredentialsFromEnvVarFunc is nil but MockLoginCredentialsManager.GetOnPremCredentialsFromEnvVar was called.")
	}

	call := struct {
	}{}

	m.calls.GetOnPremCredentialsFromEnvVar = append(m.calls.GetOnPremCredentialsFromEnvVar, call)

	return m.GetOnPremCredentialsFromEnvVarFunc()
}

// GetOnPremCredentialsFromEnvVarCalled returns true if GetOnPremCredentialsFromEnvVar was called at least once.
func (m *MockLoginCredentialsManager) GetOnPremCredentialsFromEnvVarCalled() bool {
	m.lockGetOnPremCredentialsFromEnvVar.Lock()
	defer m.lockGetOnPremCredentialsFromEnvVar.Unlock()

	return len(m.calls.GetOnPremCredentialsFromEnvVar) > 0
}

// GetOnPremCredentialsFromEnvVarCalls returns the calls made to GetOnPremCredentialsFromEnvVar.
func (m *MockLoginCredentialsManager) GetOnPremCredentialsFromEnvVarCalls() []struct {
} {
	m.lockGetOnPremCredentialsFromEnvVar.Lock()
	defer m.lockGetOnPremCredentialsFromEnvVar.Unlock()

	return m.calls.GetOnPremCredentialsFromEnvVar
}

// GetCredentialsFromConfig mocks base method by wrapping the associated func.
func (m *MockLoginCredentialsManager) GetCredentialsFromConfig(cfg *github_com_confluentinc_cli_internal_pkg_config_v1.Config) func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error) {
	m.lockGetCredentialsFromConfig.Lock()
	defer m.lockGetCredentialsFromConfig.Unlock()

	if m.GetCredentialsFromConfigFunc == nil {
		panic("mocker: MockLoginCredentialsManager.GetCredentialsFromConfigFunc is nil but MockLoginCredentialsManager.GetCredentialsFromConfig was called.")
	}

	call := struct {
		Cfg *github_com_confluentinc_cli_internal_pkg_config_v1.Config
	}{
		Cfg: cfg,
	}

	m.calls.GetCredentialsFromConfig = append(m.calls.GetCredentialsFromConfig, call)

	return m.GetCredentialsFromConfigFunc(cfg)
}

// GetCredentialsFromConfigCalled returns true if GetCredentialsFromConfig was called at least once.
func (m *MockLoginCredentialsManager) GetCredentialsFromConfigCalled() bool {
	m.lockGetCredentialsFromConfig.Lock()
	defer m.lockGetCredentialsFromConfig.Unlock()

	return len(m.calls.GetCredentialsFromConfig) > 0
}

// GetCredentialsFromConfigCalls returns the calls made to GetCredentialsFromConfig.
func (m *MockLoginCredentialsManager) GetCredentialsFromConfigCalls() []struct {
	Cfg *github_com_confluentinc_cli_internal_pkg_config_v1.Config
} {
	m.lockGetCredentialsFromConfig.Lock()
	defer m.lockGetCredentialsFromConfig.Unlock()

	return m.calls.GetCredentialsFromConfig
}

// GetCredentialsFromNetrc mocks base method by wrapping the associated func.
func (m *MockLoginCredentialsManager) GetCredentialsFromNetrc(filterParams github_com_confluentinc_cli_internal_pkg_netrc.NetrcMachineParams) func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error) {
	m.lockGetCredentialsFromNetrc.Lock()
	defer m.lockGetCredentialsFromNetrc.Unlock()

	if m.GetCredentialsFromNetrcFunc == nil {
		panic("mocker: MockLoginCredentialsManager.GetCredentialsFromNetrcFunc is nil but MockLoginCredentialsManager.GetCredentialsFromNetrc was called.")
	}

	call := struct {
		FilterParams github_com_confluentinc_cli_internal_pkg_netrc.NetrcMachineParams
	}{
		FilterParams: filterParams,
	}

	m.calls.GetCredentialsFromNetrc = append(m.calls.GetCredentialsFromNetrc, call)

	return m.GetCredentialsFromNetrcFunc(filterParams)
}

// GetCredentialsFromNetrcCalled returns true if GetCredentialsFromNetrc was called at least once.
func (m *MockLoginCredentialsManager) GetCredentialsFromNetrcCalled() bool {
	m.lockGetCredentialsFromNetrc.Lock()
	defer m.lockGetCredentialsFromNetrc.Unlock()

	return len(m.calls.GetCredentialsFromNetrc) > 0
}

// GetCredentialsFromNetrcCalls returns the calls made to GetCredentialsFromNetrc.
func (m *MockLoginCredentialsManager) GetCredentialsFromNetrcCalls() []struct {
	FilterParams github_com_confluentinc_cli_internal_pkg_netrc.NetrcMachineParams
} {
	m.lockGetCredentialsFromNetrc.Lock()
	defer m.lockGetCredentialsFromNetrc.Unlock()

	return m.calls.GetCredentialsFromNetrc
}

// GetCloudCredentialsFromPrompt mocks base method by wrapping the associated func.
func (m *MockLoginCredentialsManager) GetCloudCredentialsFromPrompt(cmd *github_com_spf13_cobra.Command, orgResourceId string) func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error) {
	m.lockGetCloudCredentialsFromPrompt.Lock()
	defer m.lockGetCloudCredentialsFromPrompt.Unlock()

	if m.GetCloudCredentialsFromPromptFunc == nil {
		panic("mocker: MockLoginCredentialsManager.GetCloudCredentialsFromPromptFunc is nil but MockLoginCredentialsManager.GetCloudCredentialsFromPrompt was called.")
	}

	call := struct {
		Cmd           *github_com_spf13_cobra.Command
		OrgResourceId string
	}{
		Cmd:           cmd,
		OrgResourceId: orgResourceId,
	}

	m.calls.GetCloudCredentialsFromPrompt = append(m.calls.GetCloudCredentialsFromPrompt, call)

	return m.GetCloudCredentialsFromPromptFunc(cmd, orgResourceId)
}

// GetCloudCredentialsFromPromptCalled returns true if GetCloudCredentialsFromPrompt was called at least once.
func (m *MockLoginCredentialsManager) GetCloudCredentialsFromPromptCalled() bool {
	m.lockGetCloudCredentialsFromPrompt.Lock()
	defer m.lockGetCloudCredentialsFromPrompt.Unlock()

	return len(m.calls.GetCloudCredentialsFromPrompt) > 0
}

// GetCloudCredentialsFromPromptCalls returns the calls made to GetCloudCredentialsFromPrompt.
func (m *MockLoginCredentialsManager) GetCloudCredentialsFromPromptCalls() []struct {
	Cmd           *github_com_spf13_cobra.Command
	OrgResourceId string
} {
	m.lockGetCloudCredentialsFromPrompt.Lock()
	defer m.lockGetCloudCredentialsFromPrompt.Unlock()

	return m.calls.GetCloudCredentialsFromPrompt
}

// GetOnPremCredentialsFromPrompt mocks base method by wrapping the associated func.
func (m *MockLoginCredentialsManager) GetOnPremCredentialsFromPrompt(cmd *github_com_spf13_cobra.Command) func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error) {
	m.lockGetOnPremCredentialsFromPrompt.Lock()
	defer m.lockGetOnPremCredentialsFromPrompt.Unlock()

	if m.GetOnPremCredentialsFromPromptFunc == nil {
		panic("mocker: MockLoginCredentialsManager.GetOnPremCredentialsFromPromptFunc is nil but MockLoginCredentialsManager.GetOnPremCredentialsFromPrompt was called.")
	}

	call := struct {
		Cmd *github_com_spf13_cobra.Command
	}{
		Cmd: cmd,
	}

	m.calls.GetOnPremCredentialsFromPrompt = append(m.calls.GetOnPremCredentialsFromPrompt, call)

	return m.GetOnPremCredentialsFromPromptFunc(cmd)
}

// GetOnPremCredentialsFromPromptCalled returns true if GetOnPremCredentialsFromPrompt was called at least once.
func (m *MockLoginCredentialsManager) GetOnPremCredentialsFromPromptCalled() bool {
	m.lockGetOnPremCredentialsFromPrompt.Lock()
	defer m.lockGetOnPremCredentialsFromPrompt.Unlock()

	return len(m.calls.GetOnPremCredentialsFromPrompt) > 0
}

// GetOnPremCredentialsFromPromptCalls returns the calls made to GetOnPremCredentialsFromPrompt.
func (m *MockLoginCredentialsManager) GetOnPremCredentialsFromPromptCalls() []struct {
	Cmd *github_com_spf13_cobra.Command
} {
	m.lockGetOnPremCredentialsFromPrompt.Lock()
	defer m.lockGetOnPremCredentialsFromPrompt.Unlock()

	return m.calls.GetOnPremCredentialsFromPrompt
}

// GetPrerunCredentialsFromConfig mocks base method by wrapping the associated func.
func (m *MockLoginCredentialsManager) GetPrerunCredentialsFromConfig(cfg *github_com_confluentinc_cli_internal_pkg_config_v1.Config) func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error) {
	m.lockGetPrerunCredentialsFromConfig.Lock()
	defer m.lockGetPrerunCredentialsFromConfig.Unlock()

	if m.GetPrerunCredentialsFromConfigFunc == nil {
		panic("mocker: MockLoginCredentialsManager.GetPrerunCredentialsFromConfigFunc is nil but MockLoginCredentialsManager.GetPrerunCredentialsFromConfig was called.")
	}

	call := struct {
		Cfg *github_com_confluentinc_cli_internal_pkg_config_v1.Config
	}{
		Cfg: cfg,
	}

	m.calls.GetPrerunCredentialsFromConfig = append(m.calls.GetPrerunCredentialsFromConfig, call)

	return m.GetPrerunCredentialsFromConfigFunc(cfg)
}

// GetPrerunCredentialsFromConfigCalled returns true if GetPrerunCredentialsFromConfig was called at least once.
func (m *MockLoginCredentialsManager) GetPrerunCredentialsFromConfigCalled() bool {
	m.lockGetPrerunCredentialsFromConfig.Lock()
	defer m.lockGetPrerunCredentialsFromConfig.Unlock()

	return len(m.calls.GetPrerunCredentialsFromConfig) > 0
}

// GetPrerunCredentialsFromConfigCalls returns the calls made to GetPrerunCredentialsFromConfig.
func (m *MockLoginCredentialsManager) GetPrerunCredentialsFromConfigCalls() []struct {
	Cfg *github_com_confluentinc_cli_internal_pkg_config_v1.Config
} {
	m.lockGetPrerunCredentialsFromConfig.Lock()
	defer m.lockGetPrerunCredentialsFromConfig.Unlock()

	return m.calls.GetPrerunCredentialsFromConfig
}

// GetOnPremPrerunCredentialsFromEnvVar mocks base method by wrapping the associated func.
func (m *MockLoginCredentialsManager) GetOnPremPrerunCredentialsFromEnvVar() func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error) {
	m.lockGetOnPremPrerunCredentialsFromEnvVar.Lock()
	defer m.lockGetOnPremPrerunCredentialsFromEnvVar.Unlock()

	if m.GetOnPremPrerunCredentialsFromEnvVarFunc == nil {
		panic("mocker: MockLoginCredentialsManager.GetOnPremPrerunCredentialsFromEnvVarFunc is nil but MockLoginCredentialsManager.GetOnPremPrerunCredentialsFromEnvVar was called.")
	}

	call := struct {
	}{}

	m.calls.GetOnPremPrerunCredentialsFromEnvVar = append(m.calls.GetOnPremPrerunCredentialsFromEnvVar, call)

	return m.GetOnPremPrerunCredentialsFromEnvVarFunc()
}

// GetOnPremPrerunCredentialsFromEnvVarCalled returns true if GetOnPremPrerunCredentialsFromEnvVar was called at least once.
func (m *MockLoginCredentialsManager) GetOnPremPrerunCredentialsFromEnvVarCalled() bool {
	m.lockGetOnPremPrerunCredentialsFromEnvVar.Lock()
	defer m.lockGetOnPremPrerunCredentialsFromEnvVar.Unlock()

	return len(m.calls.GetOnPremPrerunCredentialsFromEnvVar) > 0
}

// GetOnPremPrerunCredentialsFromEnvVarCalls returns the calls made to GetOnPremPrerunCredentialsFromEnvVar.
func (m *MockLoginCredentialsManager) GetOnPremPrerunCredentialsFromEnvVarCalls() []struct {
} {
	m.lockGetOnPremPrerunCredentialsFromEnvVar.Lock()
	defer m.lockGetOnPremPrerunCredentialsFromEnvVar.Unlock()

	return m.calls.GetOnPremPrerunCredentialsFromEnvVar
}

// GetOnPremPrerunCredentialsFromNetrc mocks base method by wrapping the associated func.
func (m *MockLoginCredentialsManager) GetOnPremPrerunCredentialsFromNetrc(arg0 *github_com_spf13_cobra.Command, arg1 github_com_confluentinc_cli_internal_pkg_netrc.NetrcMachineParams) func() (*github_com_confluentinc_cli_internal_pkg_auth.Credentials, error) {
	m.lockGetOnPremPrerunCredentialsFromNetrc.Lock()
	defer m.lockGetOnPremPrerunCredentialsFromNetrc.Unlock()

	if m.GetOnPremPrerunCredentialsFromNetrcFunc == nil {
		panic("mocker: MockLoginCredentialsManager.GetOnPremPrerunCredentialsFromNetrcFunc is nil but MockLoginCredentialsManager.GetOnPremPrerunCredentialsFromNetrc was called.")
	}

	call := struct {
		Arg0 *github_com_spf13_cobra.Command
		Arg1 github_com_confluentinc_cli_internal_pkg_netrc.NetrcMachineParams
	}{
		Arg0: arg0,
		Arg1: arg1,
	}

	m.calls.GetOnPremPrerunCredentialsFromNetrc = append(m.calls.GetOnPremPrerunCredentialsFromNetrc, call)

	return m.GetOnPremPrerunCredentialsFromNetrcFunc(arg0, arg1)
}

// GetOnPremPrerunCredentialsFromNetrcCalled returns true if GetOnPremPrerunCredentialsFromNetrc was called at least once.
func (m *MockLoginCredentialsManager) GetOnPremPrerunCredentialsFromNetrcCalled() bool {
	m.lockGetOnPremPrerunCredentialsFromNetrc.Lock()
	defer m.lockGetOnPremPrerunCredentialsFromNetrc.Unlock()

	return len(m.calls.GetOnPremPrerunCredentialsFromNetrc) > 0
}

// GetOnPremPrerunCredentialsFromNetrcCalls returns the calls made to GetOnPremPrerunCredentialsFromNetrc.
func (m *MockLoginCredentialsManager) GetOnPremPrerunCredentialsFromNetrcCalls() []struct {
	Arg0 *github_com_spf13_cobra.Command
	Arg1 github_com_confluentinc_cli_internal_pkg_netrc.NetrcMachineParams
} {
	m.lockGetOnPremPrerunCredentialsFromNetrc.Lock()
	defer m.lockGetOnPremPrerunCredentialsFromNetrc.Unlock()

	return m.calls.GetOnPremPrerunCredentialsFromNetrc
}

// SetCloudClient mocks base method by wrapping the associated func.
func (m *MockLoginCredentialsManager) SetCloudClient(client *github_com_confluentinc_ccloud_sdk_go_v1.Client) {
	m.lockSetCloudClient.Lock()
	defer m.lockSetCloudClient.Unlock()

	if m.SetCloudClientFunc == nil {
		panic("mocker: MockLoginCredentialsManager.SetCloudClientFunc is nil but MockLoginCredentialsManager.SetCloudClient was called.")
	}

	call := struct {
		Client *github_com_confluentinc_ccloud_sdk_go_v1.Client
	}{
		Client: client,
	}

	m.calls.SetCloudClient = append(m.calls.SetCloudClient, call)

	m.SetCloudClientFunc(client)
}

// SetCloudClientCalled returns true if SetCloudClient was called at least once.
func (m *MockLoginCredentialsManager) SetCloudClientCalled() bool {
	m.lockSetCloudClient.Lock()
	defer m.lockSetCloudClient.Unlock()

	return len(m.calls.SetCloudClient) > 0
}

// SetCloudClientCalls returns the calls made to SetCloudClient.
func (m *MockLoginCredentialsManager) SetCloudClientCalls() []struct {
	Client *github_com_confluentinc_ccloud_sdk_go_v1.Client
} {
	m.lockSetCloudClient.Lock()
	defer m.lockSetCloudClient.Unlock()

	return m.calls.SetCloudClient
}

// Reset resets the calls made to the mocked methods.
func (m *MockLoginCredentialsManager) Reset() {
	m.lockGetCloudCredentialsFromEnvVar.Lock()
	m.calls.GetCloudCredentialsFromEnvVar = nil
	m.lockGetCloudCredentialsFromEnvVar.Unlock()
	m.lockGetOnPremCredentialsFromEnvVar.Lock()
	m.calls.GetOnPremCredentialsFromEnvVar = nil
	m.lockGetOnPremCredentialsFromEnvVar.Unlock()
	m.lockGetCredentialsFromConfig.Lock()
	m.calls.GetCredentialsFromConfig = nil
	m.lockGetCredentialsFromConfig.Unlock()
	m.lockGetCredentialsFromNetrc.Lock()
	m.calls.GetCredentialsFromNetrc = nil
	m.lockGetCredentialsFromNetrc.Unlock()
	m.lockGetCloudCredentialsFromPrompt.Lock()
	m.calls.GetCloudCredentialsFromPrompt = nil
	m.lockGetCloudCredentialsFromPrompt.Unlock()
	m.lockGetOnPremCredentialsFromPrompt.Lock()
	m.calls.GetOnPremCredentialsFromPrompt = nil
	m.lockGetOnPremCredentialsFromPrompt.Unlock()
	m.lockGetPrerunCredentialsFromConfig.Lock()
	m.calls.GetPrerunCredentialsFromConfig = nil
	m.lockGetPrerunCredentialsFromConfig.Unlock()
	m.lockGetOnPremPrerunCredentialsFromEnvVar.Lock()
	m.calls.GetOnPremPrerunCredentialsFromEnvVar = nil
	m.lockGetOnPremPrerunCredentialsFromEnvVar.Unlock()
	m.lockGetOnPremPrerunCredentialsFromNetrc.Lock()
	m.calls.GetOnPremPrerunCredentialsFromNetrc = nil
	m.lockGetOnPremPrerunCredentialsFromNetrc.Unlock()
	m.lockSetCloudClient.Lock()
	m.calls.SetCloudClient = nil
	m.lockSetCloudClient.Unlock()
}
