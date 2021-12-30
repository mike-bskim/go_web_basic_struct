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
	parms              []postData
	expectedStatusCode int
}{
	{"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"reservation-summary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	{"make-reservation", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "kimbs"},
		{key: "last_name", value: "bskim"},
		{key: "email", value: "kimbs@kimbs.com"},
		{key: "phone", value: "333-333-3333"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	// r, err := getSession()
	// if err != nil {
	// 	t.Error(err)
	// }
	// reservation := models.Reservation{
	// 	FirstName: "Hong",
	// 	LastName:  "Gildong",
	// 	Email:     "hong@gildong.com",
	// 	Phone:     "123-4567",
	// }
	// app.Session.Put(r.Context(), "reservation", reservation)

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("(GET) page:%s, expected %d, get %d\n", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.parms {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("(POST) name:%s, expected %d, but get %d\n", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func getSession() (*http.Request, error) {

	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}
