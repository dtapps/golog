// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"go.dtapp.net/golog/internal/ent/ginlog"
)

// GinLogCreate is the builder for creating a GinLog entity.
type GinLogCreate struct {
	config
	mutation *GinLogMutation
	hooks    []Hook
}

// SetRequestTime sets the "request_time" field.
func (glc *GinLogCreate) SetRequestTime(t time.Time) *GinLogCreate {
	glc.mutation.SetRequestTime(t)
	return glc
}

// SetNillableRequestTime sets the "request_time" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableRequestTime(t *time.Time) *GinLogCreate {
	if t != nil {
		glc.SetRequestTime(*t)
	}
	return glc
}

// SetRequestHost sets the "request_host" field.
func (glc *GinLogCreate) SetRequestHost(s string) *GinLogCreate {
	glc.mutation.SetRequestHost(s)
	return glc
}

// SetNillableRequestHost sets the "request_host" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableRequestHost(s *string) *GinLogCreate {
	if s != nil {
		glc.SetRequestHost(*s)
	}
	return glc
}

// SetRequestPath sets the "request_path" field.
func (glc *GinLogCreate) SetRequestPath(s string) *GinLogCreate {
	glc.mutation.SetRequestPath(s)
	return glc
}

// SetNillableRequestPath sets the "request_path" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableRequestPath(s *string) *GinLogCreate {
	if s != nil {
		glc.SetRequestPath(*s)
	}
	return glc
}

// SetRequestQuery sets the "request_query" field.
func (glc *GinLogCreate) SetRequestQuery(s string) *GinLogCreate {
	glc.mutation.SetRequestQuery(s)
	return glc
}

// SetNillableRequestQuery sets the "request_query" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableRequestQuery(s *string) *GinLogCreate {
	if s != nil {
		glc.SetRequestQuery(*s)
	}
	return glc
}

// SetRequestMethod sets the "request_method" field.
func (glc *GinLogCreate) SetRequestMethod(s string) *GinLogCreate {
	glc.mutation.SetRequestMethod(s)
	return glc
}

// SetNillableRequestMethod sets the "request_method" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableRequestMethod(s *string) *GinLogCreate {
	if s != nil {
		glc.SetRequestMethod(*s)
	}
	return glc
}

// SetRequestScheme sets the "request_scheme" field.
func (glc *GinLogCreate) SetRequestScheme(s string) *GinLogCreate {
	glc.mutation.SetRequestScheme(s)
	return glc
}

// SetNillableRequestScheme sets the "request_scheme" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableRequestScheme(s *string) *GinLogCreate {
	if s != nil {
		glc.SetRequestScheme(*s)
	}
	return glc
}

// SetRequestContentType sets the "request_content_type" field.
func (glc *GinLogCreate) SetRequestContentType(s string) *GinLogCreate {
	glc.mutation.SetRequestContentType(s)
	return glc
}

// SetNillableRequestContentType sets the "request_content_type" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableRequestContentType(s *string) *GinLogCreate {
	if s != nil {
		glc.SetRequestContentType(*s)
	}
	return glc
}

// SetRequestBody sets the "request_body" field.
func (glc *GinLogCreate) SetRequestBody(s string) *GinLogCreate {
	glc.mutation.SetRequestBody(s)
	return glc
}

// SetNillableRequestBody sets the "request_body" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableRequestBody(s *string) *GinLogCreate {
	if s != nil {
		glc.SetRequestBody(*s)
	}
	return glc
}

// SetRequestClientIP sets the "request_client_ip" field.
func (glc *GinLogCreate) SetRequestClientIP(s string) *GinLogCreate {
	glc.mutation.SetRequestClientIP(s)
	return glc
}

// SetNillableRequestClientIP sets the "request_client_ip" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableRequestClientIP(s *string) *GinLogCreate {
	if s != nil {
		glc.SetRequestClientIP(*s)
	}
	return glc
}

// SetRequestUserAgent sets the "request_user_agent" field.
func (glc *GinLogCreate) SetRequestUserAgent(s string) *GinLogCreate {
	glc.mutation.SetRequestUserAgent(s)
	return glc
}

// SetNillableRequestUserAgent sets the "request_user_agent" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableRequestUserAgent(s *string) *GinLogCreate {
	if s != nil {
		glc.SetRequestUserAgent(*s)
	}
	return glc
}

// SetRequestHeader sets the "request_header" field.
func (glc *GinLogCreate) SetRequestHeader(s string) *GinLogCreate {
	glc.mutation.SetRequestHeader(s)
	return glc
}

// SetNillableRequestHeader sets the "request_header" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableRequestHeader(s *string) *GinLogCreate {
	if s != nil {
		glc.SetRequestHeader(*s)
	}
	return glc
}

