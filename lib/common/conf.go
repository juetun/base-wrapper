/**
* @Author:changjiang
* @Description:
* @File:conf
* @Version: 1.0.0
* @Date 2020/3/18 6:25 下午
 */
package common

var Conf = &Config{
	TagListKey:            "all:tag",
	CateListKey:           "all:cate:sort",
	ArchivesKey:           "index:archives:list",
	LinkIndexKey:          "index:all:link:list",
	PostIndexKey:          "index:all:post:list",
	SystemIndexKey:        "index:all:system:list",
	TagPostIndexKey:       "index:all:tag:post:list",
	CatePostIndexKey:      "index:all:cate:post:list",
	PostDetailIndexKey:    "index:post:detail",
	DataCacheTimeDuration: 720,
	ThemeJs:               "/static/home/assets/js",
	ThemeCss:              "/static/home/assets/css",
	ThemeImg:              "/static/home/assets/img",
	ThemeFancyboxCss:      "/static/home/assets/fancybox",
	ThemeFancyboxJs:       "/static/home/assets/fancybox",
	ThemeHLightCss:        "/static/home/assets/highlightjs",
	ThemeHLightJs:         "/static/home/assets/highlightjs",
	ThemeShareCss:         "/static/home/assets/css",
	ThemeShareJs:          "/static/home/assets/js",
	ThemeArchivesJs:       "/static/home/assets/js",
	ThemeArchivesCss:      "/static/home/assets/css",
	ThemeNiceImg:          "/static/home/assets/img",
	ThemeAllCss:           "/static/home/assets/css",
	ThemeIndexImg:         "/static/home/assets/img",
	ThemeCateImg:          "/static/home/assets/img",
	ThemeTagImg:           "/static/home/assets/img",
	OtherScript:           "<script type=\"text/javascript\"></script>",
	ImgUploadUrl:          "",
	ImgUploadDst:          "./static/uploads/images/",
	ImgUploadBoth:         true,
	DefaultIndexLimit:     "10",

	// 七牛相关设置
	QiNiuUploadImg: true,
	QiNiuHostName:  "",
	QiNiuAccessKey: "",
	QiNiuSecretKey: "",
	QiNiuBucket:    "",
	QiNiuZone:      "HUABEI", // you can use "HUADONG","HUABEI","BEIMEI","HUANAN","XINJIAPO"
	AppImgUrl:      "",       //  # 一般默认是 AppUrl + /static/uploads/images/

}

