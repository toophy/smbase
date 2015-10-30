package app

import (
	"github.com/toophy/##[AppName]##/help"
	"time"
)

const (
	Evt_gap_time  = 16     // 心跳时间(毫秒)
	Evt_gap_bit   = 4      // 心跳时间对应得移位(快速运算使用)
	Evt_lay1_time = 160000 // 第一层事件池最大支持时间(毫秒)
)

const (
	UpdateCurrTimeCount = 32 // 刷新时间戳变更上线
)

// 投递定时器事件
func (this *AppBase) PostEvent(a help.IEvent) bool {
	check_name := len(a.GetName()) > 0
	if check_name {
		if _, ok := this.evt_names[a.GetName()]; ok {
			return false
		}
	}

	if a.GetTouchTime() < 0 {
		return false
	}

	// 计算放在那一层
	pos := (a.GetTouchTime() + Evt_gap_time - 1) >> Evt_gap_bit
	if pos < 0 {
		pos = 1
	}

	var header *help.DListNode

	if pos < this.evt_lay1Size {
		new_pos := this.evt_lay1Cursor + pos
		if new_pos >= this.evt_lay1Size {
			new_pos = new_pos - this.evt_lay1Size
		}
		pos = new_pos
		header = &this.evt_lay1[pos]
	} else {
		if _, ok := this.evt_lay2[pos]; !ok {
			this.evt_lay2[pos] = new(help.DListNode)
			this.evt_lay2[pos].Init(nil)
		}
		header = this.evt_lay2[pos]
	}

	if header == nil {
		return false
	}

	n := &help.DListNode{}
	n.Init(a)

	if !a.AddNode(n) {
		return false
	}

	old_pre := header.Pre

	header.Pre = n
	n.Next = header
	n.Pre = old_pre
	old_pre.Next = n

	if check_name {
		this.evt_names[a.GetName()] = a
	}

	return true
}

// 通过别名获取事件
func (this *AppBase) GetEvent(name string) help.IEvent {
	if _, ok := this.evt_names[name]; ok {
		return this.evt_names[name]
	}
	return nil
}

func (this *AppBase) RemoveEvent(e help.IEvent) {
	delete(this.evt_names, e.GetName())
	e.Destroy()
}

// 运行一次定时器事件(一个线程心跳可以处理多次)
func (this *AppBase) runEvents() {
	all_time := (this.last_time - this.start_time) / int64(time.Millisecond)

	all_count := uint64((all_time + Evt_gap_time - 1) >> Evt_gap_bit)

	for i := this.evt_lastRunCount; i <= all_count; i++ {
		// 执行第一层事件
		this.runExec(&this.evt_lay1[this.evt_lay1Cursor])

		// 执行第二层事件
		if _, ok := this.evt_lay2[this.evt_currRunCount]; ok {
			this.runExec(this.evt_lay2[this.evt_currRunCount])
			delete(this.evt_lay2, this.evt_currRunCount)
		}

		this.evt_currRunCount++
		this.evt_lay1Cursor++
		if this.evt_lay1Cursor >= this.evt_lay1Size {
			this.evt_lay1Cursor = 0
		}
	}

	this.evt_lastRunCount = this.evt_currRunCount
}

// 运行一条定时器事件链表, 每次都执行第一个事件, 直到链表为空
func (this *AppBase) runExec(header *help.DListNode) {
	for {
		// 每次得到链表第一个事件(非)
		n := header.Next
		if n.IsEmpty() {
			break
		}

		d := n.Data.(help.IEvent)

		// 执行事件, 返回true, 删除这个事件, 返回false表示用户自己处理
		if d.Exec() {
			this.RemoveEvent(d)
		} else if header.Next == n {
			// 防止使用者没有删除使用过的事件, 造成死循环, 该事件, 用户要么重新投递到其他链表, 要么删除
			this.RemoveEvent(d)
		}
	}
}

// 节点池 : 新建节点
func (this *AppBase) newDlinkNode() *help.DListNode {

	if this.node_free.IsEmpty() {
		this.LogFatal("newDlinkNode nil(%d)", this.node_alloc_count)
		return nil
	}

	this.node_alloc_count++
	free := this.node_free.Next
	free.Pop()

	return free
}

// 节点池 : 释放节点
func (this *AppBase) releaseDlinkNode(d *help.DListNode) {
	if d == nil || d.Next == nil {
		return
	}

	// 释放一串
	if !d.IsEmpty() {
		header_pre := d.Pre
		header_next := d.Next

		d.Init(nil)

		old_pre := this.node_free.Pre

		this.node_free.Pre = header_pre
		header_pre.Next = &this.node_free

		header_next.Pre = old_pre
		old_pre.Next = header_next
	}
}

// 节点池 : 增加自由节点
func (this *AppBase) addFreeDlinkNode(n *help.DListNode) {
	old_pre := this.node_free.Pre

	this.node_free.Pre = n
	n.Next = &this.node_free
	n.Pre = old_pre
	old_pre.Next = n
}

// 获取当前时间戳(毫秒)
func (this *AppBase) GetCurrTime() int64 {
	this.get_curr_time_count++
	if this.get_curr_time_count > UpdateCurrTimeCount {
		this.get_curr_time_count = 1
		this.curr_time = time.Now().UnixNano() / int64(time.Millisecond)
	}

	return this.curr_time
}
