package types

type VoiceMessage struct {
	Data   []int32 `json:"data"`
	Sender string  `json:"sender"`
	RoomID string  `json:"roomID"`
}

type SubscriptionRequest struct {
	RoomID string `json:"roomID"`
	Sender string `json:"sender"`
}

type VoiceMessageCreateRequest struct {
	Data   string `json:"data"`
	Type   string `json:"type"`
	Sender string `json:"sender"`
	RoomID string `json:"roomID"`
}
