package help

type DListNode struct {
	Pre    *DListNode  // 前一个
	Next   *DListNode  // 后一个
	SrcTid int32       // 源线程
	Data   interface{} // 事件对象
}

func (this *DListNode) Init(d interface{}) {
	this.Pre = this
	this.Next = this
	this.Data = d
}

func (this *DListNode) Pop() {
	if this.Pre != nil {
		this.Pre.Next = this.Next
	}
	if this.Next != nil {
		this.Next.Pre = this.Pre
	}

	this.Pre = nil
	this.Next = nil
}

func (this *DListNode) IsEmpty() bool {
	return (this.Pre == this.Next) && (this.Pre == nil || this.Pre == this)
}
