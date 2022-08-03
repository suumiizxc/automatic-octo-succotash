package content

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/suumiizxc/raw_rest1/config"
	"github.com/suumiizxc/raw_rest1/response"
)

type ContentAuthor struct {
	ID            uint      `json:"id"`
	AuthorType    *string   `json:"author_type"`
	Name          *string   `json:"name"`
	Image         *string   `json:"image"`
	Description   *string   `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedUserID *uint     `json:"created_user_id"`
	UpdatedUserID *uint     `json:"updated_user_id"`
}

func (c *ContentAuthor) CreateContentAuthor(w http.ResponseWriter, r *http.Request) {
	var ca ContentAuthor
	err := json.NewDecoder(r.Body).Decode(&ca)
	resp := response.Response{}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err
		resp.Message = "Failed in request body"
		w.Write(resp.ConvertByte())
		return
	}
	sqlStatement := `
	insert into content_author 
	(
		author_type, "name", image, description,
		created_user_id, updated_user_id
	)
	values
	(
		$1, $2, $3, $4,
		$5, $6
	) returning id
	`
	var lastID int
	err = config.DB.QueryRow(sqlStatement,
		ca.AuthorType, ca.Name, ca.Image, ca.Description,
		ca.CreatedUserID, ca.UpdatedUserID,
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

func (c *ContentAuthor) GetContentAuthorList(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	var ca []ContentAuthor
	sqlStatement := `
		select id, author_type, name, image, description, created_at, updated_at, created_user_id, updated_user_id
		from content_author
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
		var ca1 ContentAuthor
		rows.Scan(
			&ca1.ID,
			&ca1.AuthorType,
			&ca1.Name,
			&ca1.Image,
			&ca1.Description,
			&ca1.CreatedAt,
			&ca1.UpdatedAt,
			&ca1.CreatedUserID,
			&ca1.UpdatedUserID,
		)
		ca = append(ca, ca1)
	}
	resp.Data = ca
	resp.Message = "Successfully fetch authors"
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}

func (c *ContentAuthor) GetContentAuthorByID(w http.ResponseWriter, r *http.Request) {
	// resp := response.Response{}
	vars := mux.Vars(r)
	id := vars["id"]
	resp := response.Response{}
	ca := []ContentAuthor{}
	sqlStatement := `
		select id, author_type, name, image, description, created_at, updated_at, created_user_id, updated_user_id
		from content_author where id = $1
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
		var ca1 ContentAuthor
		rows.Scan(
			&ca1.ID,
			&ca1.AuthorType,
			&ca1.Name,
			&ca1.Image,
			&ca1.Description,
			&ca1.CreatedAt,
			&ca1.UpdatedAt,
			&ca1.CreatedUserID,
		)
		ca = append(ca, ca1)
	}
	resp.Data = ca
	resp.Message = "Successfully get author"
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}
