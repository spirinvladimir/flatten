package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFlatten(t *testing.T) {
	input := `[
		"((10 ; 20 ; 30) ; 40)",
		"((A ; 20 ; (B)) ; 40)",
		"((10 ; ((20 ; (30))) ; (40)))",
		"(♣ ; ♦ ; ♥)"
	]`
	expected := `[{"flat":"10 ; 20 ; 30 ; 40","depth":1},{"flat":"A ; 20 ; B ; 40","depth":2},{"flat":"10 ; 20 ; 30 ; 40","depth":4},{"flat":"â£ ; â¦ ; â¥","depth":0}]`
	req, _ := http.NewRequest(http.MethodPost, "/flatten", bytes.NewBufferString(input))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Flatten)
	handler.ServeHTTP(rr, req)
	got := strings.TrimRight(rr.Body.String(), "\n")
	if got != expected {
		t.Errorf("Unexpected response: got %v want %v", got, expected)
	}
}

func TestHistory(t *testing.T) {
	extraRequests := 3
	for i := 0; i < HistorySize + extraRequests; i++ {
		req, _ := http.NewRequest(http.MethodPost, "/flatten", bytes.NewBufferString(`["(1)"]`))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Flatten)
		handler.ServeHTTP(rr, req)
	}
	req, _ := http.NewRequest(http.MethodPost, "/history", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(History)
	handler.ServeHTTP(rr, req)
	var history []Request
	body, _ := ioutil.ReadAll(rr.Body)
	json.Unmarshal(body, &history)
	if len(history) != HistorySize {
		t.Errorf("History size %v should be %v", len(history), HistorySize)
	}
}
