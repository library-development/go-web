package web

// type Table struct {
// 	LocalPath string
// 	Lock      *sync.Mutex
// }

// // newID generates a new unique ID for the given scope.
// func (t *Table) newID() string {
// 	t.Lock.Lock()
// 	defer t.Lock.Unlock()
// 	return ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()
// }

// // Get returns a file from the table.
// func (t *Table) Get(id string) (*File, error) {
// 	path := filepath.Join(t.LocalPath, id)
// 	b, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var f *File
// 	if err := json.Unmarshal(b, &f); err != nil {
// 		return nil, err
// 	}
// 	return f, nil
// }

// // Search returns all files in the table that match the search string.
// func (t *Table) Search(search string) ([]*File, error) {
// 	filter := func(f *File) bool {
// 		if strings.Contains(f.Metadata.Name, search) {
// 			return true
// 		}
// 		if strings.Contains(f.Metadata.Doc, search) {
// 			return true
// 		}
// 		return false
// 	}
// 	return t.List(filter)
// }

// // List returns all files in the table that match the filter.
// func (t *Table) List(filter func(*File) bool) ([]*File, error) {
// 	var files []*File
// 	err := filepath.Walk(t.LocalPath, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if info.IsDir() {
// 			// TODO: build and return special file for directories
// 			return nil
// 		}
// 		b, err := os.ReadFile(path)
// 		if err != nil {
// 			return err
// 		}
// 		var f *File
// 		if err := json.Unmarshal(b, &f); err != nil {
// 			return err
// 		}
// 		if filter != nil && !filter(f) {
// 			return nil
// 		}
// 		files = append(files, f)
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return files, nil
// }

// func (t *Table) Type() golang.Ident {
// 	path := filepath.Join(t.LocalPath, "type")
// 	b, err := os.ReadFile(path)
// 	var id golang.Ident
// 	if err != nil {
// 		panic(err)
// 	}
// 	json.Unmarshal(b, &golang.Ident{})
// 	return id
// }

// // Post adds a new file to the table.
// // A random name is given to the file and returned.
// func (t *Table) Post(owners map[string]bool, data []byte) (string, error) {
// 	var f *File
// 	f.Metadata.Type = t.Type().String()
// 	f.Metadata.Owners = owners
// 	f.Metadata.Name = t.newID()
// 	f.Metadata.CreatedAt = time.Now().UnixNano()
// 	f.Metadata.UpdatedAt = f.Metadata.CreatedAt
// 	f.Data = data
// 	path := filepath.Join(t.LocalPath, f.Metadata.Name)
// 	os.MkdirAll(filepath.Dir(path), os.ModePerm)
// 	b, err := json.Marshal(f)
// 	if err != nil {
// 		return "", err
// 	}
// 	if err := os.WriteFile(path, b, os.ModePerm); err != nil {
// 		return "", err
// 	}
// 	return f.Metadata.Name, nil
// }

// // Put updates an existing file in the table.
// func (t *Table) Put(id string, value []byte) error {
// 	f, err := t.Get(id)
// 	if err != nil {
// 		return err
// 	}
// 	f.Data = value
// 	f.Metadata.UpdatedAt = time.Now().UnixNano()
// 	path := filepath.Join(t.LocalPath, f.Metadata.Name)
// 	b, err := json.Marshal(f)
// 	if err != nil {
// 		return err
// 	}
// 	if err := os.WriteFile(path, b, os.ModePerm); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Delete removes an existing file from the table.
// func (t *Table) Delete(id string) error {
// 	path := filepath.Join(t.LocalPath, id)
// 	if err := os.Remove(path); err != nil {
// 		return err
// 	}
// 	return nil
// }
