package sqlxid

import (
	"log"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestDiscoveryCache(t *testing.T) {
	db, err := sqlx.Open("postgres", "host=localhost user=postgres dbname=pugit sslmode=disable")
	if err != nil {
		log.Println(err)
	}

	dc := CreateNewStore(db).DiscoveryCache

	// Put some initial values
	dc.Put("foo", &DiscoveredInfo{opEndpoint: "a", opLocalID: "b", claimedID: "c"})

	// Make sure we can retrieve them
	if di := dc.Get("foo"); di == nil {
		t.Errorf("Expected a result, got nil")
	} else if di.OpEndpoint() != "a" || di.OpLocalID() != "b" || di.ClaimedID() != "c" {
		t.Errorf("Expected a b c, got %v %v %v", di.OpEndpoint(), di.OpLocalID(), di.ClaimedID())
	}

	// Attempt to get a non-existent value
	if di := dc.Get("bar"); di != nil {
		t.Errorf("Expected nil, got %v", di)
	}
}
