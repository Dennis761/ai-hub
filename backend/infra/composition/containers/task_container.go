package container

import (
	taskapp "ai_hub.com/app/core/app/task/taskapp"

	uowmongo "ai_hub.com/app/infra/db/mongoose/uow"

	projectrepo "ai_hub.com/app/infra/db/mongoose/repos/projectrepo"
	taskrepo "ai_hub.com/app/infra/db/mongoose/repos/taskrepo"

	rediscache "ai_hub.com/app/infra/cache/redis"

	"ai_hub.com/app/infra/adapters/auth"
	"ai_hub.com/app/infra/adapters/idgen"
	"ai_hub.com/app/infra/config"

	"ai_hub.com/app/infra/http/gin/controllers"
	"ai_hub.com/app/infra/http/gin/middlewares"
	"ai_hub.com/app/infra/http/gin/routes"

	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func BuildTaskModule(
	client *mongo.Client,
	dbName string,
	rdb *redis.Client,
) (*routes.TaskRouter, *controllers.TaskController, *uowmongo.MongoUnitOfWork) {

	db := client.Database(dbName)

	// Unit of Work
	uow := uowmongo.NewMongoUnitOfWork(client)

	// Collections
	tasks := db.Collection("tasks")
	projects := db.Collection("projects")

	// Repos
	taskRead := taskrepo.NewTaskReadRepoMongo(tasks)
	taskWrite := taskrepo.NewTaskWriteRepoMongo(tasks)
	projectRead := projectrepo.NewProjectReadRepoMongo(projects)

	// Cache TTLs
	editTTLsec, _ := strconv.Atoi(config.Env.RedisProjectEditTTL)
	cacheTTLsec, _ := strconv.Atoi(config.Env.RedisProjectCacheTTL)
	cache := rediscache.NewProjectCacheAdapter(rdb, rediscache.ProjectCacheTTLs{
		EditTTL:  time.Duration(editTTLsec) * time.Second,
		CacheTTL: time.Duration(cacheTTLsec) * time.Second,
	})

	// Adapters
	idGen := idgen.NewUUIDGenerator()

	// App Handlers (core)
	createH := taskapp.NewCreateTaskHandler(
		taskWrite,
		uow,
		projectRead,
		cache,
		idGen,
	)

	updateH := taskapp.NewUpdateTaskHandler(
		taskRead,
		taskWrite,
		uow,
		projectRead,
		cache,
	)

	deleteH := taskapp.NewDeleteTaskHandler(
		taskRead,
		taskWrite,
		uow,
		projectRead,
		cache,
	)

	// Queries
	byProjectH := taskapp.NewGetTasksByProjectHandler(taskRead, projectRead)

	// Controller
	controller := controllers.NewTaskController(
		createH,
		updateH,
		deleteH,
		byProjectH,
	)

	//JWT verifier (implements middlewares.TokenVerifier)
	verifier := auth.NewJwtTokenIssuer(config.Env.JWTSecret, config.Env.JWTExpiresIn)

	// compile check for interface compliance:
	var _ middlewares.TokenVerifier = verifier

	// Pass verifier to the router
	router := routes.NewTaskRoutes(controller, verifier)

	return router, controller, uow
}
