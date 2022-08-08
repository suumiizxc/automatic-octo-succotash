package content

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
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
	// var cfi ContentFile
	resp := response.Response{}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	// Create the uploads folder if it doesn't
	// already exist
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	file_name := time.Now().UnixNano()
	file_url := fmt.Sprintf("./uploads/%d%s", file_name, filepath.Ext(fileHeader.Filename))
	dst, err := os.Create(file_url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentID := r.FormValue("content_id")
	attachment := r.FormValue("attachment")
	file_type := strings.Replace(filepath.Ext(fileHeader.Filename), ".", "", 1)
	source_url := fmt.Sprintf(":/download/%d", file_name)
	created_user_id := r.FormValue("created_user_id")

	errs := []string{}
	if contentID == "" {
		errs = append(errs, "content_id required")
	}
	if attachment == "" {
		errs = append(errs, "attachment required")
	}
	if file_type == "" {
		errs = append(errs, "file_type required")
	}
	if source_url == "" {
		errs = append(errs, "source_url required")
	}
	if created_user_id == "" {
		errs = append(errs, "created_user_id required")
	}
	if len(errs) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = errs
		resp.Message = "Fill this fields"
		w.Write(resp.ConvertByte())
		return
	}

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
		contentID, attachment, file_type, file_url, source_url, time.Now(), created_user_id,
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

func (cf *ContentFile) GetContentFileByID(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	id := mux.Vars(r)["id"]
	cfi := []ContentFile{}
	sqlStatement := `
		select id, content_id, attachment, file_type, file_url, source_url, created_at, updated_at, created_user_id
		from content_file
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
		var cfi1 ContentFile
		rows.Scan(
			&cfi1.ID,
			&cfi1.ContentID,
			&cfi1.Attachment,
			&cfi1.FileType,
			&cfi1.FileUrl,
			&cfi1.SourceUrl,
			&cfi1.CreatedAt,
			&cfi1.UpdatedAt,
			&cfi1.CreatedUserID,
		)
		cfi = append(cfi, cfi1)
	}
	resp.Data = cfi
	resp.Message = "Successfully get content file by id"
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}

func (cf *ContentFile) GetContentFileByContentID(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	content_id := mux.Vars(r)["content_id"]
	cfi := []ContentFile{}
	sqlStatement := `
		select id, content_id, attachment, file_type, file_url, source_url, created_at, updated_at, created_user_id
		from content_file
		where content_id = $1
	`
	rows, err := config.DB.Query(sqlStatement, content_id)
	if err != nil {
		resp.Error = err
		resp.Message = "Failed in query"
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp.ConvertByte())
		return
	}
	for rows.Next() {
		var cfi1 ContentFile
		rows.Scan(
			&cfi1.ID,
			&cfi1.ContentID,
			&cfi1.Attachment,
			&cfi1.FileType,
			&cfi1.FileUrl,
			&cfi1.SourceUrl,
			&cfi1.CreatedAt,
			&cfi1.UpdatedAt,
			&cfi1.CreatedUserID,
		)
		cfi = append(cfi, cfi1)
	}
	resp.Data = cfi
	resp.Message = "Successfully get content file by content id"
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}
