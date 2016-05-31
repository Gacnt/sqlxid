package models

import "time"

var NonceSchema = `
CREATE TABLE IF NOT EXISTS nonces (
  T 	   TIMESTAMP WITH TIME ZONE NOT NULL,
  S 	   TEXT                     NOT NULL,
  ENDPOINT TEXT   	            NOT NULL
);
`

type Nonce struct {
	T        time.Time
	S        string
	Endpoint string
}
