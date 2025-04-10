package server

import (
	"fmt"
	"net/http"
	"time" // Added import for time

	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/agent" // Import agent package
)

func SetupRouter(agentCore *agent.DefaultAgent /*, mcpHost *mcp.Host */) *gin.Engine {
	fmt.Println("Setting up Gin router...")
	router := gin.Default()


	router.GET("/is_alive", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive", "version": "go-agent-runtime-dev"})
	})

	router.POST("/create_session", func(c *gin.Context) {
		// }
		sessionID := "session-" + fmt.Sprintf("%d", time.Now().UnixNano()) // Placeholder
		fmt.Printf("Creating session: %s\n", sessionID)
		c.JSON(http.StatusOK, gin.H{"session_id": sessionID}) // Return runtime.CreateSessionResponse equivalent
	})

	router.POST("/run_in_session", func(c *gin.Context) {
		// }
		var actionData map[string]interface{}
		if err := c.ShouldBindJSON(&actionData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action format: " + err.Error()})
			return
		}
		fmt.Printf("Running in session: %+v\n", actionData) // Placeholder
		c.JSON(http.StatusOK, gin.H{ // Return runtime.Observation equivalent
			"observation": "Placeholder observation for action",
			"exit_code":   0,
		})
	})

	router.POST("/close_session", func(c *gin.Context) {
		// }
		var closeReq map[string]interface{}
		if err := c.ShouldBindJSON(&closeReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid close request format: " + err.Error()})
			return
		}
		fmt.Printf("Closing session: %+v\n", closeReq) // Placeholder
		c.JSON(http.StatusOK, gin.H{"status": "closed"}) // Return runtime.CloseResponse equivalent
	})

	router.POST("/execute", func(c *gin.Context) {
		// }
		var commandData map[string]interface{}
		if err := c.ShouldBindJSON(&commandData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid command format: " + err.Error()})
			return
		}
		fmt.Printf("Executing command: %+v\n", commandData) // Placeholder
		c.JSON(http.StatusOK, gin.H{ // Return runtime.Observation equivalent
			"observation": "Placeholder observation for command",
			"exit_code":   0,
		})
	})

	router.POST("/read_file", func(c *gin.Context) {
		// }
		var readReq map[string]interface{}
		if err := c.ShouldBindJSON(&readReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid read request format: " + err.Error()})
			return
		}
		fmt.Printf("Reading file: %+v\n", readReq) // Placeholder
		c.JSON(http.StatusOK, gin.H{"content": "Placeholder file content"}) // Return runtime.FileReadResponse equivalent
	})

	router.POST("/write_file", func(c *gin.Context) {
		// }
		var writeReq map[string]interface{}
		if err := c.ShouldBindJSON(&writeReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid write request format: " + err.Error()})
			return
		}
		fmt.Printf("Writing file: %+v\n", writeReq) // Placeholder
		c.JSON(http.StatusOK, gin.H{"status": "written"}) // Return runtime.FileWriteResponse equivalent
	})

	router.POST("/upload", func(c *gin.Context) {
		fmt.Println("Handling file upload...") // Placeholder
		c.JSON(http.StatusOK, gin.H{"status": "uploaded"}) // Return runtime.UploadResponse equivalent
	})

	router.POST("/close", func(c *gin.Context) {
		fmt.Println("Closing runtime...") // Placeholder
		c.JSON(http.StatusOK, gin.H{"status": "closing"}) // Return runtime.CloseResponse equivalent
	})

	router.POST("/agent/run", func(c *gin.Context) {
		var runReq map[string]interface{} // Placeholder for request structure
		if err := c.ShouldBindJSON(&runReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid run request format: " + err.Error()})
			return
		}
		fmt.Printf("Received request to run agent: %+v\n", runReq)

		go func() { // Run agent in background goroutine
			result, err := agentCore.Run(/* Pass env, problemStatement, outputDir */)
			if err != nil {
				fmt.Printf("Agent run failed: %v\n", err)
			} else {
				fmt.Printf("Agent run completed successfully: %+v\n", result.Info)
			}
		}()
		c.JSON(http.StatusOK, gin.H{"message": "Agent run initiated"})
	})

	router.GET("/agent/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "idle", "info": agentCore.Info}) // Placeholder status
	})

	return router
}
