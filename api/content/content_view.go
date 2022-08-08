package content

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/suumiizxc/raw_rest1/config"
	"github.com/suumiizxc/raw_rest1/response"
)

type ContentView struct {
	ID            uint       `json:"id"`
	Session       *string    `json:"session"`
	IpAddress     *string    `json:"ip_address"`
	ContentID     *uint      `json:"content_id"`
	CreatedAt     *time.Time `json:"created_at"`
	CreatedUserID *uint      `json:"created_user_id"`
}

func (c *ContentView) CreateContentView(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	cca := ContentView{}
	err := json.NewDecoder(r.Body).Decode(&cca)
	if err != nil {
		resp.Message = "Failed in request body"
		resp.Error = err
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp.ConvertByte())
		return
	}
	sqlStatement := `
	insert into content_view
	(
		session, ip_address, content_id, created_at, created_user_id
	)
	values
	(
		$1, $2, $3, $4, $5
	)
	returning id
	`
	var lastID int
	err = config.DB.QueryRow(sqlStatement, cca.Session, cca.IpAddress, cca.ContentID, time.Now(), cca.CreatedUserID).Scan(&lastID)
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

func (c *ContentView) GetContentViewList(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	cca := []ContentView{}
	sqlStatement := `
	select id, session, ip_address, content_id, created_at, created_user_id
	from content_view
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
		var cca1 ContentView
		rows.Scan(
			&cca1.ID,
			&cca1.Session,
			&cca1.IpAddress,
			&cca1.ContentID,
			&cca1.CreatedAt,
			&cca1.CreatedUserID,
		)
		cca = append(cca, cca1)
	}

	resp.Message = "Successfully get content view list"
	resp.Data = cca
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}
