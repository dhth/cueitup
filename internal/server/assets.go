package server

import (
	_ "embed"
	"log"
	"net/http"
)

const (
	contentTypeTextHTML   = "text/html; charset=utf-8"
	contentTypeCSS        = "text/css; charset=utf-8"
	contentTypeJavascript = "text/javascript; charset=utf-8"
	contentTypeImage      = "image/png"
)

//go:embed _client/index.html
var indexHTMLBytes []byte

//go:embed _client/priv/static/cueitup.css
var cssBytes []byte

//go:embed _client/priv/static/custom.css
var customCSSBytes []byte

//go:embed _client/priv/static/favicon.png
var faviconBytes []byte

//go:embed _client/priv/static/cueitup.mjs
var jsBytes []byte

func getIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(contentType, contentTypeTextHTML)
	if _, err := w.Write(indexHTMLBytes); err != nil {
		log.Printf("failed to write bytes to HTTP connection: %s", err.Error())
	}
}

func getFavicon(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(contentType, contentTypeImage)
	if _, err := w.Write(faviconBytes); err != nil {
		log.Printf("failed to write bytes to HTTP connection: %s", err.Error())
	}
}

func getCSS(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(contentType, contentTypeCSS)
	if _, err := w.Write(cssBytes); err != nil {
		log.Printf("failed to write bytes to HTTP connection: %s", err.Error())
	}
}

func getCustomCSS(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(contentType, contentTypeCSS)
	if _, err := w.Write(customCSSBytes); err != nil {
		log.Printf("failed to write bytes to HTTP connection: %s", err.Error())
	}
}

func getJS(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(contentType, contentTypeJavascript)
	if _, err := w.Write(jsBytes); err != nil {
		log.Printf("failed to write bytes to HTTP connection: %s", err.Error())
	}
}
