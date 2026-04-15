package model

import "github.com/example/aichat/backend/models"

const TableNameBetaApplication = "beta_application"

const (
	BetaApplicationStatusSubmitted = "submitted"
	BetaApplicationMailPending     = "pending"
	BetaApplicationMailSent        = "sent"
	BetaApplicationMailFailed      = "failed"
	BetaApplicationMailSkipped     = "skipped"
)

// BetaApplication 记录官网内测申请。
type BetaApplication struct {
	models.Model
	ProductInterest string `json:"productInterest" gorm:"column:product_interest;type:varchar(64);not null" db:"product_interest"`
	ContactType     string `json:"contactType" gorm:"column:contact_type;type:varchar(32);not null" db:"contact_type"`
	ContactValue    string `json:"contactValue" gorm:"column:contact_value;type:varchar(128);not null" db:"contact_value"`
	UseCase         string `json:"useCase" gorm:"column:use_case;type:text;not null" db:"use_case"`
	Note            string `json:"note" gorm:"column:note;type:text" db:"note"`
	SourcePage      string `json:"sourcePage" gorm:"column:source_page;type:varchar(255)" db:"source_page"`
	UserAgent       string `json:"userAgent" gorm:"column:user_agent;type:varchar(512)" db:"user_agent"`
	RemoteAddr      string `json:"remoteAddr" gorm:"column:remote_addr;type:varchar(128)" db:"remote_addr"`
	Status          string `json:"status" gorm:"column:status;type:varchar(32);not null" db:"status"`
	MailStatus      string `json:"mailStatus" gorm:"column:mail_status;type:varchar(32);not null" db:"mail_status"`
	MailError       string `json:"mailError" gorm:"column:mail_error;type:text" db:"mail_error"`
}

func (*BetaApplication) TableName() string {
	return TableNameBetaApplication
}
