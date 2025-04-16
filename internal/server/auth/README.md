# Casbin Authentication for Kled.io Framework

This package provides a robust, domain-aware authentication and authorization system for the Kled.io Framework using Casbin.

## Features

### Domain-Based Access Control
- Multi-domain RBAC (Role-Based Access Control)
- Domain-specific policies and permissions
- Hierarchical role inheritance within domains
- Fine-grained access control for API endpoints

### Flexible Policy Management
- File-based policy storage with CSV format
- Runtime policy modification and enforcement
- Policy filtering and querying capabilities
- Support for allow/deny authorizations

### Integration with Gin Web Framework
- Middleware for authentication and authorization
- Domain extraction from request context
- Support for custom domain extractors
- Seamless integration with existing routes

## Usage

### Basic Authorization

```go
// Create a new enforcer
enforcer, err := auth.NewEnforcer("config/auth/model.conf", "config/auth/policy.csv")
if err != nil {
    log.Fatalf("Failed to create enforcer: %v", err)
}

// Use the authorization middleware
router.Use(middleware.AuthorizeMiddleware(enforcer))
```

### Domain-Based Authorization

```go
// Create a new domain enforcer
enforcer, err := auth.NewDomainEnforcer("config/auth/model.conf", "config/auth/policy.csv")
if err != nil {
    log.Fatalf("Failed to create domain enforcer: %v", err)
}

// Use the domain authorization middleware with default domain extractor
router.Use(middleware.DomainAuthorizeMiddleware(enforcer))

// Or use a custom domain extractor
customExtractor := func(c *gin.Context) string {
    return c.GetHeader("X-Domain")
}
router.Use(middleware.CustomDomainAuthorizeMiddleware(enforcer, customExtractor))
```

### Policy Management

```go
// Add a policy
enforcer.AddPolicy("user", "workspace", "/api/v1/workspaces/*", "GET")

// Check if a policy exists
exists := enforcer.HasPolicy("user", "workspace", "/api/v1/workspaces/*", "GET")

// Remove a policy
enforcer.RemovePolicy("user", "workspace", "/api/v1/workspaces/*", "GET")

// Save policies to disk
enforcer.SavePolicy()
```

### Role Management

```go
// Add a role for a user in a domain
enforcer.AddRoleForUserInDomain("alice", "admin", "workspace")

// Get roles for a user in a domain
roles := enforcer.GetRolesForUserInDomain("alice", "workspace")

// Delete a role for a user in a domain
enforcer.DeleteRoleForUserInDomain("alice", "admin", "workspace")
```

## Model Configuration

The default model configuration supports domain-based RBAC with the following components:

- Request definition: `r = sub, dom, obj, act`
- Policy definition: `p = sub, dom, obj, act`
- Role definition: `g = _, _, _`
- Policy effect: `e = some(where (p.eft == allow))`
- Matchers: `m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act`

This configuration allows for defining policies and roles specific to domains, enabling fine-grained access control across different parts of the application.
