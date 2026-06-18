package core

import (
	"fmt"
	"html"
	"net/http"
	"strconv"
)

type UI string

const (
	SwaggerUI UI = "swagger" // default
	ReDoc     UI = "redoc"
	Scalar    UI = "scalar"
	RapiDoc   UI = "rapidoc"
)

var uiRenderers = map[UI]func(url string) string{
	SwaggerUI: func(url string) string { return fmt.Sprintf(swaggerUIHTML, strconv.Quote(url)) },
	ReDoc:     func(url string) string { return fmt.Sprintf(redocHTML, html.EscapeString(url)) },
	Scalar:    func(url string) string { return fmt.Sprintf(scalarHTML, html.EscapeString(url)) },
	RapiDoc:   func(url string) string { return fmt.Sprintf(rapidocHTML, html.EscapeString(url)) },
}

func UIHandler(ui UI, url string) http.Handler {
	render, ok := uiRenderers[ui]
	if !ok {
		render = uiRenderers[SwaggerUI]
	}
	body := render(url)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(body))
	})
}

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Tonic API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.32.6/swagger-ui.css">
  <style>
    html { box-sizing: border-box; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin: 0; background: #fafafa; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.32.6/swagger-ui-bundle.js" crossorigin></script>
  <script src="https://unpkg.com/swagger-ui-dist@5.32.6/swagger-ui-standalone-preset.js" crossorigin></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: %s,
      dom_id: '#swagger-ui',
      deepLinking: true,
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIStandalonePreset
      ],
      layout: "StandaloneLayout"
    });
  </script>
</body>
</html>
`

const redocHTML = `<!DOCTYPE html>
<html lang="en">
  <head>
    <title>Tonic API Docs</title>
    <!-- needed for adaptive design -->
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <redoc spec-url='%s'></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>
  </body>
</html>
`

const scalarHTML = `<!DOCTYPE html>
<html lang="en">
  <head>
    <title>Tonic API Docs</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <div id="app"></div>
	<!-- Initialize the API Reference -->
    <script id="api-reference" data-url="%s"></script>
    <!-- Load the Script -->
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`

const rapidocHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Tonic API Docs</title>
  <script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
</head>
<body>
  <rapi-doc spec-url="%s" theme="light"></rapi-doc>
</body>
</html>
`
