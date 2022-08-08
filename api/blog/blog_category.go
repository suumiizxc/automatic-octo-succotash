package blog

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/suumiizxc/raw_rest1/config"
	"github.com/suumiizxc/raw_rest1/response"
)

type BlogCategory struct {
	ID            uint       `json:"id"`
	NameMn        *string    `json:"name_mn"`
	NameEn        *string    `json:"name_en"`
	IsFeatured    *uint      `json:"is_featured"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	CreatedUserID *uint      `json:"created_user_id"`
	UpdatedUserID *uint      `json:"updated_user_id"`
}

func (bc *BlogCategory) CreateBlogCategory(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	var bla BlogCategory
	err := json.NewDecoder(r.Body).Decode(&bla)
	if err != nil {
		resp.Error = err
		resp.Message = "Failed in request body"
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp.ConvertByte())
		return
	}
	sqlStatement := `
		insert into blog_category
		(name_mn, name_en, is_featured, created_at, created_user_id)
		values
		(
			$1, $2, $3, $4, $5
		) returning id
	`
	var lastID int
	err = config.DB.QueryRow(sqlStatement,
		&bla.NameMn,
		&bla.NameEn,
		&bla.IsFeatured,
		time.Now(),
		&bla.CreatedUserID,
	).Scan(&lastID)
	if err != nil {
		resp.Message = "Failed in query"
		resp.Error = err
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp.ConvertByte())
		return
	}
}

func (bc *BlogCategory) GetBlogCategoryList(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	var bla []BlogCategory
	sqlStatement := `
		select id, name_mn, name_en, is_featured, created_at, updated_at,
			created_user_id, updated_user_id
		from blog_category
	`
	rows, err := config.DB.Query(sqlStatement)
	if err != nil {
		resp.Error = err
		resp.Message = "Failed in query"
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp.ConvertByte())
		return
	}
	for rows.Next() {
		var bla1 BlogCategory
		rows.Scan(
			&bla1.ID,
			&bla1.NameMn,
			&bla1.NameEn,
			&bla1.IsFeatured,
			&bla1.CreatedAt,
			&bla1.UpdatedAt,
			&bla1.CreatedUserID,
			&bla1.UpdatedUserID,
		)

		bla = append(bla, bla1)
	}
	resp.Data = bla
	resp.Message = "Successfully get blog category list"
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}
