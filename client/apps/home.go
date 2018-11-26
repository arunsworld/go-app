package apps

import (
	"github.com/arunsworld/go-js-dom"
)

var homeHTML = `
<div class="row justify-content-sm-center">
    <div class="col-sm-7">
        <h4>Introduction</h4>
        <p><b>Go App</b> is a starter app with a web front-end and a Go back-end designed to build out simple
            apps quickly.</p>
        <p>It is NOT a framework and uses no frameworks. It is designed to demonstrate how an app can be
            built simply and you may use this app as a starting point. To build a new app; download the package and
            start
            modifying it to your needs.</p>
        <p>Go App has a few goals:
            <ul>
                <li>Independence from any frameworks. This let's your code be agnostic and portable.</li>
                <li>Compile into a single binary.</li>
                <li>Be secure enough to deploy directly in the wild.</li>
            </ul>
        </p>
        <hr />
        <h4>Dependencies</h4>
        <p>While the goal of Go App is to minimize dependencies some dependencies like the programming
            language we cannot get away from. All the dependencies are listed and categorized below.</p>
        <h5>Key Dependencies</h5>
        <ol>
            <li>The <a href="https://golang.org">Go</a> programming language. This package uses Go for programming
                both the back-end and the front-end.</li>
            <li><a href="https://github.com/gopherjs/gopherjs">GopherJS</a> library and utility is used to
                produce Javascript to run on the browser. Install with:</li>
            <code>
            go get -u github.com/gopherjs/gopherjs
            </code>
            <li>To compile into a single binary I'm using <a href="https://github.com/gobuffalo/packr/tree/master/v2">packr2</a>.
            You may replace this with your own choice of library to package static files into the binary.
            </li>
        </ol>
    </div>
</div>
`

// HomeContentProducer produces a Div element containing the contents of Home
var HomeContentProducer = func() *dom.HTMLDivElement {
	_, result := divElementFromContent(homeHTML)
	return result
}
