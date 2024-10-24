package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"tripleS/pkg"
	"tripleS/pkg/utils"
)

func CreatBucketHandler(w http.ResponseWriter, r *http.Request, directory *string) {
	bucketName := ""
	filePath := strings.Split(r.URL.Path[1:], "/")

	if len(filePath) > 0 {
		bucketName = filePath[0]
	}
	isBucketExist, err := utils.IsBucketExist(directory, &bucketName)
	if err != nil {
		utils.HandlerXML(w, pkg.Response{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	if len(filePath) > 1 {
		objectKey := filePath[1]

		if !isBucketExist {
			utils.HandlerXML(w, pkg.Response{Status: http.StatusBadRequest, Message: "bucket does not exist"})
			return
		}

		err := ObjectHandler(r, directory, &bucketName, &objectKey)
		if err != nil {
			utils.HandlerXML(w, pkg.Response{Status: http.StatusInternalServerError, Message: err.Error()})
			return
		}

		err = utils.HandlerCSV(*directory+"/"+bucketName+"/"+bucketName+".csv", "POST objects.csv", fmt.Sprint(r.ContentLength), r.Header.Get("Content-Type"), "", objectKey)
		if err != nil {
			return
		}
		utils.HandlerXML(w, pkg.Response{Status: http.StatusOK, Message: "Object added successfully"})
		return
	}

	// error handling
	if !utils.IsValid(bucketName) {
		utils.HandlerXML(w, pkg.Response{Status: http.StatusBadRequest, Message: "Invalid Bucket Name"})
		return
	}
	if isBucketExist {
		utils.HandlerXML(w, pkg.Response{Status: http.StatusBadRequest, Message: "bucket already exist"})
		return
	}

	err = os.MkdirAll(*directory+"/"+bucketName, 0o775)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = utils.HandlerCSV(*directory+"/"+bucketName+"/"+bucketName+".csv", "PUT objects.csv", "", "", "", "")
	if err != nil {
		utils.HandlerXML(w, pkg.Response{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	err = utils.HandlerCSV(*directory+"/buckets.csv", "POST buckets.csv", "", "", bucketName, "")
	if err != nil {
		utils.HandlerXML(w, pkg.Response{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	utils.HandlerXML(w, pkg.Response{Status: http.StatusOK, Message: "Bucket created successfully"})
}
func ListBucketHandler(w http.ResponseWriter, r *http.Request, directory *string) {
	filePath := strings.Split(r.URL.Path[1:], "/")

	// Handle object retrieval if the path has two parts (e.g., bucket/object)
	if len(filePath) == 2 {
		utils.HandlerXML(w, GetObjectHandler(w, r, directory))
		return
	}

	bucketName := r.URL.Path[1:]

	// Check if the bucket name is provided and if not, return an error
	if len(bucketName) > 0 {
		utils.HandlerXML(w, pkg.Response{Status: http.StatusBadRequest, Message: "invalid request"})
		return
	}

	// Read the bucket metadata from the CSV file
	records, err := utils.ReadFile(*directory + "/buckets.csv")
	if err != nil {
		utils.HandlerXML(w, pkg.Response{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	// Accumulate bucket metadata for a single XML response
	var bucketList []pkg.BucketMetadata
	first := true
	for _, record := range records {
		if first {
			first = false
			continue
		}
		response := pkg.BucketMetadata{
			Name:         record[0],
			CreationTime: record[1],
			LastModified: record[2],
			Status:       record[3],
		}
		bucketList = append(bucketList, response)
	}

	// Send all accumulated bucket data at once
	utils.HandlerData(w, bucketList)
}

func DeleteBucketHandler(w http.ResponseWriter, r *http.Request, directory *string) {

	bucketName := ""
	objectKey := ""
	filePath := r.URL.Path[1:]
	path := strings.Split(filePath, "/")
	if len(path) > 0 {
		bucketName = path[0]
	}

	// checking if Bucket exist
	isBucketExist, err := utils.IsBucketExist(directory, &bucketName)
	if err != nil {
		utils.HandlerXML(w, pkg.Response{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	if !isBucketExist {
		utils.HandlerXML(w, pkg.Response{Status: http.StatusBadRequest, Message: "bucket does not exist"})
		return
	}

	if len(path) > 1 {
		objectKey = path[1]
		records, err := utils.ReadFile(*directory + "/" + bucketName + "/" + bucketName + ".csv")
		if err != nil {
			utils.HandlerXML(w, pkg.Response{Status: http.StatusInternalServerError, Message: err.Error()})
			return
		}

		found := false
		for _, record := range records {
			if record[0] == objectKey {
				found = true
				break
			}
		}
		if !found {
			utils.HandlerXML(w, pkg.Response{Status: http.StatusBadRequest, Message: "Object does not exist"})
			return
		}
		err = ObjectHandler(r, directory, &bucketName, &objectKey)
		if err != nil {
			utils.HandlerXML(w, pkg.Response{Status: http.StatusInternalServerError, Message: err.Error()})
			return
		}
		utils.HandlerXML(w, pkg.Response{Status: http.StatusOK, Message: "Object deleted successfully"})
		return
	}

	// Check if the bucket directory is empty
	records, err := utils.ReadFile(*directory + "/buckets.csv")

	for _, record := range records {
		if record[0] == bucketName {
			if record[3] == "active" {
				utils.HandlerXML(w, pkg.Response{Status: http.StatusBadRequest, Message: "Bucket is not empty"})
				return
			}
			break
		}
	}

	// Delete the bucket directory
	err = os.RemoveAll(*directory + "/" + bucketName)
	if err != nil {
		utils.HandlerXML(w, pkg.Response{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	// Remove the bucket from the CSV
	err = utils.HandlerCSV(*directory+"/buckets.csv", "DELETE buckets.csv", "", "", bucketName, "")
	if err != nil {
		utils.HandlerXML(w, pkg.Response{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	// Successfully deleted the bucket
	utils.HandlerXML(w, pkg.Response{Status: http.StatusOK, Message: "Bucket deleted successfully"})
}
