package store
import("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{*sql.DB}
type Task struct{ID int64 `json:"id"`;Title string `json:"title"`;Body string `json:"body"`;From string `json:"from"`;To string `json:"to"`;Status string `json:"status"`;CreatedAt time.Time `json:"created_at"`;UpdatedAt time.Time `json:"updated_at"`}
func Open(d string)(*DB,error){os.MkdirAll(d,0755);dsn:=filepath.Join(d,"courier.db")+"?_journal_mode=WAL&_busy_timeout=5000";db,err:=sql.Open("sqlite",dsn);if err!=nil{return nil,fmt.Errorf("open: %w",err)};db.SetMaxOpenConns(1);migrate(db);return &DB{db},nil}
func migrate(db *sql.DB){db.Exec(`CREATE TABLE IF NOT EXISTS tasks(id INTEGER PRIMARY KEY AUTOINCREMENT,title TEXT NOT NULL,body TEXT DEFAULT '',from_person TEXT DEFAULT '',to_person TEXT DEFAULT '',status TEXT DEFAULT 'open',created_at DATETIME DEFAULT CURRENT_TIMESTAMP,updated_at DATETIME DEFAULT CURRENT_TIMESTAMP)`)}
func(db *DB)Create(t *Task)error{res,err:=db.Exec(`INSERT INTO tasks(title,body,from_person,to_person)VALUES(?,?,?,?)`,t.Title,t.Body,t.From,t.To);if err!=nil{return err};t.ID,_=res.LastInsertId();return nil}
func(db *DB)List(status string)([]Task,error){q:=`SELECT id,title,body,from_person,to_person,status,created_at,updated_at FROM tasks`;if status!=""{q+=` WHERE status=?`};q+=` ORDER BY created_at DESC`;var rows *sql.Rows;if status!=""{rows,_=db.Query(q,status)}else{rows,_=db.Query(q)};defer rows.Close();var out[]Task;for rows.Next(){var t Task;rows.Scan(&t.ID,&t.Title,&t.Body,&t.From,&t.To,&t.Status,&t.CreatedAt,&t.UpdatedAt);out=append(out,t)};return out,nil}
func(db *DB)UpdateStatus(id int64,status string){db.Exec(`UPDATE tasks SET status=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`,status,id)}
func(db *DB)Delete(id int64){db.Exec(`DELETE FROM tasks WHERE id=?`,id)}
func(db *DB)Stats()(map[string]interface{},error){var open,done int;db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status='open'`).Scan(&open);db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status='done'`).Scan(&done);return map[string]interface{}{"open":open,"done":done},nil}
