// This file is autogenerated. DO NOT EDIT.

package model

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jxskiss/dbgen"
	"github.com/jxskiss/errors"
	"github.com/jxskiss/gopkg/v2/sqlutil"
	"gorm.io/gorm"
)

var _ context.Context
var _ time.Time
var _ proto.Message
var _ errors.ErrorGroup
var _ dbgen.Opt
var _ sqlutil.Bitmap
var _ gorm.DB

type Equipment struct {
	Id                int64  `db:"id" gorm:"column:id;primaryKey"`                        // int(11)
	Name              string `db:"name" gorm:"column:name"`                               // varchar(45)
	EquipmentTypeId   int    `db:"equipment_type_id" gorm:"column:equipment_type_id"`     // int(11)
	CurrentLocationId int    `db:"current_location_id" gorm:"column:current_location_id"` // int(11)
}

type EquipmentList []*Equipment

func (p EquipmentList) ToIdMap() map[int64]*Equipment {
	out := make(map[int64]*Equipment, len(p))
	for _, x := range p {
		out[x.Id] = x
	}
	return out
}

func (p EquipmentList) PluckIds() []int64 {
	out := make([]int64, 0, len(p))
	for _, x := range p {
		out = append(out, x.Id)
	}
	return out
}