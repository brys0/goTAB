package os

import "github.com/shirou/gopsutil/host"

type OS struct {
	//PrettyName string `json:"pretty_name"`
	Name      string `json:"name"`
	VersionID string `json:"version_id"`
	Version   string `json:"version"`
	//VersionCodename string `json:"version_codename"`
	//Id string `json:"id"`
	//HomeURL string `json:"home_url"`
	//SupportURL string `json:"support_url"`
	//BugReportURL string `json:"bug_report_url"`
}

func GetOSInfo() (*OS, error) {
	os, err := host.Info()

	if err != nil {
		return nil, err
	}

	return &OS{
		Name:      os.OS,
		Version:   os.Platform,
		VersionID: os.PlatformVersion,
	}, nil
}
