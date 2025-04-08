// ent/schema/errorlog.go

package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// 此表对应的是，错误落盘
type AppLog struct {
	ent.Schema
}

// Fields of the ErrorLog.
func (AppLog) Fields() []ent.Field {
	return []ent.Field{
		field.Time("time"),
		field.String("level"),
		field.String("msg"),
		field.String("kind").
			Optional(),
		field.String("component").
			Optional(),
		field.String("operation").
			Optional(),
		field.Int64("user_id").
			Optional(),
		field.String("trace_id").
			Optional(),
		field.String("span_id").
			Optional(),
		field.String("args").
			Optional(),
		field.Int("code").
			Optional(),
		field.String("reason").
			Optional(),
		field.String("stack").
			Optional(),
		field.Float("latency").
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.JSON("extra", map[string]any{}),
		field.String("app_name").Default(""),
	}
}

// Indexes of the ErrorLog.
func (AppLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("time"),
		index.Fields("trace_id"),
		index.Fields("code"),
	}
}
