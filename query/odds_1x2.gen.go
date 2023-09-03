// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"sports-book.com/model"
)

func newOdds1x2(db *gorm.DB, opts ...gen.DOOption) odds1x2 {
	_odds1x2 := odds1x2{}

	_odds1x2.odds1x2Do.UseDB(db, opts...)
	_odds1x2.odds1x2Do.UseModel(&model.Odds1x2{})

	tableName := _odds1x2.odds1x2Do.TableName()
	_odds1x2.ALL = field.NewAsterisk(tableName)
	_odds1x2.ID = field.NewInt32(tableName, "id")
	_odds1x2.Bookmaker = field.NewString(tableName, "bookmaker")
	_odds1x2.Match = field.NewInt32(tableName, "match")
	_odds1x2.HomeWin = field.NewFloat64(tableName, "home_win")
	_odds1x2.Draw = field.NewFloat64(tableName, "draw")
	_odds1x2.AwayWin = field.NewFloat64(tableName, "away_win")

	_odds1x2.fillFieldMap()

	return _odds1x2
}

type odds1x2 struct {
	odds1x2Do

	ALL       field.Asterisk
	ID        field.Int32
	Bookmaker field.String
	Match     field.Int32
	HomeWin   field.Float64
	Draw      field.Float64
	AwayWin   field.Float64

	fieldMap map[string]field.Expr
}

func (o odds1x2) Table(newTableName string) *odds1x2 {
	o.odds1x2Do.UseTable(newTableName)
	return o.updateTableName(newTableName)
}

func (o odds1x2) As(alias string) *odds1x2 {
	o.odds1x2Do.DO = *(o.odds1x2Do.As(alias).(*gen.DO))
	return o.updateTableName(alias)
}

func (o *odds1x2) updateTableName(table string) *odds1x2 {
	o.ALL = field.NewAsterisk(table)
	o.ID = field.NewInt32(table, "id")
	o.Bookmaker = field.NewString(table, "bookmaker")
	o.Match = field.NewInt32(table, "match")
	o.HomeWin = field.NewFloat64(table, "home_win")
	o.Draw = field.NewFloat64(table, "draw")
	o.AwayWin = field.NewFloat64(table, "away_win")

	o.fillFieldMap()

	return o
}

func (o *odds1x2) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := o.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (o *odds1x2) fillFieldMap() {
	o.fieldMap = make(map[string]field.Expr, 6)
	o.fieldMap["id"] = o.ID
	o.fieldMap["bookmaker"] = o.Bookmaker
	o.fieldMap["match"] = o.Match
	o.fieldMap["home_win"] = o.HomeWin
	o.fieldMap["draw"] = o.Draw
	o.fieldMap["away_win"] = o.AwayWin
}

func (o odds1x2) clone(db *gorm.DB) odds1x2 {
	o.odds1x2Do.ReplaceConnPool(db.Statement.ConnPool)
	return o
}

func (o odds1x2) replaceDB(db *gorm.DB) odds1x2 {
	o.odds1x2Do.ReplaceDB(db)
	return o
}

type odds1x2Do struct{ gen.DO }

