// Package cmd_server 提供一些小工具
package cmd_server

// 没添加一种回调，就添加一个CMD，一个回调，一个数据类型 如果假如回调的话需要增加回调对应的数据结构

type CMD uint64 // 指令头
type AnyData interface{}
type FUNC func(data AnyData) AnyData // 回调函数，当没有返回值时则返回nil

// 参数是协议与提供返回值的channel(避免过于复杂的回调函数设计), 也可以不提供，或者提供回调函数以供调用。或者都不提供让这边通过某种方式发一个返回
const (
// END  CMD = 2<<64 - 1
// STOP CMD = 2<<64 - 2
)

type Plt struct {
	Cmd  CMD          // 命令
	Data AnyData      // 数据
	Ret  chan AnyData // 返回channel
}

type CMDServer struct {
	funcMap     map[CMD]FUNC
	inputC      chan Plt
	adminInputC chan Plt // 优先
	stop        chan bool
}

// Update 为主线程
func (s *CMDServer) Update() {
	for {
		select {
		case <-s.stop:
			return
		default:
			select {
			case plt := <-s.adminInputC:
				s.proPlt(plt)
			default:
				// 优先处理优先命令
				select {
				case plt := <-s.inputC:
					s.proPlt(plt)
				}
			}
		}

	}
}

func (s *CMDServer) proPlt(plt Plt) {
	cmd := plt.Cmd
	d := s.funcMap[cmd](plt.Data)
	if d != nil {
		go ret(plt, d)
	}
}

func ret(plt Plt, d AnyData) {
	plt.Ret <- d
}

func MakeCmdServer(funcMap map[CMD]FUNC, adminCNUm int, cNum int) CMDServer {
	return CMDServer{
		funcMap:     funcMap,
		inputC:      make(chan Plt, cNum),
		adminInputC: make(chan Plt, adminCNUm),
		stop:        make(chan bool),
	}
}

func (s *CMDServer) Start() {
	go s.Update()
}

func (s *CMDServer) Stop() {
	s.stop <- true
}

//Send 处理发来的请求
//视情况可能阻塞，低负荷情况下非阻塞
func (s *CMDServer) Send(plt Plt) {
	s.inputC <- plt
}

//AdminSend 将优先处理请求
//视情况可能阻塞，低负荷情况下非阻塞
func (s *CMDServer) AdminSend(plt Plt) {
	s.adminInputC <- plt
}
