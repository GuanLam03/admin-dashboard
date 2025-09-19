package models

import "github.com/goravel/framework/database/orm"

type User struct {
    orm.Model
    Name     string `json:"name"`
    Email    string `json:"email" gorm:"unique"`
    Password string `json:"password"`

    TwoFactorSecret string `gorm:"column:two_factor_secret"`
    TwoFactorEnabled bool  `gorm:"column:two_factor_enabled"`
}
