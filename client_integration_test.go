package glifxyz_test

import (
	"context"
	"os"
	"testing"

	glifxyz "github.com/iamwavecut/go-glifxyz-api-client"
)

const TestEnvAPITokenName = "GLIF_API_TOKEN_TEST"

func TestIntegration(t *testing.T) {
	os.Setenv(TestEnvAPITokenName, "439a4776d0d33b141ec33443a21b5a3d")
	if os.Getenv(TestEnvAPITokenName) == "" {
		t.Skipf("GLIF API test token is not set, skipping test")
	}
	client := glifxyz.NewGlifClient(glifxyz.WithEnvToken(TestEnvAPITokenName))
	ctx := context.Background()

	glifID := "clkmq6q1w000ele08f75it9xf" // nsfw hamster utility glif
	run, err := client.RunSimple(ctx, glifID)
	if err != nil {
		t.Fatalf("Error running model: %v", err)
	}

	t.Logf("Run response: %v", run)
}
