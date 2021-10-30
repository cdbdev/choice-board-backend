package ctrl

import (
	"net/http"
	"encoding/json"
)

type ResultData struct {
	Result 		string		`json:"result"`
}

func GenerateResult(info string) []byte {
	result := ResultData{ Result: info }
	data, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	
	return data
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found!"))
}