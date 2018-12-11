package apps

import (
	"github.com/arunsworld/go-js-dom"
	"github.com/gopherjs/gopherjs/js"
)

var qrcodeContentHTML = `
	<div id="qrcode" width="512" height="512" style="display:none;"></div>
	<div class="row justify-content-sm-center">
		<div class="col-sm-4">
			<div class="card" style="width: 16rem;">
				<img id="qrcode-img" class="card-img-top" width="256" height="256"></img>
				<div class="card-body">
					<form>
						<div class="form-group">
							<input type="password" class="form-control" id="qrcode-text" placeholder="Password: transmit via QR">
						</div>
					</form>
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

	return result
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
