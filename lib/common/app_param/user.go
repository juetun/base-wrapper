package app_param

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
)

// 获取用户信息对应的表
const (
	UserDataTypeEmail  = "user_email"
	UserDataTypeMain   = "user_main"
	UserDataTypeInfo   = "user_info"
	UserDataTypeMobile = "user_mobile"
)

// 获取用户信息的响应参数结构
type (
	ResultUser struct {
		List map[int64]ResultUserItem `json:"list"`
	}
	ResultUserItem struct {
		UserHid          int64      `json:"user_hid,omitempty"`  // 用户ID
		Portrait         string     `json:"portrait,omitempty"`  // 头像
		NickName         string     `json:"nick_name,omitempty"` // 昵称
		UserName         string     `json:"user_name,omitempty"` // 用户名
		RealName         string     `json:"real_name"`           // 真实姓名
		Gender           int        `json:"gender,omitempty"`    //
		Status           int        `json:"status,omitempty"`    //
		Score            int        `json:"score,omitempty"`     //
		AuthDesc         string     `json:"auth_desc,omitempty"` // 认证描述
		IsV              int        `json:"is_v,omitempty"`      // 用户头像加V
		Remark           string     `json:"remark" `             // 个性签名
		Signature        string     `json:"signature,omitempty"`
		RegisterChannel  string     `json:"register_channel,omitempty"`
		CountryCode      string     `json:"country_code,omitempty"`
		Mobile           string     `json:"mobile,omitempty"`
		MobileVerifiedAt *time.Time `json:"mobile_verified_at,omitempty"`
		Email            string     `json:"email,omitempty"`
		EmailVerifiedAt  *time.Time `json:"email_verified_at,omitempty"`
		ShopId           string     `json:"shop_id"`
	}
	RequestUser struct {
		UserHid           int64           `json:"user_hid" form:"user_hid"`
		UserMobileIndex   string          `json:"user_mobile_index" form:"user_mobile_index"`
		UserEmailIndex    string          `json:"user_email_index" form:"user_email_index"`
		Portrait          string          `json:"portrait" form:"portrait"`
		NickName          string          `json:"nick_name" form:"nick_name"`
		UserName          string          `json:"user_name" form:"user_name"`
		Gender            int             `json:"gender" form:"gender"`
		Status            int             `json:"status" form:"status"`
		Score             int             `json:"score" form:"score"`
		RememberToken     string          `json:"remember_token" form:"remember_token"`
		MsgReadTimeCursor base.TimeNormal `json:"msg_read_time_cursor" form:"msg_read_time_cursor"`
		ShopId            string          `json:"shop_id"`
	}

	User struct {
		UserHid    int64       `json:"user_hid"`
		UserIndex  *UserIndex  `json:"user_index,omitempty"`
		UserMain   *UserMain   `json:"user_main,omitempty"`
		UserEmail  *UserEmail  `json:"user_email,omitempty"`
		UserInfo   *UserInfo   `json:"user_info,omitempty"`
		UserMobile *UserMobile `json:"user_mobile,omitempty"`
	}
	UserIndex struct {
		ID         int64            `gorm:"column:id;primary_key" json:"-"`
		UserName   string           `gorm:"column:user_name;not null;type:varchar(50) COLLATE utf8mb4_bin;uniqueIndex;comment:用户名" json:"user_name" `
		TmpAccount string           `gorm:"column:tmp_account;not null;type:varchar(200) COLLATE utf8mb4_bin;comment:注册时临时账号" json:"tmp_account" `
		IsUse      int              `json:"is_use" gorm:"column:is_use;type:tinyint(1);default:0;comment:是否启用 0-启用 大于0-已启用"`
		CreatedAt  base.TimeNormal  `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;" `
		UpdatedAt  base.TimeNormal  `json:"updated_at" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP;" `
		DeletedAt  *base.TimeNormal `json:"deleted_at" gorm:"column:deleted_at;" `
	}
	UserInfo struct {
		ID                int             `gorm:"column:id;primary_key" json:"id"`
		RealName          string          `gorm:"column:real_name;type:varchar(60);not null;comment:真实姓名"  json:"real_name"`
		UserIndexHid      string          `gorm:"column:user_index_hid;type:varchar(60);not null;comment:user_main表位置" json:"user_index_hid"`
		UserHid           int64           `gorm:"column:user_hid;not null;uniqueIndex:idx_userhid;default:0;type:bigint(20) COLLATE utf8mb4_bin" json:"user_hid"`
		RememberToken     string          `gorm:"column:remember_token;not null;default:'';size:500;comment:登录的token" json:"remember_token"`
		MsgReadTimeCursor base.TimeNormal `gorm:"column:msg_read_time_cursor;not null;default:CURRENT_TIMESTAMP;comment:最近一次读取系统公告时间" json:"msg_read_time_cursor"`
		Level             string          `gorm:"column:level;not null;type:tinyint(2);default:0;comment:用户等级0-普通用户" json:"level"`
		Remark            string          `json:"remark" gorm:"column:remark;not null;type:varchar(150);default:'';comment:个性签名"` // 个性签名
		Password          string          `gorm:"column:password;not null;type:varchar(256) COLLATE utf8mb4_general_ci;comment:密码" json:"password"`
		IdCard            string          `gorm:"column:id_card;not null;type:varchar(50) COLLATE utf8mb4_general_ci;comment:身份证号加密串" json:"id_card"`
		IdCardSuffix      string          `gorm:"column:id_card_suffix;not null;type:varchar(50) COLLATE utf8mb4_general_ci;comment:后缀后6位字符" json:"id_card_suffix"`
		QQ                string          `gorm:"column:qq;not null;type:varchar(50) COLLATE utf8mb4_general_ci;comment:qq" json:"qq"`
		WeiXin            string          `gorm:"column:wei_xin;not null;type:varchar(50) COLLATE utf8mb4_general_ci;comment:微信账号" json:"wei_xin"`
		DingDing          string          `gorm:"column:ding_ding;not null;type:varchar(50) COLLATE utf8mb4_general_ci;comment:钉钉账号" json:"ding_ding"`
		WeiBo             string          `gorm:"column:wei_bo;not null;type:varchar(50) COLLATE utf8mb4_general_ci;comment:微博账号" json:"wei_bo"`
		Signature         string          `gorm:"column:signature;not null;type:varchar(256) COLLATE utf8mb4_general_ci;comment:用户签名" json:"signature"`
		RegisterChannel   string          `gorm:"column:register_channel;not null;type:varchar(50) COLLATE utf8mb4_general_ci;comment:账号注册渠道" json:"register_channel"`
		InviteCode        int64           `gorm:"column:invite_code;not null;default:0;type:int(10);comment:邀请码" json:"invite_code"`
	}
	UserMain struct {
		ID              int              `gorm:"column:id;primary_key" json:"id"`
		UserHid         int64            `gorm:"uniqueIndex:idx_user_hid;column:user_hid;not null;default:0;type:bigint(20) COLLATE utf8mb4_bin" json:"user_hid"` // sql:"unique_index" 创建表时生成唯一索引
		AuthDesc        string           `json:"auth_desc" gorm:"column:auth_desc;not null;type:varchar(30);default:'';comment:认证描述"`                             // 认证描述
		UserMobileIndex string           `gorm:"column:user_mobile_index;not null;type:varchar(60) COLLATE utf8mb4_bin;default:'';comment:手机号索引" json:"-" `
		UserEmailIndex  string           `gorm:"column:user_email_index;not null;type:varchar(60) COLLATE utf8mb4_bin;default:'';comment:邮箱索引" json:"-" `
		Portrait        string           `gorm:"column:portrait;not null;type:varchar(1000);default:'';comment:头图地址;" json:"portrait"`
		PortraitStatus  int              `gorm:"column:portrait_status;not null;type:varchar(10);default:'';comment:用户审核状态从右向左每位依次昵称-头像;" json:"portrait_status"`
		NickName        string           `gorm:"column:nick_name;not null;type:varchar(30);default:'';comment:昵称" json:"nick_name"`
		UserName        string           `gorm:"column:user_name;not null;size:30;default:'';comment:用户名" json:"user_name" `
		Gender          int              `gorm:"column:gender;not null;type:tinyint(1);default:0;comment:性别 0-男 1-女" json:"gender"`
		Status          int              `gorm:"column:status;not null;type:tinyint(1);default:0;comment:状态 0-可用 1-不可用" json:"status"`
		Score           int              `gorm:"column:score;not null;type:int(10);default:0;comment:用户积分" json:"score"`
		Balance         float64          `gorm:"column:balance;not null;type:decimal(10,2);default:0;comment:用户账户余额" json:"balance"`
		CurrentShopId   int64            `gorm:"column:current_shop_id;not null;default:0;comment:当前店铺ID" json:"current_shop_id"`
		Country         string           `gorm:"column:country;not null;type:varchar(30) COLLATE utf8mb4_general_ci;comment:国籍" json:"country"`
		CityId          int              `gorm:"column:city_id;not null;type:varchar(30) COLLATE utf8mb4_general_ci;comment:所在城市" json:"city_id"`
		OrgCode         string           `gorm:"column:org_code;not null;type:varchar(180) COLLATE utf8mb4_bin;comment:机构号" json:"org_code"`
		OrgRoot         string           `gorm:"column:org_root;not null;type:varchar(32) COLLATE utf8mb4_bin;comment:机构号" json:"org_root"`
		IsV             int              `json:"is_v" gorm:"column:is_v;not null;type:tinyint(1);default:0;comment:用户头像加V 0-不加 1-加"` // 用户头像加V
		CreatedAt       base.TimeNormal  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at" `
		UpdatedAt       base.TimeNormal  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at" `
		DeletedAt       *base.TimeNormal `gorm:"column:deleted_at;" json:"deleted_at"`
	}

	UserEmail struct {
		ID              int        `gorm:"column:id;primary_key" json:"-"`
		UserHid         int64      `json:"user_hid" gorm:"column:user_hid;uniqueIndex:idx_userhid,priority:1;type:bigint(20);default:0;not null;"`
		UserIndexHid    string     `json:"user_index_hid" gorm:"column:user_index_hid;type:varchar(60);not null;comment:user_main表位置"`
		Email           string     `gorm:"column:email;uniqueIndex:idx_email,priority:1;not null;type:varchar(100);default:0;comment:邮箱" json:"-"`
		EmailVerifiedAt *time.Time `gorm:"column:email_verified_at;not null;uniqueIndex:idx_email,priority:3;type:datetime;default:'2000-01-01 00:00:00';comment:认证时间;" json:"-"`
		IsDel           int        `json:"is_del" gorm:"column:is_del;uniqueIndex:idx_userhid,priority:2;uniqueIndex:idx_email,priority:2;type:tinyint(2);default:0;comment:是否删除0-未删除 大于0-已删除"`
	}
	UserMobile struct {
		ID               int        `gorm:"column:id;primary_key" json:"-"`
		UserHid          int64      `json:"user_hid" gorm:"column:user_hid;uniqueIndex:inx_userhid,priority:1;not null;default:0;type:bigint(20) COLLATE utf8mb4_bin;"`
		UserIndexHid     string     `json:"user_index_hid" gorm:"column:user_index_hid;not null;default:'';type:varchar(60) COLLATE utf8mb4_bin;comment:user_main表位置"`
		CountryCode      string     `gorm:"column:country_code;uniqueIndex:idx_mobile,priority:2;type:varchar(15) COLLATE utf8mb4_bin;not null;comment:手机国别默认86" json:"country_code"`
		Mobile           string     `gorm:"column:mobile;not null;default:'';uniqueIndex:idx_mobile,priority:1;type:varchar(20) COLLATE utf8mb4_bin;comment:手机号" json:"-"`
		MobileVerifiedAt *time.Time `json:"mobile_verified_at" gorm:"column:mobile_verified_at;not null;uniqueIndex:idx_mobile,priority:4;default:'2000-01-01 00:00:00'"`
		IsDel            int        `json:"is_del" gorm:"column:is_del;type:tinyint(2);uniqueIndex:inx_userhid,priority:2;uniqueIndex:idx_mobile,priority:3;not null;idx_mobile,priority:1;default:0;comment:是否删除0-未删除 大于0-已删除"`
	}
)

