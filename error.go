package main

import (
  "template"
  "log"
  "io"
)

var tmpl *template.Template

func init() {
  tmpl = template.Must(template.New("error").Parse("\n\n<!DOCTYPE html><html><head><title>Error</title></head><body style='font-style:sans-sarief;font-size:11pt;background:#eee'><div style='border:1px solid #999;background:white;margin: 50px auto;padding:1em 3em;width:400px'>{{. | html}}</div></body></html>\n" ))
}

func error(client io.WriteCloser, s string) {
	log.Printf("Error: %s", s)
	client.Write([]byte("HTTP/1.0 500 Internal Error\n"))
	client.Write([]byte("Content-Type: text/html\n"))
  tmpl.Execute(client, s)
	client.Close()
}
