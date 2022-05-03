package models

type (
	Data struct {
		UserHid string `json:"user_hid" gorm:"column:user_hid"`
	}
	DataChildren struct {
		Data
		UserId string `json:"user_id"  gorm:"column:user_id"`
	}
)

func (d DataChildren) TableName() string {

	return  "data_children"
}

func (d DataChildren) GetTableComment() (res string) {

	return  "测试批量"
}