// SetRequestCostTime sets the "request_cost_time" field.
func (glc *GinLogCreate) SetRequestCostTime(i int64) *GinLogCreate {
	glc.mutation.SetRequestCostTime(i)
	return glc
}

// SetNillableRequestCostTime sets the "request_cost_time" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableRequestCostTime(i *int64) *GinLogCreate {
	if i != nil {
		glc.SetRequestCostTime(*i)
	}
	return glc
}

// SetResponseTime sets the "response_time" field.
func (glc *GinLogCreate) SetResponseTime(t time.Time) *GinLogCreate {
	glc.mutation.SetResponseTime(t)
	return glc
}

// SetNillableResponseTime sets the "response_time" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableResponseTime(t *time.Time) *GinLogCreate {
	if t != nil {
		glc.SetResponseTime(*t)
	}
	return glc
}

// SetResponseHeader sets the "response_header" field.
func (glc *GinLogCreate) SetResponseHeader(s string) *GinLogCreate {
	glc.mutation.SetResponseHeader(s)
	return glc
}

// SetNillableResponseHeader sets the "response_header" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableResponseHeader(s *string) *GinLogCreate {
	if s != nil {
		glc.SetResponseHeader(*s)
	}
	return glc
}

// SetResponseStatusCode sets the "response_status_code" field.
func (glc *GinLogCreate) SetResponseStatusCode(i int) *GinLogCreate {
	glc.mutation.SetResponseStatusCode(i)
	return glc
}

// SetNillableResponseStatusCode sets the "response_status_code" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableResponseStatusCode(i *int) *GinLogCreate {
	if i != nil {
		glc.SetResponseStatusCode(*i)
	}
	return glc
}

// SetResponseData sets the "response_data" field.
func (glc *GinLogCreate) SetResponseData(s string) *GinLogCreate {
	glc.mutation.SetResponseData(s)
	return glc
}

// SetNillableResponseData sets the "response_data" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableResponseData(s *string) *GinLogCreate {
	if s != nil {
		glc.SetResponseData(*s)
	}
	return glc
}

// SetGoVersion sets the "go_version" field.
func (glc *GinLogCreate) SetGoVersion(s string) *GinLogCreate {
	glc.mutation.SetGoVersion(s)
	return glc
}

// SetNillableGoVersion sets the "go_version" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableGoVersion(s *string) *GinLogCreate {
	if s != nil {
		glc.SetGoVersion(*s)
	}
	return glc
}

// SetSdkVersion sets the "sdk_version" field.
func (glc *GinLogCreate) SetSdkVersion(s string) *GinLogCreate {
	glc.mutation.SetSdkVersion(s)
	return glc
}

// SetNillableSdkVersion sets the "sdk_version" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableSdkVersion(s *string) *GinLogCreate {
	if s != nil {
		glc.SetSdkVersion(*s)
	}
	return glc
}

// SetSystemInfo sets the "system_info" field.
func (glc *GinLogCreate) SetSystemInfo(s string) *GinLogCreate {
	glc.mutation.SetSystemInfo(s)
	return glc
}

// SetNillableSystemInfo sets the "system_info" field if the given value is not nil.
func (glc *GinLogCreate) SetNillableSystemInfo(s *string) *GinLogCreate {
	if s != nil {
		glc.SetSystemInfo(*s)
	}
	return glc
}

// SetID sets the "id" field.
func (glc *GinLogCreate) SetID(i int64) *GinLogCreate {
	glc.mutation.SetID(i)
	return glc
}

// Mutation returns the GinLogMutation object of the builder.
func (glc *GinLogCreate) Mutation() *GinLogMutation {
	return glc.mutation
}

