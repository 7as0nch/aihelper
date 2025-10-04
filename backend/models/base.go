package models

import "time"

/* *
 * @Author: chengjiang
 * @Date: 2025-10-04 17:14:27
 * @Description:
**/

type Model struct {
	ID         int64     `json:"id"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	DeleteFlag bool      `json:"delete_flag"`
}
