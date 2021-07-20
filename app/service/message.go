package service

import (
	"battery/app/dao"
	"battery/app/model"
	"battery/library/push/getui"
	"battery/library/wf"
	"context"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

var MessageService = messageService{}

type messageService struct {
	wf *wf.WorkFlow
}

// Create 创建消息
func (s *messageService) Create(ctx context.Context, userId uint64, t uint, title, summary string, detail model.MessageDetail) (uint64, error) {
	id, err := dao.Message.Ctx(ctx).
		Fields(
			dao.Message.Columns.Type,
			dao.Message.Columns.UserId,
			dao.Message.Columns.Title,
			dao.Message.Columns.Summary,
			dao.Message.Columns.Detail,
		).
		InsertAndGetId(g.Map{
			dao.Message.Columns.Type:    t,
			dao.Message.Columns.UserId:  userId,
			dao.Message.Columns.Title:   title,
			dao.Message.Columns.Summary: summary,
			dao.Message.Columns.Detail:  detail,
		})
	if err == nil && t != model.MessageTypeSystem {
		s.SendWorkFlow(uint64(id))
	}
	return uint64(id), err
}

// Detail 消息详情
func (*messageService) Detail(ctx context.Context, messageId uint64) (model.Message, error) {
	var msg model.Message
	err := dao.Message.Ctx(ctx).WherePri(messageId).Scan(&msg)
	return msg, err
}

// ListUser 用户消息列表
func (s *messageService) ListUser(ctx context.Context, userId uint64, page model.Page) []*model.Message {
	var list []*model.Message
	_ = dao.Message.Ctx(ctx).
		WhereIn(dao.Message.Columns.UserId, []interface{}{userId, 0}).
		WhereBetween(dao.Message.Columns.Type, 100, 499).
		OrderDesc(dao.Message.Columns.Id).
		Page(page.PageIndex, page.PageLimit).
		Scan(&list)
	s.ListIsRead(list, userId, 1)
	return list
}

// ListShopManager 店长消息列表
func (s *messageService) ListShopManager(ctx context.Context, userId uint64, page model.Page) []*model.Message {
	var list []*model.Message
	_ = dao.Message.Ctx(ctx).
		WhereIn(dao.Message.Columns.UserId, []interface{}{userId, 0}).
		WhereBetween(dao.Message.Columns.Type, 500, 999).
		OrderDesc(dao.Message.Columns.Id).
		Page(page.PageIndex, page.PageLimit).
		Scan(&list)
	s.ListIsRead(list, userId, 2)
	return list
}

func (*messageService) ListIsRead(list []*model.Message, userId uint64, userType uint) {
	ids := make([]uint64, len(list))
	readIds := make(map[uint64]bool, len(list))
	for key, msg := range list {
		ids[key] = msg.Id
		readIds[msg.Id] = false
	}
	var readList []model.MessageRead
	_ = dao.MessageRead.WhereIn(dao.MessageRead.Columns.MessageId, ids).
		Where(dao.MessageRead.Columns.UserId, userId).
		Where(dao.MessageRead.Columns.UserType, userType).
		Scan(&readList)
	for _, read := range readList {
		readIds[read.MessageId] = true
	}
	for _, msg := range list {
		msg.IsRead = readIds[msg.Id]
	}
}

// Read 标记消息为已读
func (*messageService) Read(ctx context.Context, userId uint64, userType uint, messageIds []uint64) error {
	data := make(g.List, len(messageIds))
	for key, id := range messageIds {
		data[key] = g.Map{
			dao.MessageRead.Columns.MessageId: id,
			dao.MessageRead.Columns.UserType:  userType,
			dao.MessageRead.Columns.UserId:    userId,
		}
	}
	_, err := dao.MessageRead.Ctx(ctx).Save(data)
	return err
}

func (s *messageService) SendWorkFlowInit() {
	s.wf = wf.Start("Push-Message", 4, 20)
}

type MessagePayload struct {
	MessageId uint64
}

func (p *MessagePayload) Play() {
	message, err := MessageService.Detail(context.TODO(), p.MessageId)
	if err != nil {
		return
	}
	if message.Type == model.MessageTypeSystem {
		return
	}
	user := UserService.Detail(context.TODO(), message.UserId)
	res, _err := getui.Service().PushSingleCid(getui.PushSingleRequest{
		RequestId: fmt.Sprintf("%d-%d", gtime.Timestamp(), message.Id),
		Audience:  getui.PushSingleRequestAudience{Cid: []string{user.DeviceToken}},
		PushMessage: getui.PushSingleRequestPushMessage{Notification: getui.PushMessageNotification{
			Title:     message.Title,
			Body:      message.Summary,
			ClickType: "none",
			//TODO 跳转
		}},
	})
	if _err != nil {
		g.Log().Errorf("个推推送错误：%v, %v", _err.Error(), res)
	}
}

func (s *messageService) SendWorkFlow(messageId uint64) {
	s.wf.AddJob(wf.Job{
		Payload: &MessagePayload{
			MessageId: messageId,
		},
	})
}
