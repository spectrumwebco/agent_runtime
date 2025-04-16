# Casbin Integration Plan for Kled.io Framework

This document outlines the plan for deeply integrating Casbin's domain-aware RBAC (Role-Based Access Control) into the agent_runtime codebase to enhance the Kled.io Framework's authentication and authorization capabilities.

## 1. Current Implementation Assessment

The agent_runtime repository already has a basic Casbin integration:

- **Basic RBAC Model**: Simple role-based access control without domain support
- **File-based Policy Storage**: Policies stored in CSV format
- **Basic Gin Middleware**: Simple middleware for request authorization
- **Limited Role Management**: Basic role assignment without hierarchy

### Existing Files:
- `/config/auth/policy.csv`: Contains basic RBAC policies
- `/config/auth/rbac_model.conf`: Contains basic RBAC model
- `/internal/server/auth/casbin.go`: Basic Casbin enforcer implementation
- `/internal/server/middleware/auth.go`: Basic authorization middleware

## 2. Enhancement Goals

The enhanced Casbin integration will provide:

- **Domain-aware RBAC**: Multi-tenancy support with domain-specific roles and policies
- **Advanced Policy Models**: Support for URL pattern matching and regex
- **Dynamic Policy Management**: API endpoints for policy management
- **Role Hierarchy**: Support for role inheritance within domains
- **Integration with Authentication**: Connection with JWT or OAuth providers

## 3. Implementation Plan

### 3.1 Core Components

#### Domain Enforcer
Create a domain-aware Casbin enforcer that supports multi-tenancy:

```go
// internal/server/auth/domain_enforcer.go
package auth

import (
    "fmt"
    "github.com/casbin/casbin/v2"
    "github.com/casbin/casbin/v2/model"
    fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

// DomainEnforcer extends the basic Enforcer with domain support
type DomainEnforcer struct {
    enforcer *casbin.Enforcer
}

// NewDomainEnforcer creates a new domain-aware Casbin enforcer
func NewDomainEnforcer(modelPath, policyPath string) (*DomainEnforcer, error) {
    m, err := model.NewModelFromFile(modelPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load model: %w", err)
    }

    adapter := fileadapter.NewAdapter(policyPath)
    enforcer, err := casbin.NewEnforcer(m, adapter)
    if err != nil {
        return nil, fmt.Errorf("failed to create enforcer: %w", err)
    }

    return &DomainEnforcer{
        enforcer: enforcer,
    }, nil
}

// Enforce checks if a subject has permission to access an object in a domain
func (e *DomainEnforcer) Enforce(sub, dom, obj, act string) (bool, error) {
    return e.enforcer.Enforce(sub, dom, obj, act)
}

// AddPolicy adds a policy rule to the enforcer
func (e *DomainEnforcer) AddPolicy(sub, dom, obj, act string) (bool, error) {
    return e.enforcer.AddPolicy(sub, dom, obj, act)
}

// RemovePolicy removes a policy rule from the enforcer
func (e *DomainEnforcer) RemovePolicy(sub, dom, obj, act string) (bool, error) {
    return e.enforcer.RemovePolicy(sub, dom, obj, act)
}

// AddRoleForUserInDomain adds a role for a user in a domain
func (e *DomainEnforcer) AddRoleForUserInDomain(user, role, domain string) (bool, error) {
    return e.enforcer.AddNamedGroupingPolicy("g", user, role, domain)
}

// DeleteRoleForUserInDomain deletes a role for a user in a domain
func (e *DomainEnforcer) DeleteRoleForUserInDomain(user, role, domain string) (bool, error) {
    return e.enforcer.RemoveNamedGroupingPolicy("g", user, role, domain)
}

// GetRolesForUserInDomain gets the roles that a user has in a domain
func (e *DomainEnforcer) GetRolesForUserInDomain(user, domain string) ([]string, error) {
    return e.enforcer.GetModel()["g"]["g"].RM.GetRoles(user, domain)
}

// SavePolicy saves the current policy to the adapter
func (e *DomainEnforcer) SavePolicy() error {
    return e.enforcer.SavePolicy()
}
```

#### Domain-aware Middleware
Create middleware that extracts the domain from the request and uses it for authorization:

