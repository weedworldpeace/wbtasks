package models

import (
	"errors"
	"time"
)

var (
	ErrBadPort       = errors.New("bad port")
	ErrBadReleaseMod = errors.New("bad release mod")
	ErrNonExistId    = errors.New("non exist id")
	ErrBadEmail      = errors.New("bad email")
)

type Notification struct {
	Email        string    `json:"email" binding:"required"`
	Id           string    `json:"id"`
	CreationDate time.Time `json:"creation_date" time_format:"2006-01-02 15:04:05"`
	SendingDate  time.Time `json:"sending_date" time_format:"2006-01-02 15:04:05" binding:"required"`
	Data         string    `json:"data" binding:"required"`
}

func NewNotification() *Notification {
	return &Notification{}
}

// type BadResponse struct {
// 	Err string `json:"error"`
// }

// func NewBadResponse() *BadResponse {
// 	return &BadResponse{}
// }

// type CreateNotificationRequest struct {
// 	Notification
// }

// func NewCreateNotificationRequest() *CreateNotificationRequest {
// 	return &CreateNotificationRequest{}
// }

// type CreateNotificationResponse struct {
// 	Id string `json:"id"`
// }

// func NewCreateNotificationResponse() *CreateNotificationResponse {
// 	return &CreateNotificationResponse{}
// }

// type ReadNotificationRequest struct {
// 	Id string `form:"id"`
// }

// func NewReadNotificationRequest() *ReadNotificationRequest {
// 	return &ReadNotificationRequest{}
// }

// type ReadNotificationResponse struct {
// 	Notification
// }

// func NewReadNotificationResponse() *ReadNotificationResponse {
// 	return &ReadNotificationResponse{}
// }

// type DeleteNotificationRequest struct {
// 	Id string `form:"id"`
// }

// func NewDeleteNotificationRequest() *DeleteNotificationRequest {
// 	return &DeleteNotificationRequest{}
// }

// type DeleteNotificationResponse struct {
// 	Msg string `json:"message"`
// }

// func NewDeleteNotificationResponse() *DeleteNotificationResponse {
// 	return &DeleteNotificationResponse{}
// }
