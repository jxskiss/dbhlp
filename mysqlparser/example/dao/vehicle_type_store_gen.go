// Autogenerated by github.com/jxsiss/dbhlp/mysqlparser/gen/run.
// DO NOT EDIT.

package dao

import (
	"github.com/jxskiss/dbhlp"
	"gorm.io/gorm"
)

const tableName_VehicleType = "vehicle_type"

type VehicleTypeDAO interface {
	vehicleTypeCustomMethods
}

func GetVehicleTypeDAO(conn dbhlp.MySQLConn) VehicleTypeDAO {
	return &vehicleTypeDAOImpl{
		db: conn,
	}
}

type vehicleTypeDAOImpl struct {
	db *gorm.DB
}
