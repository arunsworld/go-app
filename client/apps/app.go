package apps

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var (
	// D is the HTMLDocument
	D = dom.GetWindow().Document().(dom.HTMLDocument)
	// H is the history of the window
	H = dom.GetWindow().History()
	// O is the JS Object constructor
	O = js.Global.Get("Object")
	// JSON is the JSON javascript object
	JSON = js.Global.Get("JSON")
)

// AppCollection is a collection of App(s)
type AppCollection interface {
	AddApp(url string, name string, contentProducer appContentProducer)
	Setup(currentURL string)
}

type appContentProducer func() *dom.HTMLDivElement

type app struct {
	url        string
	navName    string
	navElement *dom.HTMLLIElement
	content    *dom.HTMLDivElement
}

type appCollection struct {
	navChannel chan *app
	navs       []string
	navApps    map[string]*app
}

// NewAppCollection creates a new AppCollection
func NewAppCollection() AppCollection {
	ac := &appCollection{}
	ac.navChannel = make(chan *app)
	ac.navApps = make(map[string]*app)
	return ac
}

func (ac *appCollection) AddApp(url string, name string, contentProducer appContentProducer) {
	ac.navs = append(ac.navs, url)
	app := &app{url: url, navName: name, content: contentProducer()}
	ac.navApps[url] = app
}

func (ac *appCollection) Setup(currentURL string) {
	go ac.monitorNavigation()
	leftnav := D.GetElementByID("leftnav").(*dom.HTMLUListElement)
	content := D.GetElementByID("content").(*dom.HTMLDivElement)
	for _, nav := range ac.navs {
		app := ac.navApps[nav]
		li := D.CreateElement("li").(*dom.HTMLLIElement)
		app.navElement = li
		leftnav.AppendChild(li)
		link := D.CreateElement("a").(*dom.HTMLAnchorElement)
		li.AppendChild(link)

		liClass := li.Class()
		liClass.Add("nav-item")
		if nav == currentURL {
			liClass.Add("active")
			content.SetInnerHTML("")
			content.AppendChild(app.content)
		}
		link.SetAttribute("class", "nav-link")
		link.SetAttribute("href", "")
		link.SetTextContent(app.navName)
		link.AddEventListener("click", false, func(e dom.Event) {
			ac.navChannel <- app
			e.PreventDefault()
		})
	}
}

func (ac *appCollection) monitorNavigation() {
	for {
		targetApp := <-ac.navChannel
		for _, nav := range ac.navs {
			app := ac.navApps[nav]
			app.navElement.Class().Remove("active")
		}
		targetApp.navElement.Class().Add("active")
		content := D.GetElementByID("content").(*dom.HTMLDivElement)
		content.SetInnerHTML("")
		content.AppendChild(targetApp.content)
		H.PushState(struct{}{}, "", targetApp.url)
	}
}

func divElementFromContent(html string) (dom.DocumentFragment, *dom.HTMLDivElement) {
	result := D.CreateElement("div").(*dom.HTMLDivElement)
	f := D.CreateDocumentFragment()
	f.AppendChild(result)
	result.SetInnerHTML(html)
	return f, result
}
