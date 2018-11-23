package apps

import (
	"honnef.co/go/js/dom"
)

// HomeContentProducer produces a Div element containing the contents of Home
var HomeContentProducer = func() *dom.HTMLDivElement {
	_, result := divElementFromContent("<div>This is my Home!</div>")
	return result
}
