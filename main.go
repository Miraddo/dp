package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// UserDataProcessor is a data structure that encapsulates user data handling
type UserDataProcessor struct {
	userQuotas    map[string]*UserQuota
	dataStorage   map[string]bool
	requestCounts map[string]int
	mu            sync.RWMutex
}

// UserQuota is a data structure to keep track of the amount of data each user has sent
type UserQuota struct {
	TotalDataLimit int64
	UsedData       int64
}

// Data is a data structure that imitates receiving data
type Data struct {
	ID     string `json:"id"`
	UserID string `json:"userID"`
	Data   string `json:"data"`
}

// InitializeUserDataProcessor initializes a new UserDataProcessor
func InitializeUserDataProcessor() *UserDataProcessor {
	return &UserDataProcessor{
		userQuotas:    make(map[string]*UserQuota),
		dataStorage:   make(map[string]bool),
		requestCounts: make(map[string]int),
	}
}

// ProcessRequest is a function that processes incoming data
func (udp *UserDataProcessor) ProcessRequest(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(rw, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	data := &Data{}
	err := json.NewDecoder(req.Body).Decode(data)
	if err != nil {
		http.Error(rw, "Invalid data", http.StatusBadRequest)
		return
	}

	udp.mu.Lock()

	// Check for duplicate data
	if _, exists := udp.dataStorage[data.ID]; exists {
		http.Error(rw, fmt.Sprintf("Duplicate Data ID: %s", data.ID), http.StatusBadRequest)
		udp.mu.Unlock()
		return
	}

	// Check user quotas
	q, ok := udp.userQuotas[data.UserID]
	if !ok {
		// First time the user is seen, create a new quota for the user
		q = &UserQuota{TotalDataLimit: 100000, UsedData: int64(len(data.Data))}
		udp.userQuotas[data.UserID] = q
	} else {
		// User quota exists, check if the new data fits in
		if q.UsedData+int64(len(data.Data)) > q.TotalDataLimit {
			http.Error(rw, fmt.Sprintf("Data limit exceeded for User ID %s", data.UserID), http.StatusBadRequest)
			udp.mu.Unlock()
			return
		}
		q.UsedData += int64(len(data.Data))
	}

	// Limit the request frequency
	if udp.requestCounts[data.UserID] >= 5 {
		http.Error(rw, fmt.Sprintf("Too many requests from User ID %s", data.UserID), http.StatusBadRequest)
		udp.mu.Unlock()
		return
	}
	udp.requestCounts[data.UserID]++

	udp.mu.Unlock()

	// Process the data (simply print) or we can sent it to a queue (SQS, Kafka, RabitMQ,...)
	fmt.Println("Data processed:")
	fmt.Println("ID:", data.ID)
	fmt.Println("UserID:", data.UserID)
	fmt.Println("Data:", data.Data)

	// Update the data storage
	udp.mu.Lock()
	udp.dataStorage[data.ID] = true
	udp.mu.Unlock()

	log.Println("Data processed and quotas updated successfully for User ID", data.UserID)
	fmt.Fprint(rw, "Data processed successfully")
}

// resetRequestCounts is a background task that resets the user's request counts every minute
func (udp *UserDataProcessor) ResetRequestCounts() {
	for {
		time.Sleep(time.Minute)
		udp.mu.Lock()
		for k := range udp.requestCounts {
			udp.requestCounts[k] = 0
		}
		udp.mu.Unlock()
	}
}

func main() {
	processor := InitializeUserDataProcessor()
	go processor.ResetRequestCounts() // Running it as a goroutine

	http.HandleFunc("/data", processor.ProcessRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
