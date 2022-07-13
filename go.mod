module github.com/juetun/base-wrapper

go 1.15

require (
	github.com/ClickHouse/clickhouse-go v1.5.1
	github.com/Tang-RoseChild/mahonia v0.0.0-20131226213531-0eef680515cc
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/asim/go-micro/plugins/server/http/v3 v3.0.0-20210623064501-212df8e6c359
	github.com/asim/go-micro/v3 v3.5.1
	github.com/astaxie/beego v1.12.3
	github.com/elastic/go-elasticsearch/v7 v7.13.1
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-gonic/gin v1.7.2
	github.com/go-errors/errors v1.4.0
	github.com/go-redis/redis/v8 v8.10.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.2.0
	github.com/gorilla/websocket v1.4.2
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/mojocn/base64Captcha v1.3.4
	github.com/robfig/cron/v3 v3.0.1
	github.com/satori/go.uuid v1.2.0
	github.com/shopspring/decimal v1.3.1
	github.com/sirupsen/logrus v1.8.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/sony/sonyflake v1.0.0
	github.com/speps/go-hashids/v2 v2.0.1
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.7.0
	github.com/wumansgy/goEncrypt v0.0.0-20210730092718-e359121aa81e
	go.etcd.io/etcd/client/v3 v3.5.4
	go.uber.org/zap v1.18.1 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/tools v0.1.4 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.11
)

replace (
	github.com/coreos/bbolt v1.3.6 => go.etcd.io/bbolt v1.3.6
	go.etcd.io/bbolt v1.3.6 => github.com/coreos/bbolt v1.3.6
)
