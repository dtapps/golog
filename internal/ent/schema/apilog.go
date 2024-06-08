package schema

import "entgo.io/ent"

// ApiLog holds the schema definition for the ApiLog entity.
type ApiLog struct {
	ent.Schema
}

// TableName overrides the default table name used by Ent.
func (ApiLog) TableName() string {
	return "my_custom_user_table"
}

// Fields of the ApiLog.
func (ApiLog) Fields() []ent.Field {
	return nil
}

// Edges of the ApiLog.
func (ApiLog) Edges() []ent.Edge {
	return nil
}
