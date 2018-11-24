package apps

import (
	"github.com/arunsworld/gopherjs/datetimepicker"
	"github.com/arunsworld/gopherjs/dropzone"
	"github.com/arunsworld/gopherjs/http"
	"github.com/arunsworld/gopherjs/select2"
	"honnef.co/go/js/console"
	"honnef.co/go/js/xhr"

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
						<div class="form-row">
							<div class="form-group col-md-6">
								<label for="select2">Choices</label>
								<select class="form-control" id="select2">
									<option></option>
								</select>
								<div class="invalid-feedback" id="choices-feedback">A selection from Choices is required.</div>
							</div>
							<div class="form-group col-md-6">
								<label for="datetimepicker">Date &amp; Time</label>
								<input class="form-control" type="text" id="datetimepicker" required>
								<div class="invalid-feedback">Date &amp; Time is required.</div>
							</div>
						</div>
						<button class="btn btn-success" id="button">Go!</button>
						<button class="btn btn-secondary" id="reset">Reset</button>
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
					<button class="btn btn-secondary" id="clear">Clear</button>
				</div>
			</div>
		</div>
	</div>
`

// FormContentProducer produces a Div element containing a form and it's behaviors
var FormContentProducer = func() *dom.HTMLDivElement {
	f, result := divElementFromContent(formContentHTML)

	select2Elem := f.GetElementByID("select2")
	s2 := select2.NewSelect2(select2Elem, select2.Options{})
	selectionStream := make(chan *select2.Selection)
	choicesFeedback := f.GetElementByID("choices-feedback").(*dom.HTMLDivElement)
	go func() {
		for {
			selection := <-selectionStream
			if selection.ID != "" {
				choicesFeedback.Style().Set("display", "none")
			}
		}
	}()
	s2.Subscribe(selectionStream)
	go func() {
		rsp := make(chan http.Response, 1)
		http.GET("/api/choices", nil, http.Options{
			ResponseType: xhr.JSON,
		}, rsp)
		ch := <-rsp
		if ch.XhrRequest.Response == nil {
			js.Global.Call("alert", "Could not load choices... please refresh the page.")
			return
		}
		err := s2.ModifyOptions(ch.XhrRequest.Response)
		if err != nil {
			console.Log(err.Error())
			js.Global.Call("alert", "Could not load choices... please refresh the page.")
		}
	}()

	dtpElem := f.GetElementByID("datetimepicker")
	dtp := datetimepicker.NewDatetimepicker(dtpElem, datetimepicker.Options{})

	form := f.GetElementByID("form").(*dom.HTMLFormElement)
	btn := f.GetElementByID("button").(*dom.HTMLButtonElement)
	reset := f.GetElementByID("reset").(*dom.HTMLButtonElement)
	reset.AddEventListener("click", false, func(e dom.Event) {
		choicesFeedback.Style().Set("display", "none")
		resetForm(form, s2)
		e.PreventDefault()
		e.StopPropagation()
	})
	email := f.GetElementByID("email").(*dom.HTMLInputElement)
	pwd := f.GetElementByID("password").(*dom.HTMLInputElement)
	form.AddEventListener("submit", false, func(e dom.Event) {
		e.PreventDefault()
		e.StopPropagation()
		choicesFeedback.Style().Set("display", "none")
		form.Class().Add("was-validated")
		choice := s2.GetSelection()
		if len(choice) == 0 || choice[0].ID == "" {
			choicesFeedback.Style().Set("display", "block")
			return
		}
		if form.CheckValidity() {
			c := js.M{
				"email":    email.Value,
				"password": pwd.Value,
				"datetime": dtp.GetDate(),
				"choice":   choice[0].ID,
			}
			formData := JSON.Call("stringify", c).String()
			response := make(chan http.Response, 1)
			go func() {
				http.POST("/api/form-submit", formData, http.Options{}, response)
				btn.Disabled = true
				defer func() { btn.Disabled = false }()
				rsp := <-response
				if rsp.Error != nil || rsp.Status != 200 {
					js.Global.Call("alert", "Could not save...")
					return
				}
				js.Global.Call("alert", "Submitted successfully...")
				resetForm(form, s2)
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

func resetForm(f *dom.HTMLFormElement, s *select2.Select2) {
	f.Class().Remove("was-validated")
	for _, e := range f.GetElementsByTagName("input") {
		(e.(*dom.HTMLInputElement)).Value = ""
	}
	s.ResetSelection()
}
