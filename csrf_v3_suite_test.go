package csrfv3_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v3"
)

var csrfCookieName string = "csrf_"

type CookieJar struct {
	jar map[string]*http.Cookie
}

func (this CookieJar) New() *CookieJar {
	this.jar = map[string]*http.Cookie{}
	return &this
}

func (this *CookieJar) ReadCookies(resp *http.Response) {
	for _, c := range resp.Cookies() {
		this.jar[c.Name] = c
	}
}

func (this *CookieJar) WriteCookies(req *http.Request) {
	for _, c := range this.jar {
		req.AddCookie(c)
	}
}

func (this *CookieJar) GetValue(name string) string {
	return this.jar[name].Value
}

func (this *CookieJar) RemoveCookie(name string) {
	delete(this.jar, name)
}

type TestClient struct {
	cookieJar *CookieJar
	app       *fiber.App
}

func (this TestClient) New(app *fiber.App) *TestClient {
	this.cookieJar = CookieJar{}.New()
	this.app = app
	return &this
}

func (this *TestClient) GET(uri string, reqfilters ...func(*http.Request)) *http.Response {
	req := httptest.NewRequest("GET", "http://localhost"+uri, nil)
	this.cookieJar.WriteCookies(req)
	for _, f := range reqfilters {
		f(req)
	}
	resp, _ := this.app.Test(req)
	this.cookieJar.ReadCookies(resp)
	slog.Debug("GET "+uri, "Req", req, "Resp", resp)
	return resp
}

func (this *TestClient) POST(uri string, form url.Values, reqfilters ...func(*http.Request)) *http.Response {
	req := httptest.NewRequest("POST", "http://localhost"+uri, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Form = form
	req.PostForm = form
	this.cookieJar.WriteCookies(req)
	for _, f := range reqfilters {
		f(req)
	}
	resp, _ := this.app.Test(req)
	this.cookieJar.ReadCookies(resp)
	slog.Debug("POST "+uri, "Req", req, "Resp", resp)
	return resp
}

func (this *TestClient) GetCookieValue(name string) string {
	return this.cookieJar.GetValue(name)
}
