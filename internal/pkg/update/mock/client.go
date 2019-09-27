// Code generated by mocker. DO NOT EDIT.
// github.com/travisjeffery/mocker
// Source: client.go

package mock

import (
	sync "sync"
)

// Client is a mock of Client interface
type Client struct {
	lockCheckForUpdates sync.Mutex
	CheckForUpdatesFunc func(name, currentVersion string, forceCheck bool) (bool, string, error)

	lockPromptToDownload sync.Mutex
	PromptToDownloadFunc func(name, currVersion, latestVersion string, confirm bool) bool

	lockUpdateBinary sync.Mutex
	UpdateBinaryFunc func(name, version, path string) error

	calls struct {
		CheckForUpdates []struct {
			Name           string
			CurrentVersion string
			ForceCheck     bool
		}
		PromptToDownload []struct {
			Name          string
			CurrVersion   string
			LatestVersion string
			Confirm       bool
		}
		UpdateBinary []struct {
			Name    string
			Version string
			Path    string
		}
	}
}

// CheckForUpdates mocks base method by wrapping the associated func.
func (m *Client) CheckForUpdates(name, currentVersion string, forceCheck bool) (bool, string, error) {
	m.lockCheckForUpdates.Lock()
	defer m.lockCheckForUpdates.Unlock()

	if m.CheckForUpdatesFunc == nil {
		panic("mocker: Client.CheckForUpdatesFunc is nil but Client.CheckForUpdates was called.")
	}

	call := struct {
		Name           string
		CurrentVersion string
		ForceCheck     bool
	}{
		Name:           name,
		CurrentVersion: currentVersion,
		ForceCheck:     forceCheck,
	}

	m.calls.CheckForUpdates = append(m.calls.CheckForUpdates, call)

	return m.CheckForUpdatesFunc(name, currentVersion, forceCheck)
}

// CheckForUpdatesCalled returns true if CheckForUpdates was called at least once.
func (m *Client) CheckForUpdatesCalled() bool {
	m.lockCheckForUpdates.Lock()
	defer m.lockCheckForUpdates.Unlock()

	return len(m.calls.CheckForUpdates) > 0
}

// CheckForUpdatesCalls returns the calls made to CheckForUpdates.
func (m *Client) CheckForUpdatesCalls() []struct {
	Name           string
	CurrentVersion string
	ForceCheck     bool
} {
	m.lockCheckForUpdates.Lock()
	defer m.lockCheckForUpdates.Unlock()

	return m.calls.CheckForUpdates
}

// PromptToDownload mocks base method by wrapping the associated func.
func (m *Client) PromptToDownload(name, currVersion, latestVersion string, confirm bool) bool {
	m.lockPromptToDownload.Lock()
	defer m.lockPromptToDownload.Unlock()

	if m.PromptToDownloadFunc == nil {
		panic("mocker: Client.PromptToDownloadFunc is nil but Client.PromptToDownload was called.")
	}

	call := struct {
		Name          string
		CurrVersion   string
		LatestVersion string
		Confirm       bool
	}{
		Name:          name,
		CurrVersion:   currVersion,
		LatestVersion: latestVersion,
		Confirm:       confirm,
	}

	m.calls.PromptToDownload = append(m.calls.PromptToDownload, call)

	return m.PromptToDownloadFunc(name, currVersion, latestVersion, confirm)
}

// PromptToDownloadCalled returns true if PromptToDownload was called at least once.
func (m *Client) PromptToDownloadCalled() bool {
	m.lockPromptToDownload.Lock()
	defer m.lockPromptToDownload.Unlock()

	return len(m.calls.PromptToDownload) > 0
}

// PromptToDownloadCalls returns the calls made to PromptToDownload.
func (m *Client) PromptToDownloadCalls() []struct {
	Name          string
	CurrVersion   string
	LatestVersion string
	Confirm       bool
} {
	m.lockPromptToDownload.Lock()
	defer m.lockPromptToDownload.Unlock()

	return m.calls.PromptToDownload
}

// UpdateBinary mocks base method by wrapping the associated func.
func (m *Client) UpdateBinary(name, version, path string) error {
	m.lockUpdateBinary.Lock()
	defer m.lockUpdateBinary.Unlock()

	if m.UpdateBinaryFunc == nil {
		panic("mocker: Client.UpdateBinaryFunc is nil but Client.UpdateBinary was called.")
	}

	call := struct {
		Name    string
		Version string
		Path    string
	}{
		Name:    name,
		Version: version,
		Path:    path,
	}

	m.calls.UpdateBinary = append(m.calls.UpdateBinary, call)

	return m.UpdateBinaryFunc(name, version, path)
}

// UpdateBinaryCalled returns true if UpdateBinary was called at least once.
func (m *Client) UpdateBinaryCalled() bool {
	m.lockUpdateBinary.Lock()
	defer m.lockUpdateBinary.Unlock()

	return len(m.calls.UpdateBinary) > 0
}

// UpdateBinaryCalls returns the calls made to UpdateBinary.
func (m *Client) UpdateBinaryCalls() []struct {
	Name    string
	Version string
	Path    string
} {
	m.lockUpdateBinary.Lock()
	defer m.lockUpdateBinary.Unlock()

	return m.calls.UpdateBinary
}

// Reset resets the calls made to the mocked methods.
func (m *Client) Reset() {
	m.lockCheckForUpdates.Lock()
	m.calls.CheckForUpdates = nil
	m.lockCheckForUpdates.Unlock()
	m.lockPromptToDownload.Lock()
	m.calls.PromptToDownload = nil
	m.lockPromptToDownload.Unlock()
	m.lockUpdateBinary.Lock()
	m.calls.UpdateBinary = nil
	m.lockUpdateBinary.Unlock()
}
