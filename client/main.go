package main

import (
	"github.com/arunsworld/go-app/client/apps"
	"github.com/arunsworld/gopherjs/dropzone"
	"honnef.co/go/js/dom"
)

var (
	// D is the HTMLDocument
	D = dom.GetWindow().Document().(dom.HTMLDocument)
	// AC is the collection of Apps
	AC = apps.NewAppCollection()
)

func init() {
	dropzone.AutoDiscover(false)
}

func main() {
	url := currentURL()
	AC.AddApp("/", "Home", apps.HomeContentProducer)
	AC.AddApp("/form", "Form", apps.FormContentProducer)
	AC.AddApp("/chat", "Chat", apps.ChatContentProducer)
	AC.Setup(url)
}

func currentURL() string {
	l := D.CreateElement("a").(*dom.HTMLAnchorElement)
	l.Href = D.URL()
	return l.Pathname
}
