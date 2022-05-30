package users

import (
	"IPYP/database"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type User struct {
	ID            string         `json:"user_id"`
	Name          string         `json:"name"`
	Email         string         `json:"email"`
	AccountImgURL string         `json:"account_img_url"`
	Roles         map[int]string `json:"roles"`
}

func GetUser(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		user := User{
			ID: id,
		}

		row := db.Db.QueryRow(`SELECT name, email, google_picture FROM account
											WHERE user_id=$1`, id)
		err := row.Scan(&user.Name, &user.Email, &user.AccountImgURL)
		if err != nil {
			database.CheckDBErr(err.(*pq.Error), c)
			return
		}

		c.JSON(200, user)
	}
}

func GetAllUsers(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []User
		rows, err := db.Db.Query(`SELECT user_id, name, email, google_picture FROM account`)
		if err != nil {
			database.CheckDBErr(err.(*pq.Error), c)
			return
		}
		for rows.Next() {
			var user User
			err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.AccountImgURL)
			if err != nil {
				database.CheckDBErr(err.(*pq.Error), c)
				return
			}
			user.Roles, err = getUserRoles(db, user.ID)
			if err != nil {
				database.CheckDBErr(err.(*pq.Error), c)
				return
			}
			users = append(users, user)
		}

		c.JSON(200, users)
	}
}

func getUserRoles(db *database.DB, UserID string) (map[int]string, error) {
	roles := make(map[int]string)

	rows, err := db.Db.Query(`SELECT r.role_id, r.role_name FROM user_role_bridge 
    								INNER JOIN role r on r.role_id = user_role_bridge.role_id 
									WHERE user_uuid=$1;`, UserID)
	if err != nil {
		return roles, err
	}
	for rows.Next() {
		var roleID int
		var roleName string
		err = rows.Scan(&roleID, &roleName)
		if err != nil {
			return roles, err
		}

		roles[roleID] = roleName
	}

	return roles, err
}

func UpdateUser(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var user User
		err := c.BindJSON(&user)
		user.ID = id
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid request.")
			return
		}
		row := db.Db.QueryRow(`DELETE FROM user_role_bridge WHERE user_uuid=$1`, id)
		if row.Err() != nil {
			database.CheckDBErr(err.(*pq.Error), c)
			return
		}

		for role := range user.Roles {
			row = db.Db.QueryRow(`INSERT INTO user_role_bridge (user_uuid, role_id) VALUES ($1, $2)`, user.ID, role)
			if row.Err() != nil {
				database.CheckDBErr(err.(*pq.Error), c)
				return
			}
		}
		c.JSON(200, user)
	}
}

func DeleteUser(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		row := db.Db.QueryRow(`DELETE FROM account where user_id=$1 `, id)
		if row.Err() != nil {
			database.CheckDBErr(row.Err().(*pq.Error), c)
			return
		}

		c.JSON(200, nil)
	}
}
