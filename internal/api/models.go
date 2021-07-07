package api

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Debt struct {
	Lender *User `json:"lender" bson:"lender"`
	Debtor *User `json:"debtor" bson:"debtor"`
	Sum    int   `json:"sum" bson:"sum"`
}

// ChatState stores user state
type ChatState struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId       int                `json:"userId" bson:"user_id"`
	Action       Action             `json:"action" bson:"action"`
	CallbackData *CallbackData      `json:"callbackData" bson:"callback_data"`
}

// Button which is sent to the user as ReplyMarkup
type Button struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CallbackData *CallbackData      `json:"callbackData" bson:"callback_data"`
	Text         string             `json:"text" bson:"text"`
	Action       Action             `json:"action" bson:"action"`
	CreateAt     time.Time          `json:"createAt" bson:"create_at"`
}

type Action string

type CallbackData struct {
	RoomId       string             `json:"roomId" bson:"room_id,omitempty"`
	UserId       int                `json:"userId" bson:"user_id,omitempty"`
	ExternalId   string             `json:"externalId" bson:"external_id,omitempty"`
	ExternalData string             `json:"externalData" bson:"external_data,omitempty"`
	OperationId  primitive.ObjectID `json:"operationId" bson:"operation_id,omitempty"`
	Page         int                `json:"page" bson:"page,omitempty"`
}

func NewButton(action Action, data *CallbackData) *Button {
	return &Button{
		ID:           primitive.NewObjectID(),
		Action:       action,
		CallbackData: data,
		CreateAt:     time.Now(),
	}
}
