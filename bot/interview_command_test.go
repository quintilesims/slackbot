package bot

import (
	"bytes"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/mock"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestInterviewAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockSlackClient(ctrl)
	client.EXPECT().
		AddReminder("", "uid", gomock.Any(), gomock.Any()).
		Return(nil).
		Times(3)

	store := newMemoryStore(t)
	w := bytes.NewBuffer(nil)
	cmd := NewInterviewCommand(client, store, "uid", w)

	if err := runTestApp(cmd, "!interview add \"John Doe\" 12/31"); err != nil {
		t.Fatal(err)
	}

	result := models.Interviews{}
	if err := store.Read(db.InterviewsKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Interviews{
		{
			ManagerID:   "uid",
			Interviewee: "John Doe",
			Date:        time.Date(0, 12, 31, 0, 0, 0, 0, time.Local).UTC(),
		},
	}

	assert.Equal(t, expected, result)
}
