package blog

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/suumiizxc/raw_rest1/config"
	"github.com/suumiizxc/raw_rest1/response"
)

type BlogView struct {
	ID            uint       `json:"id"`
	Session       *string    `json:"session"`
	IpAddress     *string    `json:"ip_address"`
	BlogID        *uint      `json:"blog_id"`
	CreatedAt     *time.Time `json:"created_at"`
	CreatedUserID *uint      `json:"created_user_id"`
}

func (b *BlogView) CreateBlogView(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	bl := BlogView{}
	err := json.NewDecoder(r.Body).Decode(&bl)
	if err != nil {
		resp.Message = "Failed in request body"
		resp.Error = err
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp.ConvertByte())
		return
	}

	sqlStatement := `
	insert into blog_view
	(
		session, ip_address, blog_id, created_at, created_user_id
	)
	values
	(
		$1, $2, $3, $4, $5
	)
	returning id
	`
	var lastID int
	err = config.DB.QueryRow(sqlStatement,
		bl.Session,
		bl.IpAddress,
		bl.BlogID,
		time.Now(),
		bl.CreatedUserID,
	).Scan(&lastID)

	if err != nil {
		resp.Error = err
		resp.Message = "Failed in query"
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp.ConvertByte())
		return
	}
	resp.Data = lastID
	resp.Message = "Successfully created"
	w.WriteHeader(http.StatusCreated)
	w.Write(resp.ConvertByte())
}

func (b *BlogView) GetBlogViewList(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	bl := []BlogView{}
	sqlStatement := `
	select id, session, ip_address, blog_id, created_at, created_user_id
	from blog_view
	`

	rows, err := config.DB.Query(sqlStatement)
	if err != nil {
		resp.Error = err
		resp.Message = "Failed in query"
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp.ConvertByte())
		return
	}

	for rows.Next() {
		var bl1 BlogView
		rows.Scan(
			&bl1.ID,
			&bl1.Session,
			&bl1.IpAddress,
			&bl1.BlogID,
			&bl1.CreatedAt,
			&bl1.CreatedUserID,
		)
		bl = append(bl, bl1)
	}
	resp.Data = bl
	resp.Message = "Successfully get blog view list"
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}
