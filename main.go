package main

import (
	"embed"

	"github.com/amuttsch/bahnglei.se/cmd"
	"github.com/benbjohnson/hashfs"
)

//go:embed views/*
//go:embed images/*
//go:embed css/style.css
var assetFS embed.FS
var fsys = hashfs.NewFS(assetFS)

func main() {
	cmd.AssetFS = fsys
	cmd.Execute()
}
