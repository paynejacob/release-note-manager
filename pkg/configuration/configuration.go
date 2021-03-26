package configuration

type Section struct {
	Name   string
	Labels []string
}


type Configuration struct {
	includeLabel string
	excludeLabel string

	Sections []Section
}

var DefaultConfiguration = Configuration{
	includeLabel: "release-note",
	excludeLabel: "no-release-note",
	Sections:     []Section{
		{Name: "Bug Fixes", Labels: []string{"bug"}},
		{Name: "Major Bug Fixes", Labels: []string{"internal"}},
		{Name: "Enhancements", Labels: []string{"enhancement"}},
	},
}