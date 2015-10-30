package help

import (
	"fmt"
)

type IEvent interface {
	Init(name string, t uint64)      // 初始化(name可以为空, t是触发时间)
	GetName() string                 // 获取别名
	Exec() bool                      // 执行事件
	AddNode(n *DListNode) bool       // 增加节点
	Destroy()                        // 摧毁事件
	Pop()                            // 弹出事件
	GetTouchTime() uint64            // 获取定时器触发时间戳
	SetTouchTime(t uint64)           // 设置定时器时间戳
	SetDelayTime(d uint64, c uint64) // 设置定时器相对时间, c是当前时间戳
	PrintSelf()                      // 打印自己
}

const (
	Evt_node_count = 2
)

type Evt_base struct {
	Nodes      [Evt_node_count]*DListNode
	Name       string // 名称
	Touch_time uint64 // 定时器触发时间戳
}

func (this *Evt_base) Init(name string, t uint64) {
	this.Name = name
	this.Touch_time = t
}

func (this *Evt_base) GetName() string {
	return this.Name
}

func (this *Evt_base) AddNode(n *DListNode) bool {
	for i := 0; i < Evt_node_count; i++ {
		if this.Nodes[i] == nil {
			this.Nodes[i] = n
			return true
		}
	}
	return false
}

func (this *Evt_base) Destroy() {
	this.Pop()
	for i := 0; i < Evt_node_count; i++ {
		this.Nodes[i] = nil
	}
}

func (this *Evt_base) Pop() {
	for i := 0; i < Evt_node_count; i++ {
		if this.Nodes[i] != nil {
			this.Nodes[i].Pop()
		}
	}
}

func (this *Evt_base) GetTouchTime() uint64 {
	return this.Touch_time
}

func (this *Evt_base) SetTouchTime(t uint64) {
	this.Touch_time = t
}

func (this *Evt_base) SetDelayTime(d uint64, c uint64) {
	this.Touch_time = c + d
}

func (this *Evt_base) PrintSelf() {
	fmt.Println("  {E} Is evt base")
}
