package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// this is the model object of the Scheduler
type Scheduler struct {
Id string `json:"id"` // uuid
    Data string `json:"data"`// data json -- unsafe, will need to update once I figure out what it should be 
    Schedule_time int `json:"schedule_time"`//utc - time since epoch
    Created_at int `json:"created_time"`//utc - time since epoch
    Updated_at int `json:"updated_at"`//utc - time since epoch
}

type InvalidResponse struct {
    error string
}

type SchedulerInput struct {
    Data string `json:"data"`
    Schedule_at int `json:"schedule_at"`
}

func initDB() {
    var err error;
    DB, err = sql.Open("sqlite3", "./scheduler.db")
    if err != nil {
        log.Fatal(err)
    }
}

// this should be a POST request as I need the body for the data
func schedulerHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.Method)
    switch r.Method {
    case http.MethodGet:
        getSchedulers(w, r)
    case http.MethodPost:
        createScheduler(w, r)
    case http.MethodPut:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    case http.MethodDelete:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func getSchedulers(w http.ResponseWriter, r *http.Request) {
    var schedules []Scheduler
    res, err := DB.Query("SELECT id, data, schedule_time, created_at, updated_at FROM scheduler")
    if err != nil {
        return
    }
    defer res.Close()

    for res.Next() {
        s := Scheduler{}
        res.Scan(&s.Id, &s.Data, &s.Schedule_time, &s.Created_at, &s.Updated_at)
        schedules = append(schedules, s)
    }

    data, err := json.Marshal(schedules)
    if err != nil {
        log.Fatal(err)
    }
    w.Write(data)
}

// get data from body
// get schedule_at from body
func createScheduler(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var s SchedulerInput
    err := decoder.Decode(&s)
    log.Println(s)
    if err != nil {
        data, err := json.Marshal(InvalidResponse{error: "Invalid request body"})
        if err != nil {
            log.Fatal(err)
        }
        w.Write(data)
    }

    _, err = DB.Exec(`INSERT into scheduler (id, data, schedule_time, created_at, updated_at) values (?, ?, ?, ?, ?);`, uuid.New().String(), s.Data, s.Schedule_at, time.Now().Unix(), time.Now().Unix())
    if err != nil {
        data, err := json.Marshal(InvalidResponse{error: "Invalid request body"})
        if err != nil {
            log.Fatal(err)
        }
        w.Write(data)
        log.Fatal("Unable to insert into scheduler table", err)
        return
    }
}

func main() {
    initDB()
    defer DB.Close()

    http.HandleFunc("/", schedulerHandler)
    fmt.Println("Server is running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// TODO: set up the cron the read from the scheduler table -- maybe just print out what is within the data column
// Clean up this code -- add some tests. I feel like it barely hanging out; that is can break easily
