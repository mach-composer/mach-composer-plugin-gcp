package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalConfig(t *testing.T) {
	plugin := NewGcpPlugin()
	err := plugin.Configure("dev", "0.0.1")
	require.NoError(t, err)

	err = plugin.SetGlobalConfig(map[string]any{
		"project": "0123456789",
		"region":  "us-central1",
		"zone":    "us-central1-a",
	})
	require.NoError(t, err)

	result, err := plugin.RenderTerraformResources("my-site")
	require.NoError(t, err)
	assert.Contains(t, result, "provider \"google\"")
	assert.Contains(t, result, "project = \"0123456789\"")
	assert.Contains(t, result, "region = \"us-central1\"")
	assert.Contains(t, result, "zone = \"us-central1-a\"")

	result, err = plugin.RenderTerraformProviders("my-site")
	require.NoError(t, err)
	assert.Contains(t, result, "version = \"~> 0.0.1\"")

	err = plugin.SetGlobalConfig(map[string]any{
		"project": "0123456789",
		"region":  "us-central1",
		"zone":    "us-central1-a",
		"beta":    true,
	})
	require.NoError(t, err)

	result, err = plugin.RenderTerraformResources("my-site")
	require.NoError(t, err)
	assert.Contains(t, result, "provider \"google-beta\"")
	assert.Contains(t, result, "project = \"0123456789\"")
	assert.Contains(t, result, "region = \"us-central1\"")
	assert.Contains(t, result, "zone = \"us-central1-a\"")

	result, err = plugin.RenderTerraformProviders("my-site")
	require.NoError(t, err)
	assert.Contains(t, result, "google-beta = {")
}

func TestSiteConfig(t *testing.T) {
	plugin := NewGcpPlugin()
	err := plugin.SetGlobalConfig(map[string]any{
		"project": "0123456789",
		"region":  "us-central1",
		"zone":    "us-central1-a",
	})
	require.NoError(t, err)

	// check overrides
	err = plugin.SetSiteConfig("my-site", map[string]any{
		"project": "987654320",
		"region":  "europe-west1",
		"zone":    "europe-west1-b",
		"beta":    true,
	})
	require.NoError(t, err)

	err = plugin.SetComponentConfig("my-component", map[string]any{
		"integrations": []string{"gcp"},
	})
	require.NoError(t, err)

	result, err := plugin.RenderTerraformResources("my-site")
	require.NoError(t, err)
	assert.Contains(t, result, "project = \"987654320\"")
	assert.Contains(t, result, "provider \"google-beta\"")
	assert.NotContains(t, result, "us-central1")
}

func TestSetRemoteStateBackend(t *testing.T) {
	plugin := NewGcpPlugin()

	err := plugin.SetGlobalConfig(map[string]any{
		"project": "0123456789",
		"region":  "us-central1",
		"zone":    "us-central1-a",
	})
	require.NoError(t, err)

	err = plugin.SetRemoteStateBackend(map[string]any{
		"plugin": "gcp",
		"bucket": "0123456789",
		"prefix": "us-central1",
	})
	require.NoError(t, err)

	err = plugin.SetSiteConfig("my-site", map[string]any{})
	require.NoError(t, err)

	result, err := plugin.RenderTerraformStateBackend("my-site")
	require.NoError(t, err)
	assert.Contains(t, result, "backend \"gcs\"")
}

func TestRenderTerraformProviders(t *testing.T) {
	plugin := NewGcpPlugin()
	err := plugin.Configure("dev", "0.0.1")
	require.NoError(t, err)

	err = plugin.SetGlobalConfig(map[string]any{
		"project": "0123456789",
		"region":  "us-central1",
		"zone":    "us-central1-a",
	})
	require.NoError(t, err)

	result, err := plugin.RenderTerraformProviders("my-site")
	require.NoError(t, err)
	assert.Contains(t, result, "version = \"~> 0.0.1\"")
}

func TestTerraformRenderResources(t *testing.T) {
	plugin := NewGcpPlugin()
	err := plugin.Configure("dev", "0.0.1")
	require.NoError(t, err)

	err = plugin.SetGlobalConfig(map[string]any{
		"project": "0123456789",
		"region":  "us-central1",
		"zone":    "us-central1-a",
	})
	require.NoError(t, err)

	err = plugin.SetSiteConfig("my-site", map[string]any{})
	require.NoError(t, err)

	result, err := plugin.RenderTerraformResources("my-site")
	require.NoError(t, err)
	assert.Contains(t, result, "locals {")
	assert.Contains(t, result, "\ttags = {")
	assert.Contains(t, result, "\t\tSite = \"my-site\"")
	assert.Contains(t, result, "\t\tEnvironment = \"dev\"")
	assert.Contains(t, result, "\t}")
}
