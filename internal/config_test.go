package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	site := "my-site"
	prefix := "my-prefix"

	// no prefix
	state := GCSTFState{
		Bucket: "bucket",
	}
	assert.Equal(t, site, state.Key(site))

	// prefix
	state = GCSTFState{
		Bucket: "bucket",
		prefix: prefix,
	}
	assert.Equal(t, fmt.Sprintf("%s/%s", prefix, site), state.Key(site))
}

func TestProviders(t *testing.T) {
	// no beta
	config := SiteConfig{}
	assert.Equal(t, []string{
		"google = google",
		"google-beta = google-beta",
	}, config.providers())
}
