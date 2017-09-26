// apgar-server
//
// Author: Joe Block <joe.block@daqri.com>
//
// Extremely minimalist http server to serve health check status. We don't
// need or want any fancy url rewriting, cgis, just to serve up a directory
// of text files.
//
// The MIT License (MIT)
//
// Copyright (c) 2016 Daqri, LLC
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
//

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/status", healthCheck)
	panic(http.ListenAndServe(":9000", nil))
}

// Our healthcheck status is stored in /var/lib/apgar/status
// Validate that we are healthy, set proper http response if not
func healthCheck(w http.ResponseWriter, r *http.Request) {
	var healthy bool

	b, err := ioutil.ReadFile("/var/lib/apgar/status")

	// If error (such as file not found), report unhealthy
	if err != nil {
		fmt.Println(err)
		healthy = false
	} else {
		s := string(b)
		healthy = !strings.Contains(s, "UNHEALTHY")
	}

	w.Header().Set("Content-Type", "text/plain")

	if healthy {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HEALTHY\n"))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("UNHEALTHY\n"))
	}
}
