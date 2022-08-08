package blog

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/suumiizxc/raw_rest1/config"
	"github.com/suumiizxc/raw_rest1/response"
)

type BlogComment struct {
	ID            uint       `json:"id"`
	Description   *string    `json:"description"`
	BlogID        *uint      `json:"blog_id"`
	ParentID      *uint      `json:"parent_id"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	CreatedUserID *uint      `json:"created_user_id"`
}

func (b *BlogComment) CreateBlogComment(w http.ResponseWriter, r *http.Request) {
	var bl BlogComment
	resp := response.Response{}
	err := json.NewDecoder(r.Body).Decode(&bl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err
		resp.Message = "Failed in request body"
		w.Write(resp.ConvertByte())
		return
	}
	sqlStatement := `
	insert into blog_comment
	(
		description, blog_id, parent_id, created_at, created_user_id 
	)
	values
	(
		$1, $2, $3, $4, $5
	) returning id
	`
	var lastID int
	err = config.DB.QueryRow(sqlStatement).Scan(&lastID)
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

func (b *BlogComment) GetBlogCommentByID(w http.ResponseWriter, r *http.Request) {
	bl := []BlogComment{}
	resp := response.Response{}
	id := mux.Vars(r)["id"]
	sqlStatement := `
	select id, description, blog_id, parent_id, created_at,
		updated_at, created_user_id
	from blog_comment
	where id = $1
	`
	rows, err := config.DB.Query(sqlStatement, id)
	if err != nil {
		resp.Error = err
		resp.Message = "Failed in query"
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp.ConvertByte())
		return
	}

	for rows.Next() {
		var bl1 BlogComment
		rows.Scan(
			&bl1.ID,
			&bl1.Description,
			&bl1.BlogID,
			&bl1.ParentID,
			&bl1.CreatedAt,
			&bl1.UpdatedAt,
			&bl1.CreatedUserID,
		)
		bl = append(bl, bl1)
	}

	resp.Message = "Successfully created"
	resp.Data = bl
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}

func (b *BlogComment) GetBlogCommentList(w http.ResponseWriter, r *http.Request) {
	bl := []BlogComment{}
	resp := response.Response{}
	sqlStatement := `
	select id, description, blog_id, parent_id, created_at,
		updated_at, created_user_id
	from blog_comment
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
		var bl1 BlogComment
		rows.Scan(
			&bl1.ID,
			&bl1.Description,
			&bl1.BlogID,
			&bl1.ParentID,
			&bl1.CreatedAt,
			&bl1.UpdatedAt,
			&bl1.CreatedUserID,
		)

		bl = append(bl, bl1)
	}

	resp.Data = bl
	resp.Message = "Successfully get blog comment list"
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}
