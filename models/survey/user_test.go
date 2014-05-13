package survey

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSubmission(t *testing.T) {
	Convey("Test survey submission", t, func() {

		// Insert survey
		var s Survey
		s.Name = "Test Survey"
		s.Description = "Now that we know who you are, I know who I am. I'm not a mistake! It all makes sense! In a comic, you know how you can tell who the arch-villain's going to be? He's the exact opposite of the hero. And most times they're friends, like you and me! I should've known way back when... You know why, David? Because of the kids. They called me Mr Glass."
		s.UserID = 1
		err := s.Save()
		So(err, ShouldEqual, nil)

		// Insert question and answer
		q := Question{
			Question: "You think water moves fast?",
			UserID:   1,
		}

		q.Answers = []Answer{
			Answer{
				Answer:   "",
				UserID:   1,
				DataType: "input",
			},
		}
		q, err = s.AddQuestion(q)
		So(err, ShouldEqual, nil)
		So(len(q.Answers), ShouldNotEqual, 0)
	})
}
