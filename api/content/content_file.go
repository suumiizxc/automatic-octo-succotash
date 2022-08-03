package content

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/suumiizxc/raw_rest1/config"
	"github.com/suumiizxc/raw_rest1/response"
)

type ContentFile struct {
	ID            uint      `json:"id"`
	ContentID     *uint     `json:"content_id"`
	Attachment    *string   `json:"attachment"`
	FileType      *string   `json:"file_type"`
	FileUrl       *string   `json:"file_url"`
	SourceUrl     *string   `json:"source_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedUserID *uint     `json:"created_user_id"`
	UpdatedUserID *uint     `json:"updated_user_id"`
}

func (cf *ContentFile) CreateContentFile(w http.ResponseWriter, r *http.Request) {
	var cfi ContentFile
	resp := response.Response{}
	err := json.NewDecoder(r.Body).Decode(&cfi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err
		resp.Message = "Failed in request body"
		w.Write(resp.ConvertByte())
		return
	}
	var buf bytes.Buffer
	// in your case file would be fileupload
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])
	// Copy the file data to my buffer
	io.Copy(&buf, file)
	// do something with the contents...
	// I normally have a struct defined and unmarshal into a struct, but this will
	// work as an example
	contents := buf.String()
	fmt.Println(contents)
	// I reset the buffer in case I want to use it again
	// reduces memory allocations in more intense projects
	buf.Reset()

	sqlStatement := `
		insert into content_file 
		(
			content_id, attachment, file_type, file_url, source_url, created_at, created_user_id
		)
		values
		(
			$1, $2, $3, $4, $5, $6, $7
		) returning id
	`
	var lastID int
	err = config.DB.QueryRow(sqlStatement,
		cfi.ContentID, cfi.Attachment, cfi.FileType, cfi.FileUrl, cfi.SourceUrl, time.Now(), cfi.CreatedUserID,
	).Scan(&lastID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = err
		resp.Message = "Failed in query"
		w.Write(resp.ConvertByte())
		return
	}
	resp.Data = lastID
	resp.Message = "successfully created"
	w.WriteHeader(http.StatusCreated)
	w.Write(resp.ConvertByte())
}
