package validate

type (
	VisitDoctorRequest struct {
		UserId   string `json:"user_id" valiudate:"required,uuid"`
		DoctorId string `json:"doctor_id" validate:"required,uuid"`
		Slot     string `json:"slot" validate:"required"`
	}
)
