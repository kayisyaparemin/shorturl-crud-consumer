package queues

import "encoding/json"

type MailQueueModel struct {
    To              []string                `json:"to"`
    Cc              []string                `json:"cc"`
    Bcc             []string                `json:"bcc"`
    Subject         string                  `json:"subject"`
    Body            string                  `json:"body"`
    File            []FileModel             `json:"file"`
    ServiceProvider ServiceProviderEnumType `json:"serviceProvider"`
    ExtraParameters json.RawMessage         `json:"extraParameters"`
}

type FileModel struct {
    FileName  string `json:"fileName"`
    FileBase64 string `json:"fileBase64"`
}
type ServiceProviderEnumType int

const (
    Insider ServiceProviderEnumType = iota
    Emarsys
)
