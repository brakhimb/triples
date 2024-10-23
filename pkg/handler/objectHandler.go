package handler

import (
	"io"
	"net/http"
	"os"
	"time"
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
