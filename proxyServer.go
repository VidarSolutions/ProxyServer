package ProxyServer

import(

	"net/http"
	"net/http/httputil"

)

func ProxyServer(target string, res http.ResponseWriter, req *http.Request){

	 proxy := &httputil.ReverseProxy{
        Director: func(req *http.Request) {
            req.URL.Scheme = "http"
            req.URL.Host = target
            req.Host = target // set the Host header if desired
        },
    }
    proxy.ServeHTTP(res, req)

}



