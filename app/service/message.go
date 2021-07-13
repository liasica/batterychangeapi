package service

import (
	"battery/app/dao"
	"battery/app/model"
	"context"
	"github.com/gogf/gf/frame/g"
)

var MessageService = messageService{}

type messageService struct {
}

// Create 创建消息
func (*messageService) Create(ctx context.Context, userId uint64, t uint, title, summary string, detail interface{}) (uint64, error) {
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
	return uint64(id), err
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
