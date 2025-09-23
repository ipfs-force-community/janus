package orm

import "gorm.io/gorm"

// Miner represents table miner in the database
type Miner struct {
	gorm.Model
	Height    int64  `gorm:"not null"`
	Cid       string `gorm:"type:varchar(255);column:cid;not null"`
	Timestamp int64  `gorm:"not null"`
	MsgCid    string `gorm:"type:varchar(255);column:msg_cid;unique;not null"`
	From      string `gorm:"type:varchar(255);not null"`
	Cost      string `gorm:"type:varchar(255);not null"`
}
