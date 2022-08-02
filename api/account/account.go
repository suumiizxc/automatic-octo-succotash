package account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/suumiizxc/raw_rest1/config"
	"github.com/suumiizxc/raw_rest1/response"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID          uint      `json:"id"`
	Password    string    `json:"password"`
	UserName    string    `json:"user_name"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	IsActive    uint      `json:"is_active"`
	AvatarImage string    `json:"avatar_image"`
	Role        uint      `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LoginUserNameInput struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *Account) RegisterAccount(w http.ResponseWriter, r *http.Request) {
	//Test
	fmt.Println("data : ", r.Context().Value("data"))

	var aa Account
	err := json.NewDecoder(r.Body).Decode(&aa)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message" : "%s", "error" : "%s"}`, "Failed in request body", err.Error())))
		return
	}
	aa.Password, _ = HashPassword(aa.Password)
	aa.Role = 1
	aa.IsActive = 0
	currentTime := time.Now()
	aa.CreatedAt = currentTime
	sqlStatement := `
		insert into account (password, user_name, first_name, last_name, email, is_active, avatar_image, role, created_at)
		values ($1,$2,$3, $4,$5, $6,$7,$8,$9) returning id
	`

	var lastID int
	err = config.DB.QueryRow(sqlStatement, aa.Password, aa.UserName, aa.FirstName, aa.LastName, aa.Email, aa.IsActive, aa.AvatarImage, aa.Role, aa.CreatedAt).Scan(&lastID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message" : "%s", "error" : "%s"}`, "Failed in inserting", err.Error())))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"id" : "%d"}`, lastID)))
}

func (a *Account) LoginAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Println("data : ", r.Context().Value("data"))
	var input LoginUserNameInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message" : "%s", "error" : "%s"}`, "Failed in request body", err.Error())))
		return
	}
	var aa []Account
	sqlStatement := `
		select user_name, password from account where user_name = $1
	`
	rows, err := config.DB.Query(sqlStatement, input.UserName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message" : "%s", "error" : "%s"}`, "Failed internal error", err.Error())))
		return
	}
	for rows.Next() {
		var aa1 Account
		rows.Scan(&aa1.UserName, &aa1.Password)
		aa = append(aa, aa1)
	}
	fmt.Println("aa : ", aa)
	if !CheckPasswordHash(input.Password, aa[0].Password) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(fmt.Sprintf(`{"message" : "%s", "error": "%s"}`, "Password didn't match", err.Error())))
	}
	resp := response.Response{}
	resp.Data = aa
	resp.Message = "Successfully logged"

	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())

}
