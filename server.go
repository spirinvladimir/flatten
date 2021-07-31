package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Result struct {
	Flat  string `json:"flat"`
	Depth int    `json:"depth"`
}
type Input []string
type Output []Result
type Request struct {
	Input     Input    `json:"input"`
	Output    []Result `json:"output"`
	Timestamp int64    `json:"timestamp"`
}

const open byte = 40  // "("
const close byte = 41 // ")"
// Cache for avoid calculation same flattening tasks
var cache = map[string]Result{} // In case horizontal scaling ..Redis,etc
const HistorySize = 100

var history []Request

func flatten(task string) Result {
	deep := []byte(task)
	flat := ""
	level := -1
	depth := 0
	for i := 0; i < len(deep); i++ {
		char := deep[i]
		switch char {
		case open:
			level++
			if level > depth {
				depth = level
			}
		case close:
			level--
		default:
			flat += string(char)
		}
	}

	return Result{
		Flat:  flat,
		Depth: depth,
	}
}

func Flatten(w http.ResponseWriter, r *http.Request) {
	var input Input
	var output Output

	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &input)

	for i := 0; i < len(input); i++ {
		deep := input[i]
		result, exist := cache[deep]
		if !exist {
			result = flatten(deep)
			cache[deep] = result
		}
		output = append(output, result)
	}

	history = append(history, Request{
		Input:     input,
		Output:    output,
		Timestamp: time.Now().Unix(),
	})

	if len(history) > HistorySize {
		history = history[1 : HistorySize+1]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

func History(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}

func main() {
	http.HandleFunc("/flatten", Flatten)
	http.HandleFunc("/history", History)
	http.ListenAndServe(":8080", nil)
}
