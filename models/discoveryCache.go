package models

var DiscoveryCacheSchema = `
CREATE TABLE discovery_caches (
  CACHE_ID  TEXT NOT NULL,
  ENDPOINT  TEXT NOT NULL,
  LOCAL_ID   TEXT NOT NULL,
  CLAIMED_ID TEXT NOT NULL
);
`

type DiscoveryCache struct {
	CacheID   string `db:"cache_id"`
	Endpoint  string
	LocalID   string `db:"local_id"`
	ClaimedID string `db:"claimed_id"`
}
