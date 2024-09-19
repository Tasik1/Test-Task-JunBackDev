package model

type Token struct {
	User    User `gorm:"foreignkey:GUID"`
	GUID    uint
	Refresh string
}
