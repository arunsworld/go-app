//go:generate statics -i=static -o=assets/static.go -pkg=assets -group=Assets

package assets

import "github.com/go-playground/statics/static"

// newStaticAssets initializes a new *static.Files instance for use
func newStaticAssets(config *static.Config) (*static.Files, error) {

	return static.New(config, &static.DirFile{})
}
