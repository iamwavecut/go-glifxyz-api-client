package glifxyz_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	glifxyz "github.com/iamwavecut/go-glifxyz-api-client"
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

	glifAPI := glifxyz.NewGlifClient(glifxyz.WithBaseURL(server.URL))

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

	glifAPI := glifxyz.NewGlifClient(glifxyz.WithBaseURL(server.URL))

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
		json.NewEncoder(w).Encode(glifxyz.AddressList{Addresses: []string{"address1", "address2"}})
	}))
	defer server.Close()

	glifAPI := glifxyz.NewGlifClient(glifxyz.WithBaseURL(server.URL))

	addresses, err := glifAPI.GetAddresses(context.Background())
	if err != nil {
		t.Fatalf("Error getting addresses: %v", err)
	}
	if addresses == nil || len(addresses.Addresses) != 2 {
		t.Errorf("Expected 2 addresses, got: %v", addresses)
	}
}

func TestGetGlifs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/glifs" {
			t.Errorf("Expected to request '/api/glifs', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]glifxyz.GlifInfo{
			{ID: "glif1", Name: "Glif 1"},
			{ID: "glif2", Name: "Glif 2"},
		})
	}))
	defer server.Close()

	glifAPI := glifxyz.NewGlifClient(glifxyz.WithBaseURL(server.URL))

	glifs, err := glifAPI.GetGlifs(context.Background(), nil)
	if err != nil {
		t.Fatalf("Error getting glifs: %v", err)
	}
	if len(glifs) != 2 {
		t.Errorf("Expected 2 glifs, got: %d", len(glifs))
	}
}

func TestGetGlifRuns(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/runs" {
			t.Errorf("Expected to request '/api/runs', got: %s", r.URL.Path)
		}
		if r.URL.Query().Get("glifId") != "testGlifId" {
			t.Errorf("Expected glifId query parameter 'testGlifId', got: %s", r.URL.Query().Get("glifId"))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]glifxyz.GlifRun{
			{ID: "run1", GlifID: "testGlifId"},
			{ID: "run2", GlifID: "testGlifId"},
		})
	}))
	defer server.Close()

	glifAPI := glifxyz.NewGlifClient(glifxyz.WithBaseURL(server.URL))

	runs, err := glifAPI.GetGlifRuns(context.Background(), "testGlifId", nil)
	if err != nil {
		t.Fatalf("Error getting glif runs: %v", err)
	}
	if len(runs) != 2 {
		t.Errorf("Expected 2 runs, got: %d", len(runs))
	}
}

func TestGetUserInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/user" {
			t.Errorf("Expected to request '/api/user', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(glifxyz.UserInfo{ID: "user1", Username: "testuser"})
	}))
	defer server.Close()

	glifAPI := glifxyz.NewGlifClient(glifxyz.WithBaseURL(server.URL))

	userInfo, err := glifAPI.GetUserInfo(context.Background(), "testuser")
	if err != nil {
		t.Fatalf("Error getting user info: %v", err)
	}
	if userInfo.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got: %s", userInfo.Username)
	}
}

func TestGetMyInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/me" {
			t.Errorf("Expected to request '/api/me', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(glifxyz.UserInfo{ID: "myuser", Username: "myusername"})
	}))
	defer server.Close()

	glifAPI := glifxyz.NewGlifClient(glifxyz.WithBaseURL(server.URL))

	myInfo, err := glifAPI.GetMyInfo(context.Background())
	if err != nil {
		t.Fatalf("Error getting my info: %v", err)
	}
	if myInfo.Username != "myusername" {
		t.Errorf("Expected username 'myusername', got: %s", myInfo.Username)
	}
}

func TestGetSpheres(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/spheres" {
			t.Errorf("Expected to request '/api/spheres', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]glifxyz.SphereInfo{
			{ID: "sphere1", Name: "Sphere 1", Slug: "sphere-1"},
			{ID: "sphere2", Name: "Sphere 2", Slug: "sphere-2"},
		})
	}))
	defer server.Close()

	glifAPI := glifxyz.NewGlifClient(glifxyz.WithBaseURL(server.URL))

	spheres, err := glifAPI.GetSpheres(context.Background(), nil)
	if err != nil {
		t.Fatalf("Error getting spheres: %v", err)
	}
	if len(spheres) != 2 {
		t.Errorf("Expected 2 spheres, got: %d", len(spheres))
	}
}

func TestStreamRunSimple(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/run/testModelId" {
			t.Errorf("Expected to request '/api/v1/run/testModelId', got: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		for i := 0; i < 3; i++ {
			fmt.Fprintf(w, "{\"chunk\": %d}\n", i)
			w.(http.Flusher).Flush()
			time.Sleep(100 * time.Millisecond)
		}
	}))
	defer server.Close()

	glifAPI := glifxyz.NewGlifClient(glifxyz.WithBaseURL(server.URL))

	chunks := make([]int, 0)
	err := glifAPI.StreamRunSimple(context.Background(), "testModelId", map[string]interface{}{"prompt": "test"}, func(data []byte) error {
		var chunk struct {
			Chunk int `json:"chunk"`
		}
		if err := json.Unmarshal(data, &chunk); err != nil {
			return err
		}
		chunks = append(chunks, chunk.Chunk)
		return nil
	})

	if err != nil {
		t.Fatalf("Error streaming run: %v", err)
	}
	if len(chunks) != 3 {
		t.Errorf("Expected 3 chunks, got: %d", len(chunks))
	}
	for i, chunk := range chunks {
		if chunk != i {
			t.Errorf("Expected chunk %d, got: %d", i, chunk)
		}
	}
}
