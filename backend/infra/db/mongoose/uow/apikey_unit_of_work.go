package uow

import (
	"ai_hub.com/app/core/ports/apikeyports"
)

var _ apikeyports.UnitOfWorkPort = (*MongoUnitOfWork)(nil)
