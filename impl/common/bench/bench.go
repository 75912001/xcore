// 服务配置
// 服务配置文件, 用于配置服务的基本信息.
// 该配置文件与可执行程序在同一目录下.

package bench

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"runtime"
	"time"
	"xcore/impl/common"
	xconstants "xcore/lib/constants"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	xruntime "xcore/lib/runtime"
)

// 配置-主项,用户服务的基本配置

type Mgr struct {
	RootJson rootJson
	Json     benchJson
}

type rootJson struct {
	Etcd Etcd `json:"etcd"` // todo menglc [优化] 该配置,后期可改为从etcd中获取剩余配置,并覆盖本地配置.
}

func (p *rootJson) Parse(strJson string) error {
	if err := json.Unmarshal([]byte(strJson), &p); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}

	if p.Etcd.TTL == nil {
		defaultValue := common.TtlSecondDefault
		p.Etcd.TTL = &defaultValue
	}
	return nil
}

type Etcd struct {
	Addrs []string `json:"addrs"`
	TTL   *int64   `json:"ttl"` // ttl 秒 [default]: common.TtlSecondDefault 秒, e.g.:系统每10秒续约一次,该参数至少为11秒
}

type benchJson struct {
	Base       Base       `json:"base"`
	Timer      Timer      `json:"timer"`
	ServiceNet ServiceNet `json:"serviceNet"`
}

func (p *benchJson) Parse(pathFile string) error {
	file, err := os.Open(pathFile)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(p)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.Version == nil {
		defaultValue := "0.0.1.beta.2024.09.03.2034"
		p.Base.Version = &defaultValue
	}
	if p.Base.PprofHttpPort == nil {
		defaultValue := uint16(0)
		p.Base.PprofHttpPort = &defaultValue
	}
	if p.Base.LogLevel == nil {
		defaultValue := xlog.LevelOn
		p.Base.LogLevel = &defaultValue
	}
	if p.Base.LogAbsPath == nil {
		defaultValue := common.LogAbsPath
		p.Base.LogAbsPath = &defaultValue
	}
	if p.Base.GoMaxProcess == nil {
		defaultValue := runtime.NumCPU()
		p.Base.GoMaxProcess = &defaultValue
	}
	if p.Base.BusChannelCapacity == nil {
		defaultValue := xconstants.BusChannelCapacityDefault
		p.Base.BusChannelCapacity = &defaultValue
	}
	if p.Base.PacketLengthMax == nil {
		defaultValue := xconstants.PacketLengthDefault
		p.Base.PacketLengthMax = &defaultValue
	}
	if p.Base.SendChanCapacity == nil {
		defaultValue := xconstants.SendChanCapacityDefault
		p.Base.SendChanCapacity = &defaultValue
	}
	if p.Base.RunMode == nil {
		defaultValue := uint32(xruntime.RunModeRelease)
		p.Base.RunMode = &defaultValue
	}
	if p.Timer.ScanSecondDuration == nil {
		t := time.Millisecond * 100
		p.Timer.ScanSecondDuration = &t
	}
	if p.Timer.ScanMillisecondDuration == nil {
		t := time.Millisecond * 25
		p.Timer.ScanMillisecondDuration = &t
	}
	if p.ServiceNet.Domain == nil {
		defaultValue := ""
		p.ServiceNet.Domain = &defaultValue
	}
	if p.ServiceNet.IP == nil {
		defaultValue := ""
		p.ServiceNet.IP = &defaultValue
	}
	if p.ServiceNet.Port == nil {
		defaultValue := uint16(0)
		p.ServiceNet.Port = &defaultValue
	}
	if p.ServiceNet.Type == nil {
		defaultValue := "tcp"
		p.ServiceNet.Type = &defaultValue
	}
	if *p.ServiceNet.Type != "tcp" && *p.ServiceNet.Type != "udp" {
		return xerror.NotImplemented.WithExtraMessage(fmt.Sprintf("serviceNet.type must be tcp or udp. %x", xruntime.Location()))
	}
	return nil
}

type Base struct {
	Version            *string `json:"version"`            // 版本号. [default]: 0.0.1.beta.2024.09.03.2034
	PprofHttpPort      *uint16 `json:"pprofHttpPort"`      // pprof性能分析 http端口 [default]:nil 不使用
	LogLevel           *uint32 `json:"logLevel"`           // 日志等级 [default]: xlog.LevelOn
	LogAbsPath         *string `json:"logAbsPath"`         // 日志绝对路径 [default]: common.LogAbsPath
	GoMaxProcess       *int    `json:"goMaxProcess"`       // [default]: runtime.NumCPU()
	BusChannelCapacity *uint32 `json:"busChannelCapacity"` // 总线chan容量. [default]: xconstants.BusChannelCapacityDefault
	PacketLengthMax    *uint32 `json:"packetLengthMax"`    // bytes,用户 上行 每个包的最大长度. [default]:8192
	SendChanCapacity   *uint32 `json:"sendChanCapacity"`   // bytes,每个TCP链接的发送chan大小. [default]:1000
	RunMode            *uint32 `json:"runMode"`            // 运行模式 [default]:0,release
}

type Timer struct {
	// 秒级定时器 扫描间隔(纳秒) 1000*1000*100=100000000 为100毫秒 [default]:100000000
	ScanSecondDuration *time.Duration `json:"scanSecondDuration"`
	// 毫秒级定时器 扫描间隔(纳秒) 1000*1000*100=100000000 为25毫秒 [default]:25000000
	ScanMillisecondDuration *time.Duration `json:"scanMillisecondDuration"`
}

type ServiceNet struct {
	Domain *string `json:"domain"` // If domain is configured, it is used preferentially
	IP     *string `json:"ip"`
	Port   *uint16 `json:"port"`
	Type   *string `json:"type"` // [tcp, udp] [default]: tcp
}
