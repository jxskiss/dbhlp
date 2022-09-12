// This file is autogenerated. DO NOT EDIT.

package model

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jxskiss/dbhlp"
	"github.com/jxskiss/errors"
	"github.com/jxskiss/gopkg/v2/sqlutil"
	"gorm.io/gorm"
)

var _ context.Context
var _ time.Time
var _ proto.Message
var _ errors.ErrorGroup
var _ dbhlp.Opt
var _ sqlutil.Bitmap
var _ gorm.DB

type Location struct {
	Id            int64  `db:"id" gorm:"column:id;primaryKey"`              // int(11)
	StreetAddress string `db:"street_address" gorm:"column:street_address"` // varchar(100)
	City          string `db:"city" gorm:"column:city"`                     // varchar(50)
	State         string `db:"state" gorm:"column:state"`                   // varchar(2)
	Zipcode       int    `db:"zipcode" gorm:"column:zipcode"`               // int(5)
}

type LocationList []*Location

func (p LocationList) ToIdMap() map[int64]*Location {
	out := make(map[int64]*Location, len(p))
	for _, x := range p {
		out[x.Id] = x
	}
	return out
}

func (p LocationList) PluckIds() []int64 {
	out := make([]int64, 0, len(p))
	for _, x := range p {
		out = append(out, x.Id)
	}
	return out
}
