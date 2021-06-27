package storage

import (
	"context"

	"sample-project/structs"
)

func (s *Storage) GetUserUUID(u *structs.User) (string, error) {
	var uuid string
	err := s.Pool.QueryRow(context.TODO(), `SELECT uuid FROM users WHERE username = $1 AND password = $2`,
		u.Username, u.Password).
		Scan(&uuid)
	return uuid, err
}
