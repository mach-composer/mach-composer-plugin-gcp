package internal

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/mach-composer/mach-composer-plugin-helpers/helpers"
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/mitchellh/mapstructure"
)

type Plugin struct {
	provider     string
	environment  string
	remoteState  *GCSTFState
	globalConfig *GlobalConfig
	siteConfigs  map[string]*SiteConfig
	enabled      bool
}

func NewGcpPlugin() schema.MachComposerPlugin {
	state := &Plugin{
		provider:    "4",
		siteConfigs: map[string]*SiteConfig{},
	}

	return plugin.NewPlugin(&schema.PluginSchema{
		Identifier:          "gcp",
		Configure:           state.Configure,
		IsEnabled:           func() bool { return state.enabled },
		GetValidationSchema: state.GetValidationSchema,

		// Config
		SetRemoteStateBackend: state.SetRemoteStateBackend,
		SetGlobalConfig:       state.SetGlobalConfig,
		SetSiteConfig:         state.SetSiteConfig,

		// Renders
		RenderTerraformStateBackend: state.TerraformRenderStateBackend,
		RenderTerraformProviders:    state.TerraformRenderProviders,
		RenderTerraformResources:    state.TerraformRenderResources,
		RenderTerraformComponent:    state.RenderTerraformComponent,
	})
}

func (p *Plugin) Configure(environment string, provider string) error {
	p.environment = environment
	if provider != "" {
		p.provider = provider
	}
	return nil
}

func (p *Plugin) GetValidationSchema() (*schema.ValidationSchema, error) {
	result := getSchema()
	return result, nil
}

func (p *Plugin) SetRemoteStateBackend(data map[string]any) error {
	state := &GCSTFState{}
	if err := mapstructure.Decode(data, state); err != nil {
		return err
	}
	if err := defaults.Set(state); err != nil {
		return err
	}
	p.remoteState = state
	return nil
}

func (p *Plugin) SetGlobalConfig(data map[string]any) error {
	cfg := GlobalConfig{}

	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}

	p.globalConfig = &cfg
	p.enabled = true
	return nil
}

func (p *Plugin) SetSiteConfig(site string, data map[string]any) error {
	if data == nil {
		return nil
	}

	if p.globalConfig == nil {
		return fmt.Errorf("a global gcp config is required for setting per-site configuration")
	}

	cfg := SiteConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	cfg.merge(p.globalConfig)

	p.siteConfigs[site] = &cfg
	p.enabled = true
	return nil
}

func (p *Plugin) getSiteConfig(site string) *SiteConfig {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		value := SiteConfig{}
		value.merge(p.globalConfig)
		return &value
	}
	return cfg
}

func (p *Plugin) TerraformRenderStateBackend(site string) (string, error) {
	if p.remoteState == nil {
		return "", nil
	}

	templateContext := struct {
		Bucket string
		Prefix string
	}{
		Bucket: p.remoteState.Bucket,
		Prefix: p.remoteState.Key(site),
	}

	template := `
	backend "gcs" {
	  bucket  = "{{ .Bucket }}"
	  prefix = "{{ .Prefix }}"
	}
	`
	return helpers.RenderGoTemplate(template, templateContext)
}

func (p *Plugin) TerraformRenderProviders(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	var result string = fmt.Sprintf(`
	google = {
		source = "hashicorp/google"
		version = "%s"
	}
	google-beta = {
		source = "hashicorp/google-beta"
		version = "~> 4"
	}`, helpers.VersionConstraint(p.provider))

	return result, nil
}

func (p *Plugin) TerraformRenderResources(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	templateContext := struct {
		Project     string
		Region      string
		Zone        string
		SiteName    string
		Environment string
	}{
		Project:     cfg.Project,
		Region:      cfg.Region,
		Zone:        cfg.Zone,
		SiteName:    site,
		Environment: p.environment,
	}

	template := `
		provider "google" {
			{{ renderProperty "project" .Project}}
			{{ renderProperty "region" .Region}}
			{{ renderProperty "zone" .Zone}}
		}

		provider "google-beta" {
			{{ renderProperty "project" .Project}}
			{{ renderProperty "region" .Region}}
			{{ renderProperty "zone" .Zone}}
		}

		locals {
			tags = {
				Site = "{{ .SiteName }}"
				Environment = "{{ .Environment }}"
			}
		}
	`
	return helpers.RenderGoTemplate(template, templateContext)
}

func (p *Plugin) RenderTerraformComponent(site string, _ string) (*schema.ComponentSchema, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return nil, nil
	}

	result := &schema.ComponentSchema{
		Providers: cfg.providers(),
		DependsOn: []string{},
	}

	return result, nil
}
