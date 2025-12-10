/* *
 * @Author: chengjiang
 * @Date: 2025-12-10 17:10:45
 * @Description:
**/
package model

import "github.com/example/aichat/backend/models"

type AIModel struct {
	models.Model
	Category    AIModel_Category  `json:"category" gorm:"type:smallint;default:0"`     // 模型类别
	ModelType   ModelType         `json:"model_type" gorm:"type:varchar(50);not null"` // ark, openai, deepseek
	ModelName   string            `json:"model_name" gorm:"type:varchar(100);not null"`
	APIKey      string            `json:"api_key" gorm:"type:varchar(255)"` // 建议加密存储
	BaseURL     string            `json:"base_url" gorm:"type:varchar(255)"`
	MaxTokens   int               `json:"max_tokens" gorm:"type:int;default:8192"`
	Temperature float32           `json:"temperature" gorm:"type:float;default:0.7"`
	TopP        float32           `json:"top_p" gorm:"type:float;default:0.9"`
	PriceType   AIModel_PriceType `json:"price_type" gorm:"type:smallint;default:0"` // 计费方式
	Price       float32           `json:"price" gorm:"type:float;default:0.0"`       // 单价（元/1000 tokens）
	Supplier    string            `json:"supplier" gorm:"type:varchar(50)"`          // 供应商名称
	Description string            `json:"description" gorm:"type:varchar(255)"`      // 模型描述
	Status      models.Status     `json:"status" gorm:"type:smallint;default:1"`     // 状态: 是否可用
	IsDefault   models.Status     `json:"is_default" gorm:"type:smallint;default:2"` // 是否默认模型：全局有且只有一个默认模型。
}
// TODO 与用户的映射：用户可以选择使用哪些模型。

// TableName 指定表名
func (AIModel) TableName() string {
	return "ai_model"
}

type ModelType string

const (
	ARK_MODEL    ModelType = "ark"
	OPENAI_MODEL ModelType = "openai"
	DEEPSEEK     ModelType = "deepseek"
)

type AIModel_PriceType uint8

const (
	_                         AIModel_PriceType = iota // Free
	AIModel_PriceType_ByToken                          // 按token计费
	AIModel_PriceType_ByCall                           // 按次数计费
)

type AIModel_Category uint8

const (
	_ AIModel_Category = iota // 无类别
	// NLP 自然语言处理
	AIModel_Category_NLP // NLP-通用模型 「文本生成、理解、翻译等」
	// Multimodal 多模态
	AIModel_Category_Multimodal_General // 多模态-通用模型 「跨模态理解与生成」
	AIModel_Category_Multimodal_Image   // 多模态-图像模型 「图像生成、图像理解」
	AIModel_Category_Multimodal_Video   // 多模态-视频模型 「视频生成、视频理解」
	AIModel_Category_Multimodal_Audio   // 多模态-音频模型 「音频生成、音频理解、语音识别」
)
