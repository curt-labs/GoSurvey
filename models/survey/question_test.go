package survey

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestQuestion(t *testing.T) {
	Convey("Testing question interaction with survey", t, func() {
		var s Survey
		s.Name = "Test Survey"
		s.Description = "Now that we know who you are, I know who I am. I'm not a mistake! It all makes sense! In a comic, you know how you can tell who the arch-villain's going to be? He's the exact opposite of the hero. And most times they're friends, like you and me! I should've known way back when... You know why, David? Because of the kids. They called me Mr Glass."
		s.UserID = 1
		err := s.Save()
		So(err, ShouldEqual, nil)

		var q Question
		Convey("Adding empty question", func() {
			q = Question{}
			q, err = s.AddQuestion(q)
			So(err, ShouldNotEqual, nil)
		})

		Convey("Updating question", func() {
			q.Question = "So, you think water moves fast?"
			q, err = s.AddQuestion(q)
			So(err, ShouldNotEqual, nil)

			q.UserID = 1
			q, err = s.AddQuestion(q)
			So(err, ShouldEqual, nil)
		})

		Convey("Deleting question", func() {
			err := q.Delete()
			So(err, ShouldEqual, nil)
		})

		Convey("Adding question with empty user reference", func() {
			q := Question{}
			q.Question = "What is the meaning of life?"
			q, err = s.AddQuestion(q)
			So(err, ShouldNotEqual, nil)

			Convey("Delete question that failed to be added", func() {
				err = q.Delete()
				So(err, ShouldNotEqual, nil)
			})
		})

		Convey("Adding question with user reference and no answers", func() {
			q := Question{}
			q.Question = "What is the meaning of life?"
			q.UserID = 1
			q, err = s.AddQuestion(q)
			So(err, ShouldEqual, nil)
		})

	})
}
