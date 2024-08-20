package httphandler

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jakobsym/BudgetFi/api/internal/controller/budgetfi"
	"github.com/jakobsym/BudgetFi/api/pkg/model"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Handler struct {
	ctrl *budgetfi.Controller
}

func New(ctrl *budgetfi.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

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

// Google OAuth2 config
var OauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/login/auth",
	ClientID:     "YOUR_GOOGLE_CLIENT_ID",
	ClientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
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
	//var err error

	//TODO: OAuth Login to obtain user information such as name, email, google_id

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
