// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"go.dtapp.net/golog/internal/ent/hertzlog"
	"go.dtapp.net/golog/internal/ent/predicate"
)

// HertzlogQuery is the builder for querying Hertzlog entities.
type HertzlogQuery struct {
	config
	ctx        *QueryContext
	order      []hertzlog.OrderOption
	inters     []Interceptor
	predicates []predicate.Hertzlog
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the HertzlogQuery builder.
func (hq *HertzlogQuery) Where(ps ...predicate.Hertzlog) *HertzlogQuery {
	hq.predicates = append(hq.predicates, ps...)
	return hq
}

// Limit the number of records to be returned by this query.
func (hq *HertzlogQuery) Limit(limit int) *HertzlogQuery {
	hq.ctx.Limit = &limit
	return hq
}

// Offset to start from.
func (hq *HertzlogQuery) Offset(offset int) *HertzlogQuery {
	hq.ctx.Offset = &offset
	return hq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (hq *HertzlogQuery) Unique(unique bool) *HertzlogQuery {
	hq.ctx.Unique = &unique
	return hq
}

// Order specifies how the records should be ordered.
func (hq *HertzlogQuery) Order(o ...hertzlog.OrderOption) *HertzlogQuery {
	hq.order = append(hq.order, o...)
	return hq
}

// First returns the first Hertzlog entity from the query.
// Returns a *NotFoundError when no Hertzlog was found.
func (hq *HertzlogQuery) First(ctx context.Context) (*Hertzlog, error) {
	nodes, err := hq.Limit(1).All(setContextOp(ctx, hq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{hertzlog.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (hq *HertzlogQuery) FirstX(ctx context.Context) *Hertzlog {
	node, err := hq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Hertzlog ID from the query.
// Returns a *NotFoundError when no Hertzlog ID was found.
func (hq *HertzlogQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = hq.Limit(1).IDs(setContextOp(ctx, hq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{hertzlog.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (hq *HertzlogQuery) FirstIDX(ctx context.Context) int {
	id, err := hq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Hertzlog entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Hertzlog entity is found.
// Returns a *NotFoundError when no Hertzlog entities are found.
func (hq *HertzlogQuery) Only(ctx context.Context) (*Hertzlog, error) {
	nodes, err := hq.Limit(2).All(setContextOp(ctx, hq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{hertzlog.Label}
	default:
		return nil, &NotSingularError{hertzlog.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (hq *HertzlogQuery) OnlyX(ctx context.Context) *Hertzlog {
	node, err := hq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Hertzlog ID in the query.
// Returns a *NotSingularError when more than one Hertzlog ID is found.
// Returns a *NotFoundError when no entities are found.
func (hq *HertzlogQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = hq.Limit(2).IDs(setContextOp(ctx, hq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{hertzlog.Label}
	default:
		err = &NotSingularError{hertzlog.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (hq *HertzlogQuery) OnlyIDX(ctx context.Context) int {
	id, err := hq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Hertzlogs.
func (hq *HertzlogQuery) All(ctx context.Context) ([]*Hertzlog, error) {
	ctx = setContextOp(ctx, hq.ctx, "All")
	if err := hq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Hertzlog, *HertzlogQuery]()
	return withInterceptors[[]*Hertzlog](ctx, hq, qr, hq.inters)
}

// AllX is like All, but panics if an error occurs.
func (hq *HertzlogQuery) AllX(ctx context.Context) []*Hertzlog {
	nodes, err := hq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Hertzlog IDs.
func (hq *HertzlogQuery) IDs(ctx context.Context) (ids []int, err error) {
	if hq.ctx.Unique == nil && hq.path != nil {
		hq.Unique(true)
	}
	ctx = setContextOp(ctx, hq.ctx, "IDs")
	if err = hq.Select(hertzlog.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (hq *HertzlogQuery) IDsX(ctx context.Context) []int {
	ids, err := hq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (hq *HertzlogQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, hq.ctx, "Count")
	if err := hq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, hq, querierCount[*HertzlogQuery](), hq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (hq *HertzlogQuery) CountX(ctx context.Context) int {
	count, err := hq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (hq *HertzlogQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, hq.ctx, "Exist")
	switch _, err := hq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (hq *HertzlogQuery) ExistX(ctx context.Context) bool {
	exist, err := hq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the HertzlogQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (hq *HertzlogQuery) Clone() *HertzlogQuery {
	if hq == nil {
		return nil
	}
	return &HertzlogQuery{
		config:     hq.config,
		ctx:        hq.ctx.Clone(),
		order:      append([]hertzlog.OrderOption{}, hq.order...),
		inters:     append([]Interceptor{}, hq.inters...),
		predicates: append([]predicate.Hertzlog{}, hq.predicates...),
		// clone intermediate query.
		sql:  hq.sql.Clone(),
		path: hq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
func (hq *HertzlogQuery) GroupBy(field string, fields ...string) *HertzlogGroupBy {
	hq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &HertzlogGroupBy{build: hq}
	grbuild.flds = &hq.ctx.Fields
	grbuild.label = hertzlog.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
func (hq *HertzlogQuery) Select(fields ...string) *HertzlogSelect {
	hq.ctx.Fields = append(hq.ctx.Fields, fields...)
	sbuild := &HertzlogSelect{HertzlogQuery: hq}
	sbuild.label = hertzlog.Label
	sbuild.flds, sbuild.scan = &hq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a HertzlogSelect configured with the given aggregations.
func (hq *HertzlogQuery) Aggregate(fns ...AggregateFunc) *HertzlogSelect {
	return hq.Select().Aggregate(fns...)
}

func (hq *HertzlogQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range hq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, hq); err != nil {
				return err
			}
		}
	}
	for _, f := range hq.ctx.Fields {
		if !hertzlog.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if hq.path != nil {
		prev, err := hq.path(ctx)
		if err != nil {
			return err
		}
		hq.sql = prev
	}
	return nil
}

func (hq *HertzlogQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Hertzlog, error) {
	var (
		nodes = []*Hertzlog{}
		_spec = hq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Hertzlog).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Hertzlog{config: hq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, hq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (hq *HertzlogQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := hq.querySpec()
	_spec.Node.Columns = hq.ctx.Fields
	if len(hq.ctx.Fields) > 0 {
		_spec.Unique = hq.ctx.Unique != nil && *hq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, hq.driver, _spec)
}

func (hq *HertzlogQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(hertzlog.Table, hertzlog.Columns, sqlgraph.NewFieldSpec(hertzlog.FieldID, field.TypeInt))
	_spec.From = hq.sql
	if unique := hq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if hq.path != nil {
		_spec.Unique = true
	}
	if fields := hq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, hertzlog.FieldID)
		for i := range fields {
			if fields[i] != hertzlog.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := hq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := hq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := hq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := hq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (hq *HertzlogQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(hq.driver.Dialect())
	t1 := builder.Table(hertzlog.Table)
	columns := hq.ctx.Fields
	if len(columns) == 0 {
		columns = hertzlog.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if hq.sql != nil {
		selector = hq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if hq.ctx.Unique != nil && *hq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range hq.predicates {
		p(selector)
	}
	for _, p := range hq.order {
		p(selector)
	}
	if offset := hq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := hq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// HertzlogGroupBy is the group-by builder for Hertzlog entities.
type HertzlogGroupBy struct {
	selector
	build *HertzlogQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (hgb *HertzlogGroupBy) Aggregate(fns ...AggregateFunc) *HertzlogGroupBy {
	hgb.fns = append(hgb.fns, fns...)
	return hgb
}

// Scan applies the selector query and scans the result into the given value.
func (hgb *HertzlogGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, hgb.build.ctx, "GroupBy")
	if err := hgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*HertzlogQuery, *HertzlogGroupBy](ctx, hgb.build, hgb, hgb.build.inters, v)
}

func (hgb *HertzlogGroupBy) sqlScan(ctx context.Context, root *HertzlogQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(hgb.fns))
	for _, fn := range hgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*hgb.flds)+len(hgb.fns))
		for _, f := range *hgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*hgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := hgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// HertzlogSelect is the builder for selecting fields of Hertzlog entities.
type HertzlogSelect struct {
	*HertzlogQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (hs *HertzlogSelect) Aggregate(fns ...AggregateFunc) *HertzlogSelect {
	hs.fns = append(hs.fns, fns...)
	return hs
}

// Scan applies the selector query and scans the result into the given value.
func (hs *HertzlogSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, hs.ctx, "Select")
	if err := hs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*HertzlogQuery, *HertzlogSelect](ctx, hs.HertzlogQuery, hs, hs.inters, v)
}

func (hs *HertzlogSelect) sqlScan(ctx context.Context, root *HertzlogQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(hs.fns))
	for _, fn := range hs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*hs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := hs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
