package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("src")))

	http.HandleFunc("/getBankDetails", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method. Only GET requests are allowed.", http.StatusMethodNotAllowed)
			return
		}

		ifsc := r.URL.Query().Get("ifsc")

		if ifsc == "" {
			http.Error(w, "IFSC code is required. Please provide a valid IFSC code.", http.StatusBadRequest)
			return
		}

		resp, err := http.Get("https://ifsc.razorpay.com/" + ifsc)
		if err != nil {
			http.Error(w, "Failed to fetch bank details. Please try again later.", http.StatusNotFound)
			return
		}
		defer resp.Body.Close()

		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&data); err != nil {
			errorResponse := map[string]string{"error": "Failed to parse bank details. The response format is not as expected."}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			responseJSON, _ := json.Marshal(errorResponse)
			w.Write(responseJSON)
			return
		}
		

		if data["error"] != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			responseJSON := fmt.Sprintf(`{"error": "%s"}`, data["error"])
			w.Write([]byte(responseJSON))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			responseJSON, err := json.Marshal(data)
			if err != nil {
				http.Error(w, "Failed to create response. An error occurred while generating the response.", http.StatusInternalServerError)
				return
			}
			w.Write(responseJSON)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
