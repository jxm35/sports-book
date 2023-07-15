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

func newMatch(db *gorm.DB, opts ...gen.DOOption) match {
	_match := match{}

	_match.matchDo.UseDB(db, opts...)
	_match.matchDo.UseModel(&model.Match{})

	tableName := _match.matchDo.TableName()
	_match.ALL = field.NewAsterisk(tableName)
	_match.ID = field.NewInt32(tableName, "id")
	_match.Date = field.NewTime(tableName, "date")
	_match.HomeTeam = field.NewInt32(tableName, "home_team")
	_match.AwayTeam = field.NewInt32(tableName, "away_team")
	_match.Competition = field.NewInt32(tableName, "competition")
	_match.HomeGoals = field.NewInt32(tableName, "home_goals")
	_match.AwayGoals = field.NewInt32(tableName, "away_goals")
	_match.HomeExpectedGoals = field.NewFloat64(tableName, "home_expected_goals")
	_match.AwayExpectedGoals = field.NewFloat64(tableName, "away_expected_goals")
	_match.UsID = field.NewInt32(tableName, "us_id")

	_match.fillFieldMap()

	return _match
}

type match struct {
	matchDo

	ALL               field.Asterisk
	ID                field.Int32
	Date              field.Time
	HomeTeam          field.Int32
	AwayTeam          field.Int32
	Competition       field.Int32
	HomeGoals         field.Int32
	AwayGoals         field.Int32
	HomeExpectedGoals field.Float64
	AwayExpectedGoals field.Float64
	UsID              field.Int32

	fieldMap map[string]field.Expr
}

func (m match) Table(newTableName string) *match {
	m.matchDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m match) As(alias string) *match {
	m.matchDo.DO = *(m.matchDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *match) updateTableName(table string) *match {
	m.ALL = field.NewAsterisk(table)
	m.ID = field.NewInt32(table, "id")
	m.Date = field.NewTime(table, "date")
	m.HomeTeam = field.NewInt32(table, "home_team")
	m.AwayTeam = field.NewInt32(table, "away_team")
	m.Competition = field.NewInt32(table, "competition")
	m.HomeGoals = field.NewInt32(table, "home_goals")
	m.AwayGoals = field.NewInt32(table, "away_goals")
	m.HomeExpectedGoals = field.NewFloat64(table, "home_expected_goals")
	m.AwayExpectedGoals = field.NewFloat64(table, "away_expected_goals")
	m.UsID = field.NewInt32(table, "us_id")

	m.fillFieldMap()

	return m
}

func (m *match) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *match) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 10)
	m.fieldMap["id"] = m.ID
	m.fieldMap["date"] = m.Date
	m.fieldMap["home_team"] = m.HomeTeam
	m.fieldMap["away_team"] = m.AwayTeam
	m.fieldMap["competition"] = m.Competition
	m.fieldMap["home_goals"] = m.HomeGoals
	m.fieldMap["away_goals"] = m.AwayGoals
	m.fieldMap["home_expected_goals"] = m.HomeExpectedGoals
	m.fieldMap["away_expected_goals"] = m.AwayExpectedGoals
	m.fieldMap["us_id"] = m.UsID
}

func (m match) clone(db *gorm.DB) match {
	m.matchDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m match) replaceDB(db *gorm.DB) match {
	m.matchDo.ReplaceDB(db)
	return m
}

type matchDo struct{ gen.DO }

