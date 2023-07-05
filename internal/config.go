package internal

import (
	"fmt"
)

type GCSTFState struct {
	Bucket string `mapstructure:"bucket"`
	prefix string `mapstructure:"prefix"`
}

func (a GCSTFState) Key(site string) string {
	if a.prefix == "" {
		return site
	}
	return fmt.Sprintf("%s/%s", a.prefix, site)
}

type GlobalConfig struct {
	Project string `mapstructure:"project"`
	Region  string `mapstructure:"region"`
	Zone    string `mapstructure:"zone"`
	Beta    bool   `mapstructure:"beta"`
}

type SiteConfig struct {
	Project string `mapstructure:"project"`
	Region  string `mapstructure:"region"`
	Zone    string `mapstructure:"zone"`
	Beta    bool   `mapstructure:"beta"`
}

func (a *SiteConfig) merge(c *GlobalConfig) {
	if a.Project == "" {
		a.Project = c.Project
	}
	if a.Region == "" {
		a.Region = c.Region
	}
	if a.Zone == "" {
		a.Zone = c.Zone
	}
	if !a.Beta && c.Beta {
		a.Beta = c.Beta
	}
}

func (a *SiteConfig) providers() []string {
	if a.Beta {
		return []string{"google = google-beta"}
	}
	return []string{"google = google"}
}
