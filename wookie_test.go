package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"
)

//TODO Export tests from test.sh into this
// see server.go for all URL's that need testing
//
//TODO roll back all committed database queries from testing

//if you have a better idea for getting a cookie, please let me know
var cookie *http.Cookie = &http.Cookie{
	Name:    "logged-in",
	Value:   "MTM5MjY3Nzc2M3xKcVNVTnFTdmJ6MVFqVlg3XzdmX0xkNmstTXRJcERYckJwWnZmaENKbU55QjdRc3FXVmFnQ01MRmI0bjJhSmM0RUZWZC1jTEg4QWtWU010WWRpMlY2c0tofCL4IbY50NXz8rYVCFCEjaghmY6dcTxC7rq1aqQUYDkK",
	Path:    "/",
	Expires: time.Now().AddDate(20, 0, 0),
}

const URL = "http://localhost:8080"

type Reply struct {
	Success bool        `json:"success"`
	Info    interface{} `json:"info"`
	Message string      `json:"message"`
}

func getReply(resp *httptest.ResponseRecorder) (*Reply, error) {
	var r Reply
	err := json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (r *Reply) expectSuccess(t *testing.T) bool {
	if r.Success != true {
		t.Error("Expected success: true, got false, with message:", r.Message)
		return false
	}
	return true
}

func (r *Reply) expectFailure(t *testing.T) bool {
	if r.Success != false {
		t.Error("Expected success: false, got true")
		return false
	}
	return true
}

func TestLogin(t *testing.T) {
	resp := httptest.NewRecorder()

	uri := URL + "/logmein?"

	req, err := http.NewRequest("POST",
		uri+url.Values{"email": {"wookie@wookie.com"},
			"password": {"password"}}.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	login(resp, req)
	if resp.Header().Get("Location") != "/dashboard/#/main" {
		t.Error("Login was unsuccessful")
	}
}

func TestCookie(t *testing.T) {
	resp := httptest.NewRecorder()

	uri := URL + "/"
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(cookie)
	handlePage("index")(resp, req)

	r := regexp.MustCompile(`.*/dashboard/.*`)
	if !r.MatchString(resp.Body.String()) {
		t.Error("Cookie not set")
	}
}

func TestQuizList(t *testing.T) {
	resp := httptest.NewRecorder()

	uri := URL + "/quiz"

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(cookie)
	handleQuizList(resp, req)
	r, err := getReply(resp)
	if err != nil {
		t.Fatal(err)
	}
	r.expectSuccess(t)
	// TODO can check for specifics?
}

//TOD O test failures... e.g. no cookie

func TestNewClass(t *testing.T) {
	resp := httptest.NewRecorder()

	uri := URL + "/classes"

	req, err := http.NewRequest("POST", uri,
		strings.NewReader(`{"name":"myclass","students":{"123456789":"Me","12345":"You"}}`))
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/json")
	handleCreateClass(resp, req)
	r, err := getReply(resp)
	if err != nil {
		t.Fatal(err)
	}
	r.expectSuccess(t)
}

func 
