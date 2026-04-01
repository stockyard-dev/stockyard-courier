package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-courier/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){status:=r.URL.Query().Get("status");list,_:=s.db.List(status);if list==nil{list=[]store.Task{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var t store.Task;json.NewDecoder(r.Body).Decode(&t);if t.Title==""{writeError(w,400,"title required");return};s.db.Create(&t);writeJSON(w,201,t)}
func(s *Server)handleUpdate(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var req struct{Status string `json:"status"`};json.NewDecoder(r.Body).Decode(&req);if req.Status==""{req.Status="done"};s.db.UpdateStatus(id,req.Status);writeJSON(w,200,map[string]string{"status":req.Status})}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
