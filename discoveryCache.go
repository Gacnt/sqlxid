package sqlxid

import (
	"log"

	"github.com/Gacnt/sqlxid/models"
	"github.com/jmoiron/sqlx"
	"github.com/yohcop/openid-go"
)

type discoveryCache struct {
	db *sqlx.DB
}

type DiscoveredInfo struct {
	opEndpoint string
	opLocalID  string
	claimedID  string
}

func (s *DiscoveredInfo) OpEndpoint() string {
	return s.opEndpoint
}

func (s *DiscoveredInfo) OpLocalID() string {
	return s.opLocalID
}

func (s *DiscoveredInfo) ClaimedID() string {
	return s.claimedID
}

func (d *discoveryCache) Put(id string, info openid.DiscoveredInfo) {
	dC := &models.DiscoveryCache{id, info.OpEndpoint(), info.OpLocalID(), info.ClaimedID()}
	sqlQuery := `
	INSERT INTO discovery_caches (cache_id, endpoint, local_id, claimed_id)
	VALUES (:id, :endpoint, :localid, :claimedid);
	`
	_, err := d.db.NamedExec(sqlQuery, map[string]interface{}{
		"id":        id,
		"endpoint":  dC.Endpoint,
		"localid":   dC.LocalID,
		"claimedid": dC.ClaimedID,
	})
	if err != nil {
		log.Println(err)
		log.Println("[SQLXID]: Failed to save discovery cache")
	}
}

func (d *discoveryCache) Get(id string) openid.DiscoveredInfo {
	sqlQuery := `
	SELECT * FROM discovery_caches
	WHERE cache_id = $1;
	`
	dc := models.DiscoveryCache{}
	err := d.db.Get(&dc, sqlQuery, id)
	if err != nil {
		return nil
	}

	sDI := &DiscoveredInfo{dc.Endpoint, dc.LocalID, dc.ClaimedID}

	return sDI
}
