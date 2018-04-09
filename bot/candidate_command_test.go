package bot

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestCandidateAdd(t *testing.T) {
	store := newMemoryStore(t)
	cmd := NewCandidateCommand(store, ioutil.Discard)

	if err := runTestApp(cmd, "!candidate add --meta k1=v1 --meta k2=v2 \"John Doe\" <@uid>"); err != nil {
		t.Fatal(err)
	}

	result := models.Candidates{}
	if err := store.Read(db.CandidatesKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Candidates{
		{
			Name:      "John Doe",
			ManagerID: "uid",
			Meta:      map[string]string{"k1": "v1", "k2": "v2"},
		},
	}

	assert.Equal(t, expected, result)
}

func TestCandidateAddErrors(t *testing.T) {
	store := newMemoryStore(t)
	if err := store.Write(db.CandidatesKey, models.Candidates{{Name: "John"}}); err != nil {
		t.Fatal(err)
	}

	cases := map[string]string{
		"missing NAME":             "!candidate add",
		"missing MANAGER":          "!candidate add NAME",
		"missing '@' on MANAGER":   "!candidate add NAME MANAGER",
		"missing '<>' on MANAGER":  "!candidate add NAME @MANAGER",
		"missing MANAGER meta":     "!candidate add --meta NAME <@MANAGER>",
		"missing value for key":    "!candidate add --meta key NAME <@MANAGER>",
		"missing MANAGER key meta": "!candidate add --meta key:val NAME",
		"NAME exists":              "!candidate add John <@MANAGER>",
	}

	cmd := NewCandidateCommand(store, ioutil.Discard)
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

// todo: test --count and --ascending flag
func TestCandidateList(t *testing.T) {
	candidates := models.Candidates{
		{Name: "John Doe"},
		{Name: "Jane Doe"},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.CandidatesKey, candidates); err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBuffer(nil)
	cmd := NewCandidateCommand(store, w)
	if err := runTestApp(cmd, "!candidate ls"); err != nil {
		t.Fatal(err)
	}

	for _, candidate := range candidates {
		assert.Contains(t, w.String(), candidate.Name)
	}
}

func TestCandidateListErrors(t *testing.T) {
	cmd := NewCandidateCommand(newMemoryStore(t), ioutil.Discard)
	if err := runTestApp(cmd, "!candidate ls"); err == nil {
		t.Fatal("Error was nil!")
	}
}

func TestCandidateRemove(t *testing.T) {
	candidates := models.Candidates{
		{Name: "John Doe"},
		{Name: "Jane Doe"},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.CandidatesKey, candidates); err != nil {
		t.Fatal(err)
	}

	interviews := models.Interviews{
		{Candidate: "John Doe"},
		{Candidate: "John Doe"},
		{Candidate: "Jane Doe"},
	}

	if err := store.Write(db.InterviewsKey, interviews); err != nil {
		t.Fatal(err)
	}

	cmd := NewCandidateCommand(store, ioutil.Discard)
	if err := runTestApp(cmd, "!candidate rm John Doe"); err != nil {
		t.Fatal(err)
	}

	resultCandidates := models.Candidates{}
	if err := store.Read(db.CandidatesKey, &resultCandidates); err != nil {
		t.Fatal(err)
	}

	expectedCandidates := models.Candidates{
		{Name: "Jane Doe"},
	}

	assert.Equal(t, expectedCandidates, resultCandidates)

	resultInterviews := models.Interviews{}
	if err := store.Read(db.InterviewsKey, &resultInterviews); err != nil {
		t.Fatal(err)
	}

	expectedInterviews := models.Interviews{
		{Candidate: "Jane Doe"},
	}

	assert.Equal(t, expectedInterviews, resultInterviews)
}

func TestCandidateRemoveErrors(t *testing.T) {
	cases := map[string]string{
		"missing NAME":       "!candidate rm",
		"NAME doesn't exist": "!candidate rm John Doe",
	}

	cmd := NewCandidateCommand(newMemoryStore(t), ioutil.Discard)
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestCandidateShow(t *testing.T) {
	candidates := models.Candidates{
		{Name: "John Doe"},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.CandidatesKey, candidates); err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBuffer(nil)
	cmd := NewCandidateCommand(store, w)
	if err := runTestApp(cmd, "!candidate show john doe"); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, w.String(), "John Doe")
}

func TestCandidateShowErrors(t *testing.T) {
	cases := map[string]string{
		"missing NAME":       "!candidate show",
		"NAME doesn't exist": "!candidate show John Doe",
	}

	cmd := NewCandidateCommand(newMemoryStore(t), ioutil.Discard)
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestCandidateUpdate(t *testing.T) {
	candidates := models.Candidates{
		{
			Name: "John Doe",
			Meta: map[string]string{
				"k1": "v1",
				"k2": "v2",
			},
		},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.CandidatesKey, candidates); err != nil {
		t.Fatal(err)
	}

	cmd := NewCandidateCommand(store, ioutil.Discard)
	if err := runTestApp(cmd, "!candidate update \"John Doe\" k1 updated"); err != nil {
		t.Fatal(err)
	}

	result := models.Candidates{}
	if err := store.Read(db.CandidatesKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Candidates{
		{
			Name: "John Doe",
			Meta: map[string]string{
				"k1": "updated",
				"k2": "v2",
			},
		},
	}

	assert.Equal(t, expected, result)
}

func TestCandidateUpdateErrors(t *testing.T) {
	cases := map[string]string{
		"missing NAME":       "!candidate update",
		"missing KEY":        "!candidate update NAME",
		"missing VAL":        "!candidate update NAME KEY",
		"NAME doesn't exist": "!candidate update NAME KEY VAL",
	}

	cmd := NewCandidateCommand(newMemoryStore(t), ioutil.Discard)
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}
