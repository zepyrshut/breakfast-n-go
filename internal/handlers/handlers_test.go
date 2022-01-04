package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"rooms", "/rooms", "GET", []postData{}, http.StatusOK},
	{"terrace", "/terrace", "GET", []postData{}, http.StatusOK},
	{"spa", "/spa", "GET", []postData{}, http.StatusOK},
	{"availability", "/availability", "GET", []postData{}, http.StatusOK},
	{"reservation", "/reservation", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"availability-post", "/availability", "POST", []postData{
		{key: "start", value: "2022-01-01"},
		{key: "end", value: "2022-02-01"},
	}, http.StatusOK},
	{"reservation", "/reservation", "POST", []postData{
		{key: "first_name", value: "benito"},
		{key: "second_name", value: "camela"},
		{key: "last_name", value: "benitez"},
		{key: "email", value: "benito@mail.com"},
		{key: "phone", value: "555-999"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
