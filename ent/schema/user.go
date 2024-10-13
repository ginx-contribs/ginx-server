package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/ginx-contribs/ginx-server/pkg/toolset/idx"
	"github.com/ginx-contribs/ginx-server/pkg/toolset/ts"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("user info table"),
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("uid").DefaultFunc(idx.ULID).Unique(),
		field.String("username").Unique(),
		field.String("email").Unique(),
		field.String("password"),
		field.Int64("created_at").DefaultFunc(ts.UnixMicro),
		field.Int64("updated_at").DefaultFunc(ts.UnixMicro).UpdateDefault(ts.UnixMicro),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{}
}
