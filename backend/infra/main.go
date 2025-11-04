package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	container "ai_hub.com/app/infra/composition/containers"
	"ai_hub.com/app/infra/config"
	mongosvc "ai_hub.com/app/infra/services/mongoDB"
	rediscli "ai_hub.com/app/infra/services/redis"

	"ai_hub.com/app/infra/http/gin/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()

	mongoClient := mongosvc.ConnectDB()
	dbName := extractDBName(config.Env.MongoDBURI)
	if dbName == "" {
		log.Fatalf("Database name not found in MONGODB_URI. Place /<db> at the end of the URI")
	}
	log.Printf("[mongo] DB name from URI: %s", dbName)

	redisClient := rediscli.ConnectRedis()
	_ = redisClient

	engine := gin.New()

	engine.Use(gin.Logger(), gin.Recovery(), middlewares.ErrorHandler())
	engine.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}))

	//Prefix /api for all modules
	api := engine.Group("/api")

	//Initialization of modules
	adminRouter, _, _ := container.BuildAdminModule(mongoClient, dbName)
	apiKeyRouter, _, _ := container.BuildAPIKeyModule(mongoClient, dbName)
	projectRouter, _, _ := container.BuildProjectModule(mongoClient, dbName, redisClient)
	promptRouter, _, _ := container.BuildPromptModule(mongoClient, dbName, redisClient)
	taskRouter, _, _ := container.BuildTaskModule(mongoClient, dbName, redisClient)

	// Mount with the /api prefix
	adminRouter.Mount(api)
	apiKeyRouter.Mount(api)
	projectRouter.Mount(api)
	promptRouter.Mount(api)
	taskRouter.Mount(api)

	// Starting the server
	srv := &http.Server{
		Addr:    ":" + config.Env.Port,
		Handler: engine,
	}

	go func() {
		log.Printf(" Server is running on port %s\n", config.Env.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(" listen: %v", err)
		}
	}()

	gracefulShutdown(srv)
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println(" Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf(" Server forced to shutdown: %v", err)
	}

	if err := mongosvc.Disconnect(ctx); err != nil {
		log.Printf("[MongoDB] Disconnect error: %v", err)
	}

	log.Println(" Server exited properly")
}

// We take the database name from the URI format ...mongodb.net/<db>?...
// Example: mongodb+srv://.../Ai-Hub?retryWrites=true -> will return "Ai-Hub"
func extractDBName(uri string) string {
	if uri == "" {
		return ""
	}

	//cut off query (?...)
	base := uri
	if i := strings.Index(uri, "?"); i >= 0 {
		base = uri[:i]
	}

	//take everything after the last /
	if j := strings.LastIndex(base, "/"); j >= 0 && j < len(base)-1 {
		return base[j+1:]
	}
	return ""
}
