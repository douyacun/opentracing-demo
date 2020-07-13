package models

type Slot struct {
	Id     int32  `json:"id" xorm:"pk autoincr"`
	Name   string `json:"name"`
	Status int32  `json:"status"`
}

func (s *Slot) TableName() string {
	return "slot"
}
