package ProxyServer

import(

	"net/http"
	"net/httputil"

)

func ProxyServer(target string, res http.ResponseWriter, req *http.Request){

	proxy := &ReverseProxy{
	Rewrite: func(r *ProxyRequest) {
		r.SetURL(target)
		r.Out.Host = r.In.Host // if desired
	}
}

 proxy.ServeHTTP(res, req)

}