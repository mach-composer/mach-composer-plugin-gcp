package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
)

func TestValidateRemoteState(t *testing.T) {
	s := getSchema()
	data := map[string]any{
		"plugin": "gcp",
		"bucket": "something",
		"prefix": "my-prefix",
	}

	schema := gojsonschema.NewRawLoader(s.RemoteStateSchema)
	document := gojsonschema.NewRawLoader(data)

	result, err := gojsonschema.Validate(schema, document)
	require.NoError(t, err)
	assert.True(t, result.Valid())
	assert.Empty(t, result.Errors())
}

func TestValidateSiteConfig(t *testing.T) {
	s := getSchema()
	data := map[string]any{
		"zone":    "12345",
		"region":  "region",
		"project": "123456789",
	}

	schema := gojsonschema.NewRawLoader(s.SiteConfigSchema)
	document := gojsonschema.NewRawLoader(data)

	result, err := gojsonschema.Validate(schema, document)
	require.NoError(t, err)
	assert.True(t, result.Valid())
	assert.Empty(t, result.Errors())
}

func TestValidateGlobalConfig(t *testing.T) {
	s := getSchema()
	data := map[string]any{
		"zone":    "12345",
		"region":  "region",
		"project": "123456789",
	}

	schema := gojsonschema.NewRawLoader(s.GlobalConfigSchema)
	document := gojsonschema.NewRawLoader(data)

	result, err := gojsonschema.Validate(schema, document)
	require.NoError(t, err)
	assert.True(t, result.Valid())
	assert.Empty(t, result.Errors())
}
