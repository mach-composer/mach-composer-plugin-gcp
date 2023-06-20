package main

import (
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"

	"mach-composer/mach-composer-plugin-gcp/internal"
)

func main() {
	p := internal.NewGcpPlugin()
	plugin.ServePlugin(p)
}
