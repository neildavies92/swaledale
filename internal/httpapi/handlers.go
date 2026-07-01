package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/neildavies/swaledale/internal/store"
)

type API struct {
	store         store.Store
	secureCookies bool
}

func (api *API) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (api *API) summary(w http.ResponseWriter, r *http.Request) {
	summary, err := api.store.Summary(r.Context(), userFrom(r.Context()).HouseholdID)
	respond(w, summary, err)
}

func (api *API) members(w http.ResponseWriter, r *http.Request) {
	members, err := api.store.Members(r.Context(), userFrom(r.Context()).HouseholdID)
	respond(w, members, err)
}

func (api *API) memberBudget(w http.ResponseWriter, r *http.Request) {
	memberID, ok := pathInt(w, r, "memberID")
	if !ok {
		return
	}
	budget, err := api.store.MemberBudget(r.Context(), userFrom(r.Context()).HouseholdID, memberID)
	respond(w, budget, err)
}

func (api *API) updateBudgetItem(w http.ResponseWriter, r *http.Request) {
	memberID, ok := pathInt(w, r, "memberID")
	if !ok {
		return
	}
	itemID, ok := pathInt(w, r, "itemID")
	if !ok {
		return
	}
	var input store.UpdateBudgetItemInput
	if !decode(w, r, &input) {
		return
	}
	budget, err := api.store.UpdateBudgetItem(r.Context(), userFrom(r.Context()).HouseholdID, memberID, itemID, input)
	respond(w, budget, err)
}

func (api *API) jointAccount(w http.ResponseWriter, r *http.Request) {
	joint, err := api.store.JointAccount(r.Context(), userFrom(r.Context()).HouseholdID)
	respond(w, joint, err)
}

func (api *API) updateJointAccountItem(w http.ResponseWriter, r *http.Request) {
	itemID, ok := pathInt(w, r, "itemID")
	if !ok {
		return
	}
	var input store.UpdateMoneyLabelInput
	if !decode(w, r, &input) {
		return
	}
	joint, err := api.store.UpdateJointAccountItem(r.Context(), userFrom(r.Context()).HouseholdID, itemID, input)
	respond(w, joint, err)
}

func (api *API) goals(w http.ResponseWriter, r *http.Request) {
	goals, err := api.store.Goals(r.Context(), userFrom(r.Context()).HouseholdID)
	respond(w, goals, err)
}

func (api *API) updateGoal(w http.ResponseWriter, r *http.Request) {
	goalID, ok := pathInt(w, r, "goalID")
	if !ok {
		return
	}
	var input store.UpdateGoalInput
	if !decode(w, r, &input) {
		return
	}
	goals, err := api.store.UpdateGoal(r.Context(), userFrom(r.Context()).HouseholdID, goalID, input)
	respond(w, goals, err)
}

func pathInt(w http.ResponseWriter, r *http.Request, key string) (int64, bool) {
	value, err := strconv.ParseInt(chi.URLParam(r, key), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid "+key)
		return 0, false
	}
	return value, true
}

func decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	return true
}

func respond(w http.ResponseWriter, payload any, err error) {
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "resource not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
