package orm

import "gorm.io/gorm"

// Chain represents table chain in the database
type Chain struct {
	gorm.Model
	Height int64 `gorm:"not null"`
}
