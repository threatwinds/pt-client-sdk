# ThreatWinds Pentest Client SDK

Go SDK for interacting with the ThreatWinds Pentest Agent API.

## Installation

```bash
go get github.com/threatwinds/pt-client-sdk
```

## Features

- ✅ Independent HTTP client for CRUD operations
- ✅ Independent gRPC client for bidirectional streaming
- ✅ Authentication with API Key and Secret
- ✅ Fully separated architecture between HTTP and gRPC
- ✅ Schemas generated from protobuf (no duplication)

## Architecture

The SDK provides **two completely independent clients**:

### HTTPClient
Client dedicated exclusively to HTTP/REST operations:
- List pentests with pagination
- Get pentest details
- Schedule new pentests
- Download reports

### GRPCClient
Client dedicated exclusively to gRPC streaming:
- Provides direct access to gRPC client with authentication
- Users implement their own streaming logic
- Full control over data flow

### Why Separated?

This architecture allows you to:
- **Use only HTTP** if you don't need real-time streaming
- **Use only gRPC** if you only need streaming
- **Use both** independently according to your needs
- **Full control** over streaming implementation without predefined abstractions

## Usage

### HTTP Client

```go
package main

import (
    "context"
    "log"

    twpt "github.com/threatwinds/pt-client-sdk"
)

func main() {
    ctx := context.Background()

    // Create HTTP client
    httpClient := twpt.NewHTTPClient(
        "http://localhost:8000",
        twpt.Credentials{
            APIKey:    "your-api-key",
            APISecret: "your-api-secret",
        },
    )

    // List pentests
    pentests, err := httpClient.ListPentests(ctx, twpt.PaginationParams{
        Page:     1,
        PageSize: 10,
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, pt := range pentests.Pentests {
        log.Printf("Pentest ID: %s, Status: %s\n", pt.Id, pt.Status)
    }

    // Get a specific pentest
    pentest, err := httpClient.GetPentest(ctx, "pentest-id-123")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Pentest: %s\n", pentest.Id)

    // Schedule a new pentest
    req := &twpt.SchedulePentestRequest{
        Style:   twpt.Style_AGGRESSIVE,
        Exploit: true,
        Targets: []*twpt.TargetRequest{
            {
                Target: "example.com",
                Scope:  twpt.Scope_HOLISTIC,
                Type:   twpt.Type_BLACK_BOX,
            },
        },
    }

    pentestID, err := httpClient.SchedulePentest(ctx, req)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Pentest scheduled with ID: %s\n", pentestID)
}
```

### gRPC Client

```go
package main

import (
    "context"
    "io"
    "log"

    twpt "github.com/threatwinds/pt-client-sdk"
)

func main() {
    ctx := context.Background()

    // Create gRPC client
    grpcClient := twpt.NewGRPCClient(
        "localhost:50051",
        twpt.Credentials{
            APIKey:    "your-api-key",
            APISecret: "your-api-secret",
        },
    )
    defer grpcClient.Close()

    // Get gRPC client with authenticated context
    client, authCtx, err := grpcClient.GetClient(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // Start bidirectional stream
    stream, err := client.PentestStream(authCtx)
    if err != nil {
        log.Fatal(err)
    }

    // Subscribe to an existing pentest
    err = stream.Send(&twpt.ClientRequest{
        RequestType: &twpt.ClientRequest_GetPentest{
            GetPentest: &twpt.GetPentestRequest{
                PentestId: "pentest-id-123",
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    // Receive updates
    for {
        resp, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("Error: %v\n", err)
            break
        }

        // Handle different response types
        switch r := resp.ResponseType.(type) {
        case *twpt.ServerResponse_PentestData:
            log.Printf("Pentest data: %+v\n", r.PentestData)

        case *twpt.ServerResponse_StatusUpdate:
            log.Printf("Status: %s\n", r.StatusUpdate.Type)
            if r.StatusUpdate.Message != nil {
                log.Printf("  Message: %s\n", *r.StatusUpdate.Message)
            }

        case *twpt.ServerResponse_Error:
            log.Printf("Error: %s\n", r.Error.Error)
        }
    }

    stream.CloseSend()
}
```

### Schedule and Subscribe (gRPC)

```go
// Create gRPC client
grpcClient := twpt.NewGRPCClient(
    "localhost:50051",
    twpt.Credentials{
        APIKey:    "your-api-key",
        APISecret: "your-api-secret",
    },
)
defer grpcClient.Close()

// Get client with authentication
client, authCtx, err := grpcClient.GetClient(ctx)
if err != nil {
    log.Fatal(err)
}

// Start stream
stream, err := client.PentestStream(authCtx)
if err != nil {
    log.Fatal(err)
}

// Schedule pentest via gRPC
err = stream.Send(&twpt.ClientRequest{
    RequestType: &twpt.ClientRequest_SchedulePentest{
        SchedulePentest: &twpt.SchedulePentestRequest{
            Style:   twpt.Style_SAFE,
            Exploit: false,
            Targets: []*twpt.TargetRequest{
                {
                    Target: "example.com",
                    Scope:  twpt.Scope_TARGETED,
                    Type:   twpt.Type_WHITE_BOX,
                },
            },
        },
    },
})
if err != nil {
    log.Fatal(err)
}

// Receive real-time updates
for {
    resp, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Printf("Error: %v\n", err)
        break
    }

    switch r := resp.ResponseType.(type) {
    case *twpt.ServerResponse_ScheduleResponse:
        log.Printf("Pentest scheduled: %s\n", r.ScheduleResponse.PentestId)

    case *twpt.ServerResponse_StatusUpdate:
        log.Printf("Update: %s\n", r.StatusUpdate.Type)
        if r.StatusUpdate.Data != nil {
            log.Printf("  Status: %s\n", r.StatusUpdate.Data.Status)
        }
    }
}

stream.CloseSend()
```

