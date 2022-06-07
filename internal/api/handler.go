package api

import (
	"database/sql"
	"encoding/json"
	"golang-app/internal/model"
	"golang-app/internal/store"
	"golang-app/internal/validate"
	"net/http"
	"time"

	er "golang-app/internal/error"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type VisitDoctor struct {
	db    *sql.DB
	store store.Store
}

type Response struct {
	Message *string      `json:"message,omitempty"`
	Data    *interface{} `json:"data,omitempty"`
}

func NewVisitDoctor(database *sql.DB, store store.Store) *VisitDoctor {
	return &VisitDoctor{
		db:    database,
		store: store,
	}
}

func (l VisitDoctor) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := json.NewDecoder(r.Body)
	var logReq validate.VisitDoctorRequest

	err := req.Decode(&logReq)
	if err != nil {
		Err(w, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(logReq)
	if err != nil {
		Err(w, err)
		return
	}
	timeLayout := "2006-01-02 15:04:05"
	timeVal, err := time.Parse(timeLayout, logReq.Slot)
	if err != nil {
		Err(w, err)
		return
	}
	var mod model.SlotModel
	mod.Id = logReq.UserId
	mod.Time = timeVal.Format(timeLayout)

	f, err := l.store.Repository.AddSlot(logReq.DoctorId, mod)
	if err != nil {
		Err(w, err)
		return
	}

	Json(w, http.StatusOK, "Success", f)
}

func Err(w http.ResponseWriter, err error) {
	_, ok := err.(*er.RespError)
	if !ok {
		err = er.ValidationFailed(err.Error())
	}

	er, _ := err.(*er.RespError)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(er.Code)
	res := Response{
		Message: &er.Message,
	}
	json.NewEncoder(w).Encode(res)
}

func Json(w http.ResponseWriter, httpCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	res := Response{
		Message: &message,
		Data:    &data,
	}
	json.NewEncoder(w).Encode(res)
}
