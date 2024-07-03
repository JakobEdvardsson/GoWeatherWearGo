package storage

import "github.com/JakobEdvardsson/GoWeatherWearGo/types"

type PostgresStorage struct{}

func NewPostgresStorage() *PostgresStorage {
	return &PostgresStorage{}
}

func (s *PostgresStorage) Get(id int) *types.User {
	return &types.User{
		ID:   1,
		Name: "Foo",
	}
}
