package version

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type Version struct {
	Binary    string
	Name      string
	Version   string
	Commit    string
	BuildDate string
	BuildHost string
	UserAgent string // http
	ClientID  string // kafka
}

func NewVersion(cliName, version, commit, buildDate, buildHost string) *Version {
	name := GetFullCLIName(cliName)
	support := getSupport(cliName)
	return &Version{
		Binary:    cliName,
		Name:      name,
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
		BuildHost: buildHost,
		UserAgent: fmt.Sprintf("%s/%s (%s)", strings.ReplaceAll(name, " ", "-"), version, support),
		ClientID:  fmt.Sprintf("%s_%s", strings.ReplaceAll(name, " ", "-"), version),
	}
}

func (v *Version) IsReleased() bool {
	return v.Version != "0.0.0" && !strings.Contains(v.Version, "dirty") && !strings.Contains(v.Version, "-g")
}

// String returns the version in a standardized format
func (v *Version) String() string {
	return fmt.Sprintf(`%s - %s

Version:     %s
Git Ref:     %s
Build Date:  %s
Go Version:  %s (%s/%s)
Development: %s
`,
		v.Binary,
		v.Name,
		v.Version,
		v.Commit,
		v.BuildDate,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		strconv.FormatBool(!v.IsReleased()))
}

func GetFullCLIName(cliName string) string {
	if cliName == "ccloud" {
		return "Confluent Cloud CLI"
	}
	return "Confluent CLI"
}

func getSupport(cliName string) string {
	if cliName == "ccloud" {
		return "https://confluent.cloud; support@confluent.io"
	}
	return "https://confluent.io; support@confluent.io"
}
