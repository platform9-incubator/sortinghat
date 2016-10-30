package main

import (
	"log"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"flag"
	"regexp"
)

type ProxyConfig struct {
	UrlPattern string `json:"url_regex"`
	RemoteUrl  string `json:"remote_url"`
	handlerFunc http.HandlerFunc
}

type RevProxyConfig struct {
	ListenPort int `json:"listen_port"`
	ProxyList  []ProxyConfig `json:"proxy"`
}

/***
JSON Config file
{
	'listen_port':8080
	'proxy' : [ {
		'url_regex':'/static/*',
		'remote_url':'http://remote:port/foo'
	}]
}

 */
func parseConfig(configFile string) RevProxyConfig {

	fileData, e := ioutil.ReadFile(configFile)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var revProxyConfig RevProxyConfig;

	e = json.Unmarshal(fileData, &revProxyConfig)
	if e != nil {
		fmt.Printf("Error Unmarshaing: %v %s\n", e, fileData)
		os.Exit(1)
	}
	fmt.Println("RevProxyConfig ", revProxyConfig)
	return revProxyConfig
}

func Init() RevProxyConfig {

	var configFile string
	flag.StringVar(&configFile, "conf", "/Users/roopak/work/pf9-infra/whistle/revproxy.json", "Configuration file for reverse proxy")
	flag.Parse()
	return parseConfig(configFile)
}

func main() {
	revProxyConfig := Init()
	fmt.Println(revProxyConfig)
	configureHandlers(revProxyConfig)
	info(fmt.Sprintf("Listening %d", revProxyConfig.ListenPort))
	err := http.ListenAndServe(fmt.Sprintf(":%d", revProxyConfig.ListenPort), nil)
	if err != nil {
		panic(err)
	}
}


func createHttpHandler(remote url.URL) func(http.ResponseWriter, *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(&remote)
	return func(w http.ResponseWriter, r *http.Request) {
		info("To Http handler: ")
		proxy.ServeHTTP(w, r)
	}
}

func createFileHandler(remote url.URL) func(http.ResponseWriter, *http.Request) {
	info("Creating file server for "+ remote.Path)
	fileserver := http.FileServer(http.Dir(remote.Path))
	return func(w http.ResponseWriter, r *http.Request) {
		info("To file Handler "+remote.Path)
		fileserver.ServeHTTP(w, r)
	}
}

func createWSHandler(target string) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		info("To WS Handler "+target)
		dstConnection, err := net.Dial("tcp", target)
		if err != nil {
			http.Error(responseWriter, "Error contacting backend server.", 500)
			log.Printf("Error dialing websocket backend %s: %v", target, err)
			return
		}
		hj, ok := responseWriter.(http.Hijacker)
		if !ok {
			http.Error(responseWriter, "Not a hijacker?", 500)
			return
		}
		srcConnection, _, err := hj.Hijack()
		if err != nil {
			log.Printf("Hijack error: %v", err)
			return
		}
		defer srcConnection.Close()
		defer dstConnection.Close()

		err = httpRequest.Write(dstConnection)
		if err != nil {
			log.Printf("Error copying request to target: %v", err)
			return
		}

		errc := make(chan error, 2)
		cp := func(dst io.Writer, src io.Reader) {
			_, err := io.Copy(dst, src)
			errc <- err
		}
		go cp(dstConnection, srcConnection)
		go cp(srcConnection, dstConnection)
		<-errc
	}
}

func info(log string) {
	fmt.Println("INFO", log)
}

func globalHandler(revProxyConfig *RevProxyConfig) func(http.ResponseWriter, *http.Request) {

	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		for _, proxyConfig := range revProxyConfig.ProxyList {
			match, err:= regexp.MatchString(proxyConfig.UrlPattern, httpRequest.URL.Path)
			if match {
				proxyConfig.handlerFunc(responseWriter, httpRequest)
				return
			}

			if err != nil {
				info(err.Error())
			}
		}
	}
}

func configureHandlers(revProxyConfig RevProxyConfig){
    for idx, proxyConfig := range revProxyConfig.ProxyList {
		remote, err := url.Parse(proxyConfig.RemoteUrl)
		if err != nil {
			panic(err)
		}
		var handlerFunc http.HandlerFunc

		if remote.Scheme == "file" {
			info("Creating file handler")
			handlerFunc = createFileHandler(*remote)
		}
		if remote.Scheme == "ws" {
			info("Creating ws handler")
			handlerFunc = createWSHandler(remote.Host)
		}
		if remote.Scheme == "http" {
			info("Creating http handler")
			handlerFunc = createHttpHandler(*remote)
		}
		proxyConfig.handlerFunc = handlerFunc
		revProxyConfig.ProxyList[idx].handlerFunc = handlerFunc
	}
	fmt.Println("-----")
	fmt.Println("RevProxyConfig ", revProxyConfig)
	http.HandleFunc("/", globalHandler(&revProxyConfig))
}