```go
// internal/server/middleware/domain_auth.go
package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/spectrumwebco/agent_runtime/internal/server/auth"
)

// DomainAuthorizeMiddleware creates middleware that authorizes requests based on domain
func DomainAuthorizeMiddleware(enforcer *auth.DomainEnforcer) gin.HandlerFunc {
    return CustomDomainAuthorizeMiddleware(enforcer, DefaultDomainExtractor)
}

// CustomDomainAuthorizeMiddleware creates middleware with a custom domain extractor
func CustomDomainAuthorizeMiddleware(enforcer *auth.DomainEnforcer, extractor DomainExtractor) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := c.GetString("user")
        if user == "" {
            user = "anonymous"
        }

        domain := extractor(c)
        path := c.Request.URL.Path
        method := c.Request.Method

        // Check if the user is authorized
        allowed, err := enforcer.Enforce(user, domain, path, method)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                "error": "Authorization error",
            })
            return
        }

        if !allowed {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": "Forbidden",
            })
            return
        }

        c.Next()
    }
}
```

#### Domain-aware Model Configuration
Create a domain-aware RBAC model configuration:

```
# config/auth/model.conf
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
```

### 3.2 API Design

#### Policy Management API
Create API endpoints for policy management:

```go
// internal/server/routes/policy.go
package routes

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/spectrumwebco/agent_runtime/internal/server/auth"
)

// RegisterPolicyRoutes registers policy management routes
func RegisterPolicyRoutes(router *gin.RouterGroup, enforcer *auth.DomainEnforcer) {
    policyRouter := router.Group("/policies")
    {
        policyRouter.GET("", listPolicies(enforcer))
        policyRouter.POST("", addPolicy(enforcer))
        policyRouter.DELETE("", removePolicy(enforcer))
    }

    roleRouter := router.Group("/roles")
    {
        roleRouter.GET("/:user", getUserRoles(enforcer))
        roleRouter.POST("/:user", addUserRole(enforcer))
        roleRouter.DELETE("/:user", removeUserRole(enforcer))
    }
}

// listPolicies returns all policies
func listPolicies(enforcer *auth.DomainEnforcer) gin.HandlerFunc {
    return func(c *gin.Context) {
        domain := c.Query("domain")
        if domain == "" {
            domain = "system"
        }

        policies := enforcer.GetFilteredPolicy(1, domain)
        c.JSON(http.StatusOK, gin.H{
            "policies": policies,
        })
    }
}

// addPolicy adds a new policy
func addPolicy(enforcer *auth.DomainEnforcer) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            Subject string `json:"subject" binding:"required"`
            Domain  string `json:"domain" binding:"required"`
            Object  string `json:"object" binding:"required"`
            Action  string `json:"action" binding:"required"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        success, err := enforcer.AddPolicy(req.Subject, req.Domain, req.Object, req.Action)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        if success {
            enforcer.SavePolicy()
            c.JSON(http.StatusOK, gin.H{"message": "Policy added successfully"})
        } else {
            c.JSON(http.StatusConflict, gin.H{"message": "Policy already exists"})
        }
    }
}

// removePolicy removes a policy
func removePolicy(enforcer *auth.DomainEnforcer) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            Subject string `json:"subject" binding:"required"`
            Domain  string `json:"domain" binding:"required"`
            Object  string `json:"object" binding:"required"`
            Action  string `json:"action" binding:"required"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        success, err := enforcer.RemovePolicy(req.Subject, req.Domain, req.Object, req.Action)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        if success {
            enforcer.SavePolicy()
            c.JSON(http.StatusOK, gin.H{"message": "Policy removed successfully"})
        } else {
            c.JSON(http.StatusNotFound, gin.H{"message": "Policy not found"})
        }
    }
}

// getUserRoles returns all roles for a user
func getUserRoles(enforcer *auth.DomainEnforcer) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := c.Param("user")
        domain := c.Query("domain")
        if domain == "" {
            domain = "system"
        }

        roles, err := enforcer.GetRolesForUserInDomain(user, domain)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "user":  user,
            "domain": domain,
            "roles": roles,
        })
    }
}

// addUserRole adds a role for a user
func addUserRole(enforcer *auth.DomainEnforcer) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := c.Param("user")
        var req struct {
            Role   string `json:"role" binding:"required"`
            Domain string `json:"domain" binding:"required"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        success, err := enforcer.AddRoleForUserInDomain(user, req.Role, req.Domain)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        if success {
            enforcer.SavePolicy()
            c.JSON(http.StatusOK, gin.H{"message": "Role added successfully"})
        } else {
            c.JSON(http.StatusConflict, gin.H{"message": "Role already exists"})
        }
    }
}

