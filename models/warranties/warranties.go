package warranties

import (
	"errors"
	"time"

	"github.com/curt-labs/GoSurvey/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/martini-contrib/binding"
)

var (
	getAllWarranties = `select id, fname, lname, email, part, date_added from ActivatedWarranties`
	getWarranty      = `select id, fname, lname, email, part, date_added from ActivatedWarranties
											where id = ?
											limit 1`
	addWarranty = `insert into ActivatedWarranties(fname, lname, email, part)
									values(?,?,?,?)`
	deleteWarranty  = `delete from ActivatedWarranties where id = ?`
	checkPartNumber = `select partID from Part where partID = ? && brandID = 1 limit 1`
)

type Warranty struct {
	ID        int
	FirstName string `json:"fname" form:"fname"`
	LastName  string `json:"lname" form:"lname"`
	Email     string `json:"email" form:"email"`
	Part      int    `json:"part" form:"part"`
	Added     time.Time
}

// All retrieves all listed activated
// warranties from the database.
func All() ([]Warranty, error) {
	ws := make([]Warranty, 0)
	var err error

	err = database.Init()
	if err != nil {
		return ws, err
	}

	stmt, err := database.DB.Prepare(getAllWarranties)
	if err != nil {
		return ws, err
	}

	defer stmt.Close()

	res, err := stmt.Query()
	if err != nil {
		return ws, err
	}

	for res.Next() {
		var w Warranty
		err = res.Scan(&w.ID, &w.FirstName, &w.LastName, &w.Email, &w.Part, &w.Added)
		if err == nil {
			ws = append(ws, w)
		}
	}

	return ws, err
}

// Get returns a warranty based off
// the supplied ID.
func (w *Warranty) Get() error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getWarranty)
	if err != nil {
		return err
	}

	defer stmt.Close()

	id := w.ID
	w.ID = 0
	stmt.QueryRow(id).Scan(&w.ID, &w.FirstName, &w.LastName, &w.Email, &w.Part, &w.Added)
	if w.ID == 0 {
		return errors.New("no warranty found")
	}

	return nil
}

// Add inserts a new warranty.
func (w *Warranty) Add() error {

	if w.FirstName == "" {
		return errors.New("invalid first name")
	}
	if w.LastName == "" {
		return errors.New("invalid last name")
	}
	if w.Email == "" {
		return errors.New("invalid email address")
	}

	err := database.Init()
	if err != nil {
		return err
	}

	check, err := database.DB.Prepare(checkPartNumber)
	if err != nil {
		return err
	}
	defer check.Close()

	var part int
	check.QueryRow(w.Part).Scan(&part)

	if part == 0 {
		w.Part = 0
		return errors.New("invalid part number")
	}

	stmt, err := database.DB.Prepare(addWarranty)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(w.FirstName, w.LastName, w.Email, w.Part)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	w.ID = int(id)
	w.Added = time.Now()

	return nil
}

func (w *Warranty) Delete() error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(deleteWarranty)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(w.ID)
	if err != nil {
		return err
	}

	if aff, err := res.RowsAffected(); err != nil || aff == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
