package main

// Params are the parameters passed to this extension via the custom properties of the estafette stage
type Params struct {
	ReleaseVersion         string   `json:"version,omitempty" yaml:"version,omitempty"`
	CloseMilestone         bool     `json:"closeMilestone,omitempty" yaml:"closeMilestone,omitempty"`
	ReleaseTitle           string   `json:"title,omitempty" yaml:"title,omitempty"`
	IgnoreMissingMilestone bool     `json:"ignoreMissingMilestone,omitempty" yaml:"ignoreMissingMilestone,omitempty"`
	Assets                 []string `json:"assets,omitempty" yaml:"assets,omitempty"`
}

// SetDefaults fills in empty fields with convention-based defaults
func (p *Params) SetDefaults(buildVersion, gitRepoName string) {

	if p.ReleaseVersion == "" {
		p.ReleaseVersion = buildVersion
	}

	if p.ReleaseTitle == "" {
		p.ReleaseTitle = capitalize(gitRepoName)
	}
}
