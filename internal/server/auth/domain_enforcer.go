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

type DomainEnforcer struct {
	enforcer *casbin.Enforcer
}

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

func (e *DomainEnforcer) AuthorizeWithDomain(domainExtractor func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetString("user")
		if user == "" {
			user = "anonymous"
		}

		domain := domainExtractor(c)
		if domain == "" {
			domain = "system" // Default domain
		}

		path := c.Request.URL.Path
		method := c.Request.Method

		allowed, err := e.enforcer.Enforce(user, domain, path, method)
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

func (e *DomainEnforcer) GetRolesForUserInDomain(user, domain string) []string {
	return e.enforcer.GetRolesForUserInDomain(user, domain)
}

func (e *DomainEnforcer) AddRoleForUserInDomain(user, role, domain string) (bool, error) {
	return e.enforcer.AddRoleForUserInDomain(user, role, domain)
}

func (e *DomainEnforcer) DeleteRoleForUserInDomain(user, role, domain string) (bool, error) {
	return e.enforcer.DeleteRoleForUserInDomain(user, role, domain)
}

func (e *DomainEnforcer) AddPolicy(sub, dom, obj, act string) (bool, error) {
	return e.enforcer.AddPolicy(sub, dom, obj, act)
}

func (e *DomainEnforcer) RemovePolicy(sub, dom, obj, act string) (bool, error) {
	return e.enforcer.RemovePolicy(sub, dom, obj, act)
}

func (e *DomainEnforcer) HasPolicy(sub, dom, obj, act string) bool {
	return e.enforcer.HasPolicy(sub, dom, obj, act)
}

func (e *DomainEnforcer) GetAllPolicies() [][]string {
	return e.enforcer.GetPolicy()
}

func (e *DomainEnforcer) GetFilteredPolicy(fieldIndex int, fieldValues ...string) [][]string {
	return e.enforcer.GetFilteredPolicy(fieldIndex, fieldValues...)
}

func (e *DomainEnforcer) SavePolicy() error {
	return e.enforcer.SavePolicy()
}

func (e *DomainEnforcer) LoadPolicy() error {
	return e.enforcer.LoadPolicy()
}

func (e *DomainEnforcer) GetAllDomains() []string {
	policies := e.enforcer.GetPolicy()
	domains := make(map[string]bool)
	
	for _, policy := range policies {
		if len(policy) > 1 {
			domains[policy[1]] = true
		}
	}
	
	result := make([]string, 0, len(domains))
	for domain := range domains {
		result = append(result, domain)
	}
	
	return result
}

func (e *DomainEnforcer) GetUsersInDomain(domain string) []string {
	users := make(map[string]bool)
	
	policies := e.enforcer.GetFilteredPolicy(1, domain)
	for _, policy := range policies {
		if len(policy) > 0 {
			users[policy[0]] = true
		}
	}
	
	roles := e.enforcer.GetFilteredGroupingPolicy(2, domain)
	for _, role := range roles {
		if len(role) > 0 {
			users[role[0]] = true
		}
	}
	
	result := make([]string, 0, len(users))
	for user := range users {
		result = append(result, user)
	}
	
	return result
}
