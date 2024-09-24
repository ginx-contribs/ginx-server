package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// CronJob holds the schema definition for the CronJob entity.
type CronJob struct {
	ent.Schema
}

// Fields of the CronJob.
func (CronJob) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("cron"),
		field.Int("entry_id"),
		field.Int64("prev"),
		field.Int64("next"),
	}
}

// Edges of the CronJob.
func (CronJob) Edges() []ent.Edge {
	return []ent.Edge{}
}

// Annotations of the CronJob
func (CronJob) Annotations() []schema.Annotation {
	return []schema.Annotation{}
}
