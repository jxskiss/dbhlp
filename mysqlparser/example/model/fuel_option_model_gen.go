// Autogenerated by github.com/jxsiss/dbhlp/mysqlparser/gen/run.
// DO NOT EDIT.

package model

type FuelOption struct {
	Id          int64  `db:"id" gorm:"column:id;primaryKey"`        // int(11)
	Name        string `db:"name" gorm:"column:name"`               // varchar(45)
	Description string `db:"description" gorm:"column:description"` // varchar(255)
}

type FuelOptionList []*FuelOption

func (p FuelOptionList) ToIdMap() map[int64]*FuelOption {
	out := make(map[int64]*FuelOption, len(p))
	for _, x := range p {
		out[x.Id] = x
	}
	return out
}

func (p FuelOptionList) PluckIds() []int64 {
	out := make([]int64, 0, len(p))
	for _, x := range p {
		out = append(out, x.Id)
	}
	return out
}
