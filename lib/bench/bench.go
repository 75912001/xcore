// 服务配置
// 服务配置文件, 用于配置服务的基本信息.
// 该配置文件与可执行程序在同一目录下.

package bench

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"path/filepath"
	"runtime"
	"time"
	xcommon "xcore/lib/common"
	xconstants "xcore/lib/constants"
	xerror "xcore/lib/error"
	xetcd "xcore/lib/etcd"
	xlog "xcore/lib/log"
	xruntime "xcore/lib/runtime"
	xtimer "xcore/lib/timer"
)

// 配置-主项,用户服务的基本配置

type Mgr struct {
	RootJson rootJson
	Json     benchJson
}

type rootJson struct {
	Etcd Etcd `json:"etcd"`
}

func (p *rootJson) Parse(strJson string) error {
	if err := json.Unmarshal([]byte(strJson), &p); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Etcd.TTL == nil {
		defaultValue := xetcd.TtlSecondDefault
		p.Etcd.TTL = &defaultValue
	}
	return nil
}

type Etcd struct {
	Addrs []string `json:"addrs"` // [目前无效] todo [优化] 该配置,后期可改为从etcd中获取剩余配置,并覆盖本地配置.
	TTL   *int64   `json:"ttl"`   // ttl 秒 [default]: xetcd.TtlSecondDefault 秒, e.g.:系统每10秒续约一次,该参数至少为11秒
}

type benchJson struct {
	Base       Base               `json:"base"`
	Timer      Timer              `json:"timer"`
	ServiceNet xcommon.ServiceNet `json:"serviceNet"`
}

func (p *benchJson) Parse(jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), p)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.ProjectName == nil {
		defaultValue := xconstants.ProjectNameDefault
		p.Base.ProjectName = &defaultValue
	}
	if p.Base.Version == nil {
		defaultValue := xconstants.VersionDefault
		p.Base.Version = &defaultValue
	}
	if p.Base.LogLevel == nil {
		defaultValue := xlog.LevelOn
		p.Base.LogLevel = &defaultValue
	}
	if p.Base.LogAbsPath == nil {
		executablePath, err := xruntime.GetExecutablePath()
		if err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
		executablePath = filepath.Join(executablePath, "log")
		p.Base.LogAbsPath = &executablePath
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
	if p.Base.AvailableLoad == nil {
		defaultValue := xconstants.AvailableLoadDefault
		p.Base.AvailableLoad = &defaultValue
	}
	if p.Timer.ScanSecondDuration == nil {
		defaultValue := xtimer.ScanSecondDurationDefault
		p.Timer.ScanSecondDuration = &defaultValue
	}
	if p.Timer.ScanMillisecondDuration == nil {
		defaultValue := xtimer.ScanMillisecondDurationDefault
		p.Timer.ScanMillisecondDuration = &defaultValue
	}
	if p.ServiceNet.Addr == nil {
		defaultValue := ""
		p.ServiceNet.Addr = &defaultValue
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
	ProjectName        *string `json:"projectName"`        // 项目名称. [default]: xconstants.ProjectNameDefault
	Version            *string `json:"version"`            // 版本号. [default]: xconstants.VersionDefault
	PprofHttpPort      *uint16 `json:"pprofHttpPort"`      // pprof性能分析 http端口 [default]: nil 不使用
	LogLevel           *uint32 `json:"logLevel"`           // 日志等级 [default]: xlog.LevelOn
	LogAbsPath         *string `json:"logAbsPath"`         // 日志绝对路径 [default]: 当前执行的程序-绝对路径,指向启动当前进程的可执行文件-目录路径. e.g.:absPath/log
	GoMaxProcess       *int    `json:"goMaxProcess"`       // [default]: runtime.NumCPU()
	BusChannelCapacity *uint32 `json:"busChannelCapacity"` // 总线chan容量. [default]: xconstants.BusChannelCapacityDefault
	PacketLengthMax    *uint32 `json:"packetLengthMax"`    // bytes,用户 上行 每个包的最大长度. [default]: xconstants.PacketLengthDefault
	SendChanCapacity   *uint32 `json:"sendChanCapacity"`   // bytes,每个TCP链接的发送chan大小. [default]: xconstants.SendChanCapacityDefault
	RunMode            *uint32 `json:"runMode"`            // 运行模式 [default]: xruntime.RunModeRelease
	AvailableLoad      *uint32 `json:"availableLoad"`      // 剩余可用负载, 可用资源数 [default]: xconstants.AvailableLoadDefault
}

type Timer struct {
	// 秒级定时器 扫描间隔(纳秒) 1000*1000*100=100000000 为100毫秒 [default]: xtimer.ScanSecondDurationDefault
	ScanSecondDuration *time.Duration `json:"scanSecondDuration"`
	// 毫秒级定时器 扫描间隔(纳秒) 1000*1000*100=100000000 为25毫秒 [default]: xtimer.ScanMillisecondDurationDefault
	ScanMillisecondDuration *time.Duration `json:"scanMillisecondDuration"`
}
