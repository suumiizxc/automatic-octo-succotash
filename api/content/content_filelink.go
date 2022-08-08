package content

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/suumiizxc/raw_rest1/config"
	"github.com/suumiizxc/raw_rest1/response"
)

type ContentFileLink struct {
	ID            uint       `json:"id"`
	Token         *string    `json:"token"`
	ExpireDate    *time.Time `json:"expire_date"`
	ContentFileID *uint      `json:"content_file_id"`
	CreatedUserID *uint      `json:"created_user_id"`
	CreatedDate   *time.Time `json:"created_date"`
}

func (c *ContentFileLink) CreateContentFileLink(w http.ResponseWriter, r *http.Request) {
	var cf ContentFileLink
	err := json.NewDecoder(r.Body).Decode(&cf)
	resp := response.Response{}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err
		resp.Message = "Failed in request query"
		w.Write(resp.ConvertByte())
		return
	}
	sqlStatement := `
	insert into content_file_link
	(
		token, expire_date, content_file_id, created_user_id,
		created_date
	)
	values
	(
		$1, $2, $3, $4, 
		$5
	)
	returning id
	`
	var lastID int
	err = config.DB.QueryRow(sqlStatement,
		cf.Token, cf.ExpireDate, cf.ContentFileID, cf.CreatedUserID, time.Now()).Scan(&lastID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = err
		resp.Message = "Failed in query"
		w.Write(resp.ConvertByte())
		return
	}
}
