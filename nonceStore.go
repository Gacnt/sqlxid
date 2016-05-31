package sqlxid

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Gacnt/sqlxid/models"
	"github.com/jmoiron/sqlx"
)

var maxNonceAge = flag.Duration("max-nonce-age",
	60*time.Second,
	"Maximum accepted age for openid nonces. The bigger, the more"+
		"memory is needed to store used nonces.")

type nonceStore struct {
	db *sqlx.DB
	mu sync.Mutex
}

func (d *nonceStore) Accept(endpoint, nonce string) error {
	if len(nonce) < 20 || len(nonce) > 256 {
		return errors.New("Invalid nonce")
	}

	ts, err := time.Parse(time.RFC3339, nonce[0:20])
	if err != nil {
		return err
	}

	now := time.Now()
	diff := now.Sub(ts)
	if diff > *maxNonceAge {
		return fmt.Errorf("Nonce too old: %ds", diff.Seconds())
	}

	s := nonce[20:]

	d.mu.Lock()
	defer d.mu.Unlock()

	nonces := []models.Nonce{}
	sqlQuery := `
	SELECT * FROM nonces
	`
	err = d.db.Select(&nonces, sqlQuery)
	if err != nil {
		log.Println(err)
	}

	newNonces := []models.Nonce{{ts, s, endpoint}}
	if len(nonces) > 0 {
		sqlQuery = `
		DELETE FROM nonces;
		`
		_, err := d.db.Exec(sqlQuery)
		if err != nil {
			log.Println(err)
		}

		for _, n := range nonces {
			if n.T.UTC() == ts && n.S == s {
				return errors.New("Nonce already used")
			}
			if now.Sub(n.T) < *maxNonceAge {
				newNonces = append(newNonces, n)
			}
		}
		if ok := storeNonces(d.db, newNonces); !ok {
			return errors.New("Could not store nonces")
		}
	} else {
		if ok := storeNonces(d.db, newNonces); !ok {
			return errors.New("Could not store nonces")
		}
	}

	return nil
}

func storeNonces(db *sqlx.DB, nonces []models.Nonce) bool {

	tx := db.MustBegin()
	for _, n := range nonces {
		tx.MustExec("INSERT INTO nonces (t, s, endpoint) VALUES ($1, $2, $3)", n.T, n.S, n.Endpoint)
	}
	tx.Commit()

	return true
}
