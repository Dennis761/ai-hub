package uow

import (
	"ai_hub.com/app/core/ports/promptports"
)

var _ promptports.UnitOfWorkPort = (*MongoUnitOfWork)(nil)
