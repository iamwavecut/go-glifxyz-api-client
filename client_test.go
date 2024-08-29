package glifxyz

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientSimple(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/run/cm023wc6m0009k7ur9ta0g14f" {
			t.Errorf("Expected to request '/api/v1/run/cm023wc6m0009k7ur9ta0g14f', got: %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json header, got: %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"cm023wc6m0009k7ur9ta0g14f","inputs":["a happy horse","foobar"],"output":"Test response"}`))
	}))
	defer server.Close()

	glifAPI := NewGlifClient(WithBaseURL(server.URL))

	simpleArgs := []string{"a happy horse", "foobar"}
	response, err := glifAPI.RunSimple(context.Background(), "cm023wc6m0009k7ur9ta0g14f", simpleArgs)
	if err != nil {
		t.Fatalf("Error running simple: %v", err)
	}
	if response.Output != "Test response" {
		t.Errorf("Expected 'Test response', got: %s", response.Output)
	}
}

func TestClientArged(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/run/cm023wc6m0009k7ur9ta0g14f" {
			t.Errorf("Expected to request '/api/v1/run/cm023wc6m0009k7ur9ta0g14f', got: %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json header, got: %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"cm023wc6m0009k7ur9ta0g14f","inputs":{"prompt":"a happy horse","other_parameter":"foobar"},"output":"Test response"}`))
	}))
	defer server.Close()

	glifAPI := NewGlifClient(WithBaseURL(server.URL))

	namedArgs := map[string]interface{}{
		"prompt":          "a happy horse",
		"other_parameter": "foobar",
	}
	response, err := glifAPI.RunSimple(context.Background(), "cm023wc6m0009k7ur9ta0g14f", namedArgs)
	if err != nil {
		t.Fatalf("Error running arged: %v", err)
	}
	if response.Output != "Test response" {
		t.Errorf("Expected 'Test response', got: %s", response.Output)
	}
}

func TestGetAddresses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/addresses" {
			t.Errorf("Expected to request '/api/v1/addresses', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AddressList{Addresses: []string{"address1", "address2"}})
	}))
	defer server.Close()

	glifAPI := NewGlifClient(WithBaseURL(server.URL))

	addresses, err := glifAPI.GetAddresses(context.Background())
	if err != nil {
		t.Fatalf("Error getting addresses: %v", err)
	}
	if addresses == nil || len(addresses.Addresses) != 2 {
		t.Errorf("Expected 2 addresses, got: %v", addresses)
	}
}
