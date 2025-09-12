package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"subscriptionservice/interal/models"
)

// Server represents the server interface
// @Description Server interface for interacting with the server
type SubscrServer struct {
	loger *slog.Logger
	app AppSubscr
}

func NewSubscrServer(loger *slog.Logger, app AppSubscr) *SubscrServer {
	return &SubscrServer{
		loger: loger,
		app: app,
	}
}

type AppSubscr interface {
	CreateSubscr(newsubscr models.SubscrbUser) error
	DeleteSubscr(subs models.SubscrbUserSearch) error
	UpdateSubscr(subscr models.SubscrbUserSearch) error
	ReadSubscr(subscr models.SubscrbUserSearch) ([]models.SubscrbUser, error)
	GetSummSubscr(subscr models.SubscrbUserSearch) (models.SummSubscrb, error)
}

// CreateSubscr godoc
// @Summary      Create subscription
// @Description  create subscription from database
// @Tags         create
// @Accept       json
// @Produce      json
// @Param        input body models.SubscrbUser true "subscription struct"
// @Success      204 "success response"
// @Failure      400  "Bad request error"
// @Failure      404 "Not found error"
// @Failure      405 "Method not allowed"
// @Failure      500  "Internal server error"
// @Router       /create [post]
func (s *SubscrServer) CreateSubscr(w http.ResponseWriter, r *http.Request) {
	s.loger.Info("Create Subscription" + r.URL.String())

   if r.Method != "POST" {
        s.loger.Error("Error creating subscrible from server" + " method not allowed")
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	var subscr models.SubscrbUser
	if err := json.NewDecoder(r.Body).Decode(&subscr); err != nil {
		s.loger.Error("Error creating subscrible from server" + err.Error())
		http.Error(w, "Error creating subscrible from server", http.StatusInternalServerError)
		return
	}

	if err := s.app.CreateSubscr(subscr); err != nil {
		if err.Error() == "error Create Subscription - UUID, NameService, StartDate is empty" {
			s.loger.Error("Error creating subscrible from server" + err.Error())
			http.Error(w, "UUID, NameService, StartDate is empty", http.StatusBadRequest)
			return
		} else if err.Error() == "error Create Subscription - UUID, NameService, StartDate is not unique" {
			s.loger.Error("Error creating subscrible from server" + err.Error())
			http.Error(w, "Error UUID, NameService, StartDate is not unique", http.StatusInternalServerError)
			return
		} else {
			s.loger.Error("Error creating subscrible from server" + err.Error())
			http.Error(w, "Error creating subscrible from server", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
	s.loger.Info("Create Subscription Success" + r.URL.String())
}

// UpdateSubscr godoc
// @Summary      Update subscription
// @Description  update subscription from database
// @Tags         update
// @Accept       json
// @Produce      json
// @Param        input body models.SubscrbUserSearch true "update subscription struct"
// @Success      204 "success response"
// @Failure      400  "Bad request error"
// @Failure      404 "Not found error"
// @Failure      405 "Method not allowed"
// @Failure      500  "Internal server error"
// @Router       /update [post]
func (s *SubscrServer) UpdateSubscr(w http.ResponseWriter, r *http.Request) {
	s.loger.Info("Update Subscription" + r.URL.String())
	
	if r.Method != "POST" {
		s.loger.Error("Error updating subscrible from server" + " method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var subscr models.SubscrbUserSearch
	if err := json.NewDecoder(r.Body).Decode(&subscr); err != nil {
		s.loger.Error("Error update subscrible from server" + err.Error())
		http.Error(w, "Error update subscrible from server ", http.StatusInternalServerError)
		return
	}

	if err := s.app.UpdateSubscr(subscr); err != nil {
		if err.Error() == "error Update Subscription - UUID, NameService, StartDate is empty" {
			s.loger.Error("Error update subscrible from server" + err.Error())
			http.Error(w, "Error update subscrible from server", http.StatusBadRequest)
			return
		} else {
			s.loger.Error("Error update subscrible from server" + err.Error())
			http.Error(w, "Error update subscrible from server", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
	s.loger.Info("Update Subscription Success" + r.URL.String())
}

// DeleteSubscr godoc
// @Summary      Delete subscription
// @Description  delete subscription from database
// @Tags         deleted
// @Accept       json
// @Produce      json
// @Param        input body models.SubscrbUserSearch true "delete subscription struct"
// @Success      204 "success response"
// @Failure      400  "Bad request error"
// @Failure      404 "Not found error"
// @Failure      405 "Method not allowed"
// @Failure      500  "Internal server error"
// @Router       /delete [delete]
func (s *SubscrServer) DeleteSubscr(w http.ResponseWriter, r *http.Request) {
	s.loger.Info("Delete Subscription" + r.URL.String())

	if r.Method != "DELETE" {
		s.loger.Error("Error delete subscrible from server" + " method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var subscr models.SubscrbUserSearch
	if err := json.NewDecoder(r.Body).Decode(&subscr); err != nil {
		s.loger.Error("Error delete subscrible from server" + err.Error())
		http.Error(w, "Error delete subscrible from server ", http.StatusInternalServerError)
		return
	}

	if err := s.app.DeleteSubscr(subscr); err != nil {
		if err.Error() == "error Delete Subscription - UUID, NameService, StartDate is empty" {
			s.loger.Error("Error delete subscrible from server" + err.Error())
			http.Error(w, "Error delete subscrible from server", http.StatusBadRequest)
			return
		} else {
			s.loger.Error("Error delete subscrible from server" + err.Error())
			http.Error(w, "Error delete subscrible from server", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
	s.loger.Info("Delete Subscription Success" + r.URL.String())
}

// GetSubscr godoc
// @Summary      Get subscription
// @Description  get subscription from database
// @Tags         data
// @Accept       json
// @Produce      json
// @Param        input body models.SubscrbUserSearch true "filter information"
// @Success      200  {array} models.SubscrbUser
// @Failure      400  "Bad request error"
// @Failure      404 "Not found error"
// @Failure      405 "Method not allowed"
// @Failure      500  "Internal server error"
// @Router       /search [post]
func (s *SubscrServer) GetSubscr(w http.ResponseWriter, r *http.Request) {
	s.loger.Info("Get Subscription" + r.URL.String())

	if r.Method != "POST" {
		s.loger.Error("Error get subscrible from server" + " method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var subscr models.SubscrbUserSearch
	if err := json.NewDecoder(r.Body).Decode(&subscr); err != nil {
		s.loger.Error("Error get subscrible from server" + err.Error())
		http.Error(w, "Error get subscrible from server ", http.StatusInternalServerError)
		return
	}

	subres, err := s.app.ReadSubscr(subscr)
	if err != nil {
		s.loger.Error("Error get subscrible from server" + err.Error())
		http.Error(w, "Error get subscrible from server ", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(subres); err != nil {
		s.loger.Error("Error get subscrible from server" + err.Error())
		http.Error(w, "Error get subscrible from server ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	s.loger.Info("Get Subscription Success" + r.URL.String())
}


// Get Summ Subscrible godoc
// @Summary      Summ subscription
// @Description  Summ subscription from database
// @Tags         data
// @Accept       json
// @Produce      json
// @Param        input body models.SubscrbUserSearch true "filter information"
// @Success      200  {object} models.SummSubscrb
// @Failure      400  "Bad request error"
// @Failure      404 "Not found error"
// @Failure      405 "Method not allowed"
// @Failure      500  "Internal server error"
// @Router       /summsubscr [post]
func (s *SubscrServer) GetSummSubscr(w http.ResponseWriter, r *http.Request) {
	s.loger.Info("Get Summ Subscription" + r.URL.String())

	if r.Method != "POST" {
		s.loger.Error("Error get subscrible from server" + " method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var subscr models.SubscrbUserSearch
	if err := json.NewDecoder(r.Body).Decode(&subscr); err != nil {
		s.loger.Error("Error get subscrible from server" + err.Error())
		http.Error(w, "Error get subscrible from server ", http.StatusInternalServerError)
		return
	}

	ressum, err := s.app.GetSummSubscr(subscr)
	if err != nil {
		s.loger.Error("Error get subscrible from server" + err.Error())
		http.Error(w, "Error get subscrible from server ", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(ressum); err != nil {
		s.loger.Error("Error get subscrible from server" + err.Error())
		http.Error(w, "Error get subscrible from server ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	s.loger.Info("Get Summ Subscription Success" + r.URL.String())
}