package apps

import (
	"github.com/arunsworld/go-js-dom"
	"github.com/arunsworld/gopherjs/dropzone"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/console"
)

var qrcodeContentHTML = `
	<div id="qrcode" width="512" height="512" style="display:none;"></div>
	<div class="row justify-content-sm-center">
		<div class="col-md-8">
			<div class="card-deck">
				<div class="card">
					<div class="card-body">
						<img id="qrcode-img" class="card-img-top" max-width="256" max-height="256"></img>
						<hr/>
						<form onSubmit="return false">
							<div class="form-group">
								<input type="password" class="form-control" id="qrcode-text" placeholder="Enter secret...">
							</div>
						</form>
					</div>
				</div>
				<div class="card">
					<div class="card-body">
						<div id="initialize-camera" class="text-center">
							<p>Scan QR Code</p>
							<button type="submit" class="btn btn-primary" id="initialize-camera-btn">Initialize</button>
						</div>
						<div id="looking-for-camera" style="display: none;">
							<p>Looking for camera...</p>
						</div>
						<div id="camera-found" style="display: none;" class="text-center">
							<p id="scan-results" class="text-center"></p>
							<p id="camera-message"></p>
							<form class="text-center">
								<div class="btn-group" id="actions">
									<button type="submit" class="btn btn-danger shadow-none" id="scanStop">Stop</button>
								</div>
							</form>
							<br/>
							<video id="preview" class="card-img-top" max-width="256" max-height="256" style="display:none;"></video>
						</div>
						<div id="camera-not-found" style="display: none;">
							<p>No camera was found in your system...</p>
						</div>
					</div>
				</div>
				<div class="card">
					<div class="card-body">
						<img id="qrcode-upload-img" class="card-img-top" max-width="256" max-height="256"></img>
						<hr/>
						<div id="uploadFile" class="dropzone">
							<div class="dz-message" style="font-size: 2rem; color: #967ADC; text-align: center;">Upload</div>
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

	qrcodeTxt.AddEventListener("change", false, func(e dom.Event) {
		renderQRCode(qrcode, qrcodeImg, qrcodeTxt.Value)
	})

	renderQRCode(qrcode, qrcodeImg, "https://go.iapps365.com/qr")

	qrcodeUploadImg := f.GetElementByID("qrcode-upload-img")
	uploadFile := f.GetElementByID("uploadFile")

	dropzone.NewDropzone(uploadFile, dropzone.Props{
		URL: "/api/upload",
		OnSuccess: func(file *dropzone.File, response string) {
			url := "https://go.iapps365.com/uploads/" + response
			renderQRCode(qrcode, qrcodeUploadImg, url)
		},
	})

	lookingForCamera := f.GetElementByID("looking-for-camera").(*dom.HTMLDivElement)
	cameraFound := f.GetElementByID("camera-found").(*dom.HTMLDivElement)
	cameraNotFound := f.GetElementByID("camera-not-found").(*dom.HTMLDivElement)
	preview := f.GetElementByID("preview").(*dom.HTMLVideoElement)
	scanResults := f.GetElementByID("scan-results").(*dom.HTMLParagraphElement)

	initializeCamera := f.GetElementByID("initialize-camera").(*dom.HTMLDivElement)
	initializeCameraBtn := f.GetElementByID("initialize-camera-btn").(*dom.HTMLButtonElement)
	scanStop := f.GetElementByID("scanStop").(*dom.HTMLButtonElement)
	scanStop.Disabled = true
	actions := f.GetElementByID("actions").(*dom.HTMLDivElement)
	cameraMessage := f.GetElementByID("camera-message")

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
			isURL := content.Call("startsWith", "https://")
			scannedContent := content.String()
			switch isURL.Bool() {
			case true:
				scanResults.SetInnerHTML("")
				link := D.CreateElement("a").(*dom.HTMLAnchorElement)
				link.Href = scannedContent
				link.Target = "_blank"
				link.SetTextContent(scannedContent)
				scanResults.AppendChild(link)
			case false:
				scanResults.SetTextContent(scannedContent)
			}
			scanner.Call("stop").Call("then", func() {
				console.Log("camera stopped...")
				preview.Style().Set("display", "none")
				scanStop.Disabled = true
			})
		})

		js.Global.Get("Instascan").Get("Camera").Call("getCameras").Call("then", func(cameras *js.Object) {
			lookingForCamera.Style().Set("display", "none")
			if cameras.Length() > 0 {
				cameraMessage.SetTextContent(cameras.Get("length").String() + " camera(s) found.")
				cameraFound.Style().Set("display", "block")
				for i := 0; i < cameras.Length(); i++ {
					cam := cameras.Index(i)
					btn := scanButton(scanner, cam, preview, scanStop)
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
				preview.Style().Set("display", "none")
				scanStop.Disabled = true
			})
			e.PreventDefault()
		})
	})

	return result
}

func scanButton(scanner *js.Object, cam *js.Object, preview *dom.HTMLVideoElement,
	scanStop *dom.HTMLButtonElement) *dom.HTMLButtonElement {
	btn := D.CreateElement("button").(*dom.HTMLButtonElement)
	btn.SetAttribute("type", "submit")
	btn.SetClass("btn btn-success shadow-none")
	btn.SetTextContent("Scan")

	btn.AddEventListener("click", false, func(e dom.Event) {
		scanner.Call("start", cam).Call("then", func() {
			preview.Style().Set("display", "block")
			scanStop.Disabled = false
			console.Log("camera started...")
		}).Call("catch", func(e *js.Object) {
			console.Log("camera error...")
			console.Log(e)
		})
		e.PreventDefault()
	})

	return btn
}

func renderQRCode(qrcode *js.Object, qrcodeImg dom.Element, text string) {
	if text == "" {
		text = "https://go.iapps365.com/qr"
	}
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
