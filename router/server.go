package router

import (
	"context"
	"gotrue/middle_ware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func router() *gin.Engine {
	r := gin.New()
	r.Use(middle_ware.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(middle_ware.Recovery)
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
