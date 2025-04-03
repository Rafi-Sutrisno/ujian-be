package entity

type Language struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(255);not null;unique" json:"name"`
	Code string `gorm:"type:varchar(10);not null;unique" json:"code"`
}
