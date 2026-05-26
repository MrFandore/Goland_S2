package graph

import "Prak_11/internal/store"

// Resolver is the root resolver, holding application dependencies.
type Resolver struct {
	Store *store.Store
}