type Config struct {
	TagListKey            string `json:"tag_list_key" yaml:"TagListKey"`
	CateListKey           string `json:"cate_list_key" yaml:"CateListKey"`
	ArchivesKey           string `json:"archives_key" yaml:"ArchivesKey"`
	LinkIndexKey          string `json:"link_index_key" yaml:"LinkIndexKey"`
	PostIndexKey          string `json:"post_index_key" yaml:"PostIndexKey"`
	SystemIndexKey        string `json:"system_index_key" yaml:"SystemIndexKey"`
	TagPostIndexKey       string `json:"tag_post_index_key" yaml:"TagPostIndexKey"`
	CatePostIndexKey      string `json:"cate_post_index_key" yaml:"CatePostIndexKey"`
	PostDetailIndexKey    string `json:"post_detail_index_key" yaml:"PostDetailIndexKey"`
	DataCacheTimeDuration int    `json:"data_cache_time_duration" yaml:"DataCacheTimeDuration"`
	ThemeJs               string `json:"theme_js" yaml:"ThemeJs"`
	ThemeCss              string `json:"theme_css" yaml:"ThemeCss"`
	ThemeImg              string `json:"theme_img" yaml:"ThemeImg"`
	ThemeFancyboxCss      string `json:"theme_fancybox_css" yaml:"ThemeFancyboxCss"`
	ThemeFancyboxJs       string `json:"theme_fancybox_js" yaml:"ThemeFancyboxJs"`
	ThemeHLightCss        string `json:"theme_h_light_css" yaml:"ThemeHLightCss"`
	ThemeHLightJs         string `json:"theme_h_light_js" yaml:"ThemeHLightJs"`
	ThemeShareCss         string `json:"theme_share_css" yaml:"ThemeShareCss"`
	ThemeShareJs          string `json:"theme_share_js" yaml:"ThemeShareJs"`
	ThemeArchivesJs       string `json:"theme_archives_js" yaml:"ThemeArchivesJs"`
	ThemeArchivesCss      string `json:"theme_archives_css" yaml:"ThemeArchivesCss"`
	ThemeNiceImg          string `json:"theme_nice_img" yaml:"ThemeNiceImg"`
	ThemeAllCss           string `json:"theme_all_css" yaml:"ThemeAllCss"`
	ThemeIndexImg         string `json:"theme_index_img" yaml:"ThemeIndexImg"`
	ThemeCateImg          string `json:"theme_cate_img" yaml:"ThemeCateImg"`
	ThemeTagImg           string `json:"theme_tag_img" yaml:"ThemeTagImg"`
	OtherScript           string `json:"other_script" yaml:"OtherScript"`
	ImgUploadUrl          string `json:"img_upload_url" yaml:"ImgUploadUrl"`
	ImgUploadDst          string `json:"img_upload_dst" yaml:"ImgUploadDst"`
	ImgUploadBoth         bool   `json:"img_upload_both" yaml:"ImgUploadBoth"`
	DefaultIndexLimit     string `json:"default_index_limit" yaml:"DefaultIndexLimit"`

	QiNiuUploadImg bool   `json:"qi_niu_upload_img" yaml:"QiNiuUploadImg"`
	QiNiuHostName  string `json:"qi_niu_host_name" yaml:"QiNiuHostName"`
	QiNiuAccessKey string `json:"qi_niu_access_key" yaml:"QiNiuAccessKey"`
	QiNiuSecretKey string `json:"qi_niu_secret_key"yaml:"QiNiuSecretKey"`
	QiNiuBucket    string `json:"qi_niu_bucket" yaml:"QiNiuBucket"`
	QiNiuZone      string `json:"qi_niu_zone" yaml:"QiNiuZone"`

	AppImgUrl string `yaml:"AppImgUrl"`

	// Theme        int    `yaml:"Theme"`
	// Title        string `yaml:"Title"`
	// Keywords     string `yaml:"Keywords"`
	// Description  string `yaml:"Description"`
	// RecordNumber string `yaml:"RecordNumber"`
	//
	// // UserCnt int `yaml:"UserCnt"`
	//
	// // index

	//
	// // github gitment
	// GithubName         string `yaml:"GithubName"`
	// GithubRepo         string `yaml:"GithubRepo"`
	// GithubClientId     string `yaml:"GithubClientId"`
	// GithubClientSecret string `yaml:"GithubClientSecret"`
	// GithubLabels       string `yaml:"GithubLabels"`
	//
	// OtherScript string `yaml:"OtherScript"`
	//
	// ThemeNiceImg     string `yaml:"ThemeNiceImg"`
	// ThemeAllCss      string `yaml:"ThemeAllCss"`
	// ThemeIndexImg    string `yaml:"ThemeIndexImg"`
	// ThemeCateImg     string `yaml:"ThemeCateImg"`
	// ThemeTagImg      string `yaml:"ThemeTagImg"`

	// ThemeShareCss    string `yaml:"ThemeShareCss"`
	// ThemeShareJs     string `yaml:"ThemeShareJs"`
	// ThemeArchivesJs  string `yaml:"ThemeArchivesJs"`
	// ThemeArchivesCss string `yaml:"ThemeArchivesCss"`

	// AppUrl            string `yaml:"AppUrl"`
	// DefaultIndexLimit string `yaml:"DefaultIndexLimit"`
	//
	// DbUser     string `yaml:"DbUser"`
	// DbPassword string `yaml:"DbPassword"`
	// DbPort     string `yaml:"DbPort"`
	// DbDataBase string `yaml:"DbDataBase"`
	// DbHost     string `yaml:"DbHost"`
	//
	// AlarmType string `yaml:"AlarmType"`
	// MailUser  string `yaml:"MailUser"`
	// MailPwd   string `yaml:"MailPwd"`
	// MailHost  string `yaml:"MailHost"`
	//
	// HashIdSalt   string `yaml:"HashIdSalt"`
	// HashIdLength int    `yaml:"HashIdLength"`
	//
	// JwtIss       string        `yaml:"JwtIss"`
	// JwtAudience  string        `yaml:"JwtAudience"`
	// JwtJti       string        `yaml:"JwtJti"`
	// JwtSecretKey string        `yaml:"JwtSecretKey"`
	// JwtTokenLife time.Duration `yaml:"JwtTokenLife"`
	//
	// RedisAddr string `yaml:"RedisAddr"`
	// RedisPwd  string `yaml:"RedisPwd"`
	// RedisDb   int    `yaml:"RedisDb"`
	//
	// QCaptchaAid       string `yaml:"QCaptchaAid"`
	// QCaptchaSecretKey string `yaml:"QCaptchaSecretKey"`
	//
	// BackUpFilePath string `yaml:"BackUpFilePath"`
	// BackUpDuration string `yaml:"BackUpDuration"`
	// BackUpSentTo   string `yaml:"BackUpSentTo"`
	//
	// ImgUploadUrl          string `yaml:"ImgUploadUrl"`
	// ImgUploadDst          string `yaml:"ImgUploadDst"`
	// ImgUploadBoth         bool   `yaml:"ImgUploadBoth"`
	//
	// // qiniu

	//
}
