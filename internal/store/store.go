package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Message struct {
	ID string `json:"id"`
	Queue string `json:"queue"`
	Payload string `json:"payload"`
	Priority int `json:"priority"`
	Status string `json:"status"`
	Retries int `json:"retries"`
	MaxRetries int `json:"max_retries"`
	ProcessedAt string `json:"processed_at"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"courier.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS messages(id TEXT PRIMARY KEY,queue TEXT NOT NULL,payload TEXT DEFAULT '',priority INTEGER DEFAULT 0,status TEXT DEFAULT 'pending',retries INTEGER DEFAULT 0,max_retries INTEGER DEFAULT 3,processed_at TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Message)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO messages(id,queue,payload,priority,status,retries,max_retries,processed_at,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Queue,e.Payload,e.Priority,e.Status,e.Retries,e.MaxRetries,e.ProcessedAt,e.CreatedAt);return err}
func(d *DB)Get(id string)*Message{var e Message;if d.db.QueryRow(`SELECT id,queue,payload,priority,status,retries,max_retries,processed_at,created_at FROM messages WHERE id=?`,id).Scan(&e.ID,&e.Queue,&e.Payload,&e.Priority,&e.Status,&e.Retries,&e.MaxRetries,&e.ProcessedAt,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Message{rows,_:=d.db.Query(`SELECT id,queue,payload,priority,status,retries,max_retries,processed_at,created_at FROM messages ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Message;for rows.Next(){var e Message;rows.Scan(&e.ID,&e.Queue,&e.Payload,&e.Priority,&e.Status,&e.Retries,&e.MaxRetries,&e.ProcessedAt,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM messages WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM messages`).Scan(&n);return n}
