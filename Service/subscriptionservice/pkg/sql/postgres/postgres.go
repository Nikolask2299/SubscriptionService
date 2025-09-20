package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"subscriptionservice/interal/models"

	sq "github.com/Masterminds/squirrel"

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

func (p *Postgres) CreateSubscr(subscr models.SubscrbUser) (int64, error) {
    
    psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
    
    qr, args, err := psql.Insert("Services").Columns("services_name").
    Values(subscr.NameService).
    Suffix("ON CONFLICT (services_name) DO NOTHING").ToSql()
    
    if err != nil {
        return -1, err
    }

    if _, err := p.db.Exec(qr, args...); err != nil {
        return -1, err
    }
    
    ars := make([]interface{}, 0, 10)

    sb := psql.Insert("UserSubscr").
    Columns("user_id", "start_date", "end_date", "services_name", "price")
    
    ars = append(ars, subscr.UUID)
    ars = append(ars, sq.Expr("to_date(?, 'MM-YYYY')", subscr.StartDate))
    
    if subscr.EndDate == "" {
        ars = append(ars, sql.NullTime{})
    } else {
        ars = append(ars, sq.Expr("to_date(?, 'MM-YYYY')", subscr.EndDate))
    }

    ars = append(ars, subscr.NameService)
    ars = append(ars, subscr.Price)

    qr, args, err = sb.Values(ars...).
    Suffix("RETURNING subscrb_id").
    ToSql()
    if err != nil {
        return -1, err
    }

    rs, err := p.db.Query(qr, args...)
    if err != nil {
        return -1, err
    }
    var id int64  
    rs.Next()
    err = rs.Scan(&id)
    if err != nil {
        return -1, err
    }
    rs.Close()

    return id, nil
}

func (p *Postgres) DeleteSubscr(id int64, subs models.SubscrbUserSearch) error {
  psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
  
  query := psql.Delete("UserSubscr")
   
   if id > 0 {
        query.Where(sq.Eq{
            "subscrb_id": id,
        })
   }
   
   if subs.UUID != "" {
        query.Where(sq.Eq{
            "user_id":subs.UUID, 
        })
   }
   
   if subs.NameService != "" {
        query.Where(sq.Eq{
            "services_name":subs.NameService,
        })
   }

   if subs.StartDate != "" {
    query.Where(sq.Eq{
        "start_date":sq.Expr("to_date(?, 'MM-YYYY')", subs.StartDate),
    })
   }

   qr, args, err := query.ToSql()
   if err != nil {
        return err
   }

    res, err := p.db.Exec(qr, args...)
    if err != nil {
        return err
    }

    rw, err := res.RowsAffected()
    if err != nil {
        return err
    }

    if rw == 0 {
        return models.Errisnotex
    }
    
    return nil
}

func (p *Postgres) UpdateSubscr(id int64, subscr models.SubscrbUserSearch) error {
    psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

    query := psql.Update("UserSubscr")
    

    if  subscr.EndDate != "" {
       query = query.Set("end_date", sq.Expr("to_date(?, 'MM-YYYY')", subscr.EndDate))
    }
    
    if subscr.Price > 0 {
       query = query.Set("price", subscr.Price)
    }

    if id > 0 {
        if subscr.UUID != "" {
            query = query.Set("user_id", subscr.UUID)
        }

        if subscr.NameService != "" {
           query = query.Set("services_name", subscr.NameService)
        }

        if subscr.StartDate != "" {
           query = query.Set("start_date", sq.Expr("to_date(?, 'MM-YYYY')", subscr.StartDate))
        }

       query = query.Where("subscrb_id = ?", id)
    } else {
        query = query.Where("user_id = ?", subscr.UUID)
        query = query.Where("services_name = ?", subscr.NameService)
        query = query.Where("start_date = to_date(?, 'MM-YYYY')", subscr.StartDate)
    }
    
    qr, args, err := query.ToSql()
    if err != nil {
        return err
    }

    res, err := p.db.Exec(qr, args...)
    if err != nil {
        return err
    }

    rw, err := res.RowsAffected()
    if err != nil {
        return err
    }

    if rw == 0 {
        return models.Errisnotex
    }

    return err
}

