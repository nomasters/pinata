package pinata

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
		t.Error(err)
	}
}

func TestMetadata(t *testing.T) {
	m := NewMetadata()
	if m.Name != "" {
		t.Error("Metadata.Name should be an empty string")
	}
	m = NewMetadataWithName("bob")
	if m.Name != "bob" {
		t.Error("Metadata.Name should equal 'bob'")
	}

	tests := []struct {
		Key         string
		Value       interface{}
		ShouldError bool
	}{
		{"test_string", "", false},
		{"test_int", int(5), false},
		{"test_int64", int64(5), false},
		{"test_float32", float32(5.001), false},
		{"test_float64", float64(5.001), false},
		{"test_time", time.Now(), false},
		{"test_bytes", []byte("no good"), true},
		{"test_bool", true, true},
	}

	for _, test := range tests {
		err := m.SetKeyValue(test.Key, test.Value)
		if test.ShouldError {
			if err == nil {
				t.Errorf("%v : %v should have returned an error, but did not", test.Key, test.Value)
			}
		} else {
			if err != nil {
				t.Error(err)
			}
		}
	}
}
