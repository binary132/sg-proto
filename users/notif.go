package users

import "github.com/synapse-garden/sg-proto/store"

// Resource is the name of the Resource for users.
const Resource = store.Resource("users")

// Resource implements store.Resourcer on *User.
func (u *User) Resource() store.Resource {
	return Resource
}
