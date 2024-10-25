package control

// IParameters 参数
type IParameters interface {
	Set(parameters ...interface{})   // 设置参数 [全部重置]
	Get() []interface{}              // 获取参数 [全部]
	Append(parameter ...interface{}) // 追加参数
}
