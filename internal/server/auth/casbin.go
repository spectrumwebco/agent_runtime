package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gin-gonic/gin"
)

// Enforcer is the Casbin enforcer
type Enforcer struct {
	enforcer *casbin.Enforcer
}

// NewEnforcer creates a new Casbin enforcer
func NewEnforcer(modelPath, policyPath string) (*Enforcer, error) {
	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load model: %w", err)
	}

	adapter := fileadapter.NewAdapter(policyPath)
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %w", err)
	}

	return &Enforcer{
		enforcer: enforcer,
	}, nil
}

// Authorize is a middleware that authorizes a request
func (e *Enforcer) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetString("user")
		if user == "" {
			user = "anonymous"
		}

		path := c.Request.URL.Path
		method := c.Request.Method

		// Check if the user is authorized
		allowed, err := e.enforcer.Enforce(user, path, method)
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

// GetRoles returns the roles of a user
func (e *Enforcer) GetRoles(user string) []string {
	return e.enforcer.GetRolesForUser(user)
}

// AddRoleForUser adds a role for a user
func (e *Enforcer) AddRoleForUser(user, role string) (bool, error) {
	return e.enforcer.AddRoleForUser(user, role)
}

// DeleteRoleForUser deletes a role for a user
func (e *Enforcer) DeleteRoleForUser(user, role string) (bool, error) {
	return e.enforcer.DeleteRoleForUser(user, role)
}

// AddPolicy adds a policy
func (e *Enforcer) AddPolicy(sub, obj, act string) (bool, error) {
	return e.enforcer.AddPolicy(sub, obj, act)
}

// RemovePolicy removes a policy
func (e *Enforcer) RemovePolicy(sub, obj, act string) (bool, error) {
	return e.enforcer.RemovePolicy(sub, obj, act)
}

// HasPolicy checks if a policy exists
func (e *Enforcer) HasPolicy(sub, obj, act string) bool {
	return e.enforcer.HasPolicy(sub, obj, act)
}

// GetAllPolicies returns all policies
func (e *Enforcer) GetAllPolicies() [][]string {
	return e.enforcer.GetPolicy()
}

// SavePolicy saves the policy to the adapter
func (e *Enforcer) SavePolicy() error {
	return e.enforcer.SavePolicy()
}
