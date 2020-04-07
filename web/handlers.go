package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

func homeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	unameCookie, err := r.Cookie("username")
	sidCookie, err2 := r.Cookie("session")

	if err != nil || err2 != nil ||
		len(unameCookie.Value) == 0 || len(sidCookie.Value) == 0 {
		tmpl, e := template.ParseFiles("./templates/home.html")
		if e != nil {
			log.Printf("Parsing template home.html error: %s", e)
			return
		}

		tmpl.Execute(w, nil)
		return
	}
	http.Redirect(w, r, "/userhome", http.StatusFound)

}

func userHomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	unameCookie, err1 := r.Cookie("username")
	_, err2 := r.Cookie("session")

	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	funame := r.FormValue("username")

	var page *UserPage
	if len(unameCookie.Value) != 0 {
		page = &UserPage{Name: unameCookie.Value}
	} else if len(funame) != 0 {
		page = &UserPage{Name: funame}
	}

	tmpl, err := template.ParseFiles("./templates/userhome.html")
	if err != nil {
		log.Printf("Parsing userhome.html error: %s", err)
		return
	}

	tmpl.Execute(w, page)
}

func apiHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, ErrorRequestNotRecognized)
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	apiBody := &ApiBody{}
	if err := json.Unmarshal(res, apiBody); err != nil {
		sendErrorResponse(w, ErrorRequestNotRecognized)
		return
	}

	request(apiBody, w, r)
}

func proxyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u, _ := url.Parse(STREAM_SERVER_HOST)
	// 把当前主机名替换为上述主机名后直接进行代理
	// 大数据流请求时使用(此时难以解析请求体)
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}
