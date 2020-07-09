package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func auth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := os.Getenv("API_KEY")

		if len(apiKey) == 0 {
			fmt.Println("No API key set")
			http.Error(w, "Configuration error", http.StatusInternalServerError)
			return
		}

		if r.Header.Get("api-key") != apiKey {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func main() {
	http.HandleFunc("/status", auth(getStatus))
	http.HandleFunc("/listconfig", auth(listConfig))
	http.HandleFunc("/importconfig", auth(importConfig))
	http.HandleFunc("/runpipeline", auth(runPipeline))
	http.ListenAndServe(":8080", nil)
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("/app/entrypoint.sh", "status").Output()
	if err != nil {
		fmt.Fprintf(w, "Error %s while executing command: %s", err, out)
		return
	}
	fmt.Fprintf(w, "%s", out)
}

func listConfig(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("ls", "-l", "/app/config-data").Output()
	if err != nil {
		fmt.Fprintf(w, "Error %s while executing command: %s", err, out)
		return
	}

	fmt.Fprintf(w, "%s", out)
}

func importConfig(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("/app/entrypoint.sh", "import", "--dir", "/app/config-data").Output()
	if err != nil {
		fmt.Fprintf(w, "Error %s while executing command: %s", err, out)
		return
	}

	fmt.Fprintf(w, "%s", out)
}

func runPipeline(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	tapID := queryParams.Get("tap_id")
	targetID := queryParams.Get("target_id")

	if len(tapID) == 0 || len(targetID) == 0 {
		fmt.Fprint(w, "Tap id or target id not present")
		return
	}

	out, err := exec.Command("/app/entrypoint.sh", "run_tap", "--tap", tapID, "--target", targetID).Output()
	if err != nil {
		fmt.Fprintf(w, "Error %s while executing command: %s", err, out)
		return
	}

	fmt.Fprintf(w, "%s", out)
}
