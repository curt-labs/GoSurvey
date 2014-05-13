package survey

import (
	"database/sql"
	"errors"
	"github.com/curt-labs/GoSurvey/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	getAllSurveys = `select id, name, description,
										date_added, date_modifed, userID, deleted
										from Survey
										where deleted = 0
										order by date_modifed desc
										limit ?,?`
	getSurveyCount = `select count(id) as count from Survey where deleted = 0`
	getSurvey      = `select id, name, description,
										date_added, date_modifed, userID, deleted
										from Survey
										where id = ? && deleted = 0 limit 1`
	getSurveyRevisions = `select
												sv.ID as revisionID, IFNULL(sv.new_name, ""), IFNULL(sv.old_name, ""),
												sv.date, sv.changeType,
												u.id as userID, u.fname, u.lname, u.username
												from Survey_Revisions as sv
												join admin.user as u on sv.userID = u.id
												where sv.surveyID = ?
												order by date desc`
	insertSurvey = `insert into Survey(name, description, date_added, userID)
									values(?,?,NOW(), ?)`
	updateSurvey = `update Survey set name = ?, description = ?, userID = ?
									where id = ?`
	deleteSurvey = `update Survey set deleted = 1, userID = ? where id = ?`
)

type Survey struct {
	ID           int              `json:"id"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	DateAdded    time.Time        `json:"date_added"`
	DateModified time.Time        `json:"date_modified"`
	UserID       int              `json:"-"`
	Deleted      bool             `json:"-"`
	Revisions    []SurveyRevision `json:"revisions"`
	Questions    []Question       `json:"questions"`
	Completion   SurveyStatus     `json:"-"`
}

type SurveyStatus struct {
	Completed     bool `json:"completed"`
	QuestionCount int  `json:"question_count"`
}

type SurveyRevision struct {
	ID         int          `json:"id"`
	NewName    *string      `json:"new_name"`
	OldName    *string      `json:"old_name"`
	Date       time.Time    `json:"date"`
	ChangeType string       `json:"change_type"`
	User       RevisionUser `json:"user"`
}

type RevisionUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// GetSurveys will return a list of surveys in
// the database or an error if empty.
func GetSurveys(skip, take int) ([]Survey, error) {
	if take == 0 {
		take = 25
	}

	svs := make([]Survey, 0)
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return svs, err
	}

	defer db.Close()

	stmt, err := db.Prepare(getAllSurveys)
	if err != nil {
		return svs, err
	}

	defer stmt.Close()

	res, err := stmt.Query(skip, take)
	if err != nil {
		return svs, err
	}

	for res.Next() {
		var sv Survey
		err = res.Scan(&sv.ID, &sv.Name, &sv.Description, &sv.DateAdded, &sv.DateModified, &sv.UserID, &sv.Deleted)
		if err == nil {
			sv.revisions()
			sv.questions()
			svs = append(svs, sv)
		}
	}

	return svs, err
}

// SurveyCount will return the total number of
// surveys in the database that aren't marked as
// deleted.
func SurveyCount() int {
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return 0
	}

	defer db.Close()

	stmt, err := db.Prepare(getSurveyCount)
	if err != nil {
		return 0
	}

	defer stmt.Close()

	var total int
	stmt.QueryRow().Scan(&total)

	return total
}

// Get will update the values on the bound
// Survey or return an error.
func (s *Survey) Get() error {
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}

	defer db.Close()

	stmt, err := db.Prepare(getSurvey)
	if err != nil {
		return err
	}

	defer stmt.Close()

	id := s.ID
	s.ID = 0
	stmt.QueryRow(id).Scan(&s.ID, &s.Name, &s.Description, &s.DateAdded, &s.DateModified, &s.UserID, &s.Deleted)
	if s.ID == 0 {
		return errors.New("no survey found")
	}

	s.revisions()
	s.questions()

	return nil
}

// revisions will assign revision
// history to the bound Survey.
func (s *Survey) revisions() error {
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}

	defer db.Close()

	stmt, err := db.Prepare(getSurveyRevisions)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Query(s.ID)
	if err != nil {
		return err
	}

	for res.Next() {
		var r SurveyRevision
		err = res.Scan(&r.ID, &r.NewName, &r.OldName, &r.Date, &r.ChangeType,
			&r.User.ID, &r.User.FirstName, &r.User.LastName,
			&r.User.Username)
		if err == nil {
			s.Revisions = append(s.Revisions, r)
		}
	}

	return nil
}

// Add will commit the current Survey
// to the database.
func (s *Survey) Save() error {

	if s.Name == "" {
		return errors.New("survey name cannot be blank")
	}

	if s.UserID == 0 {
		return errors.New("invalid user reference")
	}

	if s.ID == 0 {
		return s.insert()
	}
	return s.update()
}

// insert will insert a new survey record.
func (s *Survey) insert() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertSurvey)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(s.Name, s.Description, s.UserID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	s.ID = int(id)

	return nil
}

// update will update an existing survey
// record.
func (s *Survey) update() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(updateSurvey)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(s.Name, s.Description, s.UserID, s.ID)

	return err
}

// Delete will remove (mark as deleted) a Survey
// from the list of returned results
func (s *Survey) Delete() error {
	if s.ID == 0 {
		return errors.New("invalid survey record")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteSurvey)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(s.UserID, s.ID)

	return err
}
