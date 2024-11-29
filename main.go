package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Server is starting on port 10001...")
	http.Handle("/", http.FileServer(http.Dir("src")))

	http.HandleFunc("/getBankDetails", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for /getBankDetails")

		if r.Method != http.MethodGet {
			log.Println("Invalid request method. Expected GET.")
			http.Error(w, "Invalid request method. Only GET requests are allowed.", http.StatusMethodNotAllowed)
			return
		}

		ifsc := r.URL.Query().Get("ifsc")
		if ifsc == "" {
			log.Println("IFSC code not provided in the request")
			errorResponse := map[string]string{"error": "IFSC code is required. Please provide a valid IFSC code."}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			responseJSON, _ := json.Marshal(errorResponse)
			w.Write(responseJSON)
			return
		}

		log.Printf("Fetching bank details for IFSC: %s\n", ifsc)
		resp, err := http.Get("https://ifsc.razorpay.com/" + ifsc)
		if err != nil {
			log.Printf("Failed to fetch bank details for IFSC %s: %v\n", ifsc, err)
			http.Error(w, "Failed to fetch bank details. Please try again later.", http.StatusNotFound)
			return
		}
		defer resp.Body.Close()

		log.Printf("Received response for IFSC %s: %v\n", ifsc, resp.StatusCode)

		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&data); err != nil {
			log.Printf("Failed to decode response for IFSC %s: %v\n", ifsc, err)
			errorResponse := map[string]string{"error": "Failed to parse bank details. The response format is not as expected."}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			responseJSON, _ := json.Marshal(errorResponse)
			w.Write(responseJSON)
			return
		}

		if data["error"] != nil {
			log.Printf("Error returned from Razorpay API for IFSC %s: %v\n", ifsc, data["error"])
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			responseJSON := fmt.Sprintf(`{"error": "%s"}`, data["error"])
			w.Write([]byte(responseJSON))
		} else {
			log.Printf("Bank details for IFSC %s: %v\n", ifsc, data)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			responseJSON, err := json.Marshal(data)
			if err != nil {
				log.Printf("Failed to marshal bank details response for IFSC %s: %v\n", ifsc, err)
				http.Error(w, "Failed to create response. An error occurred while generating the response.", http.StatusInternalServerError)
				return
			}
			w.Write(responseJSON)
		}
	})

	// Log that the application has started and is running
	log.Println("Application is started and running on port 10001")

	log.Fatal(http.ListenAndServe(":10001", nil))
}
