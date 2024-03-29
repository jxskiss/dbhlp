// Autogenerated by github.com/jxsiss/dbhlp/mysqlparser/gen/run.
// DO NOT EDIT.

package dao

import (
	"github.com/jxskiss/dbhlp"
	"gorm.io/gorm"
)

const tableName_Insurance = "insurance"

type InsuranceDAO interface {
	insuranceCustomMethods
}

func GetInsuranceDAO(conn dbhlp.DBConn) InsuranceDAO {
	return &insuranceDAOImpl{
		db: conn,
	}
}

type insuranceDAOImpl struct {
	db *gorm.DB
}
