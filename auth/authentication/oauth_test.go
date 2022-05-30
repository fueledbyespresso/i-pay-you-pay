package authentication

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"school-supply-list/database"
	"testing"
)

func createRouter() *gin.Engine {
	r := gin.Default()
	if _, err := os.Stat("../../projectvars.env"); err == nil {
		err := godotenv.Load("../../projectvars.env")
		if err != nil {
			fmt.Println("Error loading environment.env")
		}
		fmt.Println("Current environment:", os.Getenv("ENV"))
	}
	db := database.InitDBConnection()
	SStore := database.InitOauthStore()

	dbConnection := &database.DB{Db: db, SessionStore: SStore}
	r.GET("/login", handleGoogleLogin(dbConnection))
	r.GET("/callback", handleGoogleCallback(dbConnection))
	r.GET("/logout", handleGoogleLogout(dbConnection))
	r.GET("/account", getAccount(dbConnection))
	r.GET("/refresh", refreshSession(dbConnection))

	return r
}

func createSession(r *http.Request, w *httptest.ResponseRecorder, db *database.DB) {
	session, err := db.SessionStore.Get(r, "session")
	if err != nil {
		log.Fatal("Could not create session")
	}
	session.Values["GoogleId"] = "111644517051019423711"
	session.Values["Email"] = "leontyx@gmail.com"
	session.Values["Name"] = "leon"
	session.Values["Picture"] = "img.png"

	err = session.Save(r, w)

	if err != nil {
		fmt.Print("Unable to store session data")
	}
}

func TestAccount(t *testing.T) {
	r := createRouter()
	database.PerformMigrations("file://../../database/migrations")

	req, err := http.NewRequest("GET", "/account", nil)

	if err != nil {
		t.Fatalf(err.Error())
	}

	w := httptest.NewRecorder()

	createSession(req, w, &database.DB{SessionStore: database.InitOauthStore()})

	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.FailNow()
	}
	var account Account
	contents, _ := ioutil.ReadAll(w.Body)
	err = json.Unmarshal(contents, &account)

	fmt.Println("email: ", account.Email)
	fmt.Println("Name: ", account.Name)
	fmt.Println("Account picture URL: ", account.Picture)
	fmt.Println("UUID: ", account.ID)
}
