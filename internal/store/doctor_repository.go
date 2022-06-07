package store

import (
	"encoding/json"
	"golang-app/internal/model"
	"golang-app/internal/notifications"
	"sort"
	"time"
)

type Repository struct {
	store *Store
}

type SortSlotModel []model.SlotModel

func (e SortSlotModel) Less(i, j int) bool {
	timeLayout := "2006-01-02 15:04:05"
	ival, _ := time.Parse(timeLayout, e[i].Time)
	jval, _ := time.Parse(timeLayout, e[j].Time)
	return ival.Before(jval)
}

func (e SortSlotModel) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e SortSlotModel) Len() int { return len(e) }

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func (r *Repository) FindById(id string) (*model.Doctor, error) {
	doctor := &model.Doctor{}
	var tempSlot []uint8
	var tempId []uint8
	uid := []byte(id)
	if err := r.store.db.QueryRow(
		"SELECT id, name, spec, slots FROM doctors WHERE id=$1",
		uid,
	).Scan(
		&tempId,
		&doctor.Name,
		&doctor.Spec,
		&tempSlot,
	); err != nil {
		return nil, err
	}
	err := json.Unmarshal([]byte(tempSlot), &doctor.Slot)
	if err != nil {
		return nil, err
	}
	doctor.ID = string(tempId)
	return doctor, nil
}

func (r *Repository) AddSlot(doctorId string, newSlot model.SlotModel) (*model.Doctor, error) {
	doctor, err := r.FindById(doctorId)
	if err != nil {
		return nil, err
	}
	var temp []uint8
	f, err := r.CheckSlot(doctorId, newSlot)
	if err != nil {
		return nil, err
	}
	if f == false {
		return nil, err
	}
	doctor.Slot = append(doctor.Slot, newSlot)
	var sliceSort SortSlotModel = doctor.Slot
	sort.Sort(SortSlotModel(sliceSort))
	slots, err := json.Marshal(doctor.Slot)
	if err != nil {
		return nil, err
	}
	if err := r.store.db.QueryRow(
		"UPDATE doctors SET slots=$2WHERE id = $1 RETURNING slots",
		doctorId,
		slots,
	).Scan(
		&temp,
	); err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(temp), &doctor.Slot)
	if err != nil {
		return nil, err
	}
	return doctor, nil
}

func (r *Repository) CheckSlot(doctorId string, newSlot model.SlotModel) (bool, error) {
	doctorModel, err := r.FindById(doctorId)
	if err != nil {
		return false, err
	}
	if len(doctorModel.Slot) == 0 {
		return true, nil
	}
	timeLayout := "2006-01-02 15:04:05"
	timeVisit, err := time.Parse(timeLayout, newSlot.Time)
	if err != nil {
		return false, err
	}
	for i := range doctorModel.Slot {
		if i == 0 || i == len(doctorModel.Slot)-1 {
			if i == 0 {
				timeSub, err := time.Parse(timeLayout, doctorModel.Slot[i].Time)
				if err != nil {
					return false, err
				}
				timeSub = timeSub.Add(time.Minute * 20 * -1)
				if timeVisit.Before(timeSub) {
					return true, nil
				}
			}
			if i == len(doctorModel.Slot)-1 {
				timeAdd, err := time.Parse(timeLayout, doctorModel.Slot[i].Time)
				if err != nil {
					return false, err
				}
				timeAdd = timeAdd.Add(time.Minute * 20)
				if timeVisit.After(timeAdd) {
					return true, nil
				}
			}
		} else {
			timeSub, err := time.Parse(timeLayout, doctorModel.Slot[i].Time)
			if err != nil {
				return false, err
			}
			timeSub = timeSub.Add(time.Minute * 20 * -1)
			timeAdd, err := time.Parse(timeLayout, doctorModel.Slot[i-1].Time)
			if err != nil {
				return false, err
			}
			timeAdd = timeAdd.Add(time.Minute * 20)
			if timeVisit.After(timeAdd) && timeVisit.Before(timeSub) {

				return true, nil
			}
		}
	}
	return false, nil
}

func (r *Repository) Notification() error {
	timeLayout := "2006-01-02 15:04:05"
	rows, err := r.store.db.Query("SELECT id, name, spec, slots FROM doctors")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		doctor := &model.Doctor{}
		var temp []uint8
		var tempId []uint8

		err = rows.Scan(
			&tempId,
			&doctor.Name,
			&doctor.Spec,
			&temp,
		)
		if err != nil {
			return err
		}
		doctor.ID = string(tempId)
		err = json.Unmarshal([]byte(temp), &doctor.Slot)

		currentTime := time.Now()
		currentTimeStr := currentTime.Format(timeLayout)
		currentTime, err = time.Parse(timeLayout, currentTimeStr)
		if err != nil {
			return err
		}
		for i := range doctor.Slot {
			notifyTime, err := time.Parse(timeLayout, doctor.Slot[i].Time)
			if err != nil {
				return err
			}
			user, err := r.GetUsernameById(doctor.Slot[i].Id)
			if err != nil {
				return err
			}
			if currentTime.After(notifyTime.Add(-time.Hour*24)) && currentTime.Before(notifyTime.Add(-time.Hour*24+time.Minute*5)) {
				c := notifications.CreateNotification(user.Name, doctor.Spec, currentTime.Format(timeLayout))
				c.NotifyDay()
			}
			if currentTime.After(notifyTime.Add(-time.Hour*2)) && currentTime.Before(notifyTime.Add(-time.Hour*2+time.Minute*5)) {
				c := notifications.CreateNotification(user.Name, doctor.Spec, currentTime.Format(timeLayout))
				c.NotifyTwoHours()
			}
			if currentTime.After(notifyTime.Add(-time.Hour * 2)) {
				copy(doctor.Slot[i:], doctor.Slot[i+1:])
				copy(doctor.Slot[i:], doctor.Slot[i+1:])
				doctor.Slot[len(doctor.Slot)-1] = model.SlotModel{}
				doctor.Slot = doctor.Slot[:len(doctor.Slot)-1]

				r.UpdateSlot(doctor)
			}
		}
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateSlot(d *model.Doctor) (*model.Doctor, error) {
	slots, err := json.Marshal(d.Slot)
	if err != nil {
		return nil, err
	}
	var tempSlot []uint8
	uid := []byte(d.ID)
	if err := r.store.db.QueryRow(
		"UPDATE doctors SET slots=$2WHERE id = $1 RETURNING slots",
		uid,
		slots,
	).Scan(
		&tempSlot,
	); err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(tempSlot), &d.Slot)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Store) GetNewRepository() *Repository {
	s := Repository{}
	s.store = d
	return &s
}
