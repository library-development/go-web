package web

// type DB struct {
// 	LocalPath string
// 	Locks     map[string]*sync.Mutex
// }

// func (db *DB) Table(id string) *Table {
// 	return &Table{
// 		LocalPath: filepath.Join(db.LocalPath, id),
// 		Lock:      db.Locks[id],
// 	}
// }

// func (db *DB) CreateTable(id string, typ golang.Ident) error {
// 	path := filepath.Join(db.LocalPath, id)
// 	err := os.MkdirAll(path, 0755)
// 	if err != nil {
// 		return err
// 	}
// 	b, err := json.Marshal(typ)
// 	if err != nil {
// 		return err
// 	}
// 	path = filepath.Join(path, "type")
// 	err = os.WriteFile(path, b, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
