package rest

import (
	"github.com/synapse-garden/sg-proto/auth"
	"github.com/synapse-garden/sg-proto/incept"
	"github.com/synapse-garden/sg-proto/store"
	"github.com/synapse-garden/sg-proto/users"

	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
)

// Needed endpoints:
//
// Create a new user:
//  - POST /incept/:credential ("magic link")
//
// Get a new session / login:
//  - POST /session/:user_id :pwhash (returns Token to be included as Authorization: Bearer)
//  - GET  /profile (user ID inferred)
//
// Open a new chat socket
//  - GET /chat/:user_id
//
// TODOs / tasks
//  - POST /todo {bounty, due}
//  - POST /todo/:id/complete => Get bounty if before due

// API is a transform on an httprouter.Router, passing a DB for passing
// on to httprouter.Handles.
type API func(*httprouter.Router, *bolt.DB) error

// Bind binds the API on the given DB.  It sets up REST endpoints as needed.
func Bind(db *bolt.DB, source *SourceInfo) (*httprouter.Router, error) {
	if err := db.Update(store.Prep(
		incept.TicketBucket,
		users.UserBucket,
		auth.LoginBucket,
	)); err != nil {
		return nil, err
	}
	htr := httprouter.New()
	for _, api := range []API{
		Source(source),
		Incept,
	} {
		if err := api(htr, db); err != nil {
			return nil, err
		}
	}

	return htr, nil
}
