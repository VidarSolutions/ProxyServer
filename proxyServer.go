package ProxyServer

import(

	"net/http"
	"net/http/httputil"

)

func ProxyServer(target string, res http.ResponseWriter, req *http.Request){
    proxy := &httputil.ReverseProxy{
        Director: func(req *http.Request) {
            req.URL.Path = req.URL.Path // set the URL path to the incoming request's path
            req.URL.Scheme = "http"
            req.URL.Host = target
            req.Host = "" // Remove Host header to prevent redirects
        },
    }

    proxy.ServeHTTP(res, req)

}



