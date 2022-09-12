// This file is autogenerated. DO NOT EDIT.

package dao

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jxskiss/dbhlp"
	"github.com/jxskiss/errors"
	"github.com/jxskiss/gopkg/v2/sqlutil"
	"gorm.io/gorm"

	"github.com/jxskiss/dbhlp/mysqlparser/example/model"
)

var _ context.Context
var _ time.Time
var _ proto.Message
var _ errors.ErrorGroup
var _ dbhlp.Opt
var _ sqlutil.Bitmap
var _ gorm.DB

const tableName_Insurance = "insurance"

type InsuranceDAO interface {
	Get(ctx context.Context, id int64, opts ...dbhlp.Opt) (*model.Insurance, error)
	GetWhere(ctx context.Context, where string, paramsAndOpts ...interface{}) (*model.Insurance, error)
	MGet(ctx context.Context, idList []int64, opts ...dbhlp.Opt) (model.InsuranceList, error)
	MGetWhere(ctx context.Context, where string, paramsAndOpts ...interface{}) (model.InsuranceList, error)
	Update(ctx context.Context, id int64, updates map[string]interface{}, opts ...dbhlp.Opt) error
	insuranceCustomMethods
}

func GetInsuranceDAO(conn dbhlp.MySQLConn) InsuranceDAO {
	return &insuranceDAOImpl{
		db: conn,
	}
}

type insuranceDAOImpl struct {
	db *gorm.DB
}

func (p *insuranceDAOImpl) Get(ctx context.Context, id int64, opts ...dbhlp.Opt) (*model.Insurance, error) {
	conn := dbhlp.GetSession(p.db, opts...)
	tableName := tableName_Insurance
	var out = &model.Insurance{}
	err := conn.WithContext(ctx).Table(tableName).Where("id = ?", id).First(out).Error
	if err != nil {
		return nil, errors.AddStack(err)
	}
	return out, nil
}

func (p *insuranceDAOImpl) GetWhere(ctx context.Context, where string, paramsAndOpts ...interface{}) (*model.Insurance, error) {
	params, opts := dbhlp.SplitOpts(paramsAndOpts)
	conn := dbhlp.GetSession(p.db, opts...)
	tableName := tableName_Insurance
	var out = &model.Insurance{}
	err := conn.WithContext(ctx).Table(tableName).Where(where, params...).First(out).Error
	if err != nil {
		return nil, errors.AddStack(err)
	}
	return out, nil
}

func (p *insuranceDAOImpl) MGet(ctx context.Context, idList []int64, opts ...dbhlp.Opt) (model.InsuranceList, error) {
	conn := dbhlp.GetSession(p.db, opts...)
	tableName := tableName_Insurance
	var out model.InsuranceList
	err := conn.WithContext(ctx).Table(tableName).Where("id in (?)", idList).Find(&out).Error
	if err != nil {
		return nil, errors.AddStack(err)
	}
	return out, nil
}

func (p *insuranceDAOImpl) MGetWhere(ctx context.Context, where string, paramsAndOpts ...interface{}) (model.InsuranceList, error) {
	params, opts := dbhlp.SplitOpts(paramsAndOpts)
	conn := dbhlp.GetSession(p.db, opts...)
	tableName := tableName_Insurance
	var out model.InsuranceList
	err := conn.WithContext(ctx).Table(tableName).Where(where, params...).Find(&out).Error
	if err != nil {
		return nil, errors.AddStack(err)
	}
	return out, nil
}

func (p *insuranceDAOImpl) Update(ctx context.Context, id int64, updates map[string]interface{}, opts ...dbhlp.Opt) error {
	if len(updates) == 0 {
		return errors.New("programming error: empty updates map")
	}
	conn := dbhlp.GetSession(p.db, opts...)
	tableName := tableName_Insurance
	err := conn.WithContext(ctx).Table(tableName).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return errors.AddStack(err)
	}
	return nil
}
