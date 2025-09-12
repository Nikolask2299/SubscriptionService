package app

import (
	"errors"
	"log/slog"
	"strconv"
	"subscriptionservice/interal/models"
)

type SubApp struct {
	loger    *slog.Logger
	base DataBase
}

func NewApp(loger *slog.Logger, base DataBase) *SubApp {
	return &SubApp{
		loger:    loger,
		base: base,
	}
}

type DataBase interface {
	CreateSubscr(subscr models.SubscrbUser) error
	DeleteSubscr(subs models.SubscrbUserSearch) error 
	UpdateSubscr(subscr map[string]string) error
	ReadSubscr(subscr models.SubscrbUserSearch) ([]models.SubscrbUser, error)
	GetSummSubscr(subscr models.SubscrbUserSearch) (models.SummSubscrb, error)
}

func (a *SubApp) CreateSubscr(newsubscr models.SubscrbUser) error {
	log := a.loger.With(
		slog.String("OP", "Create Subscription"),
	)

	log.Info("Create Subscription", slog.String("UUID", newsubscr.UUID))


	err := a.base.CreateSubscr(newsubscr) 
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint \user_subscr_pkey\` {
			return errors.New("error Create Subscription - UUID, NameService, StartDate is not unique")
		}
		log.Error("Error Create Subscription", slog.String("Error", err.Error()))
		return errors.New("error Create Subscription")
	}

	log.Info("Create Subscription Success", slog.String("UUID", newsubscr.UUID))
	return nil
}

func (a *SubApp) DeleteSubscr(subs models.SubscrbUserSearch) error {
	log := a.loger.With(
		slog.String("OP", "Delete Subscription"),
	)
	log.Info("Delete Subscription", slog.String("UUID", subs.UUID))
	
	if (subs.UUID == "") || (subs.NameService == "") || (subs.StartDate == "") {
		log.Error("Error Delete Subscription", slog.String("Error", "Error Delete Subscription - UUID, NameService, StartDate is empty"))
		return errors.New("error Delete Subscription - UUID, NameService, StartDate is empty")
	}

	if err := a.base.DeleteSubscr(subs); err != nil {
		log.Error("Error Delete Subscription", slog.String("Error", err.Error()))
		return errors.New("error Delete Subscription")
	}

	log.Info("Delete Subscription Success", slog.String("UUID", subs.UUID))
	return nil
}

func (a *SubApp) UpdateSubscr(subscr models.SubscrbUserSearch) error {
	log := a.loger.With(
		slog.String("OP", "Update Subsription"),
	)

	log.Info("Update Subscription", slog.String("UUID", subscr.UUID))

	log.Info("Update Subscription", slog.String("UUID", subscr.UUID))
	if (subscr.UUID == "") || (subscr.NameService == "") || (subscr.StartDate == "") {
		log.Error("Error Update Subscription", slog.String("Error", "Error Update Subscription - UUID, NameService, StartDate is empty"))
		return errors.New("error Update Subscription - UUID, NameService, StartDate is empty")
	}

	subscmap := map[string]string{}
	subscmap["name_service"] = subscr.NameService
	subscmap["start_date"] = subscr.StartDate
	subscmap["end_date"] = subscr.EndDate
	if subscr.Price != 0 {
		subscmap["price"] = strconv.Itoa(subscr.Price)
	}
	subscmap["user_id"] = subscr.UUID

	if err := a.base.UpdateSubscr(subscmap); err != nil {
		log.Error("Error Update Subscription", slog.String("Error", err.Error()))
		return errors.New("error Update Subscription")
	}

	log.Info("Update Subscription Success", slog.String("UUID", subscr.UUID))
	return nil
}

func (a *SubApp) ReadSubscr(subscr models.SubscrbUserSearch) ([]models.SubscrbUser, error) {
	log := a.loger.With(
		slog.String("OP", "Read Subscription"),
	)

	log.Info("Read Subscription", slog.String("UUID", subscr.UUID), slog.String("NameService", subscr.NameService), slog.String("StartDate", subscr.StartDate))
	
	res, err := a.base.ReadSubscr(subscr)
	if err != nil {
		log.Error("Error Read Subscription", slog.String("Error", err.Error()))
		return nil, errors.New("error Read Subscription")
	}

	log.Info("Read Subscription Success", slog.String("UUID", subscr.UUID), slog.String("NameService", subscr.NameService), slog.String("StartDate", subscr.StartDate))
	return res, nil
}

func (a *SubApp) GetSummSubscr(subscr models.SubscrbUserSearch) (models.SummSubscrb, error) {
	log := a.loger.With(
		slog.String("OP", "Update Subsription"),
	)

	log.Info("Summ Subscription")
	if (subscr.EndDate == "") && (subscr.StartDate == "") {
		return models.SummSubscrb{}, errors.New("StartDate and EndDate is empty")
	}

	res, err := a.base.GetSummSubscr(subscr)
	if  err != nil {
		log.Error("Error Read Subscription", slog.String("Error", err.Error()))
		return models.SummSubscrb{}, errors.New("error Read Subscription")
	}
	
	log.Info("Summ Subscription Success")
	return res, nil
}