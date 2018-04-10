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
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Tomlmap struct {
	Webserver webserver
}

type webserver struct {
	Ipaddress string
	Port      string
}

func main() {
	raw := os.Getpid()
	myPid := []byte(strconv.Itoa(raw))
	var conf Tomlmap
	err := ioutil.WriteFile("/var/run/apgar-server.pid", myPid, 0644)
	if err != nil {
		fmt.Println("Could not write /var/run/apgar-server.pid:", err)
	}
	if _, err := toml.DecodeFile("/etc/apgar/config.toml", &conf); err != nil {
		fmt.Println("Could not open /etc/apgar/config.toml fallback to defaults", err)
	}
	webserverIP := conf.Webserver.Ipaddress
	webserverPort := conf.Webserver.Port
	if len(strings.TrimSpace(webserverIP)) == 0 {
		webserverIP = ""
	}
	if len(strings.TrimSpace(webserverPort)) == 0 {
		webserverPort = "9000"
	}
	listenAddress := fmt.Sprintf("%s:%s", webserverIP, webserverPort)
	http.HandleFunc("/status", healthCheck)
	http.HandleFunc("/", baseHandler)
	http.ListenAndServe(listenAddress, nil)
}

// Handy to allow our services to display scrapable data by writing to
// /var/lib/apgar/foo
func baseHandler(w http.ResponseWriter, r *http.Request) {
	fileName := fmt.Sprintf("/var/lib/apgar%s", r.URL)
	data, err := ioutil.ReadFile(fileName)

	w.Header().Set("Content-Type", "text/plain")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%s", err)))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
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
