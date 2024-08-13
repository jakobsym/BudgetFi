package httphandler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jakobsym/BudgetFi/api/internal/controller/budgetfi"
	"github.com/jakobsym/BudgetFi/api/pkg/model"
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

// Util
// generates [16]byte uuid
func GenerateUUID() (uuid.UUID, error) {
	return uuid.Must(uuid.NewRandom()), nil
}
