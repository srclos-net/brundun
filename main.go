package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

type CommandResponse struct {
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	Returncode int    `json:"returncode"`
}

func main() {
	http.HandleFunc("/execute", executeCommand)
	http.HandleFunc("/", healthCheck)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func executeCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Add authentication and authorization checks here

	var command string
	err := json.NewDecoder(r.Body).Decode(&struct {
		Command string `json:"command"`
	}{command})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	// Add security validation for the command here
	// e.g., check against a whitelist of allowed commands

	output, err := exec.Command("sh", "-c", command).CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error executing command: %v", err)
		return
	}

	response := CommandResponse{
		Stdout:     string(output),
		Stderr:     "", // Not used in this example
		Returncode: 0,  // Assuming success
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding response: %v", err)
		return
	}
}