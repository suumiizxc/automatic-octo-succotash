package content

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/suumiizxc/raw_rest1/config"
	"github.com/suumiizxc/raw_rest1/response"
)

type ContentCategory struct {
	ID            uint       `json:"id"`
	NameMN        *string    `json:"name_mn"`
	NameEN        *string    `json:"name_en"`
	Left          *uint      `json:"left"`
	Right         *uint      `json:"right"`
	TreeID        *uint      `json:"tree_id"`
	Level         *uint      `json:"level"`
	ParentID      *uint      `json:"parent_id"`
	IsFeatured    *uint      `json:"is_featured"`
	CoverImage    *string    `json:"cover_image"`
	IsAlwaysShow  *uint      `json:"is_always_show"`
	IsMegaMenu    *uint      `json:"is_mega_menu"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	CreatedUserID *uint      `json:"created_user_id"`
	UpdatedUserID *uint      `json:"updated_user_id"`
}

// Test document test test :P
func (c *ContentCategory) CreateContentCategory(w http.ResponseWriter, r *http.Request) {
	var cca ContentCategory
	err := json.NewDecoder(r.Body).Decode(&cca)
	resp := response.Response{}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err
		resp.Message = "Failed in request body"
		w.Write(resp.ConvertByte())
		return
	}
	sqlStatement := `
	insert into content_category
	(
		name_mn, name_en, left, right, tree_id,
		level, parent_id, is_featured, cover_image, is_always_show,
		is_mega_menu, created_at, updated_at, created_user_id, updated_user_id
	)
	values
	(
		$1, $2, $3, $4, $5,
		$6, $7, $8, $9, $10,
		$11, $12, $13, $14, $15
	) returning id
	`

	var lastID int
	err = config.DB.QueryRow(sqlStatement,
		&cca.NameMN,
		&cca.NameEN,
		&cca.Left,
		&cca.Right,
		&cca.TreeID,
		&cca.Level,
		&cca.ParentID,
		&cca.IsFeatured,
		&cca.CoverImage,
		&cca.IsAlwaysShow,
		&cca.IsMegaMenu,
		time.Now(),
		nil,
		&cca.CreatedUserID,
		nil,
	).Scan(&lastID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err
		resp.Message = "Failed in query"
		w.Write(resp.ConvertByte())
		return
	}
	resp.Data = lastID
	resp.Message = "Successfully created content category"
	w.WriteHeader(http.StatusCreated)
	w.Write(resp.ConvertByte())

}

func (c *ContentCategory) GetContentCategoryList(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	cca := []ContentCategory{}
	sqlStatement := `
	select id, name_mn, name_en, left, right,
		tree_id, level, parent_id, is_featured, cover_image, is_always_show,
		is_mega_menu, created_at, updated_at, created_user_id, updated_user_id
	from content_category
	`

	rows, err := config.DB.Query(sqlStatement)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = err
		resp.Message = "Failed in query"
		w.Write(resp.ConvertByte())
		return
	}
	for rows.Next() {
		var cca1 ContentCategory
		rows.Scan(
			&cca1.ID,
			&cca1.NameMN,
			&cca1.NameEN,
			&cca1.Left,
			&cca1.Right,
			&cca1.TreeID,
			&cca1.Level,
			&cca1.ParentID,
			&cca1.IsFeatured,
			&cca1.CoverImage,
			&cca1.IsAlwaysShow,
			&cca1.IsMegaMenu,
			&cca1.CreatedAt,
			&cca1.UpdatedAt,
			&cca1.CreatedUserID,
			&cca1.UpdatedUserID,
		)
		cca = append(cca, cca1)
	}
	resp.Data = cca
	resp.Message = "Successfully get content category list"
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())

}

func (c *ContentCategory) GetContentCategoryByID(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	cca := []ContentCategory{}
	id := mux.Vars(r)["id"]
	sqlStatement := `
	select id, name_mn, name_en, left, right,
		tree_id, level, parent_id, is_featured, cover_image, is_always_show,
		is_mega_menu, created_at, updated_at, created_user_id, updated_user_id
	from content_category
	where id = $1
	`

	rows, err := config.DB.Query(sqlStatement, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = err
		resp.Message = "Failed in query"
		w.Write(resp.ConvertByte())
		return
	}
	for rows.Next() {
		var cca1 ContentCategory
		rows.Scan(
			&cca1.ID,
			&cca1.NameMN,
			&cca1.NameEN,
			&cca1.Left,
			&cca1.Right,
			&cca1.TreeID,
			&cca1.Level,
			&cca1.ParentID,
			&cca1.IsFeatured,
			&cca1.CoverImage,
			&cca1.IsAlwaysShow,
			&cca1.IsMegaMenu,
			&cca1.CreatedAt,
			&cca1.UpdatedAt,
			&cca1.CreatedUserID,
			&cca1.UpdatedUserID,
		)

		cca = append(cca, cca1)
	}
	resp.Data = cca
	resp.Message = "Successfully get content category by id"
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}
