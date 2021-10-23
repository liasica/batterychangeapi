// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package model

import (
    "battery/app/model/internal"
    "github.com/gogf/gf/os/gtime"
)

// UserNote is the golang structure for table user_note.
type UserNote internal.UserNote

// Fill with you ideas below.

// UserNotePostReq 提交跟进请求
type UserNotePostReq struct {
    UserId  uint64 `json:"userId" v:"required"`
    Content string `json:"content" v:"required|length:3,200"`
}

// UserNoteListItem 跟进列表详情
type UserNoteListItem struct {
    UserId      uint64      `json:"userId"`      // 用户ID
    SysUserId   uint        `json:"sysUserId"`   // 操作员ID
    Content     string      `json:"content"`     // 跟进内容
    CreatedAt   *gtime.Time `json:"createdAt"`   // 创建时间
    SysUserName string      `json:"sysUserName"` // 操作员
}
