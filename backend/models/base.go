package models

import (
	"github.com/example/aichat/backend/pkg/auth"
	"github.com/example/aichat/backend/tools"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

/* *
 * @Author: chengjiang
 * @Date: 2025-10-04 17:14:27
 * @Description:
**/

type Model struct {
	ID        int64                 `gorm:"column:id;type:bigint;primaryKey" json:"id"`
	CreatedAt Time                  `gorm:"column:created_at;type:time;autoCreateTime" json:"created_at"`
	CreatedBy int64                 `gorm:"column:created_by;type:bigint" json:"created_by"` // 创建人
	UpdatedAt Time                  `gorm:"column:updated_at;type:time;autoUpdateTime" json:"updated_at"`
	UpdatedBy int64                 `gorm:"column:updated_by;type:bigint" json:"updated_by"`                             // 更新人
	IsDeleted soft_delete.DeletedAt `gorm:"column:is_deleted;type:smallint;default:0;softDelete:flag" json:"is_deleted"` // 是否删除0否1是
	DeletedAt Time                  `gorm:"column:deleted_at;type:time" json:"deleted_at"`
	DeletedBy int64                 `gorm:"column:deleted_by;type:bigint" json:"deleted_by"` // 删除人
}

func (m *Model) BeforeDelete(tx *gorm.DB) (err error) {
    ctx := tx.Statement.Context
    var userID int64 = 0
    
    if ctx != nil {
        if authUserID := ctx.Value(auth.UserId); authUserID != nil {
            if uid, ok := authUserID.(int64); ok {
                userID = uid
            } else if uid, ok := authUserID.(int); ok {
                userID = int64(uid)
            }
        }
    }
    
    m.DeletedAt = Now()
    m.DeletedBy = userID
    tx.Statement.SetColumn("deleted_at", m.DeletedAt)
    tx.Statement.SetColumn("deleted_by", m.DeletedBy)
    return
}

func (m *Model) New() {
	m.ID = tools.GetSFID()
}

type Status uint8

// iota
const (
	Status_Enabled  Status = iota + 1 // 启用
	Status_Disabled                   // 禁用
)
