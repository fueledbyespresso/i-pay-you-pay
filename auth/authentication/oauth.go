package authentication

import (
	"IPYP/database"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	GoogleOauthConfig *oauth2.Config
)

func ConfigOauth() {
	if os.Getenv("ENV") == "DEV" {
		GoogleOauthConfig = &oauth2.Config{
			RedirectURL:  "http://localhost:8000/oauth/v1/callback",
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		}
	} else {
		GoogleOauthConfig = &oauth2.Config{
			RedirectURL:  "https://" + os.Getenv("HOST") + "/oauth/v1/callback",
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		}
	}
}

type Error struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_msg"`
}

//All the routes created by the package nested in
// oauth/v1/*
func Routes(r *gin.RouterGroup, db *database.DB) {
	r.GET("/login", handleGoogleLogin(db))
	r.GET("/callback", handleGoogleCallback(db))
	r.GET("/logout", handleGoogleLogout(db))
	r.GET("/account", getAccount(db))
	r.GET("/refresh", refreshSession(db))
}

func getSeed() int64 {
	seed := time.Now().UnixNano() // A new random seed (independent from state)
	rand.Seed(seed)
	return seed
}

func handleGoogleLogin(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		state, err := db.SessionStore.Get(c.Request, "state")
		if err != nil {
			c.AbortWithStatusJSON(500, "Server was unable to connect to session database")
			return
		}

		stateString := strconv.FormatInt(getSeed(), 10)
		state.Values["state"] = stateString
		err = state.Save(c.Request, c.Writer)

		if err != nil {
			print("Unable to store state data")
			c.AbortWithStatusJSON(500, "Unable to store state data")
		}

		url := GoogleOauthConfig.AuthCodeURL(stateString)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}

}

func handleGoogleCallback(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		stateSession, err := db.SessionStore.Get(c.Request, "state")
		if err != nil {
			c.AbortWithStatusJSON(500, "The server was unable to retrieve session state")
			return
		}
		state := fmt.Sprintf("%v", stateSession.Values["state"])
		userData, err := getUserInfo(state, c.Request.FormValue("code"), c.Request)
		if err != nil {
			fmt.Println("Error getting content: " + err.Error())
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		stateSession.Options.MaxAge = -1
		_ = stateSession.Save(c.Request, c.Writer)
		// Add a user to user database if they don't exist
		// otherwise replace the previous access token field
		// with the new one

		if !userExists(userData.Email, db) {
			err = createUser(userData, db)
			if err != nil {
				database.CheckDBErr(err.(*pq.Error), c)
				return
			}
		} else {
			replaceAccessToken(userData, db)
		}

		// set the user information
		session, err := db.SessionStore.Get(c.Request, "session")
		if err != nil {
			c.AbortWithStatusJSON(500, "Server was unable to connect to session database")
		}

		session.Values["GoogleId"] = userData.GoogleId
		session.Values["Email"] = userData.Email
		session.Values["Name"] = userData.Name
		session.Values["Picture"] = userData.Picture

		err = session.Save(c.Request, c.Writer)
		if err != nil {
			fmt.Print("Unable to store session data")
			c.AbortWithStatusJSON(500, "Unable to store session data")
		}

		c.Redirect(http.StatusPermanentRedirect, "/")
	}
}

