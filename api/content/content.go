package content

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/suumiizxc/raw_rest1/config"
	"github.com/suumiizxc/raw_rest1/response"
)

type Content struct {
	ID               uint       `json:"id"`
	Name             *string    `json:"name"`
	PublishedYear    *uint      `json:"published_year"`
	PageCount        *uint      `json:"page_count"`
	ContentScope     *string    `json:"content_scope"`
	IsPurchase       *uint      `json:"is_purchase"`
	IsRent           *uint      `json:"is_rent"`
	RentPrice        *float32   `json:"rent_price"`
	Author           *string    `json:"author"`
	ContentPurpose   *string    `json:"content_purpose"`
	BriefDescription *string    `json:"brief_description"`
	Summary          *string    `json:"summary"`
	Description      *string    `json:"description"`
	TableList        *string    `json:"table_list"`
	ImageList        *string    `json:"image_list"`
	Cover            *string    `json:"cover"`
	FeatureStartDate *time.Time `json:"feature_start_date"`
	FeatureEndDate   *time.Time `json:"feature_end_date"`
	IsFeatured       *uint      `json:"is_featured"`
	FeatureImage     *string    `json:"feature_image"`
	PreviewFile      *string    `json:"preview_file"`
	Language         *string    `json:"language"`
	RentMonth        *uint      `json:"rent_month"`
	IsFree           *uint      `json:"is_free"`
	IsPublished      *uint      `json:"is_published"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	CreatedUserID    *uint      `json:"created_user_id"`
	UpdatedUserID    *uint      `json:"updated_user_id"`
}

func (c *Content) CreateContent(w http.ResponseWriter, r *http.Request) {
	var co Content
	err := json.NewDecoder(r.Body).Decode(&co)
	resp := response.Response{}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err
		resp.Message = "Failed in request body"
		w.Write(resp.ConvertByte())
		return
	}
	sqlStatement := `
	insert into content 
	(
		name, published_year, page_count, content_scope, is_purchase,
		is_rent, rent_price, author, content_purpose, brief_description,
		summary, description, table_list, image_list, cover, 
		feature_start_date, feature_end_date, is_featured, feature_image,
		preview_file, "language",rent_month, is_free, is_published,
		created_user_id, created_at
	)
	values
	(
		$1, $2, $3, $4, $5,
		$6, $7, $8, $9, $10,
		$11, $12, $13, $14, $15,
		$16, $17, $18, $19,
		$20, $21, $22, $23, $24,
		$25, $26
	) returning id
	`
	var lastID int
	err = config.DB.QueryRow(sqlStatement,
		co.Name, co.PublishedYear, co.PageCount, co.ContentScope, co.IsPurchase,
		co.IsRent, co.RentPrice, co.Author, co.ContentPurpose, co.BriefDescription,
		co.Summary, co.Description, co.TableList, co.ImageList, co.Cover,
		co.FeatureStartDate, co.FeatureEndDate, co.IsFeatured, co.FeatureImage,
		co.PreviewFile, co.Language, co.RentMonth, co.IsFree, co.IsPublished,
		co.CreatedUserID, time.Now(),
	).Scan(&lastID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = err
		resp.Message = "Failed in insert query"
		w.Write(resp.ConvertByte())
		return
	}
	w.WriteHeader(http.StatusCreated)
	resp.Data = lastID
	resp.Message = "Successfully created"
	w.Write(resp.ConvertByte())
}

func (c *Content) GetContentById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	resp := response.Response{}
	co := []Content{}
	sqlStatement := `
	select id, name, published_year, page_count, content_scope, is_purchase,
	is_rent, rent_price, author, content_purpose, brief_description,
	summary, description, table_list, image_list, cover, 
	feature_start_date, feature_end_date, is_featured, feature_image,
	preview_file, "language",rent_month, is_free, is_published,
	created_user_id
	from content where id = $1
	`
	rows, err := config.DB.Query(sqlStatement, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = err
		resp.Message = "Failed in get query"
		w.Write(resp.ConvertByte())
		return
	}
	for rows.Next() {
		var co1 Content
		rows.Scan(
			&co1.ID,
			&co1.Name,
			&co1.PublishedYear,
			&co1.PageCount,
			&co1.ContentScope,
			&co1.IsPurchase,
			&co1.IsRent,
			&co1.RentPrice,
			&co1.Author,
			&co1.ContentPurpose,
			&co1.BriefDescription,
			&co1.Summary,
			&co1.Description,
			&co1.TableList,
			&co1.ImageList,
			&co1.Cover,
			&co1.FeatureStartDate,
			&co1.FeatureEndDate,
			&co1.IsFeatured,
			&co1.FeatureImage,
			&co1.PreviewFile,
			&co1.Language,
			&co1.RentMonth,
			&co1.IsFree,
			&co1.IsPublished,
			&co1.CreatedUserID,
		)
		co = append(co, co1)
	}
	w.WriteHeader(http.StatusOK)
	resp.Data = co
	resp.Message = "Successfully get content by id"
	w.Write(resp.ConvertByte())

}
