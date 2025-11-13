package models

import "time"

/* *
 * @Author: chengjiang
 * @Date: 2025-10-04 17:14:27
 * @Description:
**/

type Model struct {
	ID         int64     `gorm:"column:id" json:"id"`
	CreateTime time.Time `gorm:"column:create_date_time;type:int unsigned;autoCreateTime" json:"create_time"`
	CreatedBy  int64     `gorm:"column:creator" json:"created_by"` // 创建人
	UpdateTime time.Time `gorm:"column:update_date_time;type:int unsigned;autoUpdateTime" json:"update_time"`
	UpdatedBy  int64     `gorm:"column:updater" json:"updated_by"` // 更新人
	IsDeleted  bool      `gorm:"column:is_soft_delete" json:"is_deleted"` // 是否删除0否1是
}
