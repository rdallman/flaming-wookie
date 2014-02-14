package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestLogin(t *testing.T) {

	fmt.Println(hashKey)
	c := gettestcookie()
	fmt.Println(c.Value)

	ts := httptest.NewServer(http.HandlerFunc(changeState))
	defer ts.Close()
	//req, err := http.NewRequest("PUT", ts.URL, )
}

func gettestcookie() *http.Cookie {

	ts := httptest.NewServer(http.HandlerFunc(login))

	v := url.Values{}
	v.Set("email", "me1@example.com")
	v.Set("password", "password")

	resp, _ := http.PostForm(ts.URL, v)
	return resp.Cookies()[0]

}
