package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"gotrue/middle_ware"
	"time"

	"github.com/gin-gonic/gin"
)

func router() *gin.Engine {
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(middle_ware.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(middle_ware.Recovery())
	r.Use(middle_ware.WrapResponse)
	apiRouter(r)
	return r
}

func ListenAndServe(port string) {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router(),
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
