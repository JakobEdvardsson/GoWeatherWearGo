package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/JakobEdvardsson/GoWeatherWearGo/types"
	"github.com/JakobEdvardsson/GoWeatherWearGo/util"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	DB *sql.DB
}

type PostgresSettings struct {
	pgHost     string
	pgPort     int
	pgUsername string
	pgPassword string
	pgDbname   string
}

var pgSettings *PostgresSettings

func init() {
	pgUsername := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PW")
	pgDbname := os.Getenv("POSTGRES_DB")

	if pgUsername == "" || pgPassword == "" || pgDbname == "" {
		// Load environment variables if they are not set
		err := util.LoadEnvFile(".env")
		if err != nil {
			log.Fatal("No env file or env vars provided!")
		}
		pgUsername = os.Getenv("POSTGRES_USER")
		pgPassword = os.Getenv("POSTGRES_PW")
		pgDbname = os.Getenv("POSTGRES_DB")
	}

	pgSettings = &PostgresSettings{
		pgHost:     "localhost",
		pgPort:     5432,
		pgUsername: pgUsername,
		pgPassword: pgPassword,
		pgDbname:   pgDbname,
	}
}

func NewPostgresStorage() *PostgresStorage {
	fmt.Println(*pgSettings)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		pgSettings.pgHost, pgSettings.pgPort, pgSettings.pgUsername, pgSettings.pgPassword, pgSettings.pgDbname)

	fmt.Println(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	return &PostgresStorage{
		DB: db,
	}
}

// Check if user exists in DB and return user, else return nil
func (s *PostgresStorage) GetUser(email string) (*types.User, error) {
	query := `SELECT * FROM "User" WHERE email = $1`
	var user types.User
	err := s.DB.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.EmailVerified, &user.Image)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *PostgresStorage) AddUser(spotifyProfile *types.SpotifyProfileResponse) (*types.User, error) {
	query := `INSERT INTO "User" (name, email, emailVerified, image) VALUES ($1, $2, $3, $4) RETURNING id, name, email, emailVerified, image;`
	var user types.User
	err := s.DB.QueryRow(query, spotifyProfile.DisplayName, spotifyProfile.Email, nil, nil).Scan(&user.ID, &user.Name, &user.Email, &user.EmailVerified, &user.Image)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
