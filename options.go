package dbgen

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
	"gorm.io/plugin/dbresolver"
)

func PrepareSession(db *gorm.DB, opts ...Opt) *gorm.DB {
	return new(Options).apply(opts...).prepareSession(db)
}

type Options struct {
	Tx           *gorm.DB
	UseMaster    bool
	ForUpdate    bool
	ForShare     bool
	InsertIgnore bool

	Debug bool

	IndexHints []clause.Expression

	OnDupDoUpdates     map[string]interface{}
	OnDupUpdateColumns []string
}

func (p *Options) apply(opts ...Opt) *Options {
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *Options) prepareSession(db *gorm.DB) *gorm.DB {
	conn := db
	if p.Tx != nil {
		conn = p.Tx
	} else if p.UseMaster {
		conn = conn.Clauses(dbresolver.Write)
	}
	if p.ForUpdate {
		conn = conn.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if p.ForShare {
		conn = conn.Clauses(clause.Locking{
			Strength: "SHARE",
			Table:    clause.Table{Name: clause.CurrentTable},
		})
	}
	if p.InsertIgnore {
		conn = conn.Clauses(clause.Insert{Modifier: "IGNORE"})
	}
	if len(p.OnDupDoUpdates) > 0 {
		conn = conn.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(p.OnDupDoUpdates),
		})
	}
	if len(p.OnDupUpdateColumns) > 0 {
		conn = conn.Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns(p.OnDupUpdateColumns),
		})
	}
	if len(p.IndexHints) > 0 {
		conn = conn.Clauses(p.IndexHints...)
	}
	if p.Debug {
		conn = conn.Debug()
	}
	return conn
}

type Opt func(*Options)

func UseTX(tx *gorm.DB) Opt {
	return func(opt *Options) {
		opt.Tx = tx
	}
}

func UseMaster(opt *Options) {
	opt.UseMaster = true
}

func ForUpdate(opt *Options) {
	opt.ForUpdate = true
}

func ForShare(opt *Options) {
	opt.ForShare = true
}

func InsertIgnore(opt *Options) {
	opt.InsertIgnore = true
}

func OnDupDoUpdates(updates map[string]interface{}) func(opt *Options) {
	return func(opt *Options) {
		opt.OnDupDoUpdates = updates
	}
}

func OnDupUpdateColumns(columns []string) func(opt *Options) {
	return func(opt *Options) {
		opt.OnDupUpdateColumns = columns
	}
}

func IndexHints(hints ...hints.IndexHint) func(opt *Options) {
	return func(opt *Options) {
		for _, h := range hints {
			opt.IndexHints = append(opt.IndexHints, h)
		}
	}
}

func Debug(opt *Options) {
	opt.Debug = true
}

func SplitOpts(paramsAndOpts []interface{}) (params []interface{}, opts []Opt) {
	var length = len(paramsAndOpts)
	var i int
	for i = 0; i < length; i++ {
		switch paramsAndOpts[i].(type) {
		case Opt, func(*Options):
			goto split
		}
	}
split:
	params = paramsAndOpts[:i]
	if i < length {
		opts = make([]Opt, 0, length-i)
		for _, x := range paramsAndOpts[i:] {
			switch opt := x.(type) {
			case Opt:
				opts = append(opts, opt)
			case func(*Options):
				opts = append(opts, opt)
			default:
				params = append(params, x)
			}
		}
	}
	return params, opts
}
