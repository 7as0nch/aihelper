/* *
 * @Author: chengjiang
 * @Date: 2025-12-10 17:01:44
 * @Description: AI 提示词模板
**/
package model

import "github.com/example/aichat/backend/models"

type AIPromptTemplate struct {
	models.Model
	Type        AI_Prompt_Type `json:"type" gorm:"type:smallint;default:1"` // system, user
	Name        string         `json:"name" gorm:"size:100"`
	Description string         `json:"description" gorm:"size:200"`
	Text        string         `json:"text" gorm:"type:text"`
}

func (AIPromptTemplate) TableName() string {
	return "ai_prompt_template"
}

type AI_Prompt_Type uint8

const (
	AI_Prompt_Type_System AI_Prompt_Type = iota + 1 // system
	AI_Prompt_Type_User                             // user
)
