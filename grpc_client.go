package pt_client_sdk

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// GRPCClient provides gRPC access to the ThreatWinds Pentest API
type GRPCClient struct {
	Address     string
	Credentials *Credentials
	conn        *grpc.ClientConn
	client      PentestServiceClient
	mu          sync.Mutex
}

// NewGRPCClient creates a new gRPC client instance
func NewGRPCClient(address string, creds Credentials) *GRPCClient {
	return &GRPCClient{
		Address:     address,
		Credentials: &creds,
	}
}

// GetClient returns the gRPC client with authenticated context.
func (c *GRPCClient) GetClient(ctx context.Context) (PentestServiceClient, context.Context, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client == nil {
		conn, err := grpc.NewClient(
			c.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
		}

		c.conn = conn
		c.client = NewPentestServiceClient(conn)
	}

	authCtx := c.getAuthContext(ctx)

	return c.client, authCtx, nil
}

// getAuthContext adds API credentials to the context
func (c *GRPCClient) getAuthContext(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{
		"api-key":    c.Credentials.APIKey,
		"api-secret": c.Credentials.APISecret,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

// Close closes the gRPC connection
func (c *GRPCClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
