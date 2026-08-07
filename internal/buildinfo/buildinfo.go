package buildinfo

type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"build_time"`
}

var (
	version   = "unknown"
	commit    = "unknown"
	buildTime = "unknown"
)

func Get() Info {
	return Info{
		Version:   version,
		Commit:    commit,
		BuildTime: buildTime,
	}
}
