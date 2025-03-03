package protocol

type LogData struct {
	Content string `json:"content"` // 日志内容
	Stream  string `json:"stream"`  // 日志流类型（stdout/stderr）
}
