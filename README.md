# Go [glif.app](https://glif.app/) API Client

[![Go Reference](https://pkg.go.dev/badge/github.com/iamwavecut/go-glifxyz-api-client.svg)](https://pkg.go.dev/github.com/iamwavecut/go-glifxyz-api-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/iamwavecut/go-glifxyz-api-client)](https://goreportcard.com/report/github.com/iamwavecut/go-glifxyz-api-client)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This is a Go client library for interacting with the Glif API. It provides a simple and efficient interface to run Glif models, retrieve addresses, and access various Glif resources.

## Features

- Run Glif models with simple and streaming options
- Retrieve Glif addresses
- Get information about Glifs, runs, users, and spheres
- Rate limiting support
- Customizable HTTP client
- Logging integration

## Installation

To install the Glif API client, use `go get`:
```shell
go get github.com/iamwavecut/go-glifxyz-api-client
```

## Usage

### Initializing the Client

To use the Glif API client, first import the package and create a new client instance:

```go
import (
    "github.com/iamwavecut/go-glifxyz-api-client"
)

client := glifxyz.NewGlifClient(
    glifxyz.WithAPIToken("your-api-token"),
    glifxyz.WithBaseURL("https://custom-glif-url.com"), // Optional: defaults to https://glif.app
)
```

### Running a Model

To run a model, use the `RunSimple` method:

```go
ctx := context.Background()
modelID := "your-model-id"
input := map[string]interface{}{
    "prompt": "Your input prompt",
}

run, err := client.RunSimple(ctx, modelID, input)
if err != nil {
    log.Fatalf("Error running model: %v", err)
}

fmt.Printf("Run ID: %s\n", run.ID)
fmt.Printf("Output: %s\n", run.Output)
```

### Streaming Model Output

To stream the output of a model run, use the `StreamRunSimple` method:

```go
err := client.StreamRunSimple(ctx, modelID, input, func(data []byte) error {
    fmt.Printf("Received data: %s\n", string(data))
    return nil
})
if err != nil {
    log.Fatalf("Error streaming model output: %v", err)
}
```

### Retrieving Addresses

To retrieve addresses, use the `GetAddresses` method:

```go
addresses, err := client.GetAddresses(ctx)
if err != nil {
    log.Fatalf("Error retrieving addresses: %v", err)
}

for _, address := range addresses.Addresses {
    fmt.Printf("Address: %s\n", address)
}
```

### Accessing Glif Information

To get information about Glifs, use the `GetGlifs` method:

```go
params := url.Values{}
params.Set("limit", "10")
glifs, err := client.GetGlifs(ctx, params)
if err != nil {
    log.Fatalf("Error retrieving Glifs: %v", err)
}

for _, glif := range glifs {
    fmt.Printf("Glif ID: %s, Name: %s\n", glif.ID, glif.Name)
}
```

### Accessing User Information

To access user information, use the `GetUserInfo` or `GetMyInfo` methods:

```go
userInfo, err := client.GetUserInfo(ctx, "username_or_id")
if err != nil {
    log.Fatalf("Error retrieving user info: %v", err)
}

fmt.Printf("User ID: %s, Username: %s\n", userInfo.ID, userInfo.Username)

// Or for the authenticated user:
myInfo, err := client.GetMyInfo(ctx)
if err != nil {
    log.Fatalf("Error retrieving my info: %v", err)
}

fmt.Printf("My ID: %s, My Username: %s\n", myInfo.ID, myInfo.Username)
```

### Accessing Sphere Information

To access sphere information, use the `GetSpheres` method:

```go
params := url.Values{}
params.Set("limit", "5")
spheres, err := client.GetSpheres(ctx, params)
if err != nil {
    log.Fatalf("Error retrieving spheres: %v", err)
}

for _, sphere := range spheres {
    fmt.Printf("Sphere ID: %s, Name: %s, Slug: %s\n", sphere.ID, sphere.Name, sphere.Slug)
}
```

## Configuration Options

The `NewGlifClient` function accepts various options to customize the client:

- `WithBaseURL(url string)`: Set a custom base URL for the Glif API
- `WithHTTPClient(client *http.Client)`: Use a custom HTTP client
- `WithLogger(logger *slog.Logger)`: Set a custom logger
- `WithRateLimit(r rate.Limit, b int)`: Set custom rate limiting parameters
- `WithAPIToken(token string)`: Set the API token for authentication

## Error Handling

All methods return errors that should be checked and handled appropriately. The client uses structured error types to provide more context about the nature of the error.

## Rate Limiting

The client implements rate limiting by default to prevent overwhelming the Glif API. You can customize the rate limit using the `WithRateLimit` option when creating the client.

## Contributing

Contributions to the Go Glif API Client are welcome! Please feel free to submit issues, fork the repository and send pull requests!

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
