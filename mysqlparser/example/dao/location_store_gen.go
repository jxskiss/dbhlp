// Autogenerated by github.com/jxsiss/dbhlp/mysqlparser/gen/run.
// DO NOT EDIT.

package dao

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jxskiss/dbhlp"
	"github.com/jxskiss/errors"
	"github.com/jxskiss/gopkg/v2/utils/sqlutil"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"

	"github.com/jxskiss/dbhlp/mysqlparser/example/model"
)

var _ context.Context
var _ json.Marshaler
var _ log.Logger
var _ time.Time
var _ dbhlp.Opt
var _ errors.ErrorGroup
var _ sqlutil.LazyBinary
var _ gorm.DB
var _ proto.Message

const tableName_Location = "location"

type LocationDAO interface {
	locationCustomMethods
}

func GetLocationDAO(conn dbhlp.MySQLConn) LocationDAO {
	return &locationDAOImpl{
		db: conn,
	}
}

type locationDAOImpl struct {
	db *gorm.DB
}
