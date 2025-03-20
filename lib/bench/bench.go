// 服务配置
// 服务配置文件, 用于配置服务的基本信息.
// 该配置文件与可执行程序在同一目录下.

package bench

import (
	"encoding/json"
	"fmt"
	xcommon "github.com/75912001/xcore/lib/common"
	xerror "github.com/75912001/xcore/lib/error"
	xetcd "github.com/75912001/xcore/lib/etcd"
	xruntime "github.com/75912001/xcore/lib/runtime"
	xtimer "github.com/75912001/xcore/lib/timer"
	"github.com/pkg/errors"
	"path/filepath"
	"runtime"
	"time"
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
	Addrs []string `json:"addrs"` // etcd地址
	TTL   *int64   `json:"ttl"`   // ttl 秒 [default]: xetcd.TtlSecondDefault 秒, e.g.:系统每10秒续约一次,该参数至少为11秒
}

type benchJson struct {
	Base      Base               `json:"base"`
	Timer     Timer              `json:"timer"`
	ServerNet xcommon.ServiceNet `json:"serverNet"`
}

func (p *benchJson) Parse(jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), p)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.ProjectName == nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.Version == nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.LogLevel == nil {
		return errors.WithMessage(err, xruntime.Location())
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
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.PacketLengthMax == nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.SendChanCapacity == nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.RunMode == nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.AvailableLoad == nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Timer.ScanSecondDuration == nil {
		defaultValue := xtimer.ScanSecondDurationDefault
		p.Timer.ScanSecondDuration = &defaultValue
	}
	if p.Timer.ScanMillisecondDuration == nil {
		defaultValue := xtimer.ScanMillisecondDurationDefault
		p.Timer.ScanMillisecondDuration = &defaultValue
	}
	if p.ServerNet.Addr == nil {
		defaultValue := ""
		p.ServerNet.Addr = &defaultValue
	}
	if p.ServerNet.Type == nil {
		defaultValue := "tcp"
		p.ServerNet.Type = &defaultValue
	}
	if *p.ServerNet.Type != "tcp" && *p.ServerNet.Type != "udp" {
		return xerror.NotImplemented.WithExtraMessage(fmt.Sprintf("serviceNet.type must be tcp or udp. %x", xruntime.Location()))
	}
	return nil
}

type Base struct {
	ProjectName        *string `json:"projectName"`        // 项目名称
	Version            *string `json:"version"`            // 版本号
	PprofHttpPort      *uint16 `json:"pprofHttpPort"`      // pprof性能分析 http端口 [default]: nil 不使用
	LogLevel           *uint32 `json:"logLevel"`           // 日志等级
	LogAbsPath         *string `json:"logAbsPath"`         // 日志绝对路径 [default]: 当前执行的程序-绝对路径,指向启动当前进程的可执行文件-目录路径. e.g.:absPath/log
	GoMaxProcess       *int    `json:"goMaxProcess"`       // [default]: runtime.NumCPU()
	BusChannelCapacity *uint32 `json:"busChannelCapacity"` // 总线chan容量
	PacketLengthMax    *uint32 `json:"packetLengthMax"`    // bytes,用户 上行 每个包的最大长度
	SendChanCapacity   *uint32 `json:"sendChanCapacity"`   // bytes,每个TCP链接的发送chan大小
	RunMode            *uint32 `json:"runMode"`            // 运行模式 [0:release 1:debug]
	AvailableLoad      *uint32 `json:"availableLoad"`      // 剩余可用负载, 可用资源数
}

type Timer struct {
	// 秒级定时器 扫描间隔(纳秒) 1000*1000*100=100000000 为100毫秒 [default]: xtimer.ScanSecondDurationDefault
	ScanSecondDuration *time.Duration `json:"scanSecondDuration"`
	// 毫秒级定时器 扫描间隔(纳秒) 1000*1000*100=100000000 为25毫秒 [default]: xtimer.ScanMillisecondDurationDefault
	ScanMillisecondDuration *time.Duration `json:"scanMillisecondDuration"`
}
