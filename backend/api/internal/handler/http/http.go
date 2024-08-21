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
	"github.com/jakobsym/BudgetFi/api/internal/controller/budgetfi"
	"github.com/jakobsym/BudgetFi/api/pkg/model"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Handler struct {
	ctrl *budgetfi.Controller
}

func New(ctrl *budgetfi.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

var env = loadOauthEnv()

// Google OAuth2 config
var OauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/login/auth",
	ClientID:     env["CLIENT_ID"],
	ClientSecret: env["CLIENT_SECRET"],
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	oauthStateString := genStateOauthCookie()
	url := OauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// when a user presses sign-in a POST request is sent with all of their OAuth credentials
func (h *Handler) OauthCallback(w http.ResponseWriter, r *http.Request) {
	var usr model.User

	// oauth code to obtain user info
	token, err := OauthConfig.Exchange(context.Background(), r.FormValue("code"))
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
		if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
			http.Error(w, "Invalid user data: "+err.Error(), http.StatusBadRequest) // 400 status code if error in request
			return
		}
		err := h.ctrl.CreateUser(r.Context(), &usr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("User creation server error: %v\n", err)
			return
		}
	}

	//TODO: Create a Session
}

// Util
// generates [16]byte uuid
func GenerateUUID() (uuid.UUID, error) {
	return uuid.Must(uuid.NewRandom()), nil
}

// Util
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
		"CLIENT_ID":     os.Getenv("CLIENT_ID"),
		"CLIENT_SECRET": os.Getenv("CLIENT_SECRET"),
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

	// TODO: Validate the user information
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
