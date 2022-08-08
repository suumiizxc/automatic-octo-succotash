package content

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/suumiizxc/raw_rest1/config"
	"github.com/suumiizxc/raw_rest1/response"
)

type ContentToContentAuthor struct {
	ID              uint  `json:"id"`
	ContentID       *uint `json:"content_id"`
	ContentAuthorID *uint `json:"content_author_id"`
}

func (c *ContentToContentAuthor) CreateContentToContentAuthor(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	cca := ContentToContentAuthor{}
	err := json.NewDecoder(r.Body).Decode(&cca)
	if err != nil {
		resp.Message = "Failed in request body"
		resp.Error = err
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp.ConvertByte())
		return
	}
	sqlStatement := `
	insert into content_to_content_author 
	(
		content_id, content_author_id
	)
	values
	(
		$1, $2
	) returning id
	`

	var lastID int

	err = config.DB.QueryRow(sqlStatement, cca.ContentID, cca.ContentAuthorID).Scan(&lastID)
	if err != nil {
		resp.Message = "Failed in query"
		resp.Error = err
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp.ConvertByte())
		return
	}
	resp.Message = "Successfully created"
	resp.Data = lastID
	w.WriteHeader(http.StatusCreated)
	w.Write(resp.ConvertByte())
}

func (c *ContentToContentAuthor) GetContentToContentAuthorList(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	cca := []ContentToContentAuthor{}
	sqlStatement := `
		select id, content_id, content_author_id
		from content_to_content_author
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
		var cca1 ContentToContentAuthor
		rows.Scan(
			&cca1.ID,
			&cca1.ContentID,
			&cca1.ContentAuthorID,
		)
		cca = append(cca, cca1)
	}
	resp.Data = cca
	resp.Message = "Successfully get content and content author list"
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}

func (c *ContentToContentAuthor) GetGontentToContentAuthorByID(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	cca := []ContentToContentAuthor{}
	vars := mux.Vars(r)
	id := vars["id"]

	sqlStatement := `
		select id, content_id, content_author_id
		from content_to_content_author
		where id = $1
	`
	rows, err := config.DB.Query(sqlStatement, id)
	if err != nil {
		resp.Message = "Failed in query"
		resp.Error = err
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp.ConvertByte())
		return
	}
	for rows.Next() {
		var cca1 ContentToContentAuthor
		rows.Scan(
			&cca1.ID,
			&cca1.ContentID,
			&cca1.ContentAuthorID,
		)
		cca = append(cca, cca1)
	}
	resp.Data = cca
	resp.Message = "Successfully get content and content author by id"
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}
