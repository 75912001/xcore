package log

//
//import (
//    "context"
//    "fmt"
//    "testing"
//    "time"
//)
//
//// todo menglc 完成覆盖率 [100%] 测试
//
//func TestGetInstance(t *testing.T) {
//    instance1 := GetInstance()
//    instance2 := GetInstance()
//
//    if instance1 != instance2 {
//        t.Errorf("GetInstance() returned different instances")
//    }
//}
//
//func TestIsEnable(t *testing.T) {
//    // 当 stdInstance 为 nil 时，isEnable 应返回 false
//    stdInstance = nil
//    if isEnable() != false {
//        t.Errorf("Expected isEnable to return false, but it returned true")
//    }
//    // 当 stdInstance 不为 nil 时， stdInstance.logChan == nil, 此时 isEnable 应返回 false
//    stdInstance = &mgr{}
//    if isEnable() != false {
//        t.Errorf("Expected isEnable to return false, but it returned true")
//    }
//    // 启动, 使 stdInstance.logChan 不为 nil, 此时 isEnable 应返回 true
//    err := stdInstance.start()
//    if err != nil {
//        t.Errorf("start() returned error: %v", err)
//    }
//    if isEnable() != true {
//        t.Errorf("Expected isEnable to return true, but it returned false")
//    }
//
//    // Reset stdInstance to nil after the test
//    stdInstance = nil
//}
//
//
//func TestNewOptions(t *testing.T) {
//    NewOptions()
//}
//
//func TestPrintErr(t *testing.T) {
//    tests := []struct {
//        name     string
//        input    []interface{}
//        preFunc  func()
//        postFunc func()
//    }{
//        {
//            name:  "日志启用",
//            input: []interface{}{"test"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelFatal)) // LevelFatal is less than LevelError
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//        {
//            name:  "日志-未启用",
//            input: []interface{}{"test"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = nil
//            },
//            postFunc: func() {
//            },
//        },
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            tt.preFunc()
//            PrintErr(tt.input...)
//            tt.postFunc()
//        })
//    }
//}
//
//func Test_entry_Debug(t *testing.T) {
//    tests := []struct {
//        name         string
//        level        int
//        time         time.Time
//        callerInfo   string
//        message      string
//        ctx          context.Context
//        extendFields extendFields
//        preFunc      func()
//        postFunc     func()
//    }{
//        {
//            name:         "normal-日志等级 < 等级",
//            level:        LevelDebug,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelInfo))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//        {
//            name:         "normal- 等级 <= 日志等级",
//            level:        LevelDebug,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelDebug))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            p := &entry{
//                level:        tt.level,
//                time:         tt.time,
//                callerInfo:   tt.callerInfo,
//                message:      tt.message,
//                ctx:          tt.ctx,
//                extendFields: tt.extendFields,
//            }
//            tt.preFunc()
//            p.Debug("this is message:", 123, "xxx", p)
//            tt.postFunc()
//        })
//    }
//}
//
//func Test_entry_Debugf(t *testing.T) {
//    tests := []struct {
//        name         string
//        level        int
//        time         time.Time
//        callerInfo   string
//        message      string
//        ctx          context.Context
//        extendFields extendFields
//        preFunc      func()
//        postFunc     func()
//    }{
//        {
//            name:         "normal-日志等级 < 等级",
//            level:        LevelDebug,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelInfo))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//        {
//            name:         "normal- 等级 <= 日志等级",
//            level:        LevelDebug,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelDebug))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            p := &entry{
//                level:        tt.level,
//                time:         tt.time,
//                callerInfo:   tt.callerInfo,
//                message:      tt.message,
//                ctx:          tt.ctx,
//                extendFields: tt.extendFields,
//            }
//            tt.preFunc()
//            p.Debugf("this is message:%v", 123)
//            tt.postFunc()
//        })
//    }
//}
//
//func Test_entry_Error(t *testing.T) {
//    tests := []struct {
//        name         string
//        level        int
//        time         time.Time
//        callerInfo   string
//        message      string
//        ctx          context.Context
//        extendFields extendFields
//        preFunc      func()
//        postFunc     func()
//    }{
//        {
//            name:         "normal-日志等级 < trace 等级",
//            level:        LevelError,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelFatal))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//        {
//            name:         "normal- trace 等级 <= 日志等级",
//            level:        LevelError,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelError))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            p := &entry{
//                level:        tt.level,
//                time:         tt.time,
//                callerInfo:   tt.callerInfo,
//                message:      tt.message,
//                ctx:          tt.ctx,
//                extendFields: tt.extendFields,
//            }
//            tt.preFunc()
//            p.Error("this is message:", 123)
//            tt.postFunc()
//        })
//    }
//}
//
//func Test_entry_Errorf(t *testing.T) {
//    tests := []struct {
//        name         string
//        level        int
//        time         time.Time
//        callerInfo   string
//        message      string
//        ctx          context.Context
//        extendFields extendFields
//        preFunc      func()
//        postFunc     func()
//    }{
//        {
//            name:         "normal-日志等级 < trace 等级",
//            level:        LevelError,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelFatal))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//        {
//            name:         "normal- trace 等级 <= 日志等级",
//            level:        LevelError,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelError))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            p := &entry{
//                level:        tt.level,
//                time:         tt.time,
//                callerInfo:   tt.callerInfo,
//                message:      tt.message,
//                ctx:          tt.ctx,
//                extendFields: tt.extendFields,
//            }
//            tt.preFunc()
//            p.Errorf("this is message:%v", 123)
//            tt.postFunc()
//        })
//    }
//}
//
//func Test_entry_Fatal(t *testing.T) {
//    tests := []struct {
//        name         string
//        level        int
//        time         time.Time
//        callerInfo   string
//        message      string
//        ctx          context.Context
//        extendFields extendFields
//        preFunc      func()
//        postFunc     func()
//    }{
//        {
//            name:         "normal-日志等级 < trace 等级",
//            level:        LevelFatal,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelOff))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//        {
//            name:         "normal- trace 等级 <= 日志等级",
//            level:        LevelFatal,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelFatal))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            p := &entry{
//                level:        tt.level,
//                time:         tt.time,
//                callerInfo:   tt.callerInfo,
//                message:      tt.message,
//                ctx:          tt.ctx,
//                extendFields: tt.extendFields,
//            }
//            tt.preFunc()
//            p.Fatal("this is message:", 123)
//            tt.postFunc()
//        })
//    }
//}
//
//func Test_entry_Fatalf(t *testing.T) {
//    tests := []struct {
//        name         string
//        level        int
//        time         time.Time
//        callerInfo   string
//        message      string
//        ctx          context.Context
//        extendFields extendFields
//        preFunc      func()
//        postFunc     func()
//    }{
//        {
//            name:         "normal-日志等级 < trace 等级",
//            level:        LevelFatal,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelOff))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//        {
//            name:         "normal- trace 等级 <= 日志等级",
//            level:        LevelFatal,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelFatal))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            p := &entry{
//                level:        tt.level,
//                time:         tt.time,
//                callerInfo:   tt.callerInfo,
//                message:      tt.message,
//                ctx:          tt.ctx,
//                extendFields: tt.extendFields,
//            }
//            tt.preFunc()
//            p.Fatalf("this is message:%v", 123)
//            tt.postFunc()
//        })
//    }
//}
//
//func Test_entry_Info(t *testing.T) {
//    tests := []struct {
//        name         string
//        level        int
//        time         time.Time
//        callerInfo   string
//        message      string
//        ctx          context.Context
//        extendFields extendFields
//        preFunc      func()
//        postFunc     func()
//    }{
//        {
//            name:         "normal-日志等级 < 等级",
//            level:        LevelInfo,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelWarn))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//        {
//            name:         "normal- 等级 <= 日志等级",
//            level:        LevelInfo,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelInfo))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            p := &entry{
//                level:        tt.level,
//                time:         tt.time,
//                callerInfo:   tt.callerInfo,
//                message:      tt.message,
//                ctx:          tt.ctx,
//                extendFields: tt.extendFields,
//            }
//            tt.preFunc()
//            p.Info("this is message:", 123, "xxx", p)
//            tt.postFunc()
//        })
//    }
//}
//
//func Test_entry_Infof(t *testing.T) {
//    tests := []struct {
//        name         string
//        level        int
//        time         time.Time
//        callerInfo   string
//        message      string
//        ctx          context.Context
//        extendFields extendFields
//        preFunc      func()
//        postFunc     func()
//    }{
//        {
//            name:         "normal-日志等级 < 等级",
//            level:        LevelInfo,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelWarn))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//        {
//            name:         "normal- 等级 <= 日志等级",
//            level:        LevelInfo,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelInfo))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            p := &entry{
//                level:        tt.level,
//                time:         tt.time,
//                callerInfo:   tt.callerInfo,
//                message:      tt.message,
//                ctx:          tt.ctx,
//                extendFields: tt.extendFields,
//            }
//            tt.preFunc()
//            p.Infof("this is message:%v", 123)
//            tt.postFunc()
//        })
//    }
//}
//
//func Test_entry_Trace(t *testing.T) {
//    tests := []struct {
//        name         string
//        level        int
//        time         time.Time
//        callerInfo   string
//        message      string
//        ctx          context.Context
//        extendFields extendFields
//        preFunc      func()
//        postFunc     func()
//    }{
//        {
//            name:         "normal-日志等级 < trace 等级",
//            level:        LevelTrace,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelDebug))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//        {
//            name:         "normal- trace 等级 <= 日志等级",
//            level:        LevelTrace,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelTrace))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            p := &entry{
//                level:        tt.level,
//                time:         tt.time,
//                callerInfo:   tt.callerInfo,
//                message:      tt.message,
//                ctx:          tt.ctx,
//                extendFields: tt.extendFields,
//            }
//            tt.preFunc()
//            p.Trace("this is message:", 123, "xxx", p)
//            tt.postFunc()
//        })
//    }
//}
//
//func Test_entry_Tracef(t *testing.T) {
//    tests := []struct {
//        name         string
//        level        int
//        time         time.Time
//        callerInfo   string
//        message      string
//        ctx          context.Context
//        extendFields extendFields
//        preFunc      func()
//        postFunc     func()
//    }{
//        {
//            name:         "normal-日志等级 < trace 等级",
//            level:        LevelTrace,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelDebug))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//        {
//            name:         "normal- trace 等级 <= 日志等级",
//            level:        LevelTrace,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelTrace))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            p := &entry{
//                level:        tt.level,
//                time:         tt.time,
//                callerInfo:   tt.callerInfo,
//                message:      tt.message,
//                ctx:          tt.ctx,
//                extendFields: tt.extendFields,
//            }
//            tt.preFunc()
//            p.Tracef("this is message:%v", 123)
//            tt.postFunc()
//        })
//    }
//}
//
//func Test_entry_Warn(t *testing.T) {
//    tests := []struct {
//        name         string
//        level        int
//        time         time.Time
//        callerInfo   string
//        message      string
//        ctx          context.Context
//        extendFields extendFields
//        preFunc      func()
//        postFunc     func()
//    }{
//        {
//            name:         "normal-日志等级 < trace 等级",
//            level:        LevelWarn,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelError))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//        {
//            name:         "normal- trace 等级 <= 日志等级",
//            level:        LevelWarn,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelWarn))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            p := &entry{
//                level:        tt.level,
//                time:         tt.time,
//                callerInfo:   tt.callerInfo,
//                message:      tt.message,
//                ctx:          tt.ctx,
//                extendFields: tt.extendFields,
//            }
//            tt.preFunc()
//            p.Warn("this is message:", 123)
//            tt.postFunc()
//        })
//    }
//}
//
//func Test_entry_Warnf(t *testing.T) {
//    tests := []struct {
//        name         string
//        level        int
//        time         time.Time
//        callerInfo   string
//        message      string
//        ctx          context.Context
//        extendFields extendFields
//        preFunc      func()
//        postFunc     func()
//    }{
//        {
//            name:         "normal-日志等级 < trace 等级",
//            level:        LevelWarn,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelError))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//        {
//            name:         "normal- trace 等级 <= 日志等级",
//            level:        LevelWarn,
//            time:         time.Now(),
//            callerInfo:   "callerInfo",
//            message:      "normal-LevelFatal-message",
//            ctx:          context.TODO(),
//            extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
//            preFunc: func() {
//                GetInstance()
//                stdInstance = &mgr{}
//                stdInstance.start(NewOptions().WithLevel(LevelWarn))
//            },
//            postFunc: func() {
//                GetInstance().Stop()
//            },
//        },
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            p := &entry{
//                level:        tt.level,
//                time:         tt.time,
//                callerInfo:   tt.callerInfo,
//                message:      tt.message,
//                ctx:          tt.ctx,
//                extendFields: tt.extendFields,
//            }
//            tt.preFunc()
//            p.Warnf("this is message:%v", 123)
//            tt.postFunc()
//        })
//    }
//}
//
//func Test_entry_WithContext(t *testing.T) {
//    stdInstance = new(mgr)
//    _ = GetInstance().Start()
//    _ = newEntry().withContext(context.Background())
//    GetInstance().Stop()
//    return
//}
//
//func Test_entry_WithExtendField(t *testing.T) {
//    stdInstance = new(mgr)
//    _ = GetInstance().Start()
//    _ = newEntry().withExtendField("key", "value")
//    GetInstance().Stop()
//    return
//}
//
//func Test_entry_WithExtendFields(t *testing.T) {
//    stdInstance = new(mgr)
//    _ = GetInstance().Start()
//    fields := extendFields{"key-1", "value-1", "key-2", "value-2"}
//    _ = newEntry().withExtendFields(fields)
//    GetInstance().Stop()
//    return
//}
//
//func Test_entry_WithMessage(t *testing.T) {
//    stdInstance = new(mgr)
//    _ = GetInstance().Start()
//    _ = newEntry().withMessage("message")
//    GetInstance().Stop()
//    return
//}
//
