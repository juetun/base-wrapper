package app_param

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	UploadDivideString = "|"
)

type (
	UploadImage struct {
		UploadCommon
	}
	UploadVideo struct {
		UploadCommon
	}
	UploadMusic struct {
		UploadCommon
	}
	UploadCommon struct {
		Channel string `json:"channel" form:"channel"`
		ID      int64  `json:"id" form:"id"`
	}
	ShowData struct {
		DefaultKey  string      `json:"default_key"`
		PlayAddress PlayAddress `json:"play_address"`
	}
	PlayAddress map[string]string
)

func (r *UploadImage) ToString() (res string) {
	res = fmt.Sprintf("%s%s%d", r.Channel, UploadDivideString, r.ID)
	return
}

// GetShowUrl 获取图片地址的播放地址
func (r *UploadImage) GetShowUrl() (res string) {
	res = ""
	return
}

func (r *UploadImage) ParseString(saveUploadString string) (err error) {
	err = r.UploadCommon.parseString(saveUploadString)
	return
}

func (r *UploadVideo) ToString() (res string) {
	res = fmt.Sprintf("%s%s%d", r.Channel, UploadDivideString, r.ID)
	return
}

// GetShowUrl 获取视频的播放地址
func (r *UploadVideo) GetShowUrl() (res ShowData) {
	res = ShowData{
		PlayAddress: map[string]string{},
	}
	return
}

func (r *UploadVideo) ParseString(saveUploadString string) (err error) {
	err = r.UploadCommon.parseString(saveUploadString)
	return
}

func (r *UploadMusic) ToString() (res string) {
	res = fmt.Sprintf("%s|%d", r.Channel, r.ID)
	return
}

func (r *UploadMusic) ParseString(saveUploadString string) (err error) {
	err = r.UploadCommon.parseString(saveUploadString)
	return
}

// GetShowUrl 获取音频的播放地址
func (r *UploadMusic) GetShowUrl() (res ShowData) {
	res = ShowData{
		PlayAddress: map[string]string{},
	}
	return
}

func (r *UploadCommon) parseString(saveUploadString string) (err error) {
	tmp := strings.Split(saveUploadString, UploadDivideString)
	switch len(tmp) {
	case 1:
		tmp[1] = "0"
	}
	r.Channel = tmp[0]
	r.ID, err = strconv.ParseInt(tmp[1], 10, 64)
	return
}
