package server

import (
	"encoding/json"
	"errors"
	"strconv"

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
	CreateSubscr(newsubscr models.SubscrbUser) (int64, error)
	DeleteSubscr(id int64, subs models.SubscrbUserSearch) error
	UpdateSubscr(id int64, subscr models.SubscrbUserSearch) error
	ReadSubscr(id int64, subscr models.SubscrbUserSearch) (models.SubscrbUser, error)
	ReadAllSubscr(off, lim int64, subscr models.SubscrbUserSearch) ([]models.SubscrbUser, error)
	GetSummSubscr(subscr models.SubscrbUserSearch) (models.SummSubscrb, error)
}

// CreateSubscr godoc
// @Summary      Create subscription
// @Description  create subscription from database
// @Tags         create
// @Accept       json
// @Produce      json
// @Param        input body models.SubscrbUser true "subscription struct"
// @Success      200 {object} models.ID
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

    id, err := s.app.CreateSubscr(subscr)
	if err != nil {
		if errors.Is(err, models.Errisempty) {
			s.loger.Error("Error Creating subscrible from server" + err.Error())
			http.Error(w, "UUID, NameService, StartDate is empty", http.StatusBadRequest)
			return
		} else if errors.Is(err, models.Errisuniqu) {
			s.loger.Error("Error Creating subscrible from server" + err.Error())
			http.Error(w, "Error UUID, NameService, StartDate is not unique", http.StatusBadRequest)
			return
		} else if errors.Is(err, models.ErrIvalidArgs) {
			s.loger.Error("Error Creating subscrible from server: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}

		s.loger.Error("Error creating subscrible from server" + err.Error())
		http.Error(w, "Error creating subscrible from server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(models.ID{ID: id}); err != nil {
		s.loger.Error("Error creating subscrible from server" + err.Error())
		http.Error(w, "Error creating subscrible from server", http.StatusInternalServerError)
		return
	}

	s.loger.Info("Create Subscription Success" + r.URL.String())
}

// UpdateSubscr godoc
// @Summary      Update subscription
// @Description  update subscription from database
// @Tags         update
// @Accept       json
// @Produce      json
// @Param        ID query int false "ID"
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

	var id int64 = -1
	qid := r.URL.Query().Get("ID")
	if qid != "" {
		id, _ = strconv.ParseInt(qid, 10, 64)
	}

	var subscr models.SubscrbUserSearch
	if err := json.NewDecoder(r.Body).Decode(&subscr); err != nil {
		s.loger.Error("Error update subscrible from server: " + err.Error())
		http.Error(w, "Error update subscrible from server", http.StatusInternalServerError)
		return
	}

	if err := s.app.UpdateSubscr(id, subscr); err != nil {
		if errors.Is(err, models.Errisempty) {
			s.loger.Error("Error update subscrible from server: " + err.Error())
			http.Error(w, "Error update UUID, NameService, StartDate is empty", http.StatusBadRequest)
			return
		} else if errors.Is(err, models.Errisnotex) {
			s.loger.Error("Error update subscrible from server: " + err.Error())
			http.Error(w, "Your record does not exist", http.StatusBadRequest)
			return
		} else if errors.Is(err, models.ErrIvalidArgs) {
			s.loger.Error("Error update subscrible from server: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}

		s.loger.Error("Error update subscrible from server: " + err.Error())
		http.Error(w, "Error update subscrible from server", http.StatusInternalServerError)
		return
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
// @Param        ID query int false "ID"
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

	var id int64 = -1
	qid := r.URL.Query().Get("ID")
	if qid != "" {
		id, _ = strconv.ParseInt(qid, 10, 64)
	}

	var subscr models.SubscrbUserSearch
	if err := json.NewDecoder(r.Body).Decode(&subscr); err != nil {
		s.loger.Error("Error delete subscrible from server" + err.Error())
		http.Error(w, "Error delete subscrible from server ", http.StatusInternalServerError)
		return
	}

	if err := s.app.DeleteSubscr(id, subscr); err != nil {
		if errors.Is(err, models.Errisnotex) {
			s.loger.Error("Error delete subscrible from server" + err.Error())
			http.Error(w, "Your record does not exist", http.StatusBadRequest)
			return
		} else if errors.Is(err, models.Errisempty) {
			s.loger.Error("Error delete subscrible from server: " + err.Error())
			http.Error(w, "Error delete: UUID, NameService, StartDate is empty", http.StatusBadRequest)
			return
		} else {
			s.loger.Error("Error delete subscrible from server: " + err.Error())
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
// @Param        ID query int false "ID"
// @Param        input body models.SubscrbUserSearch true "filter information"
// @Success      200  {object} models.SubscrbUser
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

	var id int64 = -1
	qid := r.URL.Query().Get("ID")
	if qid != "" {
		id, _ = strconv.ParseInt(qid, 10, 64)
	}

	var subscr models.SubscrbUserSearch
	if err := json.NewDecoder(r.Body).Decode(&subscr); err != nil {
		s.loger.Error("Error get subscrible from server" + err.Error())
		http.Error(w, "Error get subscrible from server ", http.StatusInternalServerError)
		return
	}

	subres, err := s.app.ReadSubscr(id, subscr)
	if err != nil {
		if errors.Is(err, models.Errisempty) {
			s.loger.Error("Error get subscrible from server: " + err.Error())
			http.Error(w, "Error get: UUID, NameService, StartDate is empty", http.StatusBadRequest)
			return
		} else if errors.Is(err, models.ErrIvalidArgs) {
			s.loger.Error("Error get subscrible from server: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}
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
		s.loger.Error("Error get summ subscrible from server" + " method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var subscr models.SubscrbUserSearch
	if err := json.NewDecoder(r.Body).Decode(&subscr); err != nil {
		s.loger.Error("Error get summ subscrible from server" + err.Error())
		http.Error(w, "Error get summ subscrible from server ", http.StatusInternalServerError)
		return
	}

	ressum, err := s.app.GetSummSubscr(subscr)
	if err != nil {
		s.loger.Error("Error get summ subscrible from server" + err.Error())
		http.Error(w, "Error get summ subscrible from server ", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(ressum); err != nil {
		s.loger.Error("Error get summ subscrible from server" + err.Error())
		http.Error(w, "Error get summ subscrible from server ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	s.loger.Info("Get Summ Subscription Success" + r.URL.String())
}

// Get All Subscrible godoc
// @Summary      All subscription
// @Description  All subscription from database
// @Tags         data
// @Accept       json
// @Produce      json
// @Param        OFFSET query int false "OFFSET"
// @Param        LIMIT query int false "LIMIT"
// @Param        input body models.SubscrbUserSearch true "filter information"
// @Success      200  {array} models.SubscrbUser
// @Failure      400  "Bad request error"
// @Failure      404 "Not found error"
// @Failure      405 "Method not allowed"
// @Failure      500  "Internal server error"
// @Router       /searchall [post]
func (s *SubscrServer) GetAllSubscr(w http.ResponseWriter, r *http.Request) {
	s.loger.Info("Get All Subscription" + r.URL.String())
	
	if r.Method != "POST" {
		s.loger.Error("Error get subscrible from server" + " method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var offset int64 = -1
	off := r.URL.Query().Get("OFFSET")
	if off != "" {
		offset, _ = strconv.ParseInt(off, 10, 64)
	}

	var limit int64 = -1
	lim := r.URL.Query().Get("LIMIT")
	if lim != "" {
		limit, _ = strconv.ParseInt(lim, 10, 64)
	}

	var subscr models.SubscrbUserSearch
	if err := json.NewDecoder(r.Body).Decode(&subscr); err != nil {
		s.loger.Error("Error get all subscrible from server" + err.Error())
		http.Error(w, "Error get all subscrible from server ", http.StatusInternalServerError)
		return
	}

	subres, err := s.app.ReadAllSubscr(offset, limit, subscr)
	if err != nil {
		if errors.Is(err, models.ErrIvalidArgs) {
			s.loger.Error("Error get all subscrible from server" + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.loger.Error("Error get all subscrible from server" + err.Error())
		http.Error(w, "Error get all subscrible from server ", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	
	if err := json.NewEncoder(w).Encode(subres); err != nil {
		s.loger.Error("Error get all subscrible from server" + err.Error())
		http.Error(w, "Error get all subscrible from server ", http.StatusInternalServerError)
		return
	}

	s.loger.Info("Get Subscription Success" + r.URL.String())
}