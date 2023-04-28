package ProxyServer

import (
	"bytes"
	"html"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

func ProxyServer(target string, res http.ResponseWriter, req *http.Request) {
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = target
			req.Host = "" // Remove Host header to prevent redirects
		},
		ModifyResponse: func(resp *http.Response) error {
			if resp.Header.Get("Content-Type") == "text/html" {
				doc, err := html.Parse(resp.Body)
				if err != nil {
					return err
				}
				resp.Body.Close()

				var traverse func(*html.Node)
				traverse = func(n *html.Node) {
					if n.Type == html.ElementNode {
						if n.Data == "a" {
							for i, attr := range n.Attr {
								if attr.Key == "href" {
									// Rewrite the link to point back to the onion site
									n.Attr[i].Val = strings.ReplaceAll(attr.Val, "http://"+target, req.URL.String())
								}
							}
						} else if n.Data == "img" {
							for i, attr := range n.Attr {
								if attr.Key == "src" {
									// Rewrite the image URL to point back to the onion site
									n.Attr[i].Val = strings.ReplaceAll(attr.Val, "http://"+target, req.URL.String())
								}
							}
						}
					}
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						traverse(c)
					}
				}
				traverse(doc)

				var buf bytes.Buffer
				if err = html.Render(&buf, doc); err != nil {
					return err
				}
				resp.Body = ioutil.NopCloser(&buf)
				resp.Header.Set("Content-Length", strconv.Itoa(buf.Len()))
			}
			return nil
		},
	}

	proxy.ServeHTTP(res, req)
}



