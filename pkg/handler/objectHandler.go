package handler

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"tripleS/pkg"
	"tripleS/pkg/utils"
)

func ObjectHandler(r *http.Request, directory, bucketName, objectKey *string) error {
	if r.Method == "PUT" {
		file, err := os.Create(*directory + "/" + *bucketName + "/" + *objectKey)
		if err != nil {
			return err
		}
		defer file.Close()

		io.Copy(file, r.Body)

		records, err := utils.ReadFile(*directory + "/buckets.csv")
		if err != nil {
			return err
		}
		for _, record := range records {
			if record[0] == *bucketName {
				path := *directory + "/buckets.csv"
				err := utils.RemoveCSV(&path, bucketName)
				if err != nil {
					return err
				}
				record[2] = time.Now().Format(time.RFC3339)
				record[3] = "active"

				err = utils.UpdateCSV(&path, &record)
				if err != nil {
					return err
				}
				break
			}
		}
	} else {
		err := os.RemoveAll(*directory + "/" + *bucketName + "/" + *objectKey)
		if err != nil {
			return err
		}

		records, err := utils.ReadFile(*directory + "/" + *bucketName + "/" + *bucketName + ".csv")
		if err != nil {
			return err
		}
		for _, record := range records {
			if record[0] == *objectKey {
				path := *directory + "/" + *bucketName + "/" + *bucketName + ".csv"
				err := utils.RemoveCSV(&path, objectKey)
				if err != nil {
					return err
				}
				break
			}
		}

		if len(records) == 2 {
			records, err := utils.ReadFile(*directory + "/buckets.csv")
			if err != nil {
				return err
			}
			for _, record := range records {
				if record[0] == *bucketName {
					path := *directory + "/buckets.csv"
					err := utils.RemoveCSV(&path, bucketName)
					if err != nil {
						return err
					}
					record[2] = time.Now().Format(time.RFC3339)
					record[3] = "inactive"

					err = utils.UpdateCSV(&path, &record)
					if err != nil {
						return err
					}
					break
				}
			}
		}
	}
	return nil
}

func GetObjectHandler(w http.ResponseWriter, r *http.Request, directory *string) pkg.Response {
	filePath := strings.Split(r.URL.Path[1:], "/")

	if len(filePath) < 2 {
		return pkg.Response{Status: http.StatusBadRequest, Message: "Bucket name or object key missing"}
	}

	bucketName := filePath[0]
	objectKey := filePath[1]

	// Check if the bucket exists
	isBucketExist, err := utils.IsBucketExist(directory, &bucketName)
	if err != nil {
		return pkg.Response{Status: http.StatusInternalServerError, Message: err.Error()}
	}
	if !isBucketExist {
		return pkg.Response{Status: http.StatusNotFound, Message: "Bucket does not exist"}
	}

	objectPath := *directory + "/" + bucketName + "/" + objectKey
	file, err := os.Open(objectPath)
	if err != nil {
		if os.IsNotExist(err) {
			return pkg.Response{Status: http.StatusNotFound, Message: "Object does not exist"}
		}
		return pkg.Response{Status: http.StatusInternalServerError, Message: err.Error()}
	}
	defer file.Close()

	// Get the file extension
	ext := filepath.Ext(objectKey)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		// If the MIME type is unknown, fall back to DetectContentType
		buffer := make([]byte, 512)
		if _, err := file.Read(buffer); err == nil {
			mimeType = http.DetectContentType(buffer)
		}
		file.Seek(0, io.SeekStart) // Reset file pointer to the beginning
	}

	// Set the correct content type
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", "inline; filename="+objectKey)

	// Write the file to the response
	if _, err := io.Copy(w, file); err != nil {
		return pkg.Response{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	return pkg.Response{Status: http.StatusOK, Message: "Object retrieved successfully"}
}
