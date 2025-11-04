package container

import (
	promptapp "ai_hub.com/app/core/app/prompt/promptapp"

	uowmongo "ai_hub.com/app/infra/db/mongoose/uow"

	apikeyrepo "ai_hub.com/app/infra/db/mongoose/repos/apikeyrepo"
	projectrepo "ai_hub.com/app/infra/db/mongoose/repos/projectrepo"
	promptrepo "ai_hub.com/app/infra/db/mongoose/repos/promptrepo"
	taskrepo "ai_hub.com/app/infra/db/mongoose/repos/taskrepo"

	"ai_hub.com/app/infra/adapters/auth"
	"ai_hub.com/app/infra/adapters/crypto"
	"ai_hub.com/app/infra/adapters/gollm"
	"ai_hub.com/app/infra/adapters/idgen"
	"ai_hub.com/app/infra/config"

	rediscache "ai_hub.com/app/infra/cache/redis"

	"ai_hub.com/app/infra/http/gin/controllers"
	"ai_hub.com/app/infra/http/gin/middlewares"
	"ai_hub.com/app/infra/http/gin/routes"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func BuildPromptModule(
	client *mongo.Client,
	dbName string,
	rdb *redis.Client,
) (*routes.PromptRouter, *controllers.PromptController, *uowmongo.MongoUnitOfWork) {
	db := client.Database(dbName)

	// Unit of Work
	uow := uowmongo.NewMongoUnitOfWork(client)

	// Collections
	prompts := db.Collection("prompts")
	tasks := db.Collection("tasks")
	projects := db.Collection("projects")
	apiKeys := db.Collection("api_keys")

	// Repositories (infra)
	promptRead := promptrepo.NewPromptReadRepoMongo(prompts)
	promptWrite := promptrepo.NewPromptWriteRepoMongo(prompts)
	taskRead := taskrepo.NewTaskReadRepoMongo(tasks)
	projectRead := projectrepo.NewProjectReadRepoMongo(projects)
	apiKeyRead := apikeyrepo.NewAPIKeyReadRepoMongo(apiKeys)
	apiKeyBalanceW := apikeyrepo.NewAPIKeyBalanceRepoMongo(apiKeys)

	// Adapters
	idGen := idgen.NewUUIDGenerator()
	crypt := crypto.NewCryptoAdapter()
	llmClient := gollm.NewGollmClientAdapter()
	billing := gollm.NewLLMBillingRunner(apiKeyRead, crypt, llmClient, apiKeyBalanceW)

	cache := rediscache.NewProjectCacheAdapterFromConfig(rdb)

	// App Handlers (core)
	createH := promptapp.NewCreatePromptHandler(
		promptWrite,
		uow,
		idGen,
		taskRead,
		projectRead,
		cache,
		billing,
	)

	updateH := promptapp.NewUpdatePromptHandler(
		promptRead,
		promptWrite,
		uow,
		taskRead,
		projectRead,
		cache,
	)

	deleteH := promptapp.NewDeletePromptHandler(
		promptRead,
		promptWrite,
		uow,
		taskRead,
		projectRead,
		cache,
	)

	rollbackH := promptapp.NewRollbackPromptHandler(
		promptRead,
		promptWrite,
		uow,
		taskRead,
		projectRead,
		cache,
	)

	reorderH := promptapp.NewReorderPromptsHandler(
		promptRead,
		promptWrite,
		uow,
		taskRead,
		projectRead,
		cache,
	)

	runH := promptapp.NewRunPromptHandler(
		promptRead,
		promptWrite,
		uow,
		billing,
		taskRead,
		projectRead,
	)

	listByTaskH := promptapp.NewGetPromptsByTaskHandler(promptRead, taskRead, projectRead)
	getByIDH := promptapp.NewGetPromptByIDHandler(promptRead, taskRead, projectRead)

	// Controller
	controller := controllers.NewPromptController(
		createH,
		updateH,
		deleteH,
		rollbackH,
		reorderH,
		runH,
		listByTaskH,
		getByIDH,
	)

	//JWT verifier (implements middlewares.TokenVerifier)
	verifier := auth.NewJwtTokenIssuer(config.Env.JWTSecret, config.Env.JWTExpiresIn)

	// compile check for interface compliance:
	var _ middlewares.TokenVerifier = verifier

	// Pass verifier to the router
	router := routes.NewPromptRoutes(controller, verifier)

	return router, controller, uow
}
