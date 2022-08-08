package blog

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/suumiizxc/raw_rest1/config"
	"github.com/suumiizxc/raw_rest1/response"
)

type Blog struct {
	ID       uint    `json:"id"`
	Language *string `json:"language"`
	Title    *string `json:"title"`
	Poster   *string `json:"poster"`
	Cover    *string `json:"cover"`

	BriefDescription *string    `json:"brief_description"`
	Description      *string    `json:"description"`
	CategoryID       *uint      `json:"category_id"`
	IsPublished      *uint      `json:"is_published"`
	IsFeatured       *uint      `json:"is_featured"`
	FeatureStartDate *time.Time `json:"feature_start_date"`
	FeatureEndDate   *time.Time `json:"feature_end_date"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	CreatedUserID    *uint      `json:"created_user_id"`
	UpdatedUserID    *uint      `json:"updated_user_id"`
}

func (b *Blog) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var bl Blog
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
	insert into blog
	(
		language, title, poster, cover, brief_description,
		description, category_id, is_published, is_featured, 
		feature_start_date, feature_end_date, created_at, created_user_id
	)
	values
	(
		$1, $2, $3, $4, $5,
		$6, $7, $8, $9,
		$10, $11, $12, $13
	)
	returning id
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
	resp.Message = "Successfully created blog"
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}

func (b *Blog) GetBlogByID(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	id := mux.Vars(r)["id"]
	bl := []Blog{}
	sqlStatement := `
		select id, language, title, poster, cover, brief_description,
			description, category_id, is_published, is_featured,
			feature_start_date, feature_end_date, created_at, 
			updated_at, created_user_id, updated_user_id 
		from blog
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
		var bl1 Blog
		rows.Scan(
			&bl1.ID,
			&bl1.Language,
			&bl1.Title,
			&bl1.Poster,
			&bl1.Cover,
			&bl1.BriefDescription,
			&bl1.Description,
			&bl1.CategoryID,
			&bl1.IsPublished,
			&bl1.IsFeatured,
			&bl1.FeatureStartDate,
			&bl1.FeatureEndDate,
			&bl1.CreatedAt,
			&bl1.UpdatedAt,
			&bl1.CreatedUserID,
			&bl1.UpdatedUserID,
		)
		bl = append(bl, bl1)
	}

	resp.Message = "Successfully created"
	resp.Data = bl
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}

func (b *Blog) GetBlogList(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	bl := []Blog{}
	sqlStatement := `
		select id, language, title, poster, cover, brief_description,
			description, category_id, is_published, is_featured,
			feature_start_date, feature_end_date, created_at,
			updated_at, created_user_id, updated_user_id
		from blog
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
		var bl1 Blog
		rows.Scan(
			&bl1.ID,
			&bl1.Language,
			&bl1.Title,
			&bl1.Poster,
			&bl1.Cover,
			&bl1.BriefDescription,
			&bl1.Description,
			&bl1.CategoryID,
			&bl1.IsPublished,
			&bl1.IsFeatured,
			&bl1.FeatureStartDate,
			&bl1.FeatureEndDate,
			&bl1.CreatedAt,
			&bl1.UpdatedAt,
			&bl1.CreatedUserID,
			&bl1.UpdatedUserID,
		)
		bl = append(bl, bl1)
	}
	resp.Message = "Successfully get blogs"
	resp.Data = bl
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())

}
