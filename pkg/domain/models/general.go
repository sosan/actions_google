package models

import (
	"time"
)

const (
	LayoutTimestamp       = time.RFC3339 // format "2006-01-02T15:04:05Z07:00"
	MaxAttempts           = 11
	MinRangeSleepDuration = 100 * time.Millisecond // min range time wait offset
	MaxRangeSleepDuration = 500 * time.Millisecond // max range time wait offset
	SleepOffset           = 50 * time.Millisecond  // offset
	SaveOffset            = 10
	MaxRowsFromDB         = 999
	MaxTimeoutContext     = 60 * time.Second
	TimeDriftForExpire    = 600 // 10 minutes
	MaxTimeForLocks       = 30 * time.Second
	TimeoutRequest        = 5 * time.Minute
	MaxSecondsGoRoutine   = 9 * time.Second
)

const (
	CredentialCreateContextKey   = "createcredential"
	CredentialExchangeContextKey = "exchangecredential"
	ActionGoogleKey              = "actiongoogle"
	ActionNotionKey              = "actionnotion"
	NotionHost                   = "notion.so/"
)

var ValidCommandTypes = map[string]bool{
	"create": true,
	"update": true,
	"delete": true,
}

var ValidActionsTypes = map[string]bool{
	"googlesheets": true,
	"notionoauth":  true,
	"notiontoken":  true,
}

var PermitedPathList = map[string]bool{
	"/api/actions/ping": true,
}
