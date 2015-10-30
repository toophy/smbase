package config

const (
	Evt_gap_time  = 16     // 心跳时间(毫秒)
	Evt_gap_bit   = 4      // 心跳时间对应得移位(快速运算使用)
	Evt_lay1_time = 160000 // 第一层事件池最大支持时间(毫秒)
)

const (
	UpdateCurrTimeCount = 32 // 刷新时间戳变更上线
)
