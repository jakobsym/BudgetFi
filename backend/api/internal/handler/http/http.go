package httphandler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jakobsym/BudgetFi/api/internal/controller/budgetfi"
	"github.com/jakobsym/BudgetFi/api/pkg/model"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type Handler struct {
	ctrl *budgetfi.Controller
}

func New(ctrl *budgetfi.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

var env = loadOauthEnv()
var store = sessions.NewCookieStore([]byte(env["SESSION_SECRET"]))

// Google OAuth2 config
var OauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth",
	ClientID:     env["CLIENT_ID"],
	ClientSecret: env["CLIENT_SECRET"],
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email", "openid"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://oauth2.googleapis.com/token",
	},
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	oauthStateString := genStateOauthCookie()
	url := OauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	//url := OauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// TODO: Condense code after testing
func (h *Handler) OauthCallback(w http.ResponseWriter, r *http.Request) {
	var usr model.User

	code := r.FormValue("code")
	token, err := OauthConfig.Exchange(context.Background(), code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error obtaining token via oauth: %v\n", err)
		return
	}
	client := OauthConfig.Client(context.Background(), token)
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("GET error on user info: %v\n", err)
		return
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&usr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error reading Oauth JSON body: %v\n", err)
	}

	// check if user is in DB via google_id
	usrUUID, err := h.ctrl.PrevUserCheck(r.Context(), &usr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error fetching from DB: %v\n", err)
		return
	}

	// New user -> Create new one
	if usrUUID == "" {
		defer r.Body.Close()
		usr.UUID, err = GenerateUUID()
		if err != nil {
			http.Error(w, "unable to gen UUID"+err.Error(), http.StatusInternalServerError)
			//log.Fatal("error generating uuid")
		}
		err := h.ctrl.CreateUser(r.Context(), &usr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("User creation server error: %v\n", err)
			return
		}
		usrUUID = string(usr.UUID[:])
	}

	// create a session
	session, err := store.Get(r, "session-name")
	session.Options.MaxAge = 86400 * 7
	session.Options.HttpOnly = true

	log.Println(session)
	if err != nil {
		http.Error(w, "unable to get session"+err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["user_id"] = usrUUID
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "unable to save session"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User logged-in"))
	// redirect to another route
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// fetch session id
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "unable to get session"+err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "unable to delete session"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User session deleted."))
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category model.Catergory
	// fetch UUID from session
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "unable to get session"+err.Error(), http.StatusInternalServerError)
		return
	}
	userID := session.Values["user_id"].(string)

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "error decoding category"+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = h.ctrl.CreateCategory(r.Context(), &category, userID)
	if err != nil {
		http.Error(w, "unable to get session"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Category Created for user: "))
	http.Redirect(w, r, "/dashbard", http.StatusSeeOther)
}

func (h *Handler) CreateExpense(w http.ResponseWriter, r *http.Request) {}

// TODO: make session creation a function?

/* Util functions */

// generates [16]byte uuid
func GenerateUUID() (uuid.UUID, error) {
	return uuid.Must(uuid.NewRandom()), nil
}

// generates oauth state cookie
func genStateOauthCookie() string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	return state
}

func loadOauthEnv() map[string]string {
	err := godotenv.Load("backend.env")
	if err != nil {
		log.Printf("unable to load oauth .env values: %v\n", err)
	}
	envMap := map[string]string{
		"CLIENT_ID":      os.Getenv("CLIENT_ID"),
		"CLIENT_SECRET":  os.Getenv("CLIENT_SECRET"),
		"SESSION_SECRET": os.Getenv("SESSION_SECRET"),
	}
	return envMap
}

/*
Deprecated
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var usr model.User
	defer r.Body.Close()
	var err error
	usr.UUID, err = GenerateUUID()
	if err != nil {
		http.Error(w, "unable to gen UUID"+err.Error(), http.StatusInternalServerError)
		//log.Fatal("error generating uuid")
	}
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		http.Error(w, "Invalid user data: "+err.Error(), http.StatusBadRequest) // 400 status code if error in request
		return
	}

	// Validate the user information
	// https://dev.to/wati_fe/how-i-validate-user-input-in-golang-c5f
	err = h.ctrl.CreateUser(r.Context(), &usr)
	//err = h.ctrl.Post(r.Context(), usr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("User creation server error: %v\n", err)
		return
	}
}
*/
