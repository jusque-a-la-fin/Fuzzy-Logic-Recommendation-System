package registry

import (
	"database/sql"
	"vehicles/packages/adapters/controller"

	"github.com/redis/go-redis/v9"
)

type registry struct {
	rdb *redis.Client
	db  *sql.DB
}

type Registry interface {
	NewAppController() controller.AppController
}

func NewRegistry(rdb *redis.Client, db *sql.DB) Registry {
	return &registry{rdb, db}
}

func (r *registry) NewAppController() controller.AppController {

	return controller.AppController{}
}
