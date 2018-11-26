package apps

import (
	"github.com/arunsworld/go-js-dom"
)

var homeHTML = `
<div class="row justify-content-sm-center">
    <div class="col-sm-7">
        <h4>Introduction</h4>
        <p><b>Go App</b> is a starter app with a web front-end and a Go back-end designed to build out
            simple apps quickly.</p>
        <p>It is NOT a framework and uses no frameworks. It is designed to demonstrate how an app can be
            built simply and you may use this app as a starting point. To build a new app; download the
            package and start modifying it to your needs (see Installation below).</p>
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
            language we cannot get away from. Others are choices to be made.</p>
        <h5>Key Dependencies</h5>
        <ol>
            <li>The <a href="https://golang.org">Go</a> programming language. Go is used for programming both the
                back-end and the front-end.</li>
            <li><a href="https://github.com/gopherjs/gopherjs">GopherJS</a> library and utility is used to
                produce Javascript to run on the browser. Install with:</li>
            <code>
                go get -u github.com/gopherjs/gopherjs
            </code>
            If you do not wish to pull this dependency in you could always write your own Javascript or use
            Typescript or a Javascript framework. However, it's quite cool to write the front-end in Go as well.
            <li>To compile into a single binary I'm using the <a href="https://github.com/go-playground/statics">statics
                    package</a>.
                You may replace this with your own choice of library to package static files into the binary.
                Interestingly I was
                using gobuffalo's packr2 previously but found that to have too many dependencies. It was a quick change
                to adapt to another package. While this dependency is key the choice of library is not. Install with:
            </li>
            <code>
                go get -u github.com/go-playground/statics
            </code>
        </ol>
        <h5>Recommended Dependencies</h5>
        <ol>
            <li>The <a href="https://github.com/gorilla/mux">gorilla/mux</a> package is a lovely URL router
                and dispatcher. I'm using it and would highly recommend it for REST APIs. However, this is
                not mandatory at all and can be easily replaced with the <code>http/ServeMux</code> multiplexer
                from the standard library.</li>
            <li>I'm using 2 packages for security: <a href="https://github.com/unrolled/secure">unrolled/secure</a>
                and <a href="https://github.com/rs/cors">rs/cors</a>. These are middlewares that help
                protect against OWASP vulnerabilities such as Cross Origin &amp; Cross-site scripting.</li>
            <li>I'm also using the <a href="https://github.com/NYTimes/gziphandler">NYTimes/gziphandler</a>
                package to compress the traffic.
                While this is entirely optional I'd highly recommend it.</li>
        </ol>
        <h5>Application specific Dependencies</h5>
        <ol>
            <li>This demo app includes a Chat service. To power that I'm using the
                <a href="https://github.com/gorilla/websocket">gorilla/websocket</a> package.
                I've paired that with a GopherJS websocket binding.</li>
            <li>In addition to Go dependencies I'm also using Bootstrap and other JS libraries
                such as Dropzone. These are entirely optional and application specific and can
                be swapped out for your use-case</li>
        </ol>
        <hr />
        <h4>Installation</h4>
        <p>First install the dependencies:</p>
        <pre>
        go get -u github.com/gopherjs/gopherjs
        go get -u github.com/go-playground/statics
        </pre>
        <p>Then install go-app:</p>
        <pre>
        go get -u github.com/arunsworld/go-app
        go get -u github.com/arunsworld/go-app/client
        </pre>
        <p>Install kick to make development easy (optional):</p>
        <pre>
        go get -u github.com/arunsworld/go-kick
        </pre>
        <hr />
        <h4>Execution</h4>
        <ol>
            <li><code>make build</code> builds the binary ready for execution &amp; deployment.</li>
            <li><code>make dockerize</code> creates a docker image (requires docker).</li>
            <li><code>make kick</code> runs the app using kick and recompiles when we see changes (requires kick).</li>
        </ol>
        <hr />
    </div>
</div>
`

// HomeContentProducer produces a Div element containing the contents of Home
var HomeContentProducer = func() *dom.HTMLDivElement {
	_, result := divElementFromContent(homeHTML)
	return result
}
