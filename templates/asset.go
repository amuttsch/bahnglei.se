package templates

import (
	"fmt"

	"github.com/benbjohnson/hashfs"
)

var AssetFS *hashfs.FS

func Asset(name string) string {
	fmt.Printf("Asset: %s", name)
	return "/assets/" + AssetFS.HashName(name)
}
