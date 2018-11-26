package assets

import (
	"log"
	"os"

	"github.com/go-playground/statics/static"
)

func runningInDev() string {
	path, found := os.LookupEnv("DEV")
	if found {
		return path
	}
	return ""
}

// StaticFiles returns a filesystem to the embedded static files
func StaticFiles() *static.Files {
	config := &static.Config{
		UseStaticFiles: true,
	}
	path := runningInDev()
	if path != "" {
		config.FallbackToDisk = true
		config.AbsPkgPath = path + "/assets/static/"
	}

	files, err := newStaticAssets(config)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

// TemplatesFiles returns a filesystem to the embedded templates files
func TemplatesFiles() *static.Files {
	config := &static.Config{
		UseStaticFiles: true,
	}
	path := runningInDev()
	if path != "" {
		config.FallbackToDisk = true
		config.AbsPkgPath = path + "/assets/templates/"
	}
	files, err := newStaticTemplates(config)
	if err != nil {
		log.Fatal(err)
	}
	return files
}
