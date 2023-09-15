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

	"sports-book.com/pkg/gorm/model"
)

func newCompetition(db *gorm.DB, opts ...gen.DOOption) competition {
	_competition := competition{}

	_competition.competitionDo.UseDB(db, opts...)
	_competition.competitionDo.UseModel(&model.Competition{})

	tableName := _competition.competitionDo.TableName()
	_competition.ALL = field.NewAsterisk(tableName)
	_competition.ID = field.NewInt32(tableName, "id")
	_competition.Code = field.NewString(tableName, "code")
	_competition.Year = field.NewInt32(tableName, "year")
	_competition.UsID = field.NewInt32(tableName, "us_id")

	_competition.fillFieldMap()

	return _competition
}

type competition struct {
	competitionDo

	ALL  field.Asterisk
	ID   field.Int32
	Code field.String
	Year field.Int32
	UsID field.Int32

	fieldMap map[string]field.Expr
}

func (c competition) Table(newTableName string) *competition {
	c.competitionDo.UseTable(newTableName)
	return c.updateTableName(newTableName)
}

func (c competition) As(alias string) *competition {
	c.competitionDo.DO = *(c.competitionDo.As(alias).(*gen.DO))
	return c.updateTableName(alias)
}

func (c *competition) updateTableName(table string) *competition {
	c.ALL = field.NewAsterisk(table)
	c.ID = field.NewInt32(table, "id")
	c.Code = field.NewString(table, "code")
	c.Year = field.NewInt32(table, "year")
	c.UsID = field.NewInt32(table, "us_id")

	c.fillFieldMap()

	return c
}

func (c *competition) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := c.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (c *competition) fillFieldMap() {
	c.fieldMap = make(map[string]field.Expr, 4)
	c.fieldMap["id"] = c.ID
	c.fieldMap["code"] = c.Code
	c.fieldMap["year"] = c.Year
	c.fieldMap["us_id"] = c.UsID
}

func (c competition) clone(db *gorm.DB) competition {
	c.competitionDo.ReplaceConnPool(db.Statement.ConnPool)
	return c
}

func (c competition) replaceDB(db *gorm.DB) competition {
	c.competitionDo.ReplaceDB(db)
	return c
}

type competitionDo struct{ gen.DO }