### Using Both Clients

```go
// HTTP client for synchronous operations
httpClient := twpt.NewHTTPClient(
    "http://localhost:8000",
    twpt.Credentials{
        APIKey:    "your-api-key",
        APISecret: "your-api-secret",
    },
)

// gRPC client for streaming
grpcClient := twpt.NewGRPCClient(
    "localhost:50051",
    twpt.Credentials{
        APIKey:    "your-api-key",
        APISecret: "your-api-secret",
    },
)
defer grpcClient.Close()

// Schedule pentest via HTTP
pentestID, err := httpClient.SchedulePentest(ctx, &twpt.SchedulePentestRequest{
    Style:   twpt.Style_AGGRESSIVE,
    Exploit: true,
    Targets: []*twpt.TargetRequest{
        {Target: "example.com", Scope: twpt.Scope_HOLISTIC, Type: twpt.Type_BLACK_BOX},
    },
})

// Subscribe to updates via gRPC
client, authCtx, _ := grpcClient.GetClient(ctx)
stream, _ := client.PentestStream(authCtx)
stream.Send(&twpt.ClientRequest{
    RequestType: &twpt.ClientRequest_GetPentest{
        GetPentest: &twpt.GetPentestRequest{PentestId: pentestID},
    },
})

// Receive updates...
```

## Data Types

All core types are defined in protobuf and automatically generated:

### Enums

```go
// Status
twpt.Status_PENDING
twpt.Status_IN_PROGRESS
twpt.Status_COMPLETED
twpt.Status_FAILED

// Phase
twpt.Phase_RECON
twpt.Phase_INITIAL_EXPLOIT
twpt.Phase_DEEP_EXPLOIT
twpt.Phase_LATERAL_MOVEMENT
twpt.Phase_REPORT

// Scope
twpt.Scope_HOLISTIC
twpt.Scope_TARGETED

// Type
twpt.Type_BLACK_BOX
twpt.Type_WHITE_BOX

// Style
twpt.Style_AGGRESSIVE
twpt.Style_SAFE

// UpdateType
twpt.UpdateType_INFO
twpt.UpdateType_ERROR
twpt.UpdateType_STATUS
twpt.UpdateType_DEBUG
```

### Main Structures

**Protobuf (shared):**
- `PentestData`: Complete pentest data
- `TargetData`: Target data
- `ClientRequest`: Client to server request (gRPC)
- `ServerResponse`: Server to client response (gRPC)
- `SchedulePentestRequest`: Request to schedule a pentest
- `TargetRequest`: Request with target data
- `GetPentestRequest`: Request to get a pentest
- `StatusUpdate`: Status update

**HTTP specific:**
- `Credentials`: API key and secret
- `PaginationParams`: Pagination parameters
- `PentestListResponse`: Paginated list response
- `SchedulePentestResponse`: Schedule pentest response
- `ReportFormat`: Report format (PDF, JSON, MD)
- `DownloadReportRequest`: Download report request

## API Reference

### HTTPClient

```go
type HTTPClient struct {
    BaseURL     string
    HTTPClient  *http.Client
    Credentials *Credentials
}

func NewHTTPClient(baseURL string, creds Credentials) *HTTPClient

func (c *HTTPClient) ListPentests(ctx context.Context, pagination PaginationParams) (*PentestListResponse, error)
func (c *HTTPClient) GetPentest(ctx context.Context, pentestID string) (*PentestData, error)
func (c *HTTPClient) SchedulePentest(ctx context.Context, req *SchedulePentestRequest) (string, error)
func (c *HTTPClient) DownloadReport(ctx context.Context, pentestID string, format ReportFormat, outputDir string) error
```

### GRPCClient

```go
type GRPCClient struct {
    Address     string
    Credentials *Credentials
}

func NewGRPCClient(address string, creds Credentials) *GRPCClient

func (c *GRPCClient) GetClient(ctx context.Context) (PentestServiceClient, context.Context, error)
func (c *GRPCClient) Close() error
```

## Development

### Generate protobuf code

```bash
chmod +x generate.sh
./generate.sh
```

This will generate the `pentest.pb.go` and `pentest_grpc.pb.go` files from the `pentest.proto` file.

### Requirements

- Go 1.21+
- protoc (Protocol Buffers compiler)
- protoc-gen-go
- protoc-gen-go-grpc

### Install protobuf tools

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Internal Architecture

```
pt-client-sdk/
├── http_client.go         # Independent HTTP client
├── grpc_client.go         # Independent gRPC client
├── schemas.go             # HTTP auxiliary types
├── pentest.proto          # Protobuf definition
├── pentest.pb.go          # Generated types
├── pentest_grpc.pb.go     # Generated gRPC client
└── generate.sh            # Script to generate proto code
```

### Data Flow

**HTTP:**
```
HTTPClient → REST API → Response
```

**gRPC:**
```
GRPCClient.GetClient() → Your code → gRPC Stream ↔ Server
```

## Authentication

- **HTTPClient**: Automatically adds `api-key` and `api-secret` headers to each request
- **GRPCClient**: Automatically adds `api-key` and `api-secret` metadata to context (via `GetClient()`)

## License

MIT
