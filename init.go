package sqlxid

import (
	"sync"

	"github.com/Gacnt/sqlxid/models"
	"github.com/jmoiron/sqlx"
)

type sqlxStore struct {
	DiscoveryCache *discoveryCache
	NonceStore     *nonceStore
}

func CreateNewStore(db *sqlx.DB) *sqlxStore {
	// Combine Schemas because only Postgres will exec multple `MustExec` others will not
	allSchemas := models.DiscoveryCacheSchema + models.NonceSchema
	db.MustExec(allSchemas)
	return &sqlxStore{
		DiscoveryCache: &discoveryCache{db},
		NonceStore:     &nonceStore{db, sync.Mutex{}},
	}
}
