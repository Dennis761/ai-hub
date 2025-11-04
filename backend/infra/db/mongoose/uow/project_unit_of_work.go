package uow

import (
	"ai_hub.com/app/core/ports/projectports"
)

var _ projectports.UnitOfWorkPort = (*MongoUnitOfWork)(nil)
