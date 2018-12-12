package apps

import (
	"github.com/arunsworld/go-js-dom"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/console"
)

var qrcodeContentHTML = `
	<div id="qrcode" width="512" height="512" style="display:none;"></div>
	<div class="row justify-content-sm-center">
		<div class="col-sm-6">
			<div class="card-deck">
				<div class="card">
					<div class="card-body">
						<img id="qrcode-img" class="card-img-top" max-width="256" max-height="256"></img>
						<hr/>
						<form>
							<div class="form-group">
								<input type="password" class="form-control" id="qrcode-text" placeholder="Password: transmit via QR">
							</div>
						</form>
					</div>
				</div>
				<div class="card">
					<div class="card-body">
						<div id="initialize-camera" class="text-center">
							<button type="submit" class="btn btn-primary" id="initialize-camera-btn">Initialize Camera</button>
						</div>
						<div id="looking-for-camera" style="display: none;">
							<p>Looking for camera...</p>
						</div>
						<div id="camera-found" style="display: none;">
							<p id="scan-results" class="text-center"></p>
							<form class="text-center">
								<div class="btn-group" id="actions">
									<button type="submit" class="btn btn-danger" id="scanStop">Stop</button>
								</div>
							</form>
							<br/>
							<video id="preview" class="card-img-top" max-width="256" max-height="256"></video>
							<br/>
							<p>Cameras:</p>
							<ul id="camera-list">
							</ul>
						</div>
						<div id="camera-not-found" style="display: none;">
							<p>No camera was found in your system...</p>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
`

// QRCodeContentProducer produces the QRCode Page
var QRCodeContentProducer = func() *dom.HTMLDivElement {
	f, result := divElementFromContent(qrcodeContentHTML)

	qrcode := js.Global.Call("jQuery", f.GetElementByID("qrcode"))
	qrcodeImg := f.GetElementByID("qrcode-img")
	qrcodeTxt := f.GetElementByID("qrcode-text").(*dom.HTMLInputElement)

	qrcodeTxt.AddEventListener("keyup", false, func(e dom.Event) {
		renderQRCode(qrcode, qrcodeImg, qrcodeTxt.Value)
	})

	renderQRCode(qrcode, qrcodeImg, "This is a QR Code Demo!")

	lookingForCamera := f.GetElementByID("looking-for-camera").(*dom.HTMLDivElement)
	cameraFound := f.GetElementByID("camera-found").(*dom.HTMLDivElement)
	cameraNotFound := f.GetElementByID("camera-not-found").(*dom.HTMLDivElement)
	preview := f.GetElementByID("preview")
	scanResults := f.GetElementByID("scan-results").(*dom.HTMLParagraphElement)

	initializeCamera := f.GetElementByID("initialize-camera").(*dom.HTMLDivElement)
	initializeCameraBtn := f.GetElementByID("initialize-camera-btn").(*dom.HTMLButtonElement)
	scanStop := f.GetElementByID("scanStop").(*dom.HTMLButtonElement)
	actions := f.GetElementByID("actions").(*dom.HTMLDivElement)
	cameraList := f.GetElementByID("camera-list").(*dom.HTMLUListElement)

	initializeCameraBtn.AddEventListener("click", false, func(e dom.Event) {
		initializeCamera.Style().Set("display", "none")
		lookingForCamera.Style().Set("display", "block")

		opts := js.M{
			"video":      preview,
			"scanPeriod": 5,
			"mirror":     false,
		}
		scanner := js.Global.Get("Instascan").Get("Scanner").New(opts)

		scanner.Call("addListener", "scan", func(content *js.Object) {
			scanResults.SetTextContent(content.String())
		})

		js.Global.Get("Instascan").Get("Camera").Call("getCameras").Call("then", func(cameras *js.Object) {
			lookingForCamera.Style().Set("display", "none")
			if cameras.Length() > 0 {
				cameraFound.Style().Set("display", "block")
				for i := 0; i < cameras.Length(); i++ {
					cam := cameras.Index(i)
					btn := scanButton(scanner, cam, cameraList)
					actions.InsertBefore(btn, scanStop)
				}
			} else {
				cameraNotFound.Style().Set("display", "block")
			}
		}).Call("catch", func(e *js.Object) {
			lookingForCamera.Style().Set("display", "none")
			cameraNotFound.Style().Set("display", "block")
			console.Log(e)
		})

		scanStop.AddEventListener("click", false, func(e dom.Event) {
			scanner.Call("stop").Call("then", func() {
				console.Log("camera stopped...")
			})
			e.PreventDefault()
		})
	})

	return result
}

func scanButton(scanner *js.Object, cam *js.Object, cameraList *dom.HTMLUListElement) *dom.HTMLButtonElement {
	btn := D.CreateElement("button").(*dom.HTMLButtonElement)
	btn.SetAttribute("type", "submit")
	btn.SetClass("btn btn-success")
	btn.SetTextContent("Scan")

	btn.AddEventListener("click", false, func(e dom.Event) {
		scanner.Call("start", cam).Call("then", func() {
			console.Log("camera started...")
		}).Call("catch", func(e *js.Object) {
			console.Log("camera error...")
			console.Log(e)
		})
		e.PreventDefault()
	})

	li := D.CreateElement("li")
	li.SetTextContent(cam.Get("name").String())
	cameraList.AppendChild(li)

	return btn
}

func renderQRCode(qrcode *js.Object, qrcodeImg dom.Element, text string) {
	qrcode.Call("empty")
	qrConfig := js.M{
		"text":   text,
		"render": "canvas",
		"width":  512,
		"height": 512,
	}
	qrcode.Call("qrcode", qrConfig)
	qrcodeCanvasData := qrcode.Call("children", "canvas").Index(0).Call("toDataURL").String()
	qrcodeImg.SetAttribute("src", qrcodeCanvasData)
}