func (p *Postgres) ReadSubscr(id int64, subscr models.SubscrbUserSearch) (models.SubscrbUser, error) {
   
    psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
   
    query := psql.Select("subscrb_id","user_id","TO_CHAR(start_date, 'MM-YYYY')","COALESCE(TO_CHAR(end_date, 'MM-YYYY'), '')","services_name","price").
    From("UserSubscr")

   if id > 0 {
       query = query.Where("subscrb_id = ?", id)
    }

    if subscr.UUID != "" {
       query = query.Where("user_id = ?", subscr.UUID)
    }
    
   if subscr.NameService != "" {
       query = query.Where("services_name = ?", subscr.NameService)
    }

    if subscr.StartDate != "" {
       query = query.Where("start_date = to_date(?, 'MM-YYYY')", subscr.StartDate)
    }

    if subscr.EndDate != "" {
       query = query.Where("end_date = to_date(?, 'MM-YYYY')", subscr.EndDate)
    }
    
    if subscr.Price > 0 {
       query = query.Where("price = ?", subscr.Price)
    }

    qr, args, err := query.ToSql()
    if err != nil {
        return models.SubscrbUser{}, err
    }

    row, err := p.db.Query(qr, args...)
    if err != nil {
        return models.SubscrbUser{}, err
    }

    var sub models.SubscrbUser
    row.Next()
    if err := row.Scan(&sub.SubscrID, &sub.UUID, &sub.StartDate, &sub.EndDate, &sub.NameService, &sub.Price); err != nil {
        return models.SubscrbUser{}, err
    }
    row.Close()

    return sub, nil
}

func (p *Postgres) GetSummSubscr(subscr models.SubscrbUserSearch) (models.SummSubscrb, error){
    res := models.SummSubscrb{
        UUID: subscr.UUID,
        StartDate: subscr.StartDate,
        EndDate: subscr.EndDate,
        ServiceName: make([]string, 0, 10),
    }

    psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
 
    query := psql.Select("SUM(price * ((EXTRACT(YEAR FROM LEAST(end_date, to_date($1, 'MM-YYYY'))) - EXTRACT(YEAR FROM GREATEST(start_date, to_date($2, 'MM-YYYY')))) * 12 + (EXTRACT(MONTH FROM LEAST(end_date, to_date($1, 'MM-YYYY'))) - EXTRACT(MONTH FROM GREATEST(start_date, to_date($2, 'MM-YYYY')))) + 1)), COALESCE(array_to_json(ARRAY_AGG(services_name)), '[]') FROM UserSubscr").
    Where(`start_date <= to_date($1, 'MM-YYYY') AND (end_date IS NULL OR end_date >= to_date($2, 'MM-YYYY'))`, subscr.EndDate, subscr.StartDate)

   if subscr.UUID != "" {
      query = query.Where("user_id = $3", subscr.UUID)
   }

   if subscr.NameService != "" {
        if subscr.UUID != "" {
            query = query.Where("services_name = $4", subscr.NameService)
        } else {
            query = query.Where("services_name = $3", subscr.NameService)
        }
   }

    qr, args, err := query.ToSql()
    if err != nil {
        return res, err
    }

    row, err := p.db.Query(qr, args...)
    if err != nil {
       return res, err
    }
    defer row.Close()

    row.Next()
    mass := []byte{}
    if err := row.Scan(&res.Summ, &mass); err != nil {
        return res, err
    }

    err = json.Unmarshal(mass, &res.ServiceName)
    if err != nil {
        return res, err
    }
    
    return res, nil
}

func (p *Postgres) ReadAllSubscr(off, lim int64, subscr models.SubscrbUserSearch) ([]models.SubscrbUser, error) {
    psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
    
    query := psql.Select("subscrb_id", "user_id", "TO_CHAR(start_date, 'MM-YYYY')", "COALESCE(TO_CHAR(end_date, 'MM-YYYY'), '')", "services_name", "price").
    From("UserSubscr")

    if subscr.UUID != "" {
      query = query.Where("user_id = ?", subscr.UUID)
    }
    
    if subscr.NameService != "" {
       query = query.Where("services_name = ?", subscr.NameService)
    }

    if subscr.StartDate != "" {
       query = query.Where("start_date = to_date(?, 'MM-YYYY')", subscr.StartDate)
    }

    if subscr.EndDate != "" {
       query = query.Where("end_date = to_date(?, 'MM-YYYY')", subscr.EndDate)
    }
    
    if subscr.Price > 0 {
       query = query.Where("price = ?", subscr.Price)
    }
    
    if off > 0 && lim > off {
       query = query.Offset(uint64(off))
    }

    if lim > 0 {
        query = query.Limit(uint64(lim))
    }

   
    qr, args, err := query.ToSql()
    if err != nil {
        return nil, err
    }

    row, err := p.db.Query(qr, args...)
    if err != nil {
       return nil, err
    }
    defer row.Close()

    res := make([]models.SubscrbUser, 0, 10)
    for row.Next() {
       var sub models.SubscrbUser
       if err := row.Scan(&sub.SubscrID, &sub.UUID, &sub.StartDate, &sub.EndDate, &sub.NameService, &sub.Price); err != nil {
            return nil, err
       }
       res = append(res, sub)
    }

    return res, nil
}