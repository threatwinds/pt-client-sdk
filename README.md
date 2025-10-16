# ThreatWinds Pentest Client SDK

Go SDK for interacting with the ThreatWinds Pentest Agent API.

## Installation

```bash
go get github.com/threatwinds/pt-client-sdk
```

## Features

- ✅ Independent HTTP client for CRUD operations with string enums
- ✅ Independent gRPC client for bidirectional streaming with int32 enums
- ✅ Authentication with API Key and Secret
- ✅ Fully separated architecture between HTTP and gRPC
- ✅ Type-safe API with protocol-specific types

## Architecture

The SDK provides **two completely independent clients** with **protocol-specific types**:

### HTTPClient
Client dedicated exclusively to HTTP/REST operations with **string-based enums** for better JSON readability:
- List pentests with pagination
- Get pentest details
- Schedule new pentests
- Download evidence

**Why string enums?** REST APIs benefit from human-readable JSON:
```json
{
  "style": "AGGRESSIVE",
  "scope": "HOLISTIC",
  "type": "BLACK_BOX"
}
```

### GRPCClient
Client dedicated exclusively to gRPC streaming with **int32 enums** for efficient binary protocol:
- Provides direct access to gRPC client with authentication
- Uses standard protobuf int32 enums
- Full control over bidirectional streaming
- Optimal performance for real-time updates

**Why int32 enums?** gRPC/protobuf uses integers for efficiency and forward compatibility.

### Type System

The SDK maintains **two separate type systems**:

1. **HTTP Types** (`schemas.go`): String-based enums
   - `HTTPScope`, `HTTPType`, `HTTPStyle`, `HTTPStatus`, etc.
   - `HTTPPentestData`, `HTTPTargetData`, etc.

2. **Protobuf Types** (`pentest.pb.go`): Int32-based enums
   - `Scope`, `Type`, `Style`, `Status`, etc.
   - `PentestData`, `TargetData`, etc.

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
        // Note: Status is a string like "COMPLETED"
        log.Printf("Pentest ID: %s, Status: %s\n", pt.ID, pt.Status)
    }

    // Get a specific pentest
    pentest, err := httpClient.GetPentest(ctx, "pentest-id-123")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Pentest: %s, Style: %s\n", pentest.ID, pentest.Style)

    // Schedule a new pentest with HTTP types
    req := &twpt.HTTPSchedulePentestRequest{
        Style:   twpt.HTTPStyleAggressive,  // String enum: "AGGRESSIVE"
        Exploit: true,
        Targets: []*twpt.HTTPTargetRequest{
            {
                Target: "example.com",
                Scope:  twpt.HTTPScopeHolistic,  // String enum: "HOLISTIC"
                Type:   twpt.HTTPTypeBlackBox,   // String enum: "BLACK_BOX"
                // Optional fields:
                // ID: &pentestID,
                // Credentials: &credsJSON,
            },
        },
    }

    pentestID, err := httpClient.SchedulePentest(ctx, req)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Pentest scheduled with ID: %s\n", pentestID)

    // Download evidence
    err = httpClient.DownloadEvidence(ctx, pentestID, "./evidence", true)
    if err != nil {
        log.Fatal(err)
    }
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
            // Note: Status is int32 (e.g., 3 for COMPLETED)
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

