package templates

import (
	"github.com/benbjohnson/hashfs"
)

var AssetFS *hashfs.FS

func Asset(name string) string {
	return "/" + AssetFS.HashName("assets/"+name)
}
