// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"go.dtapp.net/golog/internal/ent/ginlog"
	"go.dtapp.net/golog/internal/ent/predicate"
)

// GinLogQuery is the builder for querying GinLog entities.
type GinLogQuery struct {
	config
	ctx        *QueryContext
	order      []ginlog.OrderOption
	inters     []Interceptor
	predicates []predicate.GinLog
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the GinLogQuery builder.
func (glq *GinLogQuery) Where(ps ...predicate.GinLog) *GinLogQuery {
	glq.predicates = append(glq.predicates, ps...)
	return glq
}

// Limit the number of records to be returned by this query.
func (glq *GinLogQuery) Limit(limit int) *GinLogQuery {
	glq.ctx.Limit = &limit
	return glq
}

// Offset to start from.
func (glq *GinLogQuery) Offset(offset int) *GinLogQuery {
	glq.ctx.Offset = &offset
	return glq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (glq *GinLogQuery) Unique(unique bool) *GinLogQuery {
	glq.ctx.Unique = &unique
	return glq
}

// Order specifies how the records should be ordered.
func (glq *GinLogQuery) Order(o ...ginlog.OrderOption) *GinLogQuery {
	glq.order = append(glq.order, o...)
	return glq
}

// First returns the first GinLog entity from the query.
// Returns a *NotFoundError when no GinLog was found.
func (glq *GinLogQuery) First(ctx context.Context) (*GinLog, error) {
	nodes, err := glq.Limit(1).All(setContextOp(ctx, glq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{ginlog.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (glq *GinLogQuery) FirstX(ctx context.Context) *GinLog {
	node, err := glq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first GinLog ID from the query.
// Returns a *NotFoundError when no GinLog ID was found.
func (glq *GinLogQuery) FirstID(ctx context.Context) (id int64, err error) {
	var ids []int64
	if ids, err = glq.Limit(1).IDs(setContextOp(ctx, glq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{ginlog.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (glq *GinLogQuery) FirstIDX(ctx context.Context) int64 {
	id, err := glq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single GinLog entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one GinLog entity is found.
// Returns a *NotFoundError when no GinLog entities are found.
func (glq *GinLogQuery) Only(ctx context.Context) (*GinLog, error) {
	nodes, err := glq.Limit(2).All(setContextOp(ctx, glq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{ginlog.Label}
	default:
		return nil, &NotSingularError{ginlog.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (glq *GinLogQuery) OnlyX(ctx context.Context) *GinLog {
	node, err := glq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only GinLog ID in the query.
// Returns a *NotSingularError when more than one GinLog ID is found.
// Returns a *NotFoundError when no entities are found.
func (glq *GinLogQuery) OnlyID(ctx context.Context) (id int64, err error) {
	var ids []int64
	if ids, err = glq.Limit(2).IDs(setContextOp(ctx, glq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{ginlog.Label}
	default:
		err = &NotSingularError{ginlog.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (glq *GinLogQuery) OnlyIDX(ctx context.Context) int64 {
	id, err := glq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of GinLogs.
func (glq *GinLogQuery) All(ctx context.Context) ([]*GinLog, error) {
	ctx = setContextOp(ctx, glq.ctx, "All")
	if err := glq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*GinLog, *GinLogQuery]()
	return withInterceptors[[]*GinLog](ctx, glq, qr, glq.inters)
}

// AllX is like All, but panics if an error occurs.
func (glq *GinLogQuery) AllX(ctx context.Context) []*GinLog {
	nodes, err := glq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of GinLog IDs.
func (glq *GinLogQuery) IDs(ctx context.Context) (ids []int64, err error) {
	if glq.ctx.Unique == nil && glq.path != nil {
		glq.Unique(true)
	}
	ctx = setContextOp(ctx, glq.ctx, "IDs")
	if err = glq.Select(ginlog.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (glq *GinLogQuery) IDsX(ctx context.Context) []int64 {
	ids, err := glq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (glq *GinLogQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, glq.ctx, "Count")
	if err := glq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, glq, querierCount[*GinLogQuery](), glq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (glq *GinLogQuery) CountX(ctx context.Context) int {
	count, err := glq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (glq *GinLogQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, glq.ctx, "Exist")
	switch _, err := glq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (glq *GinLogQuery) ExistX(ctx context.Context) bool {
	exist, err := glq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the GinLogQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (glq *GinLogQuery) Clone() *GinLogQuery {
	if glq == nil {
		return nil
	}
	return &GinLogQuery{
		config:     glq.config,
		ctx:        glq.ctx.Clone(),
		order:      append([]ginlog.OrderOption{}, glq.order...),
		inters:     append([]Interceptor{}, glq.inters...),
		predicates: append([]predicate.GinLog{}, glq.predicates...),
		// clone intermediate query.
		sql:  glq.sql.Clone(),
		path: glq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		RequestTime time.Time `json:"request_time,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.GinLog.Query().
//		GroupBy(ginlog.FieldRequestTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (glq *GinLogQuery) GroupBy(field string, fields ...string) *GinLogGroupBy {
	glq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &GinLogGroupBy{build: glq}
	grbuild.flds = &glq.ctx.Fields
	grbuild.label = ginlog.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		RequestTime time.Time `json:"request_time,omitempty"`
//	}
//
//	client.GinLog.Query().
//		Select(ginlog.FieldRequestTime).
//		Scan(ctx, &v)
func (glq *GinLogQuery) Select(fields ...string) *GinLogSelect {
	glq.ctx.Fields = append(glq.ctx.Fields, fields...)
	sbuild := &GinLogSelect{GinLogQuery: glq}
	sbuild.label = ginlog.Label
	sbuild.flds, sbuild.scan = &glq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a GinLogSelect configured with the given aggregations.
func (glq *GinLogQuery) Aggregate(fns ...AggregateFunc) *GinLogSelect {
	return glq.Select().Aggregate(fns...)
}

func (glq *GinLogQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range glq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, glq); err != nil {
				return err
			}
		}
	}
	for _, f := range glq.ctx.Fields {
		if !ginlog.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if glq.path != nil {
		prev, err := glq.path(ctx)
		if err != nil {
			return err
		}
		glq.sql = prev
	}
	return nil
}

func (glq *GinLogQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*GinLog, error) {
	var (
		nodes = []*GinLog{}
		_spec = glq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*GinLog).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &GinLog{config: glq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, glq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (glq *GinLogQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := glq.querySpec()
	_spec.Node.Columns = glq.ctx.Fields
	if len(glq.ctx.Fields) > 0 {
		_spec.Unique = glq.ctx.Unique != nil && *glq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, glq.driver, _spec)
}

func (glq *GinLogQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(ginlog.Table, ginlog.Columns, sqlgraph.NewFieldSpec(ginlog.FieldID, field.TypeInt64))
	_spec.From = glq.sql
	if unique := glq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if glq.path != nil {
		_spec.Unique = true
	}
	if fields := glq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, ginlog.FieldID)
		for i := range fields {
			if fields[i] != ginlog.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := glq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := glq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := glq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := glq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (glq *GinLogQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(glq.driver.Dialect())
	t1 := builder.Table(ginlog.Table)
	columns := glq.ctx.Fields
	if len(columns) == 0 {
		columns = ginlog.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if glq.sql != nil {
		selector = glq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if glq.ctx.Unique != nil && *glq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range glq.predicates {
		p(selector)
	}
	for _, p := range glq.order {
		p(selector)
	}
	if offset := glq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := glq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// GinLogGroupBy is the group-by builder for GinLog entities.
type GinLogGroupBy struct {
	selector
	build *GinLogQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (glgb *GinLogGroupBy) Aggregate(fns ...AggregateFunc) *GinLogGroupBy {
	glgb.fns = append(glgb.fns, fns...)
	return glgb
}

// Scan applies the selector query and scans the result into the given value.
func (glgb *GinLogGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, glgb.build.ctx, "GroupBy")
	if err := glgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*GinLogQuery, *GinLogGroupBy](ctx, glgb.build, glgb, glgb.build.inters, v)
}

func (glgb *GinLogGroupBy) sqlScan(ctx context.Context, root *GinLogQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(glgb.fns))
	for _, fn := range glgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*glgb.flds)+len(glgb.fns))
		for _, f := range *glgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*glgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := glgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// GinLogSelect is the builder for selecting fields of GinLog entities.
type GinLogSelect struct {
	*GinLogQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (gls *GinLogSelect) Aggregate(fns ...AggregateFunc) *GinLogSelect {
	gls.fns = append(gls.fns, fns...)
	return gls
}

// Scan applies the selector query and scans the result into the given value.
func (gls *GinLogSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, gls.ctx, "Select")
	if err := gls.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*GinLogQuery, *GinLogSelect](ctx, gls.GinLogQuery, gls, gls.inters, v)
}

func (gls *GinLogSelect) sqlScan(ctx context.Context, root *GinLogQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(gls.fns))
	for _, fn := range gls.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*gls.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := gls.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
