package znet

import (
	"fmt"
	"github.com/lgc202/Zinx/ziface"
	"strconv"
)

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter // 存放每个MsgId 所对应的处理方法的map属性
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), "is not Found")
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := m.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	m.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}
