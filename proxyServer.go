package ProxyServer

import (
	"bytes"
	//"html"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	//"net/url"
	"strconv"
	"strings"
	"golang.org/x/net/html"
	"github.com/vidarsolutions/Transfer"
)

func ProxyServer(target string, res http.ResponseWriter, req *http.Request) {
	var t =Transfer.Dialer("127.0.0.1:9050")
		if err != nil {
			fmt.Println("Error creating request: ", err)
			return
		}

	resp, _ = t.Request("GET", target, nil)
		if err != nil {
			fmt.Println("Error fetching response: ", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body: ", err)
			return
		}

		doc, err := html.Parse(strings.NewReader(string(body)))
		if err != nil {
			fmt.Println("Error parsing HTML: ", err)
			return
		}

		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "a" {
				for i, a := range n.Attr {
					if a.Key == "href" {
						u, err := url.Parse(a.Val)
						if err != nil {
							continue
						}
						if u.Host == target.Host {
							u.Host = req.Host
							u.Scheme = req.URL.Scheme
							n.Attr[i].Val = u.String()
						}
					}
				}
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}

		f(doc)

		html.Render(res, doc)
}



