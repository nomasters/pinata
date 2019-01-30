package pinata

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthentication(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
	}))
	defer ts.Close()
	client := NewClient("", "")
	client.BaseURL = ts.URL
	_, err := client.TestAuthentication()
	if err != nil {
		log.Fatal(err)
	}
}
