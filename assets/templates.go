//go:generate statics -i=static -o=assets/templates.go -pkg=assets -group=Templates

package assets

import "github.com/go-playground/statics/static"

// newStaticTemplates initializes a new *static.Files instance for use
func newStaticTemplates(config *static.Config) (*static.Files, error) {

	return static.New(config, &static.DirFile{})
}
