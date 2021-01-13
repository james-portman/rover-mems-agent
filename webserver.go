package main

import (
	"io/ioutil"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
)

func runWebserver() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard // to disable web hits output to console
	// gin.DefaultWriter = colorable.NewColorableStdout()
	// gin.ForceConsoleColor()

	router := gin.Default()
	router.Use(cors.Default()) // allow all origins

	router.Use(static.Serve("/", static.LocalFile("web-static", false)))

  router.GET("/api", func(c *gin.Context) {
		globalDataOutputLock.RLock()
		defer globalDataOutputLock.RUnlock()
    c.JSON(200, gin.H{
      "faults": globalFaults,
      "connected": globalConnected,
      "ecuType": globalEcuType,
      "userCommand": globalUserCommand,
			"alert": globalAlert,
			"ecuData": globalDataOutput,
			"agentVersion": globalAgentVersion,
    })
		if globalAlert != "" {
			globalAlert = ""
		}
  })

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

  router.GET("/connected", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "connected": globalConnected,
    })
  })

  router.GET("/faults", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "faults": globalFaults,
    })
  })

	router.GET("/ecu/:name", func(c *gin.Context) {
		name := c.Param("name")
    globalEcuType = name
		c.String(http.StatusOK, "ECU type set to %s", name)
		globalAlert = "Agent confirms ECU set to "+name
	})

  router.GET("/command/:name", func(c *gin.Context) {
    name := c.Param("name")
    globalUserCommand = name
    c.String(http.StatusOK, "User command accepted %s", name)
  })

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
