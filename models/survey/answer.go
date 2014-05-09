package survey

import (
	"database/sql"
	"errors"
	"github.com/curt-labs/GoSurvey/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// SQL Statements
var (
	getQuestionAnswers = `select id, answer, data_type, date_modified,
												date_added, userID, deleted
												from SurveyAnswer
												where questionID = ? && deleted = 0
												order by date_added`
	getAnswerRevisions = `select ar.ID as revisionID, qr.new_answer,
													ar.old_answer, ar.date, ar.changeType,
													u.id as userID, u.fname, u.lname, u.username
													from SurveyAnswer_Revisions as ar
													join admin.user as u on ar.userID = u.id
													where ar.answerID = ?
													order by date desc`
	insertAnswer = `insert into SurveyAnswer(answer, data_type, date_added, userID, questionID)
									values(?,?,NOW(),?, ?)`
	updateAnswer = `update SurveyAnswer
									set answer = ?, data_type = ?, userID = ?, questionID = ?
									where id = ?`
	deleteAnswer = `update SurveyAnswer
									set deleted = 0, userID = ?
									where id = ?`
)

// Answer contains information for a possible answer
// for a given Question.
type Answer struct {
	ID           int              `json:"id"`
	Answer       string           `json:"answer"`
	DataType     string           `json:"data_type"`
	DateAdded    time.Time        `json:"date_added"`
	DateModified time.Time        `json:"date_modified"`
	UserID       int              `json:"-"`
	Deleted      bool             `json:"-"`
	Revisions    []AnswerRevision `json:"revisions"`
}

// AnswerRevision is a change record for an Answer.
type AnswerRevision struct {
	ID         int          `json:"id"`
	User       RevisionUser `json:"user"`
	NewAnswer  string       `json:"new_answer"`
	OldAnswer  string       `json:"old_answer"`
	Date       time.Time    `json:"date"`
	ChangeType string       `json:"change_type"`
}

// AddAnswer will commit a new Answer to
// a Question.
func (q *Question) AddAnswer(a Answer) error {

	if a.DataType != "input" && a.Answer == "" {
		return errors.New("only user input answers can be blank")
	}

	if a.ID == 0 {
		return a.insert(q.ID)
	}

	return a.update(q.ID)
}

// answers will push each answer for the
// given Question onto the Answers slice.
func (q *Question) answers() error {
	q.Answers = make([]Answer, 0)
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getQuestionAnswers)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(q.ID)
	if err != nil {
		return err
	}

	for res.Next() {
		var a Answer
		err = res.Scan(&a.ID, &a.Answer, &a.DataType, &a.DateModified, &a.DateAdded, &a.UserID, &a.Deleted)
		if err == nil {
			a.revisions()
			q.Answers = append(q.Answers, a)
		}
	}

	return nil
}

// insert will insert a new Answer and bind it
// to the given Question.
func (a *Answer) insert(questionID int) error {
	if questionID == 0 {
		return errors.New("invalid question reference")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(insertAnswer)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(a.Answer, a.DataType, a.UserID, questionID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = int(id)

	return nil
}

// udpate will update the answer, data_type, userID,
// and questionID properties for the given Answer.
func (a *Answer) update(questionID int) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(updateAnswer)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(a.Answer, a.DataType, a.UserID, questionID)

	return err
}

// Delete will mark the referenced Answer
// as deleted.
func (a *Answer) Delete() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(deleteAnswer)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(a.UserID, a.ID)
	return err
}

// revisions will retrieve all revisions for the
// referenced Answer and push them onto the Revisions
// slice.
func (a *Answer) revisions() error {
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAnswerRevisions)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(a.ID)
	if err != nil {
		return err
	}

	for res.Next() {
		var ar AnswerRevision
		err = res.Scan(&ar.ID, &ar.NewAnswer, &ar.OldAnswer, &ar.Date,
			&ar.ChangeType, &ar.User.ID, &ar.User.FirstName, &ar.User.LastName,
			&ar.User.Username)
		if err == nil {
			a.Revisions = append(a.Revisions, ar)
		}
	}
	return nil
}