type ICompetitionDo interface {
	gen.SubQuery
	Debug() ICompetitionDo
	WithContext(ctx context.Context) ICompetitionDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ICompetitionDo
	WriteDB() ICompetitionDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ICompetitionDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ICompetitionDo
	Not(conds ...gen.Condition) ICompetitionDo
	Or(conds ...gen.Condition) ICompetitionDo
	Select(conds ...field.Expr) ICompetitionDo
	Where(conds ...gen.Condition) ICompetitionDo
	Order(conds ...field.Expr) ICompetitionDo
	Distinct(cols ...field.Expr) ICompetitionDo
	Omit(cols ...field.Expr) ICompetitionDo
	Join(table schema.Tabler, on ...field.Expr) ICompetitionDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ICompetitionDo
	RightJoin(table schema.Tabler, on ...field.Expr) ICompetitionDo
	Group(cols ...field.Expr) ICompetitionDo
	Having(conds ...gen.Condition) ICompetitionDo
	Limit(limit int) ICompetitionDo
	Offset(offset int) ICompetitionDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ICompetitionDo
	Unscoped() ICompetitionDo
	Create(values ...*model.Competition) error
	CreateInBatches(values []*model.Competition, batchSize int) error
	Save(values ...*model.Competition) error
	First() (*model.Competition, error)
	Take() (*model.Competition, error)
	Last() (*model.Competition, error)
	Find() ([]*model.Competition, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Competition, err error)
	FindInBatches(result *[]*model.Competition, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Competition) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ICompetitionDo
	Assign(attrs ...field.AssignExpr) ICompetitionDo
	Joins(fields ...field.RelationField) ICompetitionDo
	Preload(fields ...field.RelationField) ICompetitionDo
	FirstOrInit() (*model.Competition, error)
	FirstOrCreate() (*model.Competition, error)
	FindByPage(offset int, limit int) (result []*model.Competition, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ICompetitionDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (c competitionDo) Debug() ICompetitionDo {
	return c.withDO(c.DO.Debug())
}

func (c competitionDo) WithContext(ctx context.Context) ICompetitionDo {
	return c.withDO(c.DO.WithContext(ctx))
}

func (c competitionDo) ReadDB() ICompetitionDo {
	return c.Clauses(dbresolver.Read)
}

func (c competitionDo) WriteDB() ICompetitionDo {
	return c.Clauses(dbresolver.Write)
}

func (c competitionDo) Session(config *gorm.Session) ICompetitionDo {
	return c.withDO(c.DO.Session(config))
}

func (c competitionDo) Clauses(conds ...clause.Expression) ICompetitionDo {
	return c.withDO(c.DO.Clauses(conds...))
}

func (c competitionDo) Returning(value interface{}, columns ...string) ICompetitionDo {
	return c.withDO(c.DO.Returning(value, columns...))
}

func (c competitionDo) Not(conds ...gen.Condition) ICompetitionDo {
	return c.withDO(c.DO.Not(conds...))
}

func (c competitionDo) Or(conds ...gen.Condition) ICompetitionDo {
	return c.withDO(c.DO.Or(conds...))
}

func (c competitionDo) Select(conds ...field.Expr) ICompetitionDo {
	return c.withDO(c.DO.Select(conds...))
}

func (c competitionDo) Where(conds ...gen.Condition) ICompetitionDo {
	return c.withDO(c.DO.Where(conds...))
}

func (c competitionDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) ICompetitionDo {
	return c.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (c competitionDo) Order(conds ...field.Expr) ICompetitionDo {
	return c.withDO(c.DO.Order(conds...))
}

func (c competitionDo) Distinct(cols ...field.Expr) ICompetitionDo {
	return c.withDO(c.DO.Distinct(cols...))
}

func (c competitionDo) Omit(cols ...field.Expr) ICompetitionDo {
	return c.withDO(c.DO.Omit(cols...))
}

func (c competitionDo) Join(table schema.Tabler, on ...field.Expr) ICompetitionDo {
	return c.withDO(c.DO.Join(table, on...))
}

func (c competitionDo) LeftJoin(table schema.Tabler, on ...field.Expr) ICompetitionDo {
	return c.withDO(c.DO.LeftJoin(table, on...))
}

func (c competitionDo) RightJoin(table schema.Tabler, on ...field.Expr) ICompetitionDo {
	return c.withDO(c.DO.RightJoin(table, on...))
}

func (c competitionDo) Group(cols ...field.Expr) ICompetitionDo {
	return c.withDO(c.DO.Group(cols...))
}

func (c competitionDo) Having(conds ...gen.Condition) ICompetitionDo {
	return c.withDO(c.DO.Having(conds...))
}

func (c competitionDo) Limit(limit int) ICompetitionDo {
	return c.withDO(c.DO.Limit(limit))
}

func (c competitionDo) Offset(offset int) ICompetitionDo {
	return c.withDO(c.DO.Offset(offset))
}

func (c competitionDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ICompetitionDo {
	return c.withDO(c.DO.Scopes(funcs...))
}

func (c competitionDo) Unscoped() ICompetitionDo {
	return c.withDO(c.DO.Unscoped())
}

func (c competitionDo) Create(values ...*model.Competition) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Create(values)
}

func (c competitionDo) CreateInBatches(values []*model.Competition, batchSize int) error {
	return c.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (c competitionDo) Save(values ...*model.Competition) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Save(values)
}

func (c competitionDo) First() (*model.Competition, error) {
	if result, err := c.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Competition), nil
	}
}

func (c competitionDo) Take() (*model.Competition, error) {
	if result, err := c.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Competition), nil
	}
}

func (c competitionDo) Last() (*model.Competition, error) {
	if result, err := c.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Competition), nil
	}
}

func (c competitionDo) Find() ([]*model.Competition, error) {
	result, err := c.DO.Find()
	return result.([]*model.Competition), err
}

func (c competitionDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Competition, err error) {
	buf := make([]*model.Competition, 0, batchSize)
	err = c.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (c competitionDo) FindInBatches(result *[]*model.Competition, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return c.DO.FindInBatches(result, batchSize, fc)
}

func (c competitionDo) Attrs(attrs ...field.AssignExpr) ICompetitionDo {
	return c.withDO(c.DO.Attrs(attrs...))
}

func (c competitionDo) Assign(attrs ...field.AssignExpr) ICompetitionDo {
	return c.withDO(c.DO.Assign(attrs...))
}

func (c competitionDo) Joins(fields ...field.RelationField) ICompetitionDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Joins(_f))
	}
	return &c
}

func (c competitionDo) Preload(fields ...field.RelationField) ICompetitionDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Preload(_f))
	}
	return &c
}

func (c competitionDo) FirstOrInit() (*model.Competition, error) {
	if result, err := c.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Competition), nil
	}
}

func (c competitionDo) FirstOrCreate() (*model.Competition, error) {
	if result, err := c.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Competition), nil
	}
}

func (c competitionDo) FindByPage(offset int, limit int) (result []*model.Competition, count int64, err error) {
	result, err = c.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = c.Offset(-1).Limit(-1).Count()
	return
}

func (c competitionDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = c.Count()
	if err != nil {
		return
	}

	err = c.Offset(offset).Limit(limit).Scan(result)
	return
}

func (c competitionDo) Scan(result interface{}) (err error) {
	return c.DO.Scan(result)
}

func (c competitionDo) Delete(models ...*model.Competition) (result gen.ResultInfo, err error) {
	return c.DO.Delete(models)
}

func (c *competitionDo) withDO(do gen.Dao) *competitionDo {
	c.DO = *do.(*gen.DO)
	return c
}