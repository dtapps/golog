package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// GinLog holds the schema definition for the GinLog entity.
type GinLog struct {
	ent.Schema
}

// Fields of the GinLog.
func (GinLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StorageKey("request_id").Comment("请求编号"), // 用 request_id 覆盖 框架的 id
		field.Time("request_time").Optional().Comment("请求时间"),
		field.String("request_host").Optional().Comment("请求主机"),
		field.String("request_path").Optional().Comment("请求地址"),
		field.String("request_query").Optional().Comment("请求参数"),
		field.String("request_method").Optional().Comment("请求方式"),
		field.String("request_scheme").Optional().Comment("请求协议"),
		field.String("request_content_type").Optional().Comment("请求类型"),
		field.String("request_body").Optional().Comment("请求内容"),
		field.String("request_client_ip").Optional().Comment("请求IP"),
		field.String("request_user_agent").Optional().Comment("请求UA"),
		field.String("request_header").Optional().Comment("请求头"),
		field.Int64("request_cost_time").Optional().Comment("请求消耗时长"),
		field.Time("response_time").Optional().Comment("响应时间"),
		field.String("response_header").Optional().Comment("响应头"),
		field.Int("response_status_code").Optional().Comment("响应状态"),
		field.String("response_data").Optional().Comment("响应内容"),
		field.String("go_version").Optional().Comment("Go版本"),
		field.String("sdk_version").Optional().Comment("Sdk版本"),
		field.String("system_info").Optional().Comment("系统信息"),
	}
}

// Edges of the GinLog.
func (GinLog) Edges() []ent.Edge {
	return nil
}

func (GinLog) Indexes() []ent.Index {
	return []ent.Index{
		//index.Fields("request_id"),
		index.Fields("request_time"),
		index.Fields("request_path"),
		index.Fields("request_method"),
	}
}
