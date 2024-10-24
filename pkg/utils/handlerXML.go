package utils

import (
	"encoding/xml"
	"net/http"
	"tripleS/pkg"
)

func HandlerXML(w http.ResponseWriter, response pkg.Response) {
	w.Header().Set("Content-Type", "application/xml")
	result, err := xml.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(response.Status)
	w.Write(result)
}

func HandlerData(w http.ResponseWriter, response []pkg.BucketMetadata) {
	w.Header().Set("Content-Type", "application/xml")
	result, err := xml.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	w.Write(result)
}
