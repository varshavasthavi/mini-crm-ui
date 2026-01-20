package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DepositEvent struct {
	PlayerID string  `json:"player_id"`
	Amount   float64 `json:"amount"`
}

type CampaignLog struct {
	EventID   string    `bson:"event_id" json:"event_id"`
	PlayerID  string    `bson:"player_id" json:"player_id"`
	Action    string    `bson:"action" json:"action"`
	Amount    float64   `bson:"amount" json:"amount"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

var (
	mongoCol *mongo.Collection
	chConn   clickhouse.Conn
)

func main() {
    var err error // Declare err once at the top of the function

    // 1. Setup MongoDB
    mClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil { log.Fatal(err) }
    mongoCol = mClient.Database("mini_crm").Collection("campaign_logs")

    // 2. Setup ClickHouse 
    chConn, err = clickhouse.Open(&clickhouse.Options{
        Addr: []string{"127.0.0.1:9000"},
        Auth: clickhouse.Auth{
            Database: "default",
            Username: "default",
            Password: "",
        },
    })
    if err != nil { log.Fatal("ClickHouse connection failed: ", err) }

    // Routes...
    http.HandleFunc("/ingest", handleIngest)
    http.HandleFunc("/logs", handleGetLogs)

    log.Println("Backend active on :8080")
    http.ListenAndServe(":8080", nil)
}

func handleIngest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodPost { return }

	var ev DepositEvent
	json.NewDecoder(r.Body).Decode(&ev)

	eventID := uuid.New().String()
	ts := time.Now()

	// 1. ClickHouse: Persist analytical event
	err := chConn.Exec(r.Context(), "INSERT INTO events VALUES (?, ?, ?, ?, ?)",
		eventID, ev.PlayerID, "DEPOSIT", ev.Amount, ts)
	if err != nil {
		log.Printf("ClickHouse Insert Error: %v", err)
	} else {
		log.Println("Successfully inserted into ClickHouse")
	}

	// 2. Business Logic: Eval amount >= 1000
	if ev.Amount >= 1000 {
		logEntry := CampaignLog{
			EventID: eventID, PlayerID: ev.PlayerID,
			Action: "BONUS_MESSAGE", Amount: ev.Amount, Timestamp: ts,
		}
		// 3. MongoDB: Persist log
		mongoCol.InsertOne(r.Context(), logEntry)
	}
	w.WriteHeader(http.StatusAccepted)
}

func handleGetLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	cursor, _ := mongoCol.Find(r.Context(), bson.M{})
	var logs []CampaignLog = make([]CampaignLog, 0)
	cursor.All(r.Context(), &logs)
	json.NewEncoder(w).Encode(logs)
}