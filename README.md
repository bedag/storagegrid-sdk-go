# StorageGRID SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/bedag/storagegrid-sdk-go.svg)](https://pkg.go.dev/github.com/bedag/storagegrid-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/bedag/storagegrid-sdk-go)](https://goreportcard.com/report/github.com/bedag/storagegrid-sdk-go)

> **⚠️ Community-Maintained SDK**
> 
> This SDK was created by the community due to the lack of an official NetApp StorageGRID SDK for Go. It is designed to fulfill the needs of its maintainers and contributors. If you find something missing or spot a bug, please open an [issue](https://github.com/bedag/storagegrid-sdk-go/issues) or submit a [pull request](https://github.com/bedag/storagegrid-sdk-go/pulls)! Contributions are highly encouraged.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
  - [Grid Management](#grid-management)
  - [Tenant Management](#tenant-management)
- [Examples](#examples)
- [API Coverage](#api-coverage)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Overview

`storagegrid-sdk-go` is an unofficial Go SDK for interacting with NetApp StorageGRID. It provides programmatic access to StorageGRID's REST APIs, enabling you to manage tenants, buckets, users, S3 access keys, health, regions, gateway configs, and more.

NetApp StorageGRID exposes two distinct REST API surfaces:
- **Grid Management API**: For grid administrators to manage the entire StorageGRID system
- **Tenant Management API**: For tenant users to manage their specific tenant resources

This SDK reflects this architecture with corresponding client types. For more details on StorageGRID APIs, see the [official documentation](https://docs.netapp.com/us-en/storagegrid-115/s3/storagegrid-s3-rest-api-operations.html).

## Features

### Grid Management
- **Tenants**: Create, list, update, delete, and monitor tenant usage
- **Health**: Monitor grid health status (alarms, alerts, node connectivity)
- **Regions**: List available regions for grid and tenant contexts
- **HA Groups**: Manage High Availability groups
- **Gateway Configs**: Configure load balancer endpoints
- **S3 Object Lock**: Read and update grid-wide compliance (S3 Object Lock) settings

### Tenant Management
- **Buckets**: Create, list, delete, drain buckets; manage per-bucket sub-resources (region, usage, S3 Object Lock, notifications, policy, CORS, legacy compliance)
- **Users**: Manage tenant users with password management
- **Groups**: Manage tenant groups with policies and permissions
- **S3 Access Keys**: Generate and manage S3 access keys for users
- **Regions**: List tenant-specific regions

### Additional Features
- **Auto-authentication**: Automatic token management with expiration handling
- **Context support**: All operations support Go context for cancellation and timeouts
- **Interface-based design**: Easy mocking and testing with provided mock implementations
- **SSL configuration**: Optional SSL verification skip for development environments

## Requirements
- Go 1.25 or newer (see `go.mod` for the exact version)
- Access to a NetApp StorageGRID instance with appropriate credentials

## Installation

```sh
go get github.com/bedag/storagegrid-sdk-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bedag/storagegrid-sdk-go/client"
	"github.com/bedag/storagegrid-sdk-go/models"
)

func main() {
	ctx := context.Background()

	// Create a grid client for system administration
	gridClient, err := client.NewGridClient(
		client.WithEndpoint("https://your-storagegrid-endpoint"),
		client.WithCredentials(&models.Credentials{
			Username: "admin",
			Password: "your-password",
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Check grid health
	health, err := gridClient.Health().Get(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Grid Status: %s\n", func() string {
		if health.AllGreen() {
			return "✅ Healthy"
		}
		return "⚠️  Issues detected"
	}())
}
```

## Usage

### Client Configuration

Both client types support the same configuration options:

```go
import (
	"github.com/bedag/storagegrid-sdk-go/client"
	"github.com/bedag/storagegrid-sdk-go/models"
)

// Common configuration options
opts := []client.ClientOption{
	client.WithEndpoint("https://your-storagegrid.example.com"),
	client.WithCredentials(&models.Credentials{
		Username: "your-username",
		Password: "your-password",
		// AccountId: &accountID, // Required for tenant operations only
	}),
	// client.WithSkipSSL(), // Skip SSL verification (development only)
}
```

### Grid Management

Use `GridClient` for system-wide administration operations. This requires grid administrator privileges.

<details>
<summary>🔧 <strong>Setup Grid Client</strong></summary>

```go
gridClient, err := client.NewGridClient(
	client.WithEndpoint("https://your-storagegrid.example.com"),
	client.WithCredentials(&models.Credentials{
		Username: "grid-admin",
		Password: "admin-password",
		// No AccountId needed for grid operations
	}),
)
if err != nil {
	return fmt.Errorf("failed to create grid client: %w", err)
}
```
</details>

#### Managing Tenants

```go
// List all tenant accounts
tenants, err := gridClient.Tenant().List(ctx)
if err != nil {
	return fmt.Errorf("failed to list tenants: %w", err)
}

// Create a new tenant
tenant := &models.Tenant{
	Name:         "my-tenant",
	Description:  "Some description",
	Capabilities: []string{"management", "s3"},
	Policy: &models.TenantPolicy{
		UseAccountIdentitySource: false,
		AllowPlatformServices:    true,
		QuotaObjectBytes:         100 * 1024 * 1024 * 1024, // 100GB
	},
}

createdTenant, err := gridClient.Tenant().Create(ctx, tenant)
if err != nil {
	return fmt.Errorf("failed to create tenant: %w", err)
}

fmt.Printf("Created tenant: %s (ID: %s)\n", *createdTenant.Name, createdTenant.Id)
```

#### Monitoring Grid Health

```go
health, err := gridClient.Health().Get(ctx)
if err != nil {
	return fmt.Errorf("failed to get health status: %w", err)
}

// Check overall status
if health.AllGreen() {
	log.Println("✅ Grid is healthy")
} else {
	log.Printf("⚠️  Grid has issues - Connected nodes: %d, Alerts: %d", 
		*health.Nodes.Connected, 
		*health.Alerts.Critical + *health.Alerts.Major)
}
```

### Tenant Management

Use `TenantClient` for tenant-specific operations. This requires tenant user credentials and an account ID.

<details>
<summary>🔧 <strong>Setup Tenant Client</strong></summary>

```go
accountID := "12345678901234567890"
tenantClient, err := client.NewTenantClient(
	client.WithEndpoint("https://your-storagegrid.example.com"),
	client.WithCredentials(&models.Credentials{
		Username:  "tenant-admin",
		Password:  "tenant-password",
		AccountId: &accountID, // Required for tenant operations
	}),
)
if err != nil {
	return fmt.Errorf("failed to create tenant client: %w", err)
}
```
</details>

#### Managing Buckets

```go
// Create a bucket with versioning enabled
bucket := &models.Bucket{
	Name:             "my-application-data",
	Region:           "us-east-1",
	EnableVersioning: true,
	S3ObjectLock: &models.BucketS3ObjectLockSettings{
		Enabled: false,
	},
}

createdBucket, err := tenantClient.Bucket().Create(ctx, bucket)
if err != nil {
	return fmt.Errorf("failed to create bucket: %w", err)
}

// List all buckets in the tenant
buckets, err := tenantClient.Bucket().List(ctx)
if err != nil {
	return fmt.Errorf("failed to list buckets: %w", err)
}

for _, bucket := range *buckets {
	fmt.Printf("Bucket: %s (Created: %s)\n", bucket.Name, bucket.CreationTime.Format("2006-01-02"))
}
```

#### Managing Users and Access Keys

```go
// Create a new user
user := &models.User{
	UniqueName:  "application-user", // Will be prefixed with "user/"
	DisplayName: "Application Service User",
	Disable:     false,
}

createdUser, err := tenantClient.Users().Create(ctx, user)
if err != nil {
	return fmt.Errorf("failed to create user: %w", err)
}

// Generate S3 access keys for the user
accessKey := &models.S3AccessKey{
	Expires: nil, // No expiration
}

keys, err := tenantClient.S3AccessKeys().CreateForUser(ctx, *createdUser.Id, accessKey)
if err != nil {
	return fmt.Errorf("failed to create access keys: %w", err)
}

fmt.Printf("Access Key: %s\n", *keys.AccessKey)
fmt.Printf("Secret Key: %s\n", *keys.SecretAccessKey)
```

## Examples

## Examples

Comprehensive examples are available in the [`examples/`](examples/) directory:

- **[Grid Management](examples/grid/)**: Health monitoring, tenant management, end-to-end S3 Object Lock
- **[Tenant Operations](examples/tenant/)**: Bucket operations (incl. per-bucket sub-resources)
- **[Testing](examples/testing/)**: Unit tests with mocks, integration tests

### Quick Examples

#### Health Check
```go
health, err := gridClient.Health().Get(ctx)
if err != nil {
    log.Fatalf("Health check failed: %v", err)
}
fmt.Printf("Grid Status: All Green = %v\n", health.AllGreen())
```

#### Create Tenant
```go
tenant := &models.Tenant{
    Name:         "my-tenant",
    Capabilities: []string{"s3", "management"},
}
createdTenant, err := gridClient.Tenant().Create(ctx, tenant)
```

#### Create Bucket
```go
bucket := &models.Bucket{
    Name:   "my-bucket",
    Region: "us-east-1",
}
createdBucket, err := tenantClient.Bucket().Create(ctx, bucket)
```

For complete working examples, see the [`examples/`](examples/) directory.

## API Coverage

## Testing

### Unit Testing with Mocks

The SDK provides comprehensive mock implementations for all service interfaces, making unit testing straightforward:

```go
package main

import (
	"context"
	"testing"

	"github.com/bedag/storagegrid-sdk-go/models"
	"github.com/bedag/storagegrid-sdk-go/testing"
)

func TestTenantOperations(t *testing.T) {
	ctx := context.Background()

	// Create mock tenant service
	mockService := &testing.MockTenantService{
		ListFunc: func(ctx context.Context) (*[]models.Tenant, error) {
			return &[]models.Tenant{
				{
					Id:   "tenant-123",
					Name: "Test Tenant",
				},
			}, nil
		},
		CreateFunc: func(ctx context.Context, tenant *models.Tenant) (*models.Tenant, error) {
			tenant.Id = "new-tenant-456"
			return tenant, nil
		},
	}

	// Use mock in your application code
	tenants, err := mockService.List(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(*tenants) != 1 {
		t.Fatalf("Expected 1 tenant, got %d", len(*tenants))
	}

	if (*tenants)[0].Id != "tenant-123" {
		t.Fatalf("Expected tenant ID 'tenant-123', got %s", (*tenants)[0].Id)
	}
}
```

### Integration Testing

For integration tests against a real StorageGRID instance:

```go
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	endpoint := os.Getenv("STORAGEGRID_ENDPOINT")
	username := os.Getenv("STORAGEGRID_USERNAME")
	password := os.Getenv("STORAGEGRID_PASSWORD")

	if endpoint == "" || username == "" || password == "" {
		t.Skip("Missing required environment variables for integration test")
	}

	ctx := context.Background()
	
	client, err := client.NewGridClient(
		client.WithEndpoint(endpoint),
		client.WithCredentials(&models.Credentials{
			Username: username,
			Password: password,
		}),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test actual API calls
	health, err := client.Health().Get(ctx)
	if err != nil {
		t.Fatalf("Failed to get health: %v", err)
	}

	t.Logf("Grid health status: operative=%v", health.Operative(1))
}
```

### Available Mocks

The `testing` package provides mocks for all service interfaces:

- `MockTenantService` - Grid tenant management
- `MockBucketService` - Bucket operations  
- `MockTenantUserService` - Tenant user management
- `MockTenantGroupService` - Tenant group management
- `MockS3AccessKeyService` - S3 access key management
- `MockHealthService` - Health monitoring
- `MockHAGroupService` - HA group management
- `MockGatewayConfigService` - Gateway configuration
- `MockRegionService` - Region management
- `MockS3ObjectLockService` - Grid-wide S3 Object Lock (compliance-global) settings

## API Coverage

This SDK provides access to StorageGRID's dual API architecture:

### Grid Management APIs (GridClient)
Used for system-wide administration with grid administrator credentials:

| Service            | Endpoint                   | Operations                         | Description                                           |
| ------------------ | -------------------------- | ---------------------------------- | ----------------------------------------------------- |
| **Tenants**        | `/grid/accounts`           | Create, Read, Update, Delete, List | Manage tenant accounts                                |
| **Health**         | `/grid/health`             | Read                               | Monitor grid health, alarms, alerts, node status      |
| **Regions**        | `/grid/regions`            | List                               | Manage grid-wide regions                              |
| **HA Groups**      | `/private/ha-groups`       | Create, Read, Update, Delete, List | Configure High Availability groups                    |
| **Gateways**       | `/private/gateway-configs` | Create, Read, Update, Delete, List | Manage load balancer endpoints                        |
| **S3 Object Lock** | `/grid/compliance-global`  | Read, Update                       | Manage grid-wide S3 Object Lock (compliance) settings |

### Tenant Management APIs (TenantClient)
Used for tenant-specific operations with tenant user credentials:

| Service                  | Endpoint                                                                                | Operations                         | Description                                                                                                    |
| ------------------------ | --------------------------------------------------------------------------------------- | ---------------------------------- | -------------------------------------------------------------------------------------------------------------- |
| **Buckets**              | `/org/containers`                                                                       | Create, Read, Delete, List, Drain  | Manage S3 buckets within tenant                                                                                |
| **Bucket sub-resources** | `/org/containers/{name}/{region,usage,object-lock,notification,policy,cors,compliance}` | Read, Update                       | Per-bucket configuration: region, usage, S3 Object Lock, notifications, bucket policy, CORS, legacy compliance |
| **Users**                | `/org/users`                                                                            | Create, Read, Update, Delete, List | Manage tenant users                                                                                            |
| **Groups**               | `/org/groups`                                                                           | Create, Read, Update, Delete, List | Manage tenant groups and permissions                                                                           |
| **S3 Keys**              | `/org/users/*/s3-access-keys`                                                           | Create, Read, Delete, List         | Generate and manage S3 access credentials                                                                      |
| **Regions**              | `/org/regions`                                                                          | List                               | List tenant-accessible regions                                                                                 |
| **Usage**                | `/org/usage`                                                                            | Read                               | Monitor tenant usage statistics                                                                                |

> 📚 **Official Documentation**: For comprehensive API documentation, refer to the [NetApp StorageGRID REST API Reference](https://docs.netapp.com/us-en/storagegrid-115/s3/storagegrid-s3-rest-api-operations.html).

## Project Structure

```
storagegrid-sdk-go/
├── client/             # Client implementations
│   ├── client.go       # Base HTTP client with authentication
│   ├── grid.go         # Grid administrator client
│   └── tenant.go       # Tenant client
├── models/             # Data models for API requests/responses
│   ├── auth.go         # Authentication models
│   ├── buckets.go      # Bucket-related models
│   ├── tenants.go      # Tenant models
│   ├── users.go        # User models
│   ├── health.go       # Health status models
│   └── ...             # Other model files
├── services/           # Service interfaces and implementations
│   ├── interface.go    # Base HTTP client interface
│   ├── tenant.go       # Tenant management service
│   ├── buckets.go      # Bucket management service
│   ├── health.go       # Health monitoring service
│   └── ...             # Other service files
└── testing/            # Mock implementations for testing
    ├── tenant_mock.go  # Mock tenant service
    ├── bucket_mock.go  # Mock bucket service
    └── ...             # Other mock files
```

## Contributing

We welcome contributions! Here's how you can help:

1. **Report Issues**: Found a bug or missing feature? [Open an issue](https://github.com/bedag/storagegrid-sdk-go/issues)
2. **Submit Pull Requests**: Have a fix or new feature? [Submit a PR](https://github.com/bedag/storagegrid-sdk-go/pulls)
3. **Improve Documentation**: Help make this README and code comments better
4. **Add Tests**: Increase test coverage for reliability

### Development Guidelines

- Follow Go conventions and best practices
- Add appropriate error handling and logging
- Include tests for new functionality
- Update documentation for new features
- Maintain interface compatibility when possible

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.