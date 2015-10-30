package app

import (
	"bytes"
	"github.com/toophy/##[AppName]##/help"
	lua "github.com/toophy/gopher-lua"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// 消息函数类型
type MsgFunc func(*ClientConn)
type ConnRetFunc func(string, string, int, string) bool

type AppBase struct {
	heart_time          int64                      // 心跳时间(毫秒)
	start_time          int64                      // 线程开启时间戳
	last_time           int64                      // 最近一次线程运行时间戳
	curr_time           int64                      // 当前时间戳(毫秒)
	get_curr_time_count int64                      // 索取当前时间戳次数
	heart_rate          float64                    // 本次心跳比率
	pre_stop            bool                       // 预备停止
	stop_now            bool                       // 现在就停止
	first_run           bool                       // 线程首次运行
	Listener            map[string]*ListenConn     // 本地侦听端口
	RemoteSvr           map[string]*ClientConn     // 远程服务连接
	Conns               map[int]*ClientConn        // 连接池
	ConnLast            int                        // 最后连接Id
	MsgProc             []MsgFunc                  // 消息处理函数注册表
	MsgProcCount        int                        // 消息函数数量
	log_Buffer          []byte                     // 线程日志缓冲
	log_BufferLen       int                        // 线程日志缓冲长度
	log_TimeString      string                     // 时间格式(精确到秒2015.08.13 16:33:00)
	log_Header          [LogMaxLevel]string        // 各级别日志头
	log_FileBuff        bytes.Buffer               // 日志总缓冲, Tid_world才会使用
	log_FileHandle      *os.File                   // 日志文件, Tid_world才会使用
	evt_lay1            []help.DListNode           // 第一层事件池
	evt_lay2            map[uint64]*help.DListNode // 第二层事件池
	evt_names           map[string]help.IEvent     // 别名
	evt_lay1Size        uint64                     // 第一层池容量
	evt_lay1Cursor      uint64                     // 第一层游标
	evt_lastRunCount    uint64                     // 最近一次运行次数
	evt_currRunCount    uint64                     // 当前运行次数
	node_pool           []help.DListNode           // 节点池
	node_free           help.DListNode             // 自由节点
	node_alloc_count    int                        // 节点分配数量
	luaState            *lua.LState                // Lua实体
	luaNilTable         lua.LTable                 // Lua空的Table, 供默认参数使用

}

// 程序控制核心
var app *AppBase

func GetApp() *AppBase {
	if app == nil {
		app = &AppBase{}
		app.init()
	}
	return app
}

// App初始化
func (this *AppBase) init() {

	runtime.GOMAXPROCS(1)

	this.Listener = make(map[string]*ListenConn, 10)
	this.RemoteSvr = make(map[string]*ClientConn, 10)
	this.Conns = make(map[int]*ClientConn, 1000)
	this.MsgProc = make([]MsgFunc, 8000)

	this.ConnLast = 1

	// 日志初始化
	this.log_Buffer = make([]byte, LogBuffMax)
	this.log_BufferLen = 0

	this.log_TimeString = time.Now().Format("15:04:05")
	this.MakeLogHeader()

	this.log_FileBuff.Grow(LogBuffSize)

	// 检查log目录
	if !help.IsExist(LogDir) {
		os.MkdirAll(LogDir, os.ModeDir)
	}

	if !help.IsExist(LogFileName) {
		os.Create(LogFileName)
	}
	file, err := os.OpenFile(LogFileName, os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	this.log_FileHandle = file
	this.log_FileHandle.Seek(0, 2)

	// 第一条日志
	this.LogDebug("\n          %s服务器启动\n", AppName)
}

// 程序开启
func (this *AppBase) Start(heart_time int64) bool {

	// 创建pprof文件
	f, err := os.Create(LogDir + "/" + ProfFile)
	if err != nil {
		this.LogWarn(err.Error())
		return false
	}
	pprof.StartCPUProfile(f)

	this.heart_time = heart_time * int64(time.Millisecond)
	this.start_time = time.Now().UnixNano()
	this.last_time = this.start_time

	this.log_TimeString = time.Now().Format("15:04:05")
	this.MakeLogHeader()

	// 设置当前时间戳(毫秒)
	this.get_curr_time_count = 1
	this.curr_time = this.last_time / int64(time.Millisecond)

	this.heart_rate = 1.0
	this.first_run = true

	// 初始化事件池
	this.evt_lay1Size = Evt_lay1_time >> Evt_gap_bit
	this.evt_lay1Cursor = 0
	this.evt_currRunCount = 1
	this.evt_lastRunCount = this.evt_currRunCount

	this.evt_lay1 = make([]help.DListNode, this.evt_lay1Size)
	this.evt_lay2 = make(map[uint64]*help.DListNode, 0)
	this.evt_names = make(map[string]help.IEvent, 0)

	for i := uint64(0); i < this.evt_lay1Size; i++ {
		this.evt_lay1[i].Init(nil)
	}

	// 节点初始化
	this.node_free.Init(nil)

	this.node_pool = make([]help.DListNode, 200000)
	for i := 0; i < 200000; i++ {
		this.node_pool[i].Init(nil)
		this.addFreeDlinkNode(&this.node_pool[i])
	}

	// 载入Lua脚本
	errInit := this.reloadLuaState()
	if errInit != nil {
		this.LogError(errInit.Error())
		return false
	}

	// 刷新日志到硬盘
	go func() {
		for {
			if this.log_FileBuff.Len() > 0 {
				this.log_FileHandle.Write(this.log_FileBuff.Bytes())
				this.log_FileBuff.Reset()
			}

			<-time.Tick(5 * time.Second)
		}
	}()

	// 计算心跳误差值, 决定心跳滴答(小数), heart_time, last_time, heart_rate
	// 处理线程间接收消息, 分配到水表定时器
	// 执行水表定时器
	go func() {

		next_time := time.Duration(this.heart_time)
		run_time := int64(0)

		for {

			time.Sleep(next_time)

			this.log_TimeString = time.Now().Format("15:04:05")
			this.MakeLogHeader()

			this.last_time = time.Now().UnixNano()
			// 设置当前时间戳(毫秒)
			this.get_curr_time_count = 1
			this.curr_time = this.last_time / int64(time.Millisecond)

			this.runEvents()

			// 计算下一次运行的时间
			run_time = time.Now().UnixNano() - this.last_time
			if run_time >= this.heart_time {
				run_time = this.heart_time - 10*1000*1000
			} else if run_time < 0 {
				run_time = 0
			}

			next_time = time.Duration(this.heart_time - run_time)

			if this.pre_stop {
				break
			}
		}
	}()

	this.Tolua_Common("main", "OnAppBegin")

	return true
}

// 等待协程结束
func (this *AppBase) WaitExit() {

	for {
		<-time.Tick(2 * time.Second)
		if this.stop_now {
			pprof.StopCPUProfile()
			this.LogInfo("bye bye.")
			break
		}
	}

	// 响应线程退出
	this.Tolua_Common("main", "OnAppEnd")
	if this.luaState != nil {
		this.luaState.Close()
		this.luaState = nil
	}

	// 关闭日志文件
	if this.log_FileHandle != nil {
		this.log_FileHandle.Close()
	}
}
