package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"subscriptionservice/interal/models"

	_ "github.com/lib/pq"
)

type Postgres struct {
    db *sql.DB
}

func NewPostgres(user, password, dbname, host, port string) (*Postgres, error) {
	psqlInfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)
    db, err := sql.Open("postgres", psqlInfo)
    if err!= nil {
        return nil, err
    }

    err = db.Ping()
    if err!= nil {
        return nil, err
    }

    return &Postgres{db: db}, nil
}


func (p *Postgres) Close() error {
    return p.db.Close()
}

func (p *Postgres) CreateSubscr(subscr models.SubscrbUser) error {
    query_service := `INSERT INTO Services (services_name) VALUES ($1)`
    query_subscr := `INSERT INTO UserSubscr (user_id, start_date, end_date, services_name, price) VALUES ($1, to_date($2, 'MM-YYYY'), to_date($3, 'MM-YYYY'), $4, $5)`

    if _, err := p.db.Exec(query_service, subscr.NameService); err != nil {
        return err
    }

    if _, err := p.db.Exec(query_subscr, subscr.UUID, subscr.StartDate, subscr.EndDate, subscr.NameService, subscr.Price); err != nil {
        return err;
    }

    return nil
}

func (p *Postgres) DeleteSubscr(subs models.SubscrbUserSearch) error {
    query := `DELETE FROM UserSubscr WHERE user_id = $1 AND services_name = $2 AND start_date = to_date($3, 'MM-YYYY');`
    if _, err := p.db.Exec(query, subs.UUID, subs.NameService, subs.StartDate); err != nil {
        if err == sql.ErrNoRows {
            return nil
        } else {
            return err
        }
    }
    return nil
}

func (p *Postgres) UpdateSubscr(subscr map[string]string) error {
    query := `UPDATE UserSubscr SET `
    
     if subscr["user_id"] != "" {
        query += `, user_id`
        query += ` = '` + subscr["user_id"]  + `'`
    }
    
    if subscr["services_name"] != "" {
        query += `, services_name`
        query += ` = '` + subscr["services_name"] + `'`
    }

    if  subscr["start_date"] != "" {
         query += `, start_date`
        query += ` = to_date('` + subscr["start_date"] + `, 'MM-YYYY')'`
    }

    if  subscr["end_date"] != "" {
        query += `, end_date`
        query += ` = to_date('` + subscr["end_date"] + `, 'MM-YYYY')'`
    }
    
    if subscr["price"] != "" {
        query += `, price`
        query += ` = '` + subscr["price"]  + `'`
    }
    query = strings.Replace(query, "SET ,", "SET ", -1)

    query += ` WHERE "user_id" = '` + subscr["user_id"] + `'`
    query += ` AND "services_name" = '` + subscr["services_name"] + `'`
    query += ` AND "start_date" = to_date('` + subscr["start_date"] + `', 'MM-YYYY');`

    _, err := p.db.Exec(query)
    return err
}

func (p *Postgres) ReadSubscr(subscr models.SubscrbUserSearch) ([]models.SubscrbUser, error) {
    query := `SELECT user_id, start_date, end_date, services_name, price FROM UserSubscr WHERE `
    
    if subscr.UUID != "" {
        query += ` AND user_id`
        query += ` = '` + subscr.UUID + `'`
    }
    
    if subscr.NameService != "" {
        query += ` AND services_name`
        query += ` = '` + subscr.NameService + `'`
    }

    if subscr.StartDate != "" {
         query += ` AND start_date`
        query += ` = to_date('` + subscr.StartDate + `, 'MM-YYYY')'`
    }

    if subscr.EndDate != "" {
        query += ` AND end_date`
        query += ` = to_date('` + subscr.EndDate + `, 'MM-YYYY')'`
    }
    
    if subscr.Price != 0 {
        query += ` AND price`
        pr := strconv.Itoa(subscr.Price)
        query += ` = '` + pr + `'`
    }
   
    query += ";"
    query = strings.Replace(query, "WHERE AND", "WHERE ", -1)

    rows, err := p.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    subs := make([]models.SubscrbUser, 0, 10)
    for rows.Next() {
       var sub models.SubscrbUser
       if err := rows.Scan(&sub.UUID, &sub.StartDate, &sub.EndDate, &sub.NameService, &sub.Price); err != nil {
            return nil, err
       }
       subs = append(subs, sub)
    }

    return subs, nil
}


func (p *Postgres) GetSummSubscr(subscr models.SubscrbUserSearch) (models.SummSubscrb, error){
    return models.SummSubscrb{}, nil
}