// removeUserRole removes a role from a user
func removeUserRole(enforcer *auth.DomainEnforcer) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := c.Param("user")
        var req struct {
            Role   string `json:"role" binding:"required"`
            Domain string `json:"domain" binding:"required"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        success, err := enforcer.DeleteRoleForUserInDomain(user, req.Role, req.Domain)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        if success {
            enforcer.SavePolicy()
            c.JSON(http.StatusOK, gin.H{"message": "Role removed successfully"})
        } else {
            c.JSON(http.StatusNotFound, gin.H{"message": "Role not found"})
        }
    }
}
```

#### CLI Commands for Policy Management
Create CLI commands for policy management:

```go
// cmd/cli/commands/policy.go
package commands

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/spectrumwebco/agent_runtime/internal/server/auth"
)

// NewPolicyCommand creates a new policy command
func NewPolicyCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "policy",
        Short: "Manage authorization policies",
        Long:  `Manage authorization policies for the Kled.io Framework.`,
    }

    cmd.AddCommand(newListPoliciesCommand())
    cmd.AddCommand(newAddPolicyCommand())
    cmd.AddCommand(newRemovePolicyCommand())
    cmd.AddCommand(newListRolesCommand())
    cmd.AddCommand(newAddRoleCommand())
    cmd.AddCommand(newRemoveRoleCommand())

    return cmd
}

// newListPoliciesCommand creates a command to list policies
func newListPoliciesCommand() *cobra.Command {
    var domain string

    cmd := &cobra.Command{
        Use:   "list",
        Short: "List policies",
        Long:  `List authorization policies for a domain.`,
        RunE: func(cmd *cobra.Command, args []string) error {
            enforcer, err := auth.NewDomainEnforcer("config/auth/model.conf", "config/auth/policy.csv")
            if err != nil {
                return fmt.Errorf("failed to create enforcer: %w", err)
            }

            policies := enforcer.GetFilteredPolicy(1, domain)
            fmt.Printf("Policies for domain %s:\n", domain)
            for _, policy := range policies {
                fmt.Printf("  %s, %s, %s, %s\n", policy[0], policy[1], policy[2], policy[3])
            }

            return nil
        },
    }

    cmd.Flags().StringVarP(&domain, "domain", "d", "system", "Domain to list policies for")

    return cmd
}

// newAddPolicyCommand creates a command to add a policy
func newAddPolicyCommand() *cobra.Command {
    var subject, domain, object, action string

    cmd := &cobra.Command{
        Use:   "add",
        Short: "Add a policy",
        Long:  `Add an authorization policy for a domain.`,
        RunE: func(cmd *cobra.Command, args []string) error {
            enforcer, err := auth.NewDomainEnforcer("config/auth/model.conf", "config/auth/policy.csv")
            if err != nil {
                return fmt.Errorf("failed to create enforcer: %w", err)
            }

            success, err := enforcer.AddPolicy(subject, domain, object, action)
            if err != nil {
                return fmt.Errorf("failed to add policy: %w", err)
            }

            if success {
                enforcer.SavePolicy()
                fmt.Println("Policy added successfully")
            } else {
                fmt.Println("Policy already exists")
            }

            return nil
        },
    }

    cmd.Flags().StringVarP(&subject, "subject", "s", "", "Subject (user or role)")
    cmd.Flags().StringVarP(&domain, "domain", "d", "system", "Domain")
    cmd.Flags().StringVarP(&object, "object", "o", "", "Object (resource)")
    cmd.Flags().StringVarP(&action, "action", "a", "", "Action (HTTP method)")

    cmd.MarkFlagRequired("subject")
    cmd.MarkFlagRequired("object")
    cmd.MarkFlagRequired("action")

    return cmd
}

// Additional commands for policy management...
```

### 3.3 Integration with Server

Update the server initialization to use the domain-aware enforcer:

```go
// internal/server/server.go
package server

import (
    "github.com/gin-gonic/gin"
    "github.com/spectrumwebco/agent_runtime/internal/server/auth"
    "github.com/spectrumwebco/agent_runtime/internal/server/middleware"
    "github.com/spectrumwebco/agent_runtime/internal/server/routes"
)

// Server represents the HTTP server
type Server struct {
    router   *gin.Engine
    enforcer *auth.DomainEnforcer
}

// NewServer creates a new server
func NewServer() (*Server, error) {
    router := gin.Default()

    // Create domain enforcer
    enforcer, err := auth.NewDomainEnforcer("config/auth/model.conf", "config/auth/policy.csv")
    if err != nil {
        return nil, err
    }

    server := &Server{
        router:   router,
        enforcer: enforcer,
    }

    // Set up middleware
    router.Use(middleware.DomainAuthorizeMiddleware(enforcer))

    // Set up routes
    apiRouter := router.Group("/api/v1")
    routes.RegisterAPIRoutes(apiRouter)
    routes.RegisterPolicyRoutes(apiRouter.Group("/admin"), enforcer)

    return server, nil
}

// Run starts the server
func (s *Server) Run(addr string) error {
    return s.router.Run(addr)
}
```

### 3.4 Integration with Authentication

Integrate with JWT authentication:

```go
// internal/server/middleware/jwt.go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

// JWTAuthMiddleware creates middleware that authenticates requests using JWT
func JWTAuthMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header is required",
            })
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header format must be Bearer {token}",
            })
            return
        }

        tokenString := parts[1]
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired token",
            })
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid token claims",
            })
            return
        }

        // Set user in context
        c.Set("user", claims["sub"])
        c.Set("roles", claims["roles"])

        c.Next()
    }
}
```

## 4. Database Integration

Integrate Casbin with a database adapter for policy storage:

```go
// internal/server/auth/db_adapter.go
package auth

import (
    "github.com/casbin/casbin/v2"
    "github.com/casbin/casbin/v2/model"
    gormadapter "github.com/casbin/gorm-adapter/v3"
    "github.com/spectrumwebco/agent_runtime/internal/database"
)

// NewDatabaseEnforcer creates a new Casbin enforcer with database adapter
func NewDatabaseEnforcer(modelPath string) (*DomainEnforcer, error) {
    m, err := model.NewModelFromFile(modelPath)
    if err != nil {
        return nil, err
    }

    adapter, err := gormadapter.NewAdapterByDB(database.GetDB())
    if err != nil {
        return nil, err
    }

    enforcer, err := casbin.NewEnforcer(m, adapter)
    if err != nil {
        return nil, err
    }

    return &DomainEnforcer{
        enforcer: enforcer,
    }, nil
}
```

## 5. Implementation Timeline

1. **Phase 1: Core Components**
   - Implement domain enforcer
   - Create domain-aware middleware
   - Update model configuration
   - Update policy file

2. **Phase 2: API Integration**
   - Implement policy management API
   - Create CLI commands for policy management
   - Update server initialization

3. **Phase 3: Authentication Integration**
   - Implement JWT middleware
   - Integrate with authentication system

4. **Phase 4: Database Integration**
   - Implement database adapter
   - Migrate policies to database

5. **Phase 5: Testing and Documentation**
   - Write unit tests
   - Create integration tests
   - Update documentation

## 6. Dependencies

- github.com/casbin/casbin/v2
- github.com/casbin/gorm-adapter/v3
- github.com/golang-jwt/jwt/v4

## 7. Conclusion

This integration plan provides a comprehensive approach to enhancing the Kled.io Framework with Casbin's domain-aware RBAC capabilities. The implementation will provide a robust, flexible, and scalable authorization system that can handle the framework's growing needs.

The deep integration approach ensures that the framework leverages the full capabilities of Casbin while maintaining a cohesive architecture that users can easily interact with. The domain-aware RBAC model will provide proper isolation between workspaces and users, enhancing the security and multi-tenancy capabilities of the framework.
