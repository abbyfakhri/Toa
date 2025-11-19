package models

import "time"

type EmailBatch struct {
	Id           string    `json:"id" db:"id"`
	From         string    `json:"from" db:"from"`
	EmailCount   int       `json:"emailCount" db:"email_count"`
	SuccessCount int       `json:"successCount" db:"success_count"`
	FailCount    int       `json:"failCount" db:"fail_count"`
	StartAt      time.Time `json:"startAt" db:"start_at"`
	FinishAt     time.Time `json:"finishAt" db:"finish_at"`
}

type Email struct {
	Id      int       `json:"id" db:"id"`
	BatchId string    `json:"batchId" db:"batch_id"`
	Email   string    `json:"email" db:"email"`
	IsSent  bool      `json:"isSent" db:"is_sent"`
	SentAt  time.Time `json:"sentAt" db:"sent_at"`
	Log     *string   `json:"log" db:"log"`
}

type GetEmailRequest struct {
	BatchId string `json:"batchId" db:"batch_id"`
	Email   string `json:"email" db:"email"`
}

type PostEmailRequest struct {
	Destinations []string `json:"destinations" validate:"required"`
	Subject      string   `json:"subject" validate:"required"`
	Body         string   `json:"body"`
	Template     string   `json:"template"`
}

type PostEmailRequestCsv struct {
	Subject      string `form:"subject" validate:"required"`
	Body         string `form:"body"`
	Template     string `form:"template"`
	TargetColumn string `form:"targetColumn" validate:"required"`
}

type PostEmailResponse struct {
	BatchId string `json:"emailBatchId"`
}
