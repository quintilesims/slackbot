package bot

/*
func TestCandidateAdd(t *testing.T) {
	store := newMemoryStore(t)
	if err := store.Write(db.CandidatesKey, models.Candidates{"John Doe": nil}); err != nil {
		t.Fatal(err)
	}

	cmd := NewCandidateCommand(store, ioutil.Discard)
	if err := runTestApp(cmd, "!candidate add --meta k1=v1 --meta k2=v2 Jane Doe"); err != nil {
		t.Fatal(err)
	}

	result := models.Candidates{}
	if err := store.Read(db.CandidatesKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Candidates{
		"John Doe": nil,
		"Jane Doe": map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
	}

	assert.Equal(t, expected, result)
}

func TestCandidateAddErrors(t *testing.T) {
	inputs := []string{
		"!candidate add",
		"!candidate add John Doe",
		"!candidate add --meta stuff NAME",
		"!candidate add --meta key:val NAME",
	}

	store := newMemoryStore(t)
	if err := store.Write(db.CandidatesKey, models.Candidates{"John Doe": nil}); err != nil {
		t.Fatal(err)
	}

	cmd := NewCandidateCommand(store, ioutil.Discard)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

// todo: test --count flag
func TestCandidateList(t *testing.T) {
	candidates := models.Candidates{
		"John Doe": nil,
		"Jane Doe": nil,
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

	for name := range candidates {
		assert.Contains(t, w.String(), name)
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
		"John Doe": nil,
		"Jane Doe": nil,
	}

	store := newMemoryStore(t)
	if err := store.Write(db.CandidatesKey, candidates); err != nil {
		t.Fatal(err)
	}

	cmd := NewCandidateCommand(store, ioutil.Discard)
	if err := runTestApp(cmd, "!candidate rm John Doe"); err != nil {
		t.Fatal(err)
	}

	result := models.Candidates{}
	if err := store.Read(db.CandidatesKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Candidates{
		"Jane Doe": nil,
	}

	assert.Equal(t, expected, result)
}

func TestCandidateRemoveErrors(t *testing.T) {
	inputs := []string{
		"!candidate rm",
		"!candidate rm John Doe",
	}

	cmd := NewCandidateCommand(newMemoryStore(t), ioutil.Discard)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestCandidateInfo(t *testing.T) {
	candidates := models.Candidates{
		"John Doe": map[string]string{"k1": "v1", "k2": "v2"},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.CandidatesKey, candidates); err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBuffer(nil)
	cmd := NewCandidateCommand(store, w)
	if err := runTestApp(cmd, "!candidate info John Doe"); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, w.String(), "John Doe")
	for key, val := range map[string]string{"k1": "v1", "k2": "v2"} {
		assert.Contains(t, w.String(), key)
		assert.Contains(t, w.String(), val)
	}
}

func TestCandidateInfoErrors(t *testing.T) {
	inputs := []string{
		"!candidate info",
		"!candidate info John Doe",
	}

	cmd := NewCandidateCommand(newMemoryStore(t), ioutil.Discard)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestCandidateUpdate(t *testing.T) {
	candidates := models.Candidates{
		"John Doe": map[string]string{"k1": "v1"},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.CandidatesKey, candidates); err != nil {
		t.Fatal(err)
	}

	cmd := NewCandidateCommand(store, ioutil.Discard)
	if err := runTestApp(cmd, "!candidate update \"John Doe\" k1 v2"); err != nil {
		t.Fatal(err)
	}

	result := models.Candidates{}
	if err := store.Read(db.CandidatesKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Candidates{
		"John Doe": map[string]string{
			"k1": "v2",
		},
	}

	assert.Equal(t, expected, result)
}

func TestCandidateUpdateErrors(t *testing.T) {
	inputs := []string{
		"!candidate update",
		"!candidate update \"John Doe\"",
		"!candidate update \"John Doe\" k1",
		"!candidate update \"John Doe\" k1 v1",
	}

	cmd := NewCandidateCommand(newMemoryStore(t), ioutil.Discard)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}
*/
