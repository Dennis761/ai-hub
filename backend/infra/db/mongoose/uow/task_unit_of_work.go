package uow

import (
	"ai_hub.com/app/core/ports/taskports"
)

var _ taskports.UnitOfWorkPort = (*MongoUnitOfWork)(nil)
