package apps

import (
	"github.com/arunsworld/gopherjs/dropzone"
	"github.com/arunsworld/gopherjs/http"

	"github.com/gopherjs/gopherjs/js"

	"honnef.co/go/js/dom"
)

var formContentHTML = `
	<div class="row">
		<div class="col-md-6">
			<div class="card">
				<div class="card-header">A Demo Form</div>
				<div class="card-body">
					<form id="form" class="needs-validation" novalidate>
						<div class="form-row">
							<div class="form-group col-md-6">
								<label for="email">Email</label>
								<input type="email" class="form-control" id="email" name="email" 
									placeholder="Enter email" required>
								<div class="invalid-feedback">Proper email address required.</div>
							</div>
							<div class="form-group col-md-6">
								<label for="password">Password</label>
								<input type="password" class="form-control" id="password" name="password" 
									placeholder="Password" required minlength="6">
								<div class="invalid-feedback">A password is required. It should have atleast 6 letters.</div>
							</div>
						</div>
						<button class="btn btn-success" id="button">Click Me!</button>
						<button class="btn btn-default" id="reset">Reset</button>
					</form>
				</div>
			</div>
		</div>
		<div class="col-md-6">
			<div class="card">
				<div class="card-header">File Uploads</div>
				<div class="card-body">
					<div id="uploadFile" class="dropzone">
						<div class="dz-message" style="font-size: 2rem; color: #967ADC; text-align: center;">Upload</div>
					</div>
					<br/>
					<button class="btn btn-default" id="clear">Clear</button>
				</div>
			</div>
		</div>
	</div>
`

type creds struct {
	*js.Object
	email    string `js:"email"`
	password string `js:"password"`
}

// FormContentProducer produces a Div element containing a form and it's behaviors
var FormContentProducer = func() *dom.HTMLDivElement {
	f, result := divElementFromContent(formContentHTML)

	form := f.GetElementByID("form").(*dom.HTMLFormElement)
	btn := f.GetElementByID("button").(*dom.HTMLButtonElement)
	reset := f.GetElementByID("reset").(*dom.HTMLButtonElement)
	reset.AddEventListener("click", false, func(e dom.Event) {
		resetForm(form)
		e.PreventDefault()
		e.StopPropagation()
	})
	email := f.GetElementByID("email").(*dom.HTMLInputElement)
	pwd := f.GetElementByID("password").(*dom.HTMLInputElement)
	form.AddEventListener("submit", false, func(e dom.Event) {
		e.PreventDefault()
		e.StopPropagation()
		form.Class().Add("was-validated")
		if form.CheckValidity() {
			c := &creds{Object: js.Global.Get("Object").New()}
			c.email = email.Value
			c.password = pwd.Value
			formData := js.Global.Get("JSON").Call("stringify", c).String()
			options := http.Options{}
			response := make(chan http.Response, 1)
			go func() {
				http.POST("/api/form-submit", formData, options, response)
				btn.Disabled = true
				defer func() { btn.Disabled = false }()
				rsp := <-response
				if rsp.Error != nil || rsp.Status != 200 {
					js.Global.Call("alert", "Could not save...")
					return
				}
				js.Global.Call("alert", "Submitted successfully...")
				resetForm(form)
			}()
		}
	})

	uploadFile := f.GetElementByID("uploadFile").(*dom.HTMLDivElement)
	dz := dropzone.NewDropzone(uploadFile, dropzone.Props{
		URL:    "/api/upload",
		Params: map[string]string{"param1": "value in param1"},
	})
	clear := f.GetElementByID("clear").(*dom.HTMLButtonElement)
	clear.AddEventListener("click", false, func(e dom.Event) {
		dz.RemoveFiles()
	})

	return result
}

func resetForm(f *dom.HTMLFormElement) {
	f.Class().Remove("was-validated")
	for _, e := range f.GetElementsByTagName("input") {
		(e.(*dom.HTMLInputElement)).Value = ""
	}
}
