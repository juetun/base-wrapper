package snowflake

import (
	"log"
	"strconv"
	"strings"

	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/sony/sonyflake"
)

var (
	SFlake *SnowFlake
)

// SnowFlake SnowFlake算法结构体
type SnowFlake struct {
	sFlake *sonyflake.Sonyflake
}

type OperateSnowFlakeHandler func() (uint16, error)

var SnoFlakeHandler OperateSnowFlakeHandler

func init() {
	SFlake = newSnowFlake()
	SnoFlakeHandler = getMachineID
}

// 模拟获取本机的机器ID
func getMachineID() (mID uint16, err error) {
	var (
		ipNum int64
		ip    string
	)
	if ip, err = utils.GetLocalIP(); err != nil {
		return
	}
	if ipNum, err = strconv.ParseInt(strings.ReplaceAll(ip, ".", ""), 10, 64); err != nil {
		return
	}
	mID = uint16(ipNum % 32)
	log.Printf("ip:%s,mechineId:%d\n",ip,mID)
	return
}

func newSnowFlake() *SnowFlake {
	st := sonyflake.Settings{}
	// machineID是个回调函数
	st.MachineID = SnoFlakeHandler
	return &SnowFlake{
		sFlake: sonyflake.NewSonyflake(st),
	}
}

func (s *SnowFlake) GetID() (uint64, error) {
	return s.sFlake.NextID()
}
