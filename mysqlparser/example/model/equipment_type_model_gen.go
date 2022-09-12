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

type EquipmentType struct {
	Id          int64  `db:"id" gorm:"column:id;primaryKey"`          // int(11)
	Name        string `db:"name" gorm:"column:name"`                 // varchar(45)
	RentalValue string `db:"rental_value" gorm:"column:rental_value"` // decimal(13,2) UNSIGNED
}

type EquipmentTypeList []*EquipmentType

func (p EquipmentTypeList) ToIdMap() map[int64]*EquipmentType {
	out := make(map[int64]*EquipmentType, len(p))
	for _, x := range p {
		out[x.Id] = x
	}
	return out
}

func (p EquipmentTypeList) PluckIds() []int64 {
	out := make([]int64, 0, len(p))
	for _, x := range p {
		out = append(out, x.Id)
	}
	return out
}
