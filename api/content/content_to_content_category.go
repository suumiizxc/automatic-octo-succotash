package content

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/suumiizxc/raw_rest1/config"
	"github.com/suumiizxc/raw_rest1/response"
)

type ContentToContentCategory struct {
	ID                uint `json:"id"`
	ContentID         uint `json:"content_id"`
	ContentCategoryID uint `json:"content_category_id"`
}

func (c *ContentToContentCategory) CreateContentToContentCategory(w http.ResponseWriter, r *http.Request) {
	var cca ContentToContentCategory
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
	insert into content_to_content_category
	(
		content_id, content_category_id
	)
	values
	(
		$1, $2
	)
	returning id

	`
	var lastID int
	err = config.DB.QueryRow(sqlStatement,
		&cca.ContentID,
		&cca.ContentCategoryID,
	).Scan(&lastID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = err
		resp.Message = "Failed in query"
		w.Write(resp.ConvertByte())
		return
	}
	resp.Data = lastID
	resp.Message = "Successfully created"
	w.WriteHeader(http.StatusCreated)
	w.Write(resp.ConvertByte())
}

func (c *ContentToContentCategory) GetContentToContentCategoryList(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	cca := []ContentToContentCategory{}
	sqlStatement := `
	select id, content_id, content_category_id
	from content_to_content_category
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
		var cca1 ContentToContentCategory
		rows.Scan(
			&cca1.ID,
			&cca1.ContentID,
			&cca1.ContentCategoryID,
		)
		cca = append(cca, cca1)
	}
	resp.Data = cca
	resp.Message = "Successfully get content and content category list"
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}

func (c *ContentToContentCategory) GetContentToContentCategoryByID(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	cca := []ContentToContentCategory{}
	id := mux.Vars(r)["id"]
	sqlStatement := `
	select id, content_id, content_category_id
	from content_to_content_category
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
		var cca1 ContentToContentCategory
		rows.Scan(
			&cca1.ID,
			&cca1.ContentID,
			&cca1.ContentCategoryID,
		)
		cca = append(cca, cca1)
	}
	resp.Data = cca
	resp.Message = "Successfully get content and content category by id"
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}
