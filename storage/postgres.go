package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/JakobEdvardsson/GoWeatherWearGo/types"
	"github.com/JakobEdvardsson/GoWeatherWearGo/util"
	"github.com/google/uuid"
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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		pgSettings.pgHost, pgSettings.pgPort, pgSettings.pgUsername, pgSettings.pgPassword, pgSettings.pgDbname)

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
	// TODO: Change * to explicit select
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
		return nil, err
	}

	res, err := s.DB.Exec(accountQuery, user.ID, "oauth", "spotify", spotifyProfile.ID, "bearer", "user-read-email")
	rowsAffected, _ := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return nil, err
	}

	return &user, nil
}

func (s *PostgresStorage) UpdateSpotifySession(refreshToken string, accessToken string, expiry time.Time, userId string) error {
	accountQuery := `UPDATE "Account" SET refresh_token = $1, access_token = $2, expires_at = $3 WHERE "userId" = $4`

	res, err := s.DB.Exec(accountQuery, refreshToken, accessToken, expiry.Unix(), userId)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}

	return nil
}

// TODO implement creation of user sessions in DB
func (s *PostgresStorage) CreateUserSession(token *oauth2.Token, user *types.User) (session *types.Session, err error) {
	accountQuery := `INSERT INTO "Session" ("userId", "sessionToken", "expires") VALUES ($1, $2, $3) RETURNING id, "userId", "sessionToken", "expires";`
	session = &types.Session{}

	sessionToken := uuid.New().String()
	expires := time.Now().Add(time.Hour * 24).UTC()

	err = s.DB.QueryRow(accountQuery, user.ID, sessionToken, expires).Scan(&session.ID, &session.UserID, &session.SessionToken, &session.Expires)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *PostgresStorage) GetUserSession(sessionToken string) (session *types.Session, err error) {
	query := `SELECT "id", "userId", "sessionToken", "expires" FROM "Session" WHERE "sessionToken" = $1;`
	session = &types.Session{}

	err = s.DB.QueryRow(query, sessionToken).Scan(&session.ID, &session.UserID, &session.SessionToken, &session.Expires)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *PostgresStorage) GetAccount(userID string) (account *types.Account, err error) {
	query := `SELECT id, "userId", type, provider, "providerAccountId", refresh_token, access_token, expires_at, token_type, scope, id_token, session_state FROM "Account" WHERE "userId" = $1;`
	account = &types.Account{}

	err = s.DB.QueryRow(query, userID).Scan(
		&account.ID,
		&account.UserID,
		&account.Type,
		&account.Provider,
		&account.ProviderAccountID,
		&account.RefreshToken,
		&account.AccessToken,
		&account.ExpiresAt,
		&account.TokenType,
		&account.Scope,
		&account.IDToken,
		&account.SessionState)

	if err != nil {
		return nil, err
	}

	return account, nil
}
