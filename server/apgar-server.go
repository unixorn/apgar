// apgar-server
// Copyright 2016 DAQRI, LLC.
// Author: Joe Block <joe.block@daqri.com>
//
// Extremely minimalist http server to serve health check status. We don't
// need or want any fancy url rewriting, cgis, just to serve up a couple of
// text files.

package main

import (
	"flag"
	"net/http"
)

var document_root = flag.String("document-root", "/var/lib/apgar", "Document root")
var port = flag.String("port", "9000", "port to serve apgar results on")

func main() {
	flag.Parse()
	panic(http.ListenAndServe(":"+*port, http.FileServer(http.Dir(*document_root))))
}
