package entity

// Template represents the domain entity for the template aggregate.
// The domain does not know about Postgres, HTTP, or any infrastructure.
type Template struct {
	ID string
}
