package main

import (
  "template"
  "log"
  "io"
  "regexp"
)

var tmpl *template.Template
var ex *regexp.Regexp 

func n2br(text string) string {
  log.Print("Got text: " + text)
  return ex.ReplaceAllString(text, "<br>")
}

var fmap = template.FuncMap{
	"n2br": n2br,
}

func init() {
  layout := template.New("error template")
  layout.Funcs(fmap)
  tmpl = template.Must(layout.Parse(`
    <!DOCTYPE html>
    <html>
    <head><title>Error: {{.title | html}}</title></head>
    <body style='font-style:sans-serif;font-size:10pt;background:#eee'>
    <div style='border:1px solid #999;background:white;margin: 50px auto;padding:1em 3em;width:600px'>
    <h2>{{.title | html | n2br}}</h2>
    <pre style='background:#222;color:#eee;padding:8px 5px;border:1px solid #666'>{{.error | html | n2br}}</pre>
    </div>
    </body>
    </html>`))

  ex = regexp.MustCompile("\n")
}

func error(client io.WriteCloser, title, error string) {
	log.Printf("---\nError in %s: \n\n%s\n\n", title, error)
	client.Write([]byte("HTTP/1.0 500 Internal Error\n"))
	client.Write([]byte("Content-Type: text/html\n\n"))
  tmpl.Execute(client, map[string]string{"error": error, "title": title})
 	client.Write([]byte("\n"))
	client.Close()
}
