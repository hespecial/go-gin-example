package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/hespecial/go-gin-example/docs"
	"github.com/hespecial/go-gin-example/models"
	"github.com/hespecial/go-gin-example/pkg/gredis"
	"github.com/hespecial/go-gin-example/pkg/logging"
	"github.com/hespecial/go-gin-example/pkg/setting"
	"github.com/hespecial/go-gin-example/pkg/util"
	"github.com/hespecial/go-gin-example/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()
	util.Setup()
}

// @title			Golang Gin API
// @version		1.0
// @description	An example of gin
// @termsOfService	https://github.com/hespecial/go-gin-example
// @license.name	MIT
// @license.url	https://github.com/hespecial/go-gin-example/blob/main/LICENSE
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	srv := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %v", err)
		}
	}()

	// 优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown server ... ")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}
	log.Println("Server exited.")
}
