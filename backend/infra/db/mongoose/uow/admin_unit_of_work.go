package uow

import (
	"ai_hub.com/app/core/ports/adminports"
)

var _ adminports.UnitOfWorkPort = (*MongoUnitOfWork)(nil)
