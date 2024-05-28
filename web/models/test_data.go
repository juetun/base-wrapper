package models

import "encoding/json"

type (
	Data struct {
		UserHid string `json:"user_hid" gorm:"column:user_hid"`
	}
	DataChildren struct {
		Data
		UserId string `json:"user_id"  gorm:"column:user_id"`
	}
)

func (r *DataChildren) Default() (err error) {

	return
}

func (r *DataChildren) UnmarshalBinary(data []byte) (err error) {
	err = json.Unmarshal(data, r)
	return
}

func (r *DataChildren) MarshalBinary() (data []byte, err error) {
	if r == nil {
		data = []byte("{}")
		return
	}
	data, err = json.Marshal(r)
	return
}

func (r *DataChildren) TableName() string {
	return "data_children"
}

func (r *DataChildren) GetTableComment() (res string) {

	return "测试批量"
}
