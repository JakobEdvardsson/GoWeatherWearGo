package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/JakobEdvardsson/GoWeatherWearGo/types"
	"github.com/JakobEdvardsson/GoWeatherWearGo/util"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
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
	userQuery := `INSERT INTO "User" (name, email, "emailVerified", image) VALUES ($1, $2, $3, $4) RETURNING id, name, email, "emailVerified", image;`
	accountQuery := `INSERT INTO "Account" ("userId", type, provider, "providerAccountId", refresh_token, access_token, expires_at, token_type, scope, id_token, session_state) VALUES ($1, $2, $3, $4, null, null, null, $5, $6, null, null);`

	var user types.User
	err := s.DB.QueryRow(userQuery, spotifyProfile.DisplayName, spotifyProfile.Email, nil, nil).Scan(&user.ID, &user.Name, &user.Email, &user.EmailVerified, &user.Image)
	if err != nil {
		fmt.Println("booom: ", err)
		return nil, err
	}
	fmt.Println("Added db.User")

	// insert into Account (id, userId, type, provider, providerAccountId, refresh_token, access_token, expires_at, token_type, scope, id_token, session_state) values ('1da5797a-429e-4c0c-86f6-8480461396af', '22f01aa9-99e2-477b-a001-9b21b659d77f', 'oauth', 'spotify', 'jakob.edvardsson', 'AQDQJAH7Qj4FYSjLWBNXXXn3lc31EnbhsMUvYWlcr2fEOt0tElmFw2txC3-6SOz7zlen-kkJII4jbl5rqy2bKZKKQqA6YJrgKWcHkxAA7uwI3o7LA5UM3YzrwG5-kpiptx15hg', 'BQDOoVM2GZ2sF50M7cWPKYw-vWYUeBItp1Yr-F2svEySaRQuFdS53kYGYECQV00O93LCklLxAlzsa8FXFXH0tWLH1PMA8l5MjqwQb1Mn4P7emSVE4MggQsQj1C-0RJhPvLRxKvb1_Bbh91ri5X6v8uy6Z4Uee6EB9NKPhCN76S9f8FMvAS7p4izc1ZLSJi0', 1710600239, 'bearer', 'user-read-email', null, null);
	res, err := s.DB.Exec(accountQuery, user.ID, "oauth", "spotify", spotifyProfile.ID, "bearer", "user-read-email")
	rowsAffected, _ := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return nil, err
	}
	fmt.Println("Added db.Account")

	return &user, nil
}

func (s *PostgresStorage) RefreshAccountSession(token *oauth2.Token, user *types.User) error {
	accountQuery := `UPDATE "Account" SET refresh_token = $1, access_token = $2, expires_at = $3 WHERE "userId" = $4`

	res, err := s.DB.Exec(accountQuery, token.RefreshToken, token.AccessToken, token.Expiry.Unix(), user.ID)
	if err != nil {
		fmt.Println("True: err != nil: ", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		fmt.Println("True: err != nil || rowsAffected == 0")
		return err
	}
	fmt.Println("Refreshed db.Account")

	return nil
}
