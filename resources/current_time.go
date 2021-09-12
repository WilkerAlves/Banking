package resources

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func _() {
	router := mux.NewRouter()

	router.HandleFunc("/api/time", func(w http.ResponseWriter, r *http.Request) {
		tz := r.URL.Query()["tz"]
		result := make(map[string]string)

		for _, value := range tz {

			loc, _ := time.LoadLocation(value)
			result[value] = time.Now().In(loc).String()
		}

		w.Header().Add("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(result)
		if err != nil {
			panic(err.Error())
		}

	}).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe("localhost:8000", router))
}