type IOdds1x2Do interface {
	gen.SubQuery
	Debug() IOdds1x2Do
	WithContext(ctx context.Context) IOdds1x2Do
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IOdds1x2Do
	WriteDB() IOdds1x2Do
	As(alias string) gen.Dao
	Session(config *gorm.Session) IOdds1x2Do
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IOdds1x2Do
	Not(conds ...gen.Condition) IOdds1x2Do
	Or(conds ...gen.Condition) IOdds1x2Do
	Select(conds ...field.Expr) IOdds1x2Do
	Where(conds ...gen.Condition) IOdds1x2Do
	Order(conds ...field.Expr) IOdds1x2Do
	Distinct(cols ...field.Expr) IOdds1x2Do
	Omit(cols ...field.Expr) IOdds1x2Do
	Join(table schema.Tabler, on ...field.Expr) IOdds1x2Do
	LeftJoin(table schema.Tabler, on ...field.Expr) IOdds1x2Do
	RightJoin(table schema.Tabler, on ...field.Expr) IOdds1x2Do
	Group(cols ...field.Expr) IOdds1x2Do
	Having(conds ...gen.Condition) IOdds1x2Do
	Limit(limit int) IOdds1x2Do
	Offset(offset int) IOdds1x2Do
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IOdds1x2Do
	Unscoped() IOdds1x2Do
	Create(values ...*model.Odds1x2) error
	CreateInBatches(values []*model.Odds1x2, batchSize int) error
	Save(values ...*model.Odds1x2) error
	First() (*model.Odds1x2, error)
	Take() (*model.Odds1x2, error)
	Last() (*model.Odds1x2, error)
	Find() ([]*model.Odds1x2, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Odds1x2, err error)
	FindInBatches(result *[]*model.Odds1x2, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Odds1x2) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IOdds1x2Do
	Assign(attrs ...field.AssignExpr) IOdds1x2Do
	Joins(fields ...field.RelationField) IOdds1x2Do
	Preload(fields ...field.RelationField) IOdds1x2Do
	FirstOrInit() (*model.Odds1x2, error)
	FirstOrCreate() (*model.Odds1x2, error)
	FindByPage(offset int, limit int) (result []*model.Odds1x2, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IOdds1x2Do
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (o odds1x2Do) Debug() IOdds1x2Do {
	return o.withDO(o.DO.Debug())
}

func (o odds1x2Do) WithContext(ctx context.Context) IOdds1x2Do {
	return o.withDO(o.DO.WithContext(ctx))
}

func (o odds1x2Do) ReadDB() IOdds1x2Do {
	return o.Clauses(dbresolver.Read)
}

func (o odds1x2Do) WriteDB() IOdds1x2Do {
	return o.Clauses(dbresolver.Write)
}

func (o odds1x2Do) Session(config *gorm.Session) IOdds1x2Do {
	return o.withDO(o.DO.Session(config))
}

func (o odds1x2Do) Clauses(conds ...clause.Expression) IOdds1x2Do {
	return o.withDO(o.DO.Clauses(conds...))
}

func (o odds1x2Do) Returning(value interface{}, columns ...string) IOdds1x2Do {
	return o.withDO(o.DO.Returning(value, columns...))
}

func (o odds1x2Do) Not(conds ...gen.Condition) IOdds1x2Do {
	return o.withDO(o.DO.Not(conds...))
}

func (o odds1x2Do) Or(conds ...gen.Condition) IOdds1x2Do {
	return o.withDO(o.DO.Or(conds...))
}

func (o odds1x2Do) Select(conds ...field.Expr) IOdds1x2Do {
	return o.withDO(o.DO.Select(conds...))
}

func (o odds1x2Do) Where(conds ...gen.Condition) IOdds1x2Do {
	return o.withDO(o.DO.Where(conds...))
}

func (o odds1x2Do) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IOdds1x2Do {
	return o.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (o odds1x2Do) Order(conds ...field.Expr) IOdds1x2Do {
	return o.withDO(o.DO.Order(conds...))
}

func (o odds1x2Do) Distinct(cols ...field.Expr) IOdds1x2Do {
	return o.withDO(o.DO.Distinct(cols...))
}

func (o odds1x2Do) Omit(cols ...field.Expr) IOdds1x2Do {
	return o.withDO(o.DO.Omit(cols...))
}

func (o odds1x2Do) Join(table schema.Tabler, on ...field.Expr) IOdds1x2Do {
	return o.withDO(o.DO.Join(table, on...))
}

func (o odds1x2Do) LeftJoin(table schema.Tabler, on ...field.Expr) IOdds1x2Do {
	return o.withDO(o.DO.LeftJoin(table, on...))
}

func (o odds1x2Do) RightJoin(table schema.Tabler, on ...field.Expr) IOdds1x2Do {
	return o.withDO(o.DO.RightJoin(table, on...))
}

func (o odds1x2Do) Group(cols ...field.Expr) IOdds1x2Do {
	return o.withDO(o.DO.Group(cols...))
}

func (o odds1x2Do) Having(conds ...gen.Condition) IOdds1x2Do {
	return o.withDO(o.DO.Having(conds...))
}

func (o odds1x2Do) Limit(limit int) IOdds1x2Do {
	return o.withDO(o.DO.Limit(limit))
}

func (o odds1x2Do) Offset(offset int) IOdds1x2Do {
	return o.withDO(o.DO.Offset(offset))
}

func (o odds1x2Do) Scopes(funcs ...func(gen.Dao) gen.Dao) IOdds1x2Do {
	return o.withDO(o.DO.Scopes(funcs...))
}

func (o odds1x2Do) Unscoped() IOdds1x2Do {
	return o.withDO(o.DO.Unscoped())
}

func (o odds1x2Do) Create(values ...*model.Odds1x2) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Create(values)
}

func (o odds1x2Do) CreateInBatches(values []*model.Odds1x2, batchSize int) error {
	return o.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (o odds1x2Do) Save(values ...*model.Odds1x2) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Save(values)
}

func (o odds1x2Do) First() (*model.Odds1x2, error) {
	if result, err := o.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Odds1x2), nil
	}
}