func createUser(userData User, db *database.DB) error {
	// Prepare the sql query for later
	insert, err := db.Db.Prepare(`INSERT INTO account (email, access_token, google_id, expires_in, google_picture, name) VALUES ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return err
	}

	//Execute the previous sql query using data from the
	// userData struct being passed into the function
	_, err = insert.Exec(userData.Email, userData.AccessToken, userData.GoogleId, userData.ExpiresIn, userData.Picture, userData.Name)

	if err != nil {
		return err
	}
	return nil
}

func userExists(email string, db *database.DB) bool {
	// Prepare the sql query for later
	rows, err := db.Db.Query("SELECT COUNT(*) as count FROM account WHERE email = $1", email)
	PanicOnErr(err)

	return checkCount(rows) > 0
}

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		PanicOnErr(err)
	}
	return count
}

func replaceAccessToken(userData User, db *database.DB) {
	_, err := db.Db.Query("UPDATE account SET access_token=$1, expires_in=$2, google_picture=$3, name=$4 WHERE email = $5",
		userData.AccessToken, userData.ExpiresIn, userData.Picture, userData.Name, userData.Email)
	if err != nil {
		fmt.Println("Unable to update access token", err)
	}
}

type User struct {
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Picture     string    `json:"picture"`
	GoogleId    string    `json:"id"`
	ExpiresIn   time.Time `json:"expires_in"`
	AccessToken string
}

func getUserInfo(state string, code string, r *http.Request) (User, error) {
	var userData User
	if state != r.FormValue("state") {
		return userData, fmt.Errorf("invalid oauth state")
	}

	token, err := GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return userData, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	//Send access token to google's user api in return for a users data!
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	if err != nil {
		return userData, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return userData, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	err = json.Unmarshal(contents, &userData)
	if err != nil {
		log.Println(err)
	}

	userData.ExpiresIn = token.Expiry
	userData.AccessToken = token.AccessToken

	return userData, nil
}

func handleGoogleLogout(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Attempting to expire session")

		session, err := db.SessionStore.Get(c.Request, "session")
		if err != nil {
			c.AbortWithStatusJSON(500, "The server was unable to retrieve this session")
			return
		}

		if session.ID != "" {
			session.Options.MaxAge = -1

			err = session.Save(c.Request, c.Writer)

			if err != nil {
				c.AbortWithStatusJSON(500, "The server was unable to expire this session")
			} else {
				c.JSON(200, `{"successful logout"}`)
			}

		} else {
			c.Redirect(http.StatusTemporaryRedirect, "./")
		}
	}
}

type Account struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	ID      string `json:"user_id"`
}

func refreshSession(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := db.SessionStore.Get(c.Request, "session")
		if err != nil {
			c.AbortWithStatusJSON(500, "The server was unable to retrieve this session")
			return
		}

		if session.ID != "" {
			session.Options.MaxAge = 3600

			err = session.Save(c.Request, c.Writer)
			if err != nil {
				c.AbortWithStatusJSON(500, "The server was unable to refresh this session")
			} else {
				c.JSON(200, "successful refresh")
			}
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "./login")
		}
	}
}

func getAccount(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := db.SessionStore.Get(c.Request, "session")
		if err != nil {
			c.AbortWithStatusJSON(500, "The server was unable to retrieve this session")
			return
		}

		if session.ID != "" {
			// get some session values
			Email := session.Values["Email"]
			EmailStr := fmt.Sprintf("%v", Email)
			Name := session.Values["Name"]
			NameStr := fmt.Sprintf("%v", Name)
			PictureUrl := session.Values["Picture"]
			PictureUrlStr := fmt.Sprintf("%v", PictureUrl)
			GoogleID := session.Values["GoogleId"]
			GoogleIDStr := fmt.Sprintf("%v", GoogleID)
			if err != nil {
				database.CheckDBErr(err.(*pq.Error), c)
				return
			}
			userID, err := getUUIDFromGoogleID(db, GoogleIDStr)
			if err != nil {
				database.CheckDBErr(err.(*pq.Error), c)
				return
			}

			userData := Account{EmailStr, NameStr, PictureUrlStr, userID}

			c.JSON(200, userData)
		} else {
			c.AbortWithStatusJSON(401, "Session not found. Session may be expired or non-existent")
		}
	}
}

func getUUIDFromGoogleID(db *database.DB, googleID string) (string, error) {
	var userID string
	userRoles, err := db.Db.Query(`SELECT user_id from account a where a.google_id=$1`, googleID)
	if err != nil {
		return userID, err
	}

	for userRoles.Next() {
		err = userRoles.Scan(&userID)
		if err != nil {
			return userID, err
		}
	}

	return userID, nil
}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
