package producer

import (
	"sync"
)

var (
	AsyncProducerInstance *AsyncProducer
	asyncProducerOnce     sync.Once
)

// AsyncProducer DB异步操作生产者
type AsyncProducer struct {
	workChan chan<- interface{}
	isStart  bool
}

func NewDbAsyncProducer(workerChan chan<- interface{}) *AsyncProducer {
	asyncProducerOnce.Do(func() {
		AsyncProducerInstance = &AsyncProducer{
			workChan: workerChan,
		}
	})
	return AsyncProducerInstance
}

// IsEnableAsync 是否开启异步
func IsEnableAsync() bool {
	return AsyncProducerInstance != nil && AsyncProducerInstance.isStart
}

// Start 启动
func (p *AsyncProducer) Start() {
	p.isStart = true
}

// Stop 关闭
func (p *AsyncProducer) Stop() {
	p.isStart = false
	close(p.workChan)
}

// Produce 生产
func (p *AsyncProducer) Produce(asyncFunc interface{}) {
	p.workChan <- asyncFunc
}

// IsStart 是否启动
func (p *AsyncProducer) IsStart() bool {
	return p.isStart
}
