package main

import (
	"fmt"
	"log"
	"strings"
	"time"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
  "github.com/gorilla/websocket"
)

func runWebserver() {

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard // to disable web hits output to console
	// gin.DefaultWriter = colorable.NewColorableStdout()
	// gin.ForceConsoleColor()

	router := gin.Default()
	router.Use(cors.Default()) // allow all origins

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


	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}

var wsupgrader = websocket.Upgrader {
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
		wsupgrader.CheckOrigin = func(r *http.Request) bool {
			// TODO: check for localhost/127/rovermems.com ? don't actually care though
			return true
		}
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
		data["logLines"] = globalLogLines

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
