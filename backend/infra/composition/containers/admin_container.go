package container

import (
	adminapp "ai_hub.com/app/core/app/admin/adminapp"

	adminrepo "ai_hub.com/app/infra/db/mongoose/repos/adminrepo"
	uowmongo "ai_hub.com/app/infra/db/mongoose/uow"

	"ai_hub.com/app/infra/adapters/auth"
	"ai_hub.com/app/infra/adapters/crypto"
	"ai_hub.com/app/infra/adapters/idgen"
	"ai_hub.com/app/infra/adapters/mailer"

	"ai_hub.com/app/infra/config"
	"ai_hub.com/app/infra/http/gin/controllers"
	"ai_hub.com/app/infra/http/gin/routes"

	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func BuildAdminModule(client *mongo.Client, dbName string) (*routes.AdminRouter, *controllers.AdminController, *uowmongo.MongoUnitOfWork) {
	db := client.Database(dbName)

	// Unit of Work
	uow := uowmongo.NewMongoUnitOfWork(client)

	// Collections
	admins := db.Collection("admins")

	// Repositories (infra)
	readRepo := adminrepo.NewAdminReadRepoMongo(admins)
	writeRepo := adminrepo.NewAdminWriteRepoMongo(admins)

	// Adapters
	hasher := crypto.NewBcryptHasher(10)

	// CodeTTL = 15 minutes
	codeTTL := 15 * time.Minute
	codegen := crypto.NewSimpleCodeGenerator(codeTTL)

	// Mailer (SMTP)

	mailer := mailer.NewSMTPMailer(
		config.Env.SMTPHost,
		config.Env.EmailUser,
		config.Env.EmailPass,
	)

	// JWT issuer/verifier
	tokenIssuer := auth.NewJwtTokenIssuer(config.Env.JWTSecret, config.Env.JWTExpiresIn)

	idGen := idgen.NewUUIDGenerator()

	// App Handlers (core)
	registerH := adminapp.NewCreateAdminHandler(readRepo, writeRepo, hasher, codegen, mailer, uow, idGen)
	verifyH := adminapp.NewVerifyAdminHandler(readRepo, writeRepo, uow)
	startResetH := adminapp.NewStartPasswordResetHandler(readRepo, writeRepo, codegen, mailer, uow)
	confirmResetH := adminapp.NewConfirmResetCodeHandler(readRepo, writeRepo, uow)
	changePwH := adminapp.NewChangePasswordWithCodeHandler(readRepo, writeRepo, hasher, uow)
	renameH := adminapp.NewRenameAdminHandler(readRepo, writeRepo, uow)
	deleteH := adminapp.NewDeleteAdminHandler(readRepo, writeRepo, uow)
	authH := adminapp.NewLoginAdminHandler(readRepo, hasher, tokenIssuer)

	// Controller & Router
	controller := controllers.NewAdminController(
		registerH, verifyH, startResetH, confirmResetH, changePwH, renameH, deleteH, authH,
	)

	router := routes.NewAdminRoutes(controller)

	return router, controller, uow
}
