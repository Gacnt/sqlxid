package sqlxid

import (
	"log"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestDefaultNonceStore(t *testing.T) {
	db, err := sqlx.Open("postgres", "host=localhost user=postgres dbname=pugit sslmode=disable")
	if err != nil {
		log.Println(err)
	}
	*maxNonceAge = 60 * time.Second
	now := time.Now().UTC()
	// 30 seconds ago
	now30s := now.Add(-30 * time.Second)
	// 2 minutes ago
	now2m := now.Add(-2 * time.Minute)

	now30sStr := now30s.Format(time.RFC3339)
	now2mStr := now2m.Format(time.RFC3339)

	ns := CreateNewStore(db)
	reject(t, ns.NonceStore, "1", "foo")                        // invalid nonce
	reject(t, ns.NonceStore, "1", "fooBarBazLongerThan20Chars") // invalid nonce

	accept(t, ns.NonceStore, "1", now30sStr+"asd")
	reject(t, ns.NonceStore, "1", now30sStr+"asd") // same nonce
	accept(t, ns.NonceStore, "1", now30sStr+"xxx") // different nonce
	reject(t, ns.NonceStore, "1", now30sStr+"xxx") // different nonce again to verify storage of multiple nonces per endpoint
	accept(t, ns.NonceStore, "2", now30sStr+"asd") // different endpoint

	reject(t, ns.NonceStore, "1", now2mStr+"old") // too old
	reject(t, ns.NonceStore, "3", now2mStr+"old") // too old
}

func accept(t *testing.T, ns *nonceStore, op, nonce string) {
	e := ns.Accept(op, nonce)
	if e != nil {
		t.Errorf("Should accept %s nonce %s", op, nonce)
	}
}

func reject(t *testing.T, ns *nonceStore, op, nonce string) {
	e := ns.Accept(op, nonce)
	if e == nil {
		t.Errorf("Should reject %s nonce %s", op, nonce)
	}
}
