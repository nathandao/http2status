package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	. "github.com/nathandao/http2status/http2status"
)

var header = `
<html>
<head>
<meta charset='utf-8'>
<title>⚡</title>
</head>
<body>
<p>⚡⚡⚡⚡⚡⚡⚡⚡</p>
<p>⚡ HTML/2 checker ⚡</p>
<p>⚡⚡⚡⚡⚡⚡⚡⚡</p>
`

var footer = `
</body>
</html>
`

var form = `
<form method="POST" action="/" accept-charset="UTF-8">
<input type="text" name="url" placeholder="yoursite.com">
{{ .csrfField }}
<input type="submit" value="Check for HTTP2 status!">
</form>
`

var result = `
<p>{{ .err }}</p>
<h2>{{ .siteUrl }}</h2>
<h2>{{ .status }}</h2>
<p>{{ .response }}</p>
`

var tForm = template.Must(template.New("form.tmpl").Parse(form))
var tResult = template.Must(template.New("result.tmpl").Parse(result))
var tHeader = template.Must(template.New("header.tmpl").Parse(header))
var tFooter = template.Must(template.New("footer.tmpl").Parse(footer))

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)

	// Add csrf middleware.
	http.ListenAndServe(":8000",
		csrf.Protect([]byte("32-byte-long-auth-key"), csrf.Secure(false))(r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tHeader.Execute(w, nil)

	tForm.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})

	//url := r.URL.Query().Get("url")
	url := r.PostFormValue("url")

	if url != "" {

		obj := map[string]interface{}{}

		isH2, res, sanitizedUrl, err := Http2Status(url)
		if err != nil {
			obj["err"] = err
		} else {
			if !isH2 {
				obj["status"] = "Nope. It's not 1984 anymore, time to upgrade to http2."
			} else {
				obj["status"] = "You're on HTTP2!"
			}
		}

		obj["response"] = res
		obj["siteUrl"] = sanitizedUrl

		tResult.Execute(w, obj)
	}

	tFooter.Execute(w, nil)
}
