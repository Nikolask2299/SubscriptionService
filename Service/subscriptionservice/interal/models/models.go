package models

import "errors"

var Errisempty = errors.New("error Update Subscription - UUID, NameService, StartDate is empty")
var Errisnotex = errors.New("Is not exist in database")
var Errisuniqu = errors.New("error Create Subscription - UUID, NameService, StartDate is not unique")
var Errisdatae = errors.New("StartDate and EndDate is empty")
var ErrIvalidArgs = errors.New("Invalid arguments in you request")

// Subscriptions User model info
// @Description Subscriptions information about the account
type SubscrbUser struct {
	SubscrID    int64  `json:"subscrb_id" validate:"omitempty, optional"`
	NameService string `json:"service_name" validate:"required"`
	Price       int    `json:"price" validate:"required"`
	UUID        string `json:"user_id" validate:"required"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date" validate:"omitempty, optional"`
}

// Subscriptions User model filter
// @Description Subscriptions filter information about the account
type SubscrbUserSearch struct {
	NameService string `json:"service_name"`
	Price       int    `json:"price"`
	UUID        string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

// Subscriptions summary price
// @Description Summary price subscriptions
type SummSubscrb struct {
	UUID        string   `json:"user_id"`
	Summ        int      `json:"summ_subscriptions"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
	ServiceName []string `json:"service_name"`
}

// Subscriptions ID
// @Description Subscriptions ID
type ID struct {
	ID int64 `json:"id_subscr"`
}
