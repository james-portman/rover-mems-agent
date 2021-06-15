package main

import (
	"embed"
	"fmt"
	"log"
	"strings"
	"time"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	// "github.com/gin-contrib/static"
  "github.com/gorilla/websocket"
)

//go:embed web-static/*
var web_static_content embed.FS

func runWebserver() {

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard // to disable web hits output to console
	// gin.DefaultWriter = colorable.NewColorableStdout()
	// gin.ForceConsoleColor()

	router := gin.Default()
	router.Use(cors.Default()) // allow all origins

	// router.Use(static.Serve("/", static.LocalFile("web-static", false)))

  router.GET("/api", func(c *gin.Context) {
		globalDataOutputLock.Lock()
    c.JSON(200, gin.H{
      "faults": globalFaults,
      "connected": globalConnected,
      "ecuType": globalEcuType,
      "userCommand": globalUserCommand,
			"alert": globalAlert,
			"error": globalError,
			"ecuData": globalDataOutput,
			"agentVersion": globalAgentVersion,
    })
		// clear the error and alert for next time
		if globalAlert != "" {
			globalAlert = ""
		}
		if globalError != "" {
			globalError = ""
		}
		globalDataOutputLock.Unlock()
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
		globalDataOutputLock.RLock()
    c.JSON(200, gin.H{
      "faults": globalFaults,
    })
		globalDataOutputLock.RUnlock()
  })

	router.GET("/ecu/:name", func(c *gin.Context) {
		globalDataOutputLock.Lock()
		name := c.Param("name")
    globalEcuType = name
		c.String(http.StatusOK, "ECU type set to %s", name)
		// globalAlert = "Agent confirms ECU set to "+name
		globalDataOutputLock.Unlock()
	})

	router.GET("/serialPort/:name", func(c *gin.Context) {
		globalDataOutputLock.Lock()
		name := c.Param("name")
		globalSelectedSerialPort = name
		c.String(http.StatusOK, "Serial port set to %s", name)
		// globalAlert = "Agent confirms ECU set to "+name
		globalDataOutputLock.Unlock()
	})

  router.GET("/command/:name", func(c *gin.Context) {
		globalDataOutputLock.Lock()
    name := c.Param("name")
    globalUserCommand = name
    c.String(http.StatusOK, "User command accepted %s", name)
		globalDataOutputLock.Unlock()
  })

	router.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})


	router.GET("/", myrouter)
	// router.GET("/", func(c *gin.Context) {
	// 	data, _ := web_static_content.ReadFile("web-static/index.html")
	// 	c.Writer.Header().Set("Content-Type", "text/html")
	// 	c.String(http.StatusOK, string(data))
	// 	// Cache-Control: no-cache, no-store, must-revalidate
	// 	// Pragma: no-cache
	// 	// Expires: 0
	// })
	router.GET("/:name1/:name2", myrouter)
	router.GET("/:name1", myrouter)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}

func myrouter(c *gin.Context) {
	// do this nicer somehow
	name := c.Param("name1")
	// fmt.Println(name)
	name2 := c.Param("name2")
	if name2 != "" {
		name = name + "/" + name2
	}
	if name == "" {
		name = "index.html"
	}

	data, _ := web_static_content.ReadFile("web-static/"+name)

	mimeType := ""
	fileSplit := strings.Split(name, ".")
	extension := fileSplit[len(fileSplit)-1]
	if extension == name { extension = "" }
	switch extension {
	case "html":
		mimeType = "text/html"
		break
	case "js":
		mimeType = "script/javascript"
		break
	case "css":
		mimeType = "text/css"
		break
	default:
		// mimeType = "text/plain"
		mimeType = "application/octet-stream"
		break
	}
	c.Writer.Header().Set("Content-Type", mimeType)
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Writer.Header().Set("Pragma", "no-cache")
	c.Writer.Header().Set("Expires", "0")
	c.String(http.StatusOK, string(data))
}

var wsupgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
    conn, err := wsupgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Failed to set websocket upgrade: %+v", err)
        return
    }

		iteration := 0
    for {
			// TODO: change to channel read then call function to send?
      err := wsiteration(conn, iteration)
      if err != nil {
        break
      }
			iteration++
    }
}

func wsiteration(conn *websocket.Conn, iteration int) error {

  // wait for a message from the browser (it is sending "." to request data)
  // message type, msg, err
  _, message, err := conn.ReadMessage()
  if err != nil {
    // fmt.Println("WS readmessage failed")
    return err
  }

	var data map[string]interface {} = make(map[string]interface{})

  if strings.Compare(string(message), ".") == 0 {

		globalDataOutputLock.RLock()
		data["faults"] = globalFaults
		data["connected"] = globalConnected
		data["ecuType"] = globalEcuType
		data["userCommand"] = globalUserCommand
		data["alert"] = globalAlert
		data["error"] = globalError
		data["ecuData"] = globalDataOutput
		data["agentVersion"] = globalAgentVersion
		data["timestamp"] = time.Now().String()
		data["serialPorts"] = globalSerialPorts
		data["selectedSerialPort"] = globalSelectedSerialPort
		if globalAlert != "" {
			globalAlert = ""
		}
		if globalError != "" {
			globalError = ""
		}
		globalDataOutputLock.RUnlock()

	} else {

		// must be a command if it wasn't . above
		log.Printf("recv: %s", message)

		data["command"] = "worked"
	}



  jsondata, err := json.Marshal(data)
  if err != nil {
    return err
  }
  conn.WriteMessage(websocket.TextMessage, jsondata)
  return nil
}
