package models

import "github.com/example/aichat/backend/tools"

/* *
 * @Author: chengjiang
 * @Date: 2025-10-04 17:14:27
 * @Description:
**/

type Model struct {
	ID        int64 `gorm:"column:id;type:bigint(20)" json:"id"`
	CreatedAt Time  `gorm:"column:created_at;type:time;autoUpdateTime" json:"created_at"`
	CreatedBy int64 `gorm:"column:creator" json:"created_by"` // 创建人
	UpdatedAt Time  `gorm:"column:updated_at;type:time;autoUpdateTime" json:"updated_at"`
	UpdatedBy int64 `gorm:"column:updater" json:"updated_by"`        // 更新人
	IsDeleted bool  `gorm:"column:is_soft_delete" json:"is_deleted"` // 是否删除0否1是
}

func (m *Model) New() {
	m.ID = tools.GetSnowID()
}

type Status uint8

// iota
const (
	Status_Enabled  Status = iota + 1 // 启用
	Status_Disabled                   // 禁用
)
