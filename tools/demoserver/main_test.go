package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// The recording is only as trustworthy as the stand-in behind it: every path the
// tape drives has to answer with JSON the CLI can render.
func TestDemoHandlerAnswersTheRecordedPaths(t *testing.T) {
	h := demoHandler()
	for _, path := range []string{"/control/status", "/control/stats", "/control/rewrite/list"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusOK {
			t.Errorf("GET %s = %d, want 200", path, rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("GET %s content type = %q", path, ct)
		}
		if !json.Valid(rec.Body.Bytes()) {
			t.Errorf("GET %s did not answer JSON: %s", path, rec.Body)
		}
	}
}

// A path nobody recorded is a 404, never an empty success that would look like
// a real answer on screen.
func TestDemoHandlerRefusesEverythingElse(t *testing.T) {
	h := demoHandler()
	for _, r := range []*http.Request{
		httptest.NewRequest(http.MethodGet, "/control/filtering/status", nil),
		httptest.NewRequest(http.MethodGet, "/control/rewrite/add", nil), // right path, wrong method
		httptest.NewRequest(http.MethodPost, "/control/status", strings.NewReader("{}")),
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s %s = %d, want 404", r.Method, r.URL.Path, rec.Code)
		}
	}
}

// Adding a rewrite has to show up in the next listing: a write that left no
// trace would make the recording show a result that never happened.
func TestDemoHandlerAddThenList(t *testing.T) {
	h := demoHandler()

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/control/rewrite/add",
		strings.NewReader(`{"domain":"lab.example.test","answer":"192.0.2.30"}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("POST rewrite/add = %d: %s", rec.Code, rec.Body)
	}
	var echoed rewrite
	if err := json.Unmarshal(rec.Body.Bytes(), &echoed); err != nil {
		t.Fatal(err)
	}
	if echoed.Domain != "lab.example.test" || echoed.Answer != "192.0.2.30" {
		t.Errorf("the write should echo what it was sent: %+v", echoed)
	}

	list := httptest.NewRecorder()
	h.ServeHTTP(list, httptest.NewRequest(http.MethodGet, "/control/rewrite/list", nil))
	var rewrites []rewrite
	if err := json.Unmarshal(list.Body.Bytes(), &rewrites); err != nil {
		t.Fatal(err)
	}
	if len(rewrites) != 3 || rewrites[2].Domain != "lab.example.test" {
		t.Errorf("the new rewrite should be last in the listing: %+v", rewrites)
	}
}

func TestDemoHandlerRejectsBrokenJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	demoHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/control/rewrite/add", strings.NewReader("not json")))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("a broken body = %d, want 400", rec.Code)
	}
}
