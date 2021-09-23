package model

type Call struct {
	CallID string 			`json:"call_id"`
	Queue_ID string 		`json:"queue_id"`
	DialedNumber string 	`json:"dialed_phone"`
	CallingNumber string 	`json:"calling_number"`
	CallingLevel string 	`json:"calling_level"`
}