// Save creates the GinLog in the database.
func (glc *GinLogCreate) Save(ctx context.Context) (*GinLog, error) {
	return withHooks(ctx, glc.sqlSave, glc.mutation, glc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (glc *GinLogCreate) SaveX(ctx context.Context) *GinLog {
	v, err := glc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (glc *GinLogCreate) Exec(ctx context.Context) error {
	_, err := glc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (glc *GinLogCreate) ExecX(ctx context.Context) {
	if err := glc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (glc *GinLogCreate) check() error {
	return nil
}

func (glc *GinLogCreate) sqlSave(ctx context.Context) (*GinLog, error) {
	if err := glc.check(); err != nil {
		return nil, err
	}
	_node, _spec := glc.createSpec()
	if err := sqlgraph.CreateNode(ctx, glc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int64(id)
	}
	glc.mutation.id = &_node.ID
	glc.mutation.done = true
	return _node, nil
}

func (glc *GinLogCreate) createSpec() (*GinLog, *sqlgraph.CreateSpec) {
	var (
		_node = &GinLog{config: glc.config}
		_spec = sqlgraph.NewCreateSpec(ginlog.Table, sqlgraph.NewFieldSpec(ginlog.FieldID, field.TypeInt64))
	)
	if id, ok := glc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := glc.mutation.RequestTime(); ok {
		_spec.SetField(ginlog.FieldRequestTime, field.TypeTime, value)
		_node.RequestTime = value
	}
	if value, ok := glc.mutation.RequestHost(); ok {
		_spec.SetField(ginlog.FieldRequestHost, field.TypeString, value)
		_node.RequestHost = value
	}
	if value, ok := glc.mutation.RequestPath(); ok {
		_spec.SetField(ginlog.FieldRequestPath, field.TypeString, value)
		_node.RequestPath = value
	}
	if value, ok := glc.mutation.RequestQuery(); ok {
		_spec.SetField(ginlog.FieldRequestQuery, field.TypeString, value)
		_node.RequestQuery = value
	}
	if value, ok := glc.mutation.RequestMethod(); ok {
		_spec.SetField(ginlog.FieldRequestMethod, field.TypeString, value)
		_node.RequestMethod = value
	}
	if value, ok := glc.mutation.RequestScheme(); ok {
		_spec.SetField(ginlog.FieldRequestScheme, field.TypeString, value)
		_node.RequestScheme = value
	}
	if value, ok := glc.mutation.RequestContentType(); ok {
		_spec.SetField(ginlog.FieldRequestContentType, field.TypeString, value)
		_node.RequestContentType = value
	}
	if value, ok := glc.mutation.RequestBody(); ok {
		_spec.SetField(ginlog.FieldRequestBody, field.TypeString, value)
		_node.RequestBody = value
	}
	if value, ok := glc.mutation.RequestClientIP(); ok {
		_spec.SetField(ginlog.FieldRequestClientIP, field.TypeString, value)
		_node.RequestClientIP = value
	}
	if value, ok := glc.mutation.RequestUserAgent(); ok {
		_spec.SetField(ginlog.FieldRequestUserAgent, field.TypeString, value)
		_node.RequestUserAgent = value
	}
	if value, ok := glc.mutation.RequestHeader(); ok {
		_spec.SetField(ginlog.FieldRequestHeader, field.TypeString, value)
		_node.RequestHeader = value
	}
	if value, ok := glc.mutation.RequestCostTime(); ok {
		_spec.SetField(ginlog.FieldRequestCostTime, field.TypeInt64, value)
		_node.RequestCostTime = value
	}
	if value, ok := glc.mutation.ResponseTime(); ok {
		_spec.SetField(ginlog.FieldResponseTime, field.TypeTime, value)
		_node.ResponseTime = value
	}
	if value, ok := glc.mutation.ResponseHeader(); ok {
		_spec.SetField(ginlog.FieldResponseHeader, field.TypeString, value)
		_node.ResponseHeader = value
	}
	if value, ok := glc.mutation.ResponseStatusCode(); ok {
		_spec.SetField(ginlog.FieldResponseStatusCode, field.TypeInt, value)
		_node.ResponseStatusCode = value
	}
	if value, ok := glc.mutation.ResponseData(); ok {
		_spec.SetField(ginlog.FieldResponseData, field.TypeString, value)
		_node.ResponseData = value
	}
	if value, ok := glc.mutation.GoVersion(); ok {
		_spec.SetField(ginlog.FieldGoVersion, field.TypeString, value)
		_node.GoVersion = value
	}
	if value, ok := glc.mutation.SdkVersion(); ok {
		_spec.SetField(ginlog.FieldSdkVersion, field.TypeString, value)
		_node.SdkVersion = value
	}
	if value, ok := glc.mutation.SystemInfo(); ok {
		_spec.SetField(ginlog.FieldSystemInfo, field.TypeString, value)
		_node.SystemInfo = value
	}
	return _node, _spec
}

// GinLogCreateBulk is the builder for creating many GinLog entities in bulk.
type GinLogCreateBulk struct {
	config
	err      error
	builders []*GinLogCreate
}

// Save creates the GinLog entities in the database.
func (glcb *GinLogCreateBulk) Save(ctx context.Context) ([]*GinLog, error) {
	if glcb.err != nil {
		return nil, glcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(glcb.builders))
	nodes := make([]*GinLog, len(glcb.builders))
	mutators := make([]Mutator, len(glcb.builders))
	for i := range glcb.builders {
		func(i int, root context.Context) {
			builder := glcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*GinLogMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, glcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, glcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int64(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, glcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (glcb *GinLogCreateBulk) SaveX(ctx context.Context) []*GinLog {
	v, err := glcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (glcb *GinLogCreateBulk) Exec(ctx context.Context) error {
	_, err := glcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (glcb *GinLogCreateBulk) ExecX(ctx context.Context) {
	if err := glcb.Exec(ctx); err != nil {
		panic(err)
	}
}
