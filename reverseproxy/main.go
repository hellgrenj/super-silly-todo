package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/hellgrenj/super-silly-todo/reverseproxy/features"
)

func logSettings() {
	apiURL := os.Getenv("API_URL")
	microserviceURL := os.Getenv("MICROSERVICE_URL")
	log.Printf("apiURL set to %v", apiURL)
	log.Printf("microservice url set to %v", microserviceURL)
}
func requestIsAddItems(req *http.Request) bool {
	return req.Method == "POST" && strings.Contains(req.URL.Path, "/todolist/") && strings.Contains(req.URL.Path, "/item")

}
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {

	var targetURL string
	flags := features.GetFlags()
	if requestIsAddItems(req) && flags.DelegateAddListItemToMicroservice == true {
		targetURL = os.Getenv("MICROSERVICE_URL")
	} else {
		targetURL = os.Getenv("API_URL")
	}
	url, _ := url.Parse(targetURL)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

func main() {

	logSettings()
	port := ":3000"
	// start server
	fmt.Printf("\nReverse proxy up and running at port %v\n", port)
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}

}
