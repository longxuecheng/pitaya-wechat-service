package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"pitaya-wechat-service/controller"
	"pitaya-wechat-service/middle_ware"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(middle_ware.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(middle_ware.Recovery())
	// r.GET("/", func(c *gin.Context) {
	// 	time.Sleep(5 * time.Second)
	// 	c.String(http.StatusOK, "Welcome Gin Server")
	// })
	root := r.Group("/")
	root.Use(middle_ware.WrapResponse)
	root.GET("/user/all", controller.GetUserListByConditions)
	root.GET("/category/tree", controller.GetCatogoriesTree)
	root.GET("/category/info", controller.GetCategoryInfo)
	goodsGroup := root.Group("/goods")
	goodsGroup.GET("/list", controller.GetGoodsListByCategory)
	goodsGroup.GET("/detail", controller.GetGoodsInfo)
	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
