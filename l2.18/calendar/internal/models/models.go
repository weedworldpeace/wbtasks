package models

type Event struct {
	EventId string `json:"event_id"`
	Message string `json:"event"`
	Date    string `json:"date"`
}

type UserEvent struct {
	UserId string `json:"user_id"`
	Event
}

func NewUserEvent() *UserEvent {
	return &UserEvent{}
}

type BadResponse struct {
	Err string `json:"error"`
}

func NewBadResponse(err error) *BadResponse {
	return &BadResponse{err.Error()}
}

type GoodPostResponse struct { // tipo not bad...
	Result string `json:"result"`
}

func NewGoodPostResponse(message string) *GoodPostResponse {
	return &GoodPostResponse{message}
}

type GoodGetResponse struct {
	Result []Event `json:"result"`
}

func NewGoodGetResponse(result []Event) *GoodGetResponse {
	return &GoodGetResponse{Result: result}
}