func (o odds1x2Do) Take() (*model.Odds1x2, error) {
	if result, err := o.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Odds1x2), nil
	}
}

func (o odds1x2Do) Last() (*model.Odds1x2, error) {
	if result, err := o.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Odds1x2), nil
	}
}

func (o odds1x2Do) Find() ([]*model.Odds1x2, error) {
	result, err := o.DO.Find()
	return result.([]*model.Odds1x2), err
}

func (o odds1x2Do) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Odds1x2, err error) {
	buf := make([]*model.Odds1x2, 0, batchSize)
	err = o.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (o odds1x2Do) FindInBatches(result *[]*model.Odds1x2, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return o.DO.FindInBatches(result, batchSize, fc)
}

func (o odds1x2Do) Attrs(attrs ...field.AssignExpr) IOdds1x2Do {
	return o.withDO(o.DO.Attrs(attrs...))
}

func (o odds1x2Do) Assign(attrs ...field.AssignExpr) IOdds1x2Do {
	return o.withDO(o.DO.Assign(attrs...))
}

func (o odds1x2Do) Joins(fields ...field.RelationField) IOdds1x2Do {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Joins(_f))
	}
	return &o
}

func (o odds1x2Do) Preload(fields ...field.RelationField) IOdds1x2Do {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Preload(_f))
	}
	return &o
}

func (o odds1x2Do) FirstOrInit() (*model.Odds1x2, error) {
	if result, err := o.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Odds1x2), nil
	}
}

func (o odds1x2Do) FirstOrCreate() (*model.Odds1x2, error) {
	if result, err := o.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Odds1x2), nil
	}
}

func (o odds1x2Do) FindByPage(offset int, limit int) (result []*model.Odds1x2, count int64, err error) {
	result, err = o.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = o.Offset(-1).Limit(-1).Count()
	return
}

func (o odds1x2Do) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = o.Count()
	if err != nil {
		return
	}

	err = o.Offset(offset).Limit(limit).Scan(result)
	return
}

func (o odds1x2Do) Scan(result interface{}) (err error) {
	return o.DO.Scan(result)
}

func (o odds1x2Do) Delete(models ...*model.Odds1x2) (result gen.ResultInfo, err error) {
	return o.DO.Delete(models)
}

func (o *odds1x2Do) withDO(do gen.Dao) *odds1x2Do {
	o.DO = *do.(*gen.DO)
	return o
}