package store

import (
	"golang-app/internal/model"
)

func (r *Repository) GetUsernameById(id string) (*model.User, error) {
	user := &model.User{}
	var tempId []uint8
	if err := r.store.db.QueryRow(
		"SELECT id, phone, name FROM users WHERE id=$1",
		id,
	).Scan(
		&tempId,
		&user.Phone,
		&user.Name,
	); err != nil {
		return nil, err
	}
	user.Id = string(tempId)

	return user, nil
}