// Schedule pentest via gRPC with protobuf types
err = stream.Send(&twpt.ClientRequest{
    RequestType: &twpt.ClientRequest_SchedulePentest{
        SchedulePentest: &twpt.SchedulePentestRequest{
            Style:   twpt.Style_SAFE,  // Int32 enum: 2
            Exploit: false,
            Targets: []*twpt.TargetRequest{
                {
                    Target: "example.com",
                    Scope:  twpt.Scope_TARGETED,   // Int32 enum: 2
                    Type:   twpt.Type_WHITE_BOX,   // Int32 enum: 2
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
            log.Printf("  Status: %d\n", r.StatusUpdate.Data.Status)
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
pentestID, err := httpClient.SchedulePentest(ctx, &twpt.HTTPSchedulePentestRequest{
    Style:   twpt.HTTPStyleAggressive,
    Exploit: true,
    Targets: []*twpt.HTTPTargetRequest{
        {
            Target: "example.com",
            Scope:  twpt.HTTPScopeHolistic,
            Type:   twpt.HTTPTypeBlackBox,
            // Optional: ID and Credentials
        },
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

### HTTP Enums (String-based)

```go
// Scope
twpt.HTTPScopeHolistic  // "HOLISTIC"
twpt.HTTPScopeTargeted  // "TARGETED"

// Type
twpt.HTTPTypeBlackBox  // "BLACK_BOX"
twpt.HTTPTypeWhiteBox  // "WHITE_BOX"

// Style
twpt.HTTPStyleAggressive  // "AGGRESSIVE"
twpt.HTTPStyleSafe        // "SAFE"

// Status
twpt.HTTPStatusPending       // "PENDING"
twpt.HTTPStatusInProgress    // "IN_PROGRESS"
twpt.HTTPStatusCompleted     // "COMPLETED"
twpt.HTTPStatusFailed        // "FAILED"

// Phase
twpt.HTTPPhaseRecon           // "RECON"
twpt.HTTPPhaseInitialExploit  // "INITIAL_EXPLOIT"
twpt.HTTPPhaseDeepExploit     // "DEEP_EXPLOIT"
twpt.HTTPPhaseLateralMovement // "LATERAL_MOVEMENT"
twpt.HTTPPhaseReport          // "REPORT"
twpt.HTTPPhaseFinished        // "FINISHED"

// Severity
twpt.HTTPSeverityNone      // "NONE"
twpt.HTTPSeverityLow       // "LOW"
twpt.HTTPSeverityMedium    // "MEDIUM"
twpt.HTTPSeverityHigh      // "HIGH"
twpt.HTTPSeverityCritical  // "CRITICAL"
```

### gRPC Enums (Int32-based)

```go
// Status
twpt.Status_PENDING       // 1
twpt.Status_IN_PROGRESS   // 2
twpt.Status_COMPLETED     // 3
twpt.Status_FAILED        // 4

// Phase
twpt.Phase_RECON             // 1
twpt.Phase_INITIAL_EXPLOIT   // 2
twpt.Phase_DEEP_EXPLOIT      // 3
twpt.Phase_LATERAL_MOVEMENT  // 4
twpt.Phase_REPORT            // 5
twpt.Phase_FINISHED          // 6

// Scope
twpt.Scope_HOLISTIC  // 1
twpt.Scope_TARGETED  // 2

// Type
twpt.Type_BLACK_BOX  // 1
twpt.Type_WHITE_BOX  // 2

// Style
twpt.Style_AGGRESSIVE  // 1
twpt.Style_SAFE        // 2

// Severity
twpt.Severity_NONE      // 1
twpt.Severity_LOW       // 2
twpt.Severity_MEDIUM    // 3
twpt.Severity_HIGH      // 4
twpt.Severity_CRITICAL  // 5

// UpdateType
twpt.UpdateType_INFO    // 1
twpt.UpdateType_ERROR   // 2
twpt.UpdateType_STATUS  // 3
twpt.UpdateType_DEBUG   // 4
```

### Main Structures

**HTTP Types (schemas.go):**
- `HTTPPentestData`: Complete pentest data with string enums (ID, Status, CreatedAt, StartedAt, FinishedAt, Style, Exploit, Summary, Targets, Severity, Findings)
- `HTTPTargetData`: Target data with string enums (ID, PentestID, Target, Scope, Type, Status, Phase, CreatedAt, StartedAt, FinishedAt, Credentials, Severity, Findings, Summary)
- `HTTPSchedulePentestRequest`: Request to schedule a pentest (ID, Style, Exploit, Targets)
- `HTTPTargetRequest`: Target request data (ID, Target, Scope, Type, Credentials)
- `HTTPPentestListResponse`: Paginated list response (Pentests, Total, Page, PageSize, TotalPages)
- `HTTPSchedulePentestResponse`: Schedule response (PentestID)
- `Credentials`: API key and secret (APIKey, APISecret)
- `PaginationParams`: Pagination parameters (Page, PageSize)

**Protobuf Types (pentest.pb.go):**
- `PentestData`: Complete pentest data with int32 enums
- `TargetData`: Target data with int32 enums
- `ClientRequest`: Client to server request (gRPC)
- `ServerResponse`: Server to client response (gRPC)
- `SchedulePentestRequest`: Request to schedule a pentest
- `TargetRequest`: Request with target data
- `GetPentestRequest`: Request to get a pentest
- `StatusUpdate`: Status update
- `ErrorResponse`: Error response

## API Reference

### HTTPClient

```go
type HTTPClient struct {
    BaseURL     string
    HTTPClient  *http.Client
    Credentials *Credentials
}

func NewHTTPClient(baseURL string, creds Credentials) *HTTPClient

func (c *HTTPClient) ListPentests(ctx context.Context, pagination PaginationParams) (*HTTPPentestListResponse, error)
func (c *HTTPClient) GetPentest(ctx context.Context, pentestID string) (*HTTPPentestData, error)
func (c *HTTPClient) SchedulePentest(ctx context.Context, req *HTTPSchedulePentestRequest) (string, error)
func (c *HTTPClient) DownloadEvidence(ctx context.Context, pentestID string, outputPath string, unzip bool) error
func (c *HTTPClient) GetCurrentVersion(ctx context.Context) (string, error)
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
├── pentest.proto          # Protobuf definition (source of truth)
├── pentest.pb.go          # Generated types (int32 enums)
├── pentest_grpc.pb.go     # Generated gRPC client
├── schemas.go             # HTTP types (string enums)
├── http_client.go         # HTTP client (uses HTTP types)
├── grpc_client.go         # gRPC client (uses protobuf types)
├── auth.go                # Authentication validation
├── helpers/               # Utility functions
│   ├── download.go
│   └── zip.go
└── generate.sh            # Script to generate proto code
```

### Data Flow

**HTTP:**
```
HTTPClient → HTTP Types (string enums) → REST API → JSON Response
```

**gRPC:**
```
GRPCClient → Protobuf Types (int32 enums) → gRPC Stream ↔ Server
```

## Why This Architecture?

### ✅ Protocol-Appropriate Types
- **HTTP/REST**: Uses string enums for human-readable JSON
- **gRPC**: Uses int32 enums for efficiency and protobuf standards

### ✅ No Breaking Changes
- Each protocol has its own type system
- No forced serialization issues
- Forward compatible with protobuf evolution

### ✅ Type Safety
- Compile-time checks for enum usage
- No runtime conversion errors
- Clear separation of concerns

### ✅ Flexibility
- Direct access to protocol-specific types
- Mix and match HTTP and gRPC clients
- Independent operation of each client

## Authentication

- **HTTPClient**: Automatically adds `api-key` and `api-secret` headers to each request
- **GRPCClient**: Automatically adds `api-key` and `api-secret` metadata to context (via `GetClient()`)

## License

MIT
