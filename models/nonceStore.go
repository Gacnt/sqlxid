package models

import "time"

var NonceSchema = `
CREATE TABLE nonces (
  T 	   TIME WITH TIMESTAMP ZONE NOT NULL,
  S 	   TEXT                     NOT NULL,
  ENDPOINT TEXT   	            NOT NULL
);
`

type Nonce struct {
	T        time.Time
	S        string
	Endpoint string
}
