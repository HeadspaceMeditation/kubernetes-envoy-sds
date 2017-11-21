// Copyright 2017 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type registrationHandler struct {
	kubeProxyEndpoint string
	namespace         string
}

func registrationServer(kubeProxyEndpoint string, namespace string) http.Handler {
	return &registrationHandler{kubeProxyEndpoint, namespace}
}

func (h *registrationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	service := serviceFromURL(r.URL.Path)

	sds, err := getService(h.kubeProxyEndpoint, h.namespace, service)
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(500)
		return
	}

	data, err := json.MarshalIndent(sds, "", "  ")
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(500)
		return
	}
	w.Write(data)

	elapsed := time.Since(start)
	log.Printf("%s %s", r.URL.Path, elapsed)
}

func serviceFromURL(path string) string {
	s := strings.Split(path, "/")
	if len(s) < 3 {
		return ""
	}
	return s[3]
}
