package survey

import (
	"database/sql"
	"errors"
	"time"

	"github.com/curt-labs/GoSurvey/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	insertUser = `insert into SurveyUser(fname, lname, email)
								values(?,?,?)`
	deleteUser       = `delete from SurveyUser where id = ?`
	insertUserAnswer = `insert into SurveyUserAnswer(userID, surveyID, questionID, answer)
									values(?,?,?,?)`
	getAllSubmissions = `select su.id, su.fname, su.lname, su.email, su.date_added,
												sua.answer, sq.question, sua.date_answered, sua.surveyID
												from SurveyUserAnswer sua
												join SurveyUser as su on sua.userID = su.id
												join SurveyQuestion as sq on sua.questionID = sq.id
												order by su.date_added desc`
	getAllSubmissionsBySurvey = `select su.id, su.fname, su.lname, su.email, su.date_added,
												sua.answer, sq.question, sua.date_answered, sua.surveyID
												from SurveyUserAnswer sua
												join SurveyUser as su on sua.userID = su.id
												join SurveyQuestion as sq on sua.questionID = sq.id
												where sua.surveyID = ?
												order by su.date_added desc
												limit ?,?`
	getSubmissionCount         = `select count(id) as count from SurveyUser`
	getSubmissionCountBySurvey = `select count(distinct su.id) as count from SurveyUser as su
																join SurveyUserAnswer as sua on su.id = sua.userID
																where sua.surveyID = ?`
	getSubmission = `select su.id, su.fname, su.lname, su.email, su.date_added,
												sua.answer, sq.question, sua.date_answered, sua.surveyID
												from SurveyUserAnswer sua
												join SurveyUser as su on sua.userID = su.id
												join SurveyQuestion as sq on sua.questionID = sq.id
												where su.id = ?`
)

const (
	TIME_LAYOUT = "2006-01-02 15:04:05 MST"
)

type SurveySubmission struct {
	ID        int                      `json:"id"`
	User      SurveyUser               `json:"user"`
	Questions []SurveySubmissionAnswer `json:"questions"`
	Survey    Survey                   `json:"survey"`
}

type SurveyUser struct {
	ID        int       `json:"id"`
	FirstName string    `json:"fname"`
	LastName  string    `json:"lname"`
	Email     string    `json:"email"`
	DateAdded time.Time `json:"date_added"`
}

type SurveySubmissionAnswer struct {
	ID           int       `json:"id"`
	Answer       string    `json:"answer"`
	Question     string    `json:"question"`
	DateAnswered time.Time `json:"date_answered"`
}

func GetAllSubmissions(skip, take, surveyID int) ([]SurveySubmission, error) {
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

	var stmt *sql.Stmt
	if surveyID == 0 {
		stmt, err = db.Prepare(getAllSubmissions)
	} else {
		stmt, err = db.Prepare(getAllSubmissionsBySurvey)
	}

	if err != nil {
		return submissions, err
	}
	defer stmt.Close()

	var res *sql.Rows
	if surveyID == 0 {
		res, err = stmt.Query()
	} else {
		res, err = stmt.Query(surveyID, skip, take)
	}
	if err != nil {
		return submissions, err
	}

	indexedSubmissions := make(map[int]SurveySubmission, 0)

	for res.Next() {
		var ans SurveySubmissionAnswer
		var user SurveyUser
		var s Survey

		if err = res.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateAdded, &ans.Answer, &ans.Question, &ans.DateAnswered, &s.ID); err == nil {
			if sm, _ := indexedSubmissions[user.ID]; sm.ID == 0 {

				ss := SurveySubmission{
					ID:   user.ID,
					User: user,
				}

				if err := s.Get(); err == nil {
					ss.Survey = s
				}
				indexedSubmissions[user.ID] = ss
			}

			sb := indexedSubmissions[user.ID]
			sb.Questions = append(sb.Questions, ans)
			indexedSubmissions[user.ID] = sb
		}
	}

	for _, sb := range indexedSubmissions {
		submissions = append(submissions, sb)
	}

	return submissions, nil
}

func SubmissionCount(surveyID int) int {
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return 0
	}

	defer db.Close()

	var stmt *sql.Stmt
	if surveyID == 0 {
		stmt, err = db.Prepare(getSubmissionCount)
	} else {
		stmt, err = db.Prepare(getSubmissionCountBySurvey)
	}

	if err != nil {
		return 0
	}

	defer stmt.Close()

	var total int
	if surveyID == 0 {
		stmt.QueryRow().Scan(&total)
	} else {
		stmt.QueryRow(surveyID).Scan(&total)
	}

	return total
}

func (s *SurveySubmission) Get() error {

	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getSubmission)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(s.ID)
	if err != nil {
		return err
	}

	for res.Next() {
		var ans SurveySubmissionAnswer
		var sur Survey

		if err = res.Scan(&s.User.ID, &s.User.FirstName, &s.User.LastName, &s.User.Email, &s.User.DateAdded, &ans.Answer, &ans.Question, &ans.DateAnswered, &sur.ID); err == nil {

			if sur.Name == "" {
				if err := sur.Get(); err == nil {
					s.Survey = sur
				}
			}

			s.Questions = append(s.Questions, ans)
		}
	}

	return nil
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
