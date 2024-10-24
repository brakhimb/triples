package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func CreateBucketsCSV(bucketsCSVPath *string) error {
	record := []string{
		"Name",
		"CreationTime",
		"LastModified",
		"Status",
	}
	err := UpdateCSV(bucketsCSVPath, &record)
	if err != nil {
		return err
	}
	return nil
}

func CreateObjectsCSV(objectsCSVPath *string) error {
	record := []string{
		"ObjectKey",
		"Size",
		"ContentType",
		"LastModified",
	}
	err := UpdateCSV(objectsCSVPath, &record)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCSV(path *string, record *[]string) error {
	file, err := os.OpenFile(*path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(*record)
	if err != nil {
		return err
	}
	return nil
}

func RemoveCSV(path, bucketName *string) error {
	// Open the CSV file in write mode
	file, err := os.OpenFile(*path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer file.Close() // Make sure file is closed properly after use

	records, err := ReadFile(*path)
	if err != nil {
		return err
	}

	// Filter out the bucket to be deleted
	var updatedRecords [][]string
	for _, record := range records {
		if record[0] != *bucketName {
			updatedRecords = append(updatedRecords, record)
		}
	}

	// Create a new CSV file to write updated records
	newFile, err := os.Create(*path)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer newFile.Close() // Ensure this file is closed properly as well

	// Write the updated records back to the CSV file
	writer := csv.NewWriter(newFile)
	defer writer.Flush()

	err = writer.WriteAll(updatedRecords)
	if err != nil {
		return fmt.Errorf("error writing data: %w", err)
	}

	return nil
}

func HandlerCSV(path, request, size, contentType, bucketName, objectKey string) error {

	switch request {
	case "PUT buckets.csv":
		err := CreateBucketsCSV(&path)
		if err != nil {
			return err
		}
	case "PUT objects.csv":
		err := CreateObjectsCSV(&path)
		if err != nil {
			return err
		}
	case "POST buckets.csv":
		record := []string{
			bucketName,
			time.Now().Format(time.RFC3339),
			time.Now().Format(time.RFC3339),
			"inactive",
		}
		err := UpdateCSV(&path, &record)
		if err != nil {
			return err
		}
	case "POST objects.csv":
		record := []string{
			objectKey,
			size,
			contentType,
			time.Now().Format(time.RFC3339),
		}
		err := UpdateCSV(&path, &record)
		if err != nil {
			return err
		}
	case "DELETE buckets.csv":
		err := RemoveCSV(&path, &bucketName)
		if err != nil {
			return err
		}
	case "DELETE objects.csv":
		err := RemoveCSV(&path, &objectKey)
		if err != nil {
			return err
		}
	}
	return nil
}
