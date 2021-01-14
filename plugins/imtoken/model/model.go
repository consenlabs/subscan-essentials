package model

type EventTransfer struct {
	ID         uint   `gorm:"primary_key" json:"-"`
	EventIndex string `json:"event_index"`
	BlockNum   uint64 `json:"block_num"`
	BlockHash  string `json:"block_hash"`
	Timestamp  uint64 `json:"timestamp"`
	Sender     string `json:"sender"`
	Receiver   string `json:"receiver"`
	Amount     string `json:"amount"`
	// Nonce 				int         `json:"nonce"`
	Fee           string `json:"fee"`
	ExtrinsicHash string `json:"extrinsic_hash"`
	ExtrinsicIdx  int    `json:"extrinsic_idx"`
	// CallModule 			string      `json:"call_module"`
	// CallModuleFunction 	string		`json:"call_module_function"`
	ModuleId string `json:"module_id"`
	EventId  string `json:"event_id"`
	EventIdx int    `json:"event_idx"`
}
