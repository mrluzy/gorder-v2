package builder

import (
	"github.com/mrluzy/gorder-v2/common/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Stock struct {
	ID        []int64  `json:"id"`
	ProductID []string `json:"product_id"`
	Quantity  []int32  `json:"quantity"`
	Version   []int64  `json:"version"`

	// extend fields
	OrderBy       string `json:"order_by"`
	ForUpdateLock bool   `json:"for_update_lock"`
}

func (s *Stock) FormatArg() (string, error) {
	return util.MarshalString(s)
}

func NewStock() *Stock {
	return &Stock{}
}
func (s *Stock) Fill(db *gorm.DB) *gorm.DB {
	db = s.fillWhere(db)
	if s.OrderBy != "" {
		db = db.Order(s.OrderBy)
	}
	return db
}

func (s *Stock) fillWhere(db *gorm.DB) *gorm.DB {
	if len(s.ID) > 0 {
		db = db.Where("id IN (?)", s.ID)
	}
	if len(s.ProductID) > 0 {
		db = db.Where("product_id IN (?)", s.ProductID)
	}
	if len(s.Version) > 0 {
		db = db.Where("version IN (?)", s.Version)
	}
	if len(s.Quantity) > 0 {
		db = s.fillQuantityGT(db)
	}
	if s.ForUpdateLock {
		db = db.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate})
	}
	return db
}

func (s *Stock) fillQuantityGT(db *gorm.DB) *gorm.DB {
	db = db.Where("quantity >= ?", s.Quantity)
	return db
}

func (s *Stock) IDs(v ...int64) *Stock {
	s.ID = v
	return s
}

func (s *Stock) ProductIDs(v ...string) *Stock {
	s.ProductID = v
	return s
}

func (s *Stock) Order(v string) *Stock {
	s.OrderBy = v
	return s
}

func (s *Stock) QuantityGT(v ...int32) *Stock {
	s.Quantity = v
	return s
}

func (s *Stock) Versions(v ...int64) *Stock {
	s.Version = v
	return s
}

func (s *Stock) ForUpdate() *Stock {
	s.ForUpdateLock = true
	return s
}
