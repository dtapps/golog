package schema

import "entgo.io/ent"

// Hertzlog holds the schema definition for the Hertzlog entity.
type Hertzlog struct {
	ent.Schema
}

// Fields of the Hertzlog.
func (Hertzlog) Fields() []ent.Field {
	return nil
}

// Edges of the Hertzlog.
func (Hertzlog) Edges() []ent.Edge {
	return nil
}
