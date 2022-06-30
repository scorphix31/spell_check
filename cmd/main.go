package main

import (
	"encoding/json"
	"fmt"
	"os"

	"spell_check/controllers"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("ERROR: %v\n", r)
			os.Exit(1)
		}
	}()

	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	logger, err := cfg.Build()

	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	gin.SetMode(gin.ReleaseMode)
	spellController := controllers.NewSpellController(logger)
	router := gin.New()
	router.POST("/checkText", spellController.CheckText)
	router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/pathsNotToLog/"),
		gin.Recovery(),
	)

	logger.Info("logger construction succeeded")
	logger.Info("Starting server...")

	router.Run("localhost:8080")
}