// GetRealName 获取用户的真实姓名
func (r *ResultUserItem) GetRealName(nilDefaultValue ...string) (res string) {
	if r.RealName != "" {
		res = r.RealName
		return
	}
	if len(nilDefaultValue) > 0 {
		res = nilDefaultValue[0]
	}
	return
}

func (r *ResultUserItem) InitData(item *User) {
	r.UserHid = item.UserHid
	if item.UserMain != nil {
		r.AuthDesc = item.UserMain.AuthDesc
		r.Portrait = item.UserMain.Portrait
		r.NickName = item.UserMain.NickName
		r.Gender = item.UserMain.Gender
		r.Status = item.UserMain.Status
		r.Score = item.UserMain.Score
		r.IsV = item.UserMain.IsV
	}
	if item.UserInfo != nil {
		r.Signature = item.UserInfo.Signature
		r.Remark = item.UserInfo.Remark
		r.RegisterChannel = item.UserInfo.RegisterChannel
		r.RealName = item.UserInfo.RealName
	}

	if item.UserEmail != nil {
		r.Email = item.UserEmail.Email
		r.EmailVerifiedAt = item.UserEmail.EmailVerifiedAt
	}
	if item.UserMobile != nil {
		r.Mobile = item.UserMobile.Mobile
		r.CountryCode = item.UserMobile.CountryCode
	}
}

func (r *RequestUser) InitRequestUser(c *gin.Context) (err error) {
	if r.UserHid == 0 {
		uidString := c.GetHeader(app_obj.HttpUserHid)
		if uidString == "" {
			err = fmt.Errorf("请先登录系统")
			return
		}
		r.UserHid, err = strconv.ParseInt(uidString, 10, 64)
		if err != nil {
			err = fmt.Errorf("用户信息参数格式不正确(uid:%s)", uidString)
			return
		}
	}
	return
}
