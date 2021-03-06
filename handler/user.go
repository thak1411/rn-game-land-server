package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rn-game-land-server/usecase"
	"github.com/thak1411/rn-game-land-server/util"
)

type UserHandler struct {
	uc usecase.UserUsecase
}

type UserForm struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var body UserForm
		if err := util.BindBody(r.Body, &body); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		user := model.User{
			Role:     config.RoleBasic,
			Name:     body.Name,
			Username: body.Username,
			Password: body.Password,
		}
		nameRegex, _ := regexp.Compile("^[ㄱ-ㅎ|가-힣|ㅏ-ㅣ|a-z|A-Z|0-9|]{2,6}$")
		usernameRegex, _ := regexp.Compile("^[a-z|A-Z|0-9|]{6,12}$")
		passwordRegex, _ := regexp.Compile("^[a-z|A-Z|0-9|]{6,12}$")
		if !nameRegex.MatchString(user.Name) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if !usernameRegex.MatchString(user.Username) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if !passwordRegex.MatchString(user.Password) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if err := h.uc.CreateUser(user); err != nil {
			ret := model.RnHttpStatus{
				Status:  909,
				Message: err.Error(),
			}
			if err.Error() != "duplicated username" { // TODO: Remove err.Error() & check name, username before create user
				ret.Status = 908
			}
			if err := json.NewEncoder(w).Encode(ret); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		}
		ret := model.RnHttpStatus{
			Status:  910,
			Message: "Success Create User",
		}
		if err := json.NewEncoder(w).Encode(ret); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users, err := h.uc.GetAllUser()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "User List: %+v", users)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

type UserListResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (h *UserHandler) GetUserList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users, err := h.uc.GetAllUser()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		ret := []UserListResponse{}
		for _, v := range users {
			ret = append(ret, UserListResponse{
				Id:       v.Id,
				Name:     v.Name,
				Username: v.Username,
			})
		}
		if err := json.NewEncoder(w).Encode(ret); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var body LoginForm
		if err := util.BindBody(r.Body, &body); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		flag, err := h.uc.CheckUser(body.Username, body.Password)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if flag {
			user, err := h.uc.GetUser(body.Username)
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			token, err := util.CreateToken(user)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			cookie := &http.Cookie{
				Name:     config.Session,
				Domain:   config.Domain,
				Path:     "/",
				Expires:  time.Now().Add(time.Hour * 24 * 7),
				Value:    token,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
			w.WriteHeader(200)
			return
		}
		w.WriteHeader(401)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		iToken := r.Context().Value(config.Session)
		token := iToken.(model.AuthTokenClaims)
		var _ = token
		// TODO: expire user's token //

		cookie := &http.Cookie{
			Name:     config.Session,
			Value:    "",
			Path:     "/",
			Domain:   config.Domain,
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(200)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

type RetUserProfile struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		iToken := r.Context().Value(config.Session)
		token := iToken.(model.AuthTokenClaims)

		var user string
		username, ok := r.URL.Query()["username"]
		if !ok || len(username) < 1 {
			user = token.Username
		} else {
			user = username[0]
		}
		ret, err := h.uc.GetUser(user)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		userProfile := RetUserProfile{
			Id:       ret.Id,
			Name:     ret.Name,
			Username: ret.Username,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(userProfile); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

type FriendForm struct {
	Name string `json:"name"`
}

func (h *UserHandler) AddFriend(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var body FriendForm
		if err := util.BindBody(r.Body, &body); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		iToken := r.Context().Value(config.Session)
		token := iToken.(model.AuthTokenClaims)

		if body.Name == token.Name {
			http.Error(w, "Can't Add Self", http.StatusBadRequest)
			return
		}

		userId, err := h.uc.GetUserIdByName(body.Name)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if err := h.uc.AddFriend(token.Id, userId); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var body FriendForm
		if err := util.BindBody(r.Body, &body); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		iToken := r.Context().Value(config.Session)
		token := iToken.(model.AuthTokenClaims)

		if body.Name == token.Name {
			http.Error(w, "Can't Remove Self", http.StatusBadRequest)
			return
		}

		userId, err := h.uc.GetUserIdByName(body.Name)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if err := h.uc.RemoveFriend(token.Id, userId); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

type MyFriendForm struct {
	Name string `json:"name"`
}

type RetMyFriend struct {
	Res bool `json:"res"`
}

func (h *UserHandler) IsMyFriend(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var body MyFriendForm
		if err := util.BindBody(r.Body, &body); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		iToken := r.Context().Value(config.Session)
		token := iToken.(model.AuthTokenClaims)

		if body.Name == token.Name {
			http.Error(w, "Can't Check Self", http.StatusBadRequest)
			return
		}

		id, err := h.uc.GetUserIdByName(body.Name)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		res, err := h.uc.IsMyFriend(token.Id, id)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		var ret RetMyFriend
		ret.Res = res
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ret); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

type RetFriend struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (h *UserHandler) GetFriend(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		iToken := r.Context().Value(config.Session)
		token := iToken.(model.AuthTokenClaims)

		var name string
		qname, ok := r.URL.Query()["name"]
		if !ok || len(qname) < 1 {
			name = token.Username
		} else {
			name = qname[0]
		}

		id, err := h.uc.GetUserIdByName(name)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		res, err := h.uc.GetFriend(id)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		ret := []RetFriend{}
		for _, v := range res {
			ret = append(ret, RetFriend{
				Id:       v.Id,
				Name:     v.Name,
				Username: v.Username,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ret); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func NewUser(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc}
}
