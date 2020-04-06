/**
* @Author:changjiang
* @Description:
* @File:upload
* @Version: 1.0.0
* @Date 2020/3/18 5:32 下午
 */
package common

type ConfUpload struct {
	ImgUploadUrl  string `yaml:"ImgUploadUrl"`
	ImgUploadDst  string `yaml:"ImgUploadDst"`
	ImgUploadBoth bool   `yaml:"ImgUploadBoth"`

	// qiniu
	QiNiuUploadImg bool   `yaml:"QiNiuUploadImg"`
	QiNiuHostName  string `yaml:"QiNiuHostName"`
	QiNiuAccessKey string `yaml:"QiNiuAccessKey"`
	QiNiuSecretKey string `yaml:"QiNiuSecretKey"`
	QiNiuBucket    string `yaml:"QiNiuBucket"`
	QiNiuZone      string `yaml:"QiNiuZone"`
	AppImgUrl      string `yaml:"AppImgUrl"`
}

var ConfigUpload *ConfUpload
