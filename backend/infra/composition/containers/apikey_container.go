package container

import (
	apikeyapp "ai_hub.com/app/core/app/apikey/apikeyapp"

	uowmongo "ai_hub.com/app/infra/db/mongoose/uow"

	apikeyrepo "ai_hub.com/app/infra/db/mongoose/repos/apikeyrepo"

	"ai_hub.com/app/infra/adapters/auth"
	"ai_hub.com/app/infra/adapters/crypto"
	"ai_hub.com/app/infra/adapters/gollm"
	"ai_hub.com/app/infra/adapters/idgen"
	"ai_hub.com/app/infra/config"

	"ai_hub.com/app/infra/http/gin/controllers"
	"ai_hub.com/app/infra/http/gin/middlewares"
	"ai_hub.com/app/infra/http/gin/routes"

	"go.mongodb.org/mongo-driver/mongo"
)

func BuildAPIKeyModule(client *mongo.Client, dbName string) (*routes.APIKeyRouter, *controllers.APIKeyController, *uowmongo.MongoUnitOfWork) {
	db := client.Database(dbName)

	// Unit of Work
	uow := uowmongo.NewMongoUnitOfWork(client)

	// Collections
	apiKeys := db.Collection("api_keys")

	// Repositories (infra)
	readRepo := apikeyrepo.NewAPIKeyReadRepoMongo(apiKeys)
	writeRepo := apikeyrepo.NewAPIKeyWriteRepoMongo(apiKeys)

	// Adapters
	crypt := crypto.NewCryptoAdapter()
	idGen := idgen.NewUUIDGenerator()
	modelCatalog := gollm.NewLLMGoModelCatalog()

	// App Handlers (core)
	createH := apikeyapp.NewCreateAPIKeyHandler(readRepo, writeRepo, crypt, idGen, uow, modelCatalog)
	updateH := apikeyapp.NewUpdateAPIKeyHandler(readRepo, writeRepo, crypt, uow, modelCatalog)
	deleteH := apikeyapp.NewDeleteAPIKeyHandler(readRepo, writeRepo, uow)

	getByIDH := apikeyapp.NewGetKeyByIDHandler(readRepo, crypt)
	getByOwnerH := apikeyapp.NewGetKeysByOwnerHandler(readRepo, crypt)

	// Controller
	controller := controllers.NewAPIKeyController(
		createH,
		updateH,
		deleteH,
		getByIDH,
		getByOwnerH,
	)

	//JWT verifier (implements middlewares.TokenVerifier)
	verifier := auth.NewJwtTokenIssuer(config.Env.JWTSecret, config.Env.JWTExpiresIn)

	// compile check for interface compliance:
	var _ middlewares.TokenVerifier = verifier

	// Pass verifier to the router
	router := routes.NewAPIKeyRoutes(controller, verifier)

	return router, controller, uow
}
