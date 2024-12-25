package main

import (
	"embed"

	"github.com/amuttsch/bahnglei.se/cmd"
	"github.com/amuttsch/bahnglei.se/templates"
	"github.com/benbjohnson/hashfs"
)

//go:embed assets
var assetFS embed.FS
var fsys = hashfs.NewFS(assetFS)

//go:embed translations/*.toml
var translationsFS embed.FS

func main() {
	cmd.AssetFS = fsys
	cmd.TranslationsFS = translationsFS
	templates.AssetFS = fsys
	cmd.Execute()
}
