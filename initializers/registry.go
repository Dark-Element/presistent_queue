package initializers

import (
	"database/sql"
)

type Registry struct {
	Db  *sql.DB
}


func GetRegistry() *Registry{
	registry := Registry{
		Db: database(),
		}
	return &registry
}