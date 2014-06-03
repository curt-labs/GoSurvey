package survey

import (
	"database/sql"
	"errors"
	"github.com/curt-labs/GoSurvey/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	insertUser = `insert into SurveyUser(fname, lname, email)
								values(?,?,?)`
	deleteUser       = `delete from SurveyUser where id = ?`
	insertUserAnswer = `insert into SurveyUserAnswer(userID, surveyID, questionID, answer)
									values(?,?,?,?)`
	getAllSubmissions = ``
)

type SurveySubmission struct {
	ID        int                      `json:"id"`
	User      SurveyUser               `json:"user"`
	Questions []SurveySubmissionAnswer `json:"questions"`
}

type SurveyUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	Email     string `json:"email"`
}

type SurveySubmissionAnswer struct {
	ID     int    `json:"id"`
	Answer string `json:"answer"`
}

func GetAllSubmissions(skip, take int) ([]SurveySubmission, error) {
	if take == 0 {
		take = 25
	}

	if skip > 0 {
		skip = (skip - 1) * take
	}

	submissions := make([]SurveySubmission, 0)
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return submissions, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllSubmissions)
	if err != nil {
		return submissions, err
	}
	defer stmt.Close()

	res, err := stmt.Query(skip, take)
	if err != nil {
		return submissions, err
	}

	for res.Next() {
		var sm SurveySubmission
		if err = res.Scan(); err == nil {
			submissions = append(submissions, sm)
		}

	}

	return submissions, nil
}

func (s *SurveySubmission) Submit() error {

	if len(s.Questions) == 0 {
		return errors.New("cannot submit a survey without answers")
	}

	surv := Survey{
		ID: s.ID,
	}

	if err := surv.Get(); err != nil {
		return err
	}

	if err := s.User.save(); err != nil {
		return err
	}

	ch := make(chan error)

	for _, question := range s.Questions {
		go func(ss *SurveySubmission, q SurveySubmissionAnswer) {
			ch <- q.save(ss.User.ID, ss.ID)
		}(s, question)
	}

	for _, _ = range s.Questions {
		if err := <-ch; err != nil {
			s.User.delete()
			return err
		}
	}

	return nil
}

func (s *SurveyUser) save() error {
	if s.FirstName == "" {
		return errors.New("first name cannot be blank")
	}
	if s.LastName == "" {
		return errors.New("last name cannot be blank")
	}
	if s.Email == "" {
		return errors.New("e-mail name cannot be blank")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertUser)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(s.FirstName, s.LastName, s.Email)
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

func (s *SurveySubmissionAnswer) save(userID, surveyID int) error {
	if s.Answer == "" {
		return errors.New("answers cannot be blank")
	}
	if userID == 0 {
		return errors.New("failed to set name and/or email")
	}
	if s.ID == 0 {
		return errors.New("failed to assign question")
	}
	if surveyID == 0 {
		return errors.New("failed to assign survey")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertUserAnswer)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, surveyID, s.ID, s.Answer)
	return err
}

func (s *SurveyUser) delete() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteUser)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(s.ID)
	return err
}
