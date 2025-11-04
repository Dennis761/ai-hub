package container

import (
	projectapp "ai_hub.com/app/core/app/project/projectapp"
	projectrepo "ai_hub.com/app/infra/db/mongoose/repos/projectrepo"
	uowmongo "ai_hub.com/app/infra/db/mongoose/uow"
	"github.com/redis/go-redis/v9"

	"ai_hub.com/app/infra/adapters/auth"
	"ai_hub.com/app/infra/adapters/crypto"
	"ai_hub.com/app/infra/adapters/idgen"

	rediscache "ai_hub.com/app/infra/cache/redis"
	"ai_hub.com/app/infra/config"

	"ai_hub.com/app/infra/http/gin/controllers"
	"ai_hub.com/app/infra/http/gin/middlewares"
	"ai_hub.com/app/infra/http/gin/routes"

	"go.mongodb.org/mongo-driver/mongo"
)

func BuildProjectModule(
	client *mongo.Client,
	dbName string,
	rdb *redis.Client,
) (*routes.ProjectRouter, *controllers.ProjectController, *uowmongo.MongoUnitOfWork) {
	db := client.Database(dbName)

	// Unit of Work
	uow := uowmongo.NewMongoUnitOfWork(client)

	// Collections
	projects := db.Collection("projects")

	// Repositories (infra)
	readRepo := projectrepo.NewProjectReadRepoMongo(projects)
	writeRepo := projectrepo.NewProjectWriteRepoMongo(projects)

	// Adapters
	hasher := crypto.NewBcryptHasher(10)
	idGen := idgen.NewUUIDGenerator()
	cache := rediscache.NewProjectCacheAdapterFromConfig(rdb)

	// App Handlers (core)
	createH := projectapp.NewCreateProjectHandler(readRepo, writeRepo, uow, idGen, hasher)
	updateH := projectapp.NewUpdateProjectHandler(readRepo, writeRepo, uow, cache, hasher)
	deleteH := projectapp.NewDeleteProjectHandler(readRepo, writeRepo, uow, cache)
	joinH := projectapp.NewJoinProjectByNameHandler(readRepo, writeRepo, uow, hasher)

	getMyH := projectapp.NewGetMyProjectsHandler(readRepo)
	getByIDH := projectapp.NewGetProjectByIDHandler(readRepo, cache)

	// Controller
	controller := controllers.NewProjectController(
		createH,
		updateH,
		deleteH,
		joinH,
		getMyH,
		getByIDH,
	)

	//JWT verifier (implements middlewares.TokenVerifier)
	verifier := auth.NewJwtTokenIssuer(config.Env.JWTSecret, config.Env.JWTExpiresIn)

	// compile check for interface compliance:
	var _ middlewares.TokenVerifier = verifier

	// Pass verifier to the router
	router := routes.NewProjectRoutes(controller, verifier)

	return router, controller, uow
}
