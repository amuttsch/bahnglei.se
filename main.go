package main

import (
	"embed"

	"github.com/amuttsch/bahnglei.se/cmd"
)

//go:embed views/*
//go:embed images/*
//go:embed css/*
var assetFS embed.FS

func main() {
	cmd.AssetFS = assetFS
	cmd.Execute()
}
