package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/server/auth"
)

type DomainExtractor func(*gin.Context) string

func DefaultDomainExtractor(c *gin.Context) string {
	path := c.Request.URL.Path
	parts := strings.Split(path, "/")
	
	var validParts []string
	for _, part := range parts {
		if part != "" {
			validParts = append(validParts, part)
		}
	}
	
	if len(validParts) >= 3 && validParts[0] == "api" && validParts[1] == "v1" {
		switch validParts[2] {
		case "workspaces":
			return "workspace"
		case "agents":
			return "agent"
		case "users":
			return "user"
		case "tools":
			return "tool"
		case "modules":
			return "module"
		case "auth":
			return "auth"
		}
	}
	
	return "system" // Default domain
}

func HeaderDomainExtractor(headerName string) DomainExtractor {
	return func(c *gin.Context) string {
		domain := c.GetHeader(headerName)
		if domain == "" {
			return "system" // Default domain
		}
		return domain
	}
}

func QueryDomainExtractor(paramName string) DomainExtractor {
	return func(c *gin.Context) string {
		domain := c.Query(paramName)
		if domain == "" {
			return "system" // Default domain
		}
		return domain
	}
}

func DomainAuthMiddleware(authHeader string, extractor DomainExtractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header must be in the format 'Bearer {token}'",
			})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			return
		}

		if token != authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		domain := extractor(c)
		
		c.Set("user", "user")
		c.Set("domain", domain)

		c.Next()
	}
}

func DomainAuthorizeMiddleware(enforcer *auth.DomainEnforcer) gin.HandlerFunc {
	return enforcer.AuthorizeWithDomain(DefaultDomainExtractor)
}

func CustomDomainAuthorizeMiddleware(enforcer *auth.DomainEnforcer, extractor DomainExtractor) gin.HandlerFunc {
	return enforcer.AuthorizeWithDomain(extractor)
}
