package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Payload struct {
    Mag       []string `json:"mag"`
    StartDate string   `json:"startDate"`
    EndDate   string   `json:"endDate"`
    States    []string `json:"states"`
}

type TornadoData struct {
    ID        int     `json:"id"`
    Om        string  `json:"om"`
    Yr        string  `json:"yr"`
    Mo        string  `json:"mo"`
    Dy        string  `json:"dy"`
    Date      string  `json:"date"`
    Time      string  `json:"time"`
    Tz        string  `json:"tz"`
    St        string  `json:"st"`
    Stf       string  `json:"stf"`
    Stn       string  `json:"stn"`
    Mag       string  `json:"mag"`
    Inj       string  `json:"inj"`
    Fat       string  `json:"fat"`
    Loss      string  `json:"loss"`
    Closs     string  `json:"closs"`
    Slat      string  `json:"slat"`
    Slon      string  `json:"slon"`
    Elat      string  `json:"elat"`
    Elon      string  `json:"elon"`
    Len       string  `json:"len"`
    Wid       string  `json:"wid"`
    Ns        string  `json:"ns"`
    Sn        string  `json:"sn"`
    Sg        string  `json:"sg"`
    F1        string  `json:"f1"`
    F2        string  `json:"f2"`
    F3        string  `json:"f3"`
    F4        string  `json:"f4"`
    Fc        string  `json:"fc"`
}

type APIResponse struct {
    Data          []TornadoData `json:"data"`
    TotalEntries  int           `json:"total_entries"`
    Timestamp     string        `json:"timestamp"`
    StatusCode    int           `json:"statusCode"`
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    http.HandleFunc("/", serveIndex)
    http.HandleFunc("/api", handleAPIRequest)

    log.Println("Server started at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("index.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

func handleAPIRequest(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Error parsing form data", http.StatusBadRequest)
        return
    }

    payload := Payload{
        StartDate: r.FormValue("startDate"),
        EndDate:   r.FormValue("endDate"),
        Mag:       r.Form["mag"],
        States:    r.Form["states"],
    }

    url := os.Getenv("API_URL")
    if url == "" {
        http.Error(w, "API URL is not set", http.StatusInternalServerError)
        return
    }

    data, err := json.Marshal(payload)
    if err != nil {
        http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
        return
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
    if err != nil {
        http.Error(w, "Error creating request", http.StatusInternalServerError)
        return
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        http.Error(w, "Error making API request", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var apiResponse APIResponse
    err = json.NewDecoder(resp.Body).Decode(&apiResponse)
    if err != nil {
        http.Error(w, "Error decoding response", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(apiResponse)
}
