package sqlite

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Filename = ""
)

type Userdata struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

func openConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", Filename)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func exists(username string) int {
	username = strings.ToLower(username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	stmt := fmt.Sprintf(`SELECT ID FROM Users where Username = '%s'`, username)
	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Println("db query", err)
		return -1
	}
	defer rows.Close()

	userID := 1
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("exists() Scan", err)
			return -1
		}

		userID = id
	}

	return userID
}

func AddUser(d Userdata) int {
	d.Username = strings.ToLower(d.Username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := exists(d.Username)
	if userID != -1 {
		fmt.Println("User already exists:", d.Username)
		return -1
	}

	insertStmt := `INSERT INTO Users values (NULL,?)`
	_, err = db.Exec(insertStmt, d.Username)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	userID = exists(d.Username)
	if userID == -1 {
		return userID
	}

	insertStmt = `INSERT INTO Userdata values (?, ?, ?, ?)`
	_, err = db.Exec(insertStmt, userID, d.Name, d.Surname, d.Description)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}

	return userID
}

func DeleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt := fmt.Sprintf(`SELECT Username FROM Users WHERE ID = %d`, id)
	rows, err := db.Query(stmt)
	defer rows.Close()

	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return err
		}
	}

	if exists(username) != id {
		return fmt.Errorf("User with ID %d does not exist", id)
	}

	deleteStmt := `DELETE FROM Userdata WHERE UserID = ?`
	_, err = db.Exec(deleteStmt, id)
	if err != nil {
		return err
	}

	return nil
}

func ListUsers() ([]Userdata, error) {
	Data := []Userdata{}
	db, err := openConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT ID, Username, Name, Surname, Description FROM Users, Userdata WHERE Users.ID = Userdata.UserID`)
	defer rows.Close()
	if err != nil {
		return Data, err
	}

	for rows.Next() {
		var id int
		var username string
		var name string
		var surname string
		var desc string
		err = rows.Scan(&id, &username, &name, &surname, &desc)
		temp := Userdata{ID: id, Username: username, Name: name, Surname: surname, Description: desc}
		Data = append(Data, temp)
		if err != nil {
			return nil, err
		}
	}

	return Data, nil
}
