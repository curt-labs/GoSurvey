package prizes

import (
	"github.com/curt-labs/GoSurvey/models/prize"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
)

type PrizeError struct {
	Message string `json:"error"`
}

type PrizeResponse struct {
	Prizes       []prize.Prize `json:"prizes"`
	TotalPrizes  int           `json:"total_prizes"`
	CurrentPage  int           `json:"current_page"`
	TotalResults int           `json:"total_results"`
}

func All(rw http.ResponseWriter, req *http.Request, r render.Render) {
	params := req.URL.Query()
	var take int
	var page int
	var err error
	total := make(chan int, 0)

	go func() {
		total <- prize.PrizeCount()
	}()

	take, err = strconv.Atoi(params.Get("count"))
	page, err = strconv.Atoi(params.Get("page"))

	prizes, err := prize.All(page, take)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if page == 0 {
		page = 1
	}

	pr := PrizeResponse{
		Prizes:       prizes,
		CurrentPage:  page,
		TotalResults: len(prizes),
		TotalPrizes:  <-total,
	}

	r.JSON(200, pr)
}

func Get(rw http.ResponseWriter, req *http.Request, r render.Render, params martini.Params) {
	var p prize.Prize
	var err error

	if p.ID, err = strconv.Atoi(params["id"]); err != nil {
		r.JSON(500, PrizeError{err.Error()})
		return
	}

	if err := p.Get(); err != nil {
		r.JSON(500, PrizeError{err.Error()})
		return
	}

	r.JSON(200, p)
}

func Current(rw http.ResponseWriter, req *http.Request, r render.Render) {
	var p prize.Prize
	var err error

	p, err = prize.Current()
	if err != nil {
		r.JSON(500, PrizeError{err.Error()})
		return
	}

	if p.ID == 0 {
		r.JSON(500, PrizeError{"no prize is currently setup"})
		return
	}

	r.JSON(200, p)
}
