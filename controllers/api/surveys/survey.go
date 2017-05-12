package surveys

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/curt-labs/GoSurvey/helpers/email"
	"github.com/curt-labs/GoSurvey/models/survey"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type SurveyError struct {
	Message string `json:"error"`
}

type SurveyResponse struct {
	Surveys      []survey.Survey `json:"surveys"`
	TotalSurveys int             `json:"total_surveys"`
	CurrentPage  int             `json:"current_page"`
	TotalResults int             `json:"total_results"`
}

func All(rw http.ResponseWriter, req *http.Request, r render.Render) {
	params := req.URL.Query()
	var take int
	var page int
	var err error
	total := make(chan int, 0)

	go func() {
		total <- survey.SurveyCount()
	}()

	take, err = strconv.Atoi(params.Get("count"))
	page, err = strconv.Atoi(params.Get("page"))

	svs, err := survey.GetSurveys(page, take)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if page == 0 {
		page = 1
	}

	sr := SurveyResponse{
		Surveys:      svs,
		CurrentPage:  page,
		TotalResults: len(svs),
		TotalSurveys: <-total,
	}

	r.JSON(200, sr)
}

func Get(rw http.ResponseWriter, req *http.Request, r render.Render, params martini.Params) {
	var sv survey.Survey
	var err error

	if sv.ID, err = strconv.Atoi(params["id"]); err != nil {
		r.JSON(500, SurveyError{err.Error()})
		return
	}

	if err := sv.Get(); err != nil {
		r.JSON(500, SurveyError{err.Error()})
		return
	}

	r.JSON(200, sv)
}

func Submit(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	dec := json.NewDecoder(req.Body)
	var s survey.SurveySubmission
	err := dec.Decode(&s)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.Submit()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	// send email to user filling out form.
	if s.User.Email != "" {
		tos := []string{s.User.Email}
		subject := "Warranty Confirmation"
		body := "<html><h3>Thank you for filling out your Warranty</h3>"
		body += "<p><strong>Provided below is a copy of your warranty questionnaire:</strong></p>"
		body += "<table><tr><th>Queston</th><th>Answer</th></tr>"
		for _, surv := range s.Questions {
			q := survey.Question{}
			q.ID = surv.ID
			err = q.GetRevisions()
			if err != nil {
				log.Println("error finding revision for question:", surv)
				continue
			}
			if len(q.Revisions) > 0 {
				body += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", q.Revisions[0].NewQuestion, surv.Answer)
			}
		}
		body += "</table>"

		err = email.Send(tos, subject, body, true)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	success := struct {
		Success bool `json:"success"`
	}{
		true,
	}

	js, _ := json.Marshal(success)

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}
