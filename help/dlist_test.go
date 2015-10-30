package help

import (
	"fmt"
	"testing"
)

type EvtPool struct {
	header DListNode
}

func (this *EvtPool) Init() {
	this.header.Init(nil)
}

func (this *EvtPool) Eat(name string) {
	fmt.Printf("吃%s\n", name)
}

func (this *EvtPool) Post(d IEvent) bool {

	n := &DListNode{}
	n.Init(d)

	if !d.AddNode(n) {
		return false
	}

	old_pre := this.header.Pre

	this.header.Pre = n
	n.Next = &this.header
	n.Pre = old_pre
	old_pre.Next = n

	return true
}

func (this *EvtPool) Run() {
	for {
		if this.header.IsEmpty() {
			break
		}

		n := this.header.Next

		n.Data.(IEvent).Exec(this)

		n.Data.(IEvent).Destroy()
	}
}

type Evt_eat struct {
	Evt_base
	FoodName string
}

func (this *Evt_eat) Exec() bool {
	return true
}

func TestDlist(t *testing.T) {

	var g_Pool EvtPool
	g_Pool.Init()

	g_Pool.Post(&Evt_eat{FoodName: "西瓜"})
	g_Pool.Post(&Evt_eat{FoodName: "葡萄"})
	g_Pool.Post(&Evt_eat{FoodName: "黄瓜"})
	g_Pool.Post(&Evt_eat{FoodName: "大蒜"})

	g_Pool.Run()
}