type IMatchDo interface {
	gen.SubQuery
	Debug() IMatchDo
	WithContext(ctx context.Context) IMatchDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IMatchDo
	WriteDB() IMatchDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IMatchDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IMatchDo
	Not(conds ...gen.Condition) IMatchDo
	Or(conds ...gen.Condition) IMatchDo
	Select(conds ...field.Expr) IMatchDo
	Where(conds ...gen.Condition) IMatchDo
	Order(conds ...field.Expr) IMatchDo
	Distinct(cols ...field.Expr) IMatchDo
	Omit(cols ...field.Expr) IMatchDo
	Join(table schema.Tabler, on ...field.Expr) IMatchDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IMatchDo
	RightJoin(table schema.Tabler, on ...field.Expr) IMatchDo
	Group(cols ...field.Expr) IMatchDo
	Having(conds ...gen.Condition) IMatchDo
	Limit(limit int) IMatchDo
	Offset(offset int) IMatchDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IMatchDo
	Unscoped() IMatchDo
	Create(values ...*model.Match) error
	CreateInBatches(values []*model.Match, batchSize int) error
	Save(values ...*model.Match) error
	First() (*model.Match, error)
	Take() (*model.Match, error)
	Last() (*model.Match, error)
	Find() ([]*model.Match, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Match, err error)
	FindInBatches(result *[]*model.Match, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Match) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IMatchDo
	Assign(attrs ...field.AssignExpr) IMatchDo
	Joins(fields ...field.RelationField) IMatchDo
	Preload(fields ...field.RelationField) IMatchDo
	FirstOrInit() (*model.Match, error)
	FirstOrCreate() (*model.Match, error)
	FindByPage(offset int, limit int) (result []*model.Match, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IMatchDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (m matchDo) Debug() IMatchDo {
	return m.withDO(m.DO.Debug())
}

func (m matchDo) WithContext(ctx context.Context) IMatchDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m matchDo) ReadDB() IMatchDo {
	return m.Clauses(dbresolver.Read)
}

func (m matchDo) WriteDB() IMatchDo {
	return m.Clauses(dbresolver.Write)
}

func (m matchDo) Session(config *gorm.Session) IMatchDo {
	return m.withDO(m.DO.Session(config))
}

func (m matchDo) Clauses(conds ...clause.Expression) IMatchDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m matchDo) Returning(value interface{}, columns ...string) IMatchDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m matchDo) Not(conds ...gen.Condition) IMatchDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m matchDo) Or(conds ...gen.Condition) IMatchDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m matchDo) Select(conds ...field.Expr) IMatchDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m matchDo) Where(conds ...gen.Condition) IMatchDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m matchDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IMatchDo {
	return m.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (m matchDo) Order(conds ...field.Expr) IMatchDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m matchDo) Distinct(cols ...field.Expr) IMatchDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m matchDo) Omit(cols ...field.Expr) IMatchDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m matchDo) Join(table schema.Tabler, on ...field.Expr) IMatchDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m matchDo) LeftJoin(table schema.Tabler, on ...field.Expr) IMatchDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m matchDo) RightJoin(table schema.Tabler, on ...field.Expr) IMatchDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m matchDo) Group(cols ...field.Expr) IMatchDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m matchDo) Having(conds ...gen.Condition) IMatchDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m matchDo) Limit(limit int) IMatchDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m matchDo) Offset(offset int) IMatchDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m matchDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IMatchDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m matchDo) Unscoped() IMatchDo {
	return m.withDO(m.DO.Unscoped())
}

func (m matchDo) Create(values ...*model.Match) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m matchDo) CreateInBatches(values []*model.Match, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m matchDo) Save(values ...*model.Match) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m matchDo) First() (*model.Match, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Match), nil
	}
}

func (m matchDo) Take() (*model.Match, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Match), nil
	}
}

func (m matchDo) Last() (*model.Match, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Match), nil
	}
}

func (m matchDo) Find() ([]*model.Match, error) {
	result, err := m.DO.Find()
	return result.([]*model.Match), err
}

func (m matchDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Match, err error) {
	buf := make([]*model.Match, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m matchDo) FindInBatches(result *[]*model.Match, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m matchDo) Attrs(attrs ...field.AssignExpr) IMatchDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m matchDo) Assign(attrs ...field.AssignExpr) IMatchDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m matchDo) Joins(fields ...field.RelationField) IMatchDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m matchDo) Preload(fields ...field.RelationField) IMatchDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m matchDo) FirstOrInit() (*model.Match, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Match), nil
	}
}

func (m matchDo) FirstOrCreate() (*model.Match, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Match), nil
	}
}

func (m matchDo) FindByPage(offset int, limit int) (result []*model.Match, count int64, err error) {
	result, err = m.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = m.Offset(-1).Limit(-1).Count()
	return
}

func (m matchDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m matchDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m matchDo) Delete(models ...*model.Match) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *matchDo) withDO(do gen.Dao) *matchDo {
	m.DO = *do.(*gen.DO)
	return m
}
