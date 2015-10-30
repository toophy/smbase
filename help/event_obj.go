package help

type EventObj struct {
	NodeObj DListNode
}

func (this *EventObj) InitEventHeader() {
	this.NodeObj.Init(nil)
}

func (this *EventObj) GetEventHeader() *DListNode {
	return &this.NodeObj
}

func (this *EventObj) AddEvent(e IEvent) bool {
	n := &DListNode{}
	n.Init(e)

	if !e.AddNode(n) {
		return false
	}

	old_pre := this.NodeObj.Pre

	this.NodeObj.Pre = n
	n.Next = &this.NodeObj
	n.Pre = old_pre
	old_pre.Next = n

	return true
}
