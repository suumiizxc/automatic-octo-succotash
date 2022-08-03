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

type AccountProfile struct {
	ID          uint   `json:"id"`
	UserName    string `json:"user_name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	IsActive    uint   `json:"is_active"`
	AvatarImage string `json:"avatar_image"`
	Role        uint   `json:"role"`
	Token       string `json:"token"`
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

type RequestIDKey struct{}

func (a *Account) ProfileAccount(w http.ResponseWriter, r *http.Request) {

	var userID uint = r.Context().Value("user").(uint)
	fmt.Println("data : ", r.Context().Value("user"))

	resp := response.Response{}

	var ap []AccountProfile
	sqlStatement := `
		select id, user_name, first_name, last_name, email, is_active, avatar_image, role
		from account where id = $1
	`
	rows, err := config.DB.Query(sqlStatement, userID)

	if err != nil {
		resp.Error = err
		resp.Message = "Failed in query"
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp.ConvertByte())
		return
	}
	for rows.Next() {
		var ap1 AccountProfile
		rows.Scan(
			&ap1.ID,
			&ap1.UserName,
			&ap1.FirstName,
			&ap1.LastName,
			&ap1.Email,
			&ap1.IsActive,
			&ap1.AvatarImage,
			&ap1.Role,
		)
		ap = append(ap, ap1)
	}
	resp.Data = ap[0]
	resp.Message = "Successfully fetch profile"
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())
}

func (a *Account) LoginAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Println("data : ", r.Context().Value("data"))
	var input LoginUserNameInput
	resp := response.Response{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		resp.Error = err
		resp.Message = "Failed in request body"
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp.ConvertByte())
		return
	}
	var aa []Account
	var ap AccountProfile
	sqlStatement := `
		select id, user_name, password, first_name, last_name, email, is_active, avatar_image, role
		from account where user_name = $1
	`
	rows, err := config.DB.Query(sqlStatement, input.UserName)
	if err != nil {
		resp.Error = err
		resp.Message = "Failed in query"
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp.ConvertByte())
		return
	}
	for rows.Next() {
		var aa1 Account
		rows.Scan(
			&aa1.ID,
			&aa1.UserName,
			&aa1.Password,
			&aa1.FirstName,
			&aa1.LastName,
			&aa1.Email,
			&aa1.IsActive,
			&aa1.AvatarImage,
			&aa1.Role,
		)
		aa = append(aa, aa1)
	}
	fmt.Println("password : ", aa[0])

	if !CheckPasswordHash(input.Password, aa[0].Password) {
		resp.Error = err
		resp.Message = "Password did not match"
		w.WriteHeader(http.StatusNotImplemented)
		w.Write(resp.ConvertByte())
		return
	}
	ap.ID = aa[0].ID
	ap.UserName = aa[0].UserName
	ap.FirstName = aa[0].FirstName
	ap.LastName = aa[0].LastName
	ap.Email = aa[0].Email
	ap.IsActive = aa[0].IsActive
	ap.Role = aa[0].Role
	ap.AvatarImage = aa[0].AvatarImage
	ap.Token = "mongol"

	resp.Data = ap
	resp.Message = "Successfully logged"
	jsonAccount, _ := json.Marshal(ap)

	err = config.RS.Set(ap.Token, jsonAccount, 0).Err()
	if err != nil {
		fmt.Println("errorRedis : ", err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ConvertByte())

}
