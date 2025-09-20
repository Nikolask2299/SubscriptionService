package app

import (
	"errors"
	"log/slog"
	"subscriptionservice/interal/models"

	"github.com/lib/pq"
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
	CreateSubscr(subscr models.SubscrbUser) (int64, error)
	DeleteSubscr(id int64, subs models.SubscrbUserSearch) error 
	UpdateSubscr(id int64, subscr models.SubscrbUserSearch) error
	ReadSubscr(id int64, subscr models.SubscrbUserSearch) (models.SubscrbUser, error)
	GetSummSubscr(subscr models.SubscrbUserSearch) (models.SummSubscrb, error)
	ReadAllSubscr(off, lim int64, subscr models.SubscrbUserSearch) ([]models.SubscrbUser, error)
}

func (a *SubApp) CreateSubscr(newsubscr models.SubscrbUser) (int64, error) {
	log := a.loger.With(
		slog.String("OP", "Create Subscription"),
	)

	log.Info("Create Subscription", slog.String("UUID", newsubscr.UUID))


	id, err := a.base.CreateSubscr(newsubscr) 
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Error("Error Create Subscription: " + pqErr.Error())
			if pqErr.Code.Name() == "unique_violation" {
				return -1, models.Errisuniqu
			} else if pqErr.Code.Name() == "invalid_text_representation" {
				log.Error("Error Read Subscription", slog.String("Error", err.Error()))
				return -1, models.ErrIvalidArgs
			}

			return -1, errors.New("error Create Subscription")
		}
			
		log.Error("Error Create Subscription: " + err.Error())
		return -1, errors.New("error Create Subscription")
	}

	log.Info("Create Subscription Success", slog.String("UUID", newsubscr.UUID))
	return id, nil
}

func (a *SubApp) DeleteSubscr(id int64, subs models.SubscrbUserSearch) error {
	log := a.loger.With(
		slog.String("OP", "Delete Subscription"),
	)
	
	log.Info("Delete Subscription", slog.String("UUID", subs.UUID))
	
	if id <= 0 {
		if (subs.UUID == "") || (subs.NameService == "") || (subs.StartDate == "") {
			log.Error("Error Delete Subscription: " + models.Errisempty.Error())
			return models.Errisempty
		}
	}

	if err := a.base.DeleteSubscr(id, subs); err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			log.Error("Error Delete Subscription: " + pqErr.Code.Name())
			return errors.New("error Delete Subscription")
		}

		if err.Error() == "Is not exist in database" {
			return models.Errisnotex
		}

		log.Error("Error Delete Subscription: ", slog.String("Error", err.Error()))
		return errors.New("error Delete Subscription")
	}

	log.Info("Delete Subscription Success", slog.String("UUID", subs.UUID))
	return nil
}

func (a *SubApp) UpdateSubscr(id int64, subscr models.SubscrbUserSearch) error {
	log := a.loger.With(
		slog.String("OP", "Update Subsription"),
	)

	log.Info("Update Subscription", slog.String("UUID", subscr.UUID))

	log.Info("Update Subscription", slog.String("UUID", subscr.UUID))

	if (subscr.UUID == "") || (subscr.NameService == "") || (subscr.StartDate == "") {
		if id <= 0 {
			log.Error("Error Update Subscription: " + models.Errisempty.Error())
			return models.Errisempty
		}
	}

	if err := a.base.UpdateSubscr(id, subscr); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Error("Error Update Subscription: " + pqErr.Error())
			if pqErr.Code.Name() == "no_data" {
				return models.Errisnotex
			}

			return errors.New("error Update Subscription")
		}

		if err.Error() == "Is not exist in database" {
			return models.Errisnotex
		}

		if err.Error() == "update statements must have at least one Set clause" {
			return models.ErrIvalidArgs
		}

		log.Error("Error Update Subscription", slog.String("Error", err.Error()))
		return errors.New("error Update Subscription")
	}

	log.Info("Update Subscription Success", slog.String("UUID", subscr.UUID))
	return nil
}

func (a *SubApp) ReadSubscr(id int64, subscr models.SubscrbUserSearch) (models.SubscrbUser, error) {
	log := a.loger.With(
		slog.String("OP", "Read Subscription"),
	)

	log.Info("Read Subscription", slog.Any("ID", id), slog.String("UUID", subscr.UUID), slog.String("NameService", subscr.NameService), slog.String("StartDate", subscr.StartDate))
	
	if id <= 0 {
		if (subscr.UUID == "") || (subscr.NameService == "") || (subscr.StartDate == "") {
			log.Error("Error Read Subscription: " + models.Errisempty.Error())
			return models.SubscrbUser{}, models.Errisempty
		}
	}

	res, err := a.base.ReadSubscr(id, subscr)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "invalid_text_representation" {
				log.Error("Error Read Subscription", slog.String("Error", err.Error()))
				return models.SubscrbUser{}, models.ErrIvalidArgs
			}
		}

		if err.Error() == "Is not exist in database" {
			return models.SubscrbUser{}, models.Errisnotex
		}

		log.Error("Error Read Subscription", slog.String("Error", err.Error()))
		return models.SubscrbUser{}, errors.New("error Read Subscription")
	}

	log.Info("Read Subscription Success", slog.Any("ID", id), slog.String("UUID", subscr.UUID), slog.String("NameService", subscr.NameService), slog.String("StartDate", subscr.StartDate))
	return res, nil
}

func (a *SubApp) GetSummSubscr(subscr models.SubscrbUserSearch) (models.SummSubscrb, error) {
	log := a.loger.With(
		slog.String("OP", "Summ Subsription"),
	)

	
	if (subscr.EndDate == "") && (subscr.StartDate == "") {
		return models.SummSubscrb{}, models.Errisdatae
	}
	
	log.Info("Summ Subscription", slog.String("Start Date", subscr.StartDate), slog.String("End Date", subscr.EndDate))
	
	res, err := a.base.GetSummSubscr(subscr)
	if  err != nil {
		log.Error("Error Summ Subscription", slog.String("Error", err.Error()))
		return models.SummSubscrb{}, errors.New("error Summ Subscription")
	}
	
	log.Info("Summ Subscription Success")
	return res, nil
}

func (a *SubApp) ReadAllSubscr(off, lim int64, subscr models.SubscrbUserSearch) ([]models.SubscrbUser, error) {
		log := a.loger.With(
		slog.String("OP", "Read Subscription"),
	)

	log.Info("Read Subscription", slog.String("UUID", subscr.UUID), slog.String("NameService", subscr.NameService), slog.String("StartDate", subscr.StartDate))
	
	res, err := a.base.ReadAllSubscr(off, lim, subscr)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "invalid_text_representation" {
				log.Error("Error Read All Subscription", slog.String("Error", err.Error()))
				return nil, models.ErrIvalidArgs
			}

			log.Error("Error Read All Subscription", slog.String("Error", err.Error()))
			return nil, errors.New("error Read All Subscription")
		}


		log.Error("Error Read All Subscription", slog.String("Error", err.Error()))
		return nil, errors.New("error Read All Subscription")
	}

	log.Info("Read All Subscription Success", slog.String("UUID", subscr.UUID), slog.String("NameService", subscr.NameService), slog.String("StartDate", subscr.StartDate))
	return res, nil
}