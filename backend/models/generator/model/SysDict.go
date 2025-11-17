/* *
 * @Author: chengjiang
 * @Date: 2025-11-17 14:13:32
 * @Description: 字典，前端各种状态枚举的存储。
**/
package model

import "github.com/example/aichat/backend/models"

// "dictCode": "5",
//
//	"dictSort": 2,
//	"dictLabel": "隐藏",
//	"dictValue": "1",
//	"dictType": "sys_show_hide",
//	"cssClass": "",
//	"listClass": "danger",
//	"isDefault": false,
//	"status": models.StatusDisabled,
//	"remark": "隐藏菜单",
//	"createTime": 1684048781000
type SysDict struct {
	models.Model
	DictCode  string        `json:"dictCode" gorm:"dict_code" db:"dict_code"`
	DictSort  int           `json:"dictSort" gorm:"dict_sort" db:"dict_sort"`
	DictLabel string        `json:"dictLabel" gorm:"dict_label" db:"dict_label"`
	DictValue string        `json:"dictValue" gorm:"dict_value" db:"dict_value"`
	DictType  string        `json:"dictType" gorm:"dict_type" db:"dict_type"`
	CssClass  string        `json:"cssClass" gorm:"css_class" db:"css_class"`
	ListClass string        `json:"listClass" gorm:"list_class" db:"list_class"`
	IsDefault bool          `json:"isDefault" gorm:"is_default;type:bool;default:false" db:"is_default"`
	Status    models.Status `json:"status" gorm:"status;type:tinyint(1);default:2" db:"status"`
	Remark    string        `json:"remark" gorm:"remark;type:varchar(255)" db:"remark"`
}
