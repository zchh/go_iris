package models

import "github.com/jinzhu/gorm"

type Userinfo struct {
	gorm.Model
	UID         string `gorm:"primary_key"`
	Username    string `gorm:"type:varchar(64);"`
	Department  string `gorm:"type:varchar(64);"`
	Created     string `gorm:"type:date;"`
	Deleted_at     string `gorm:"type:date;"`
}

