package utils

func IsValid(name string) bool {
	if len(name) < 3 || len(name) > 63 {
		return false
	}
	if name[0] == '.' || name[0] == '-' || name[len(name)-1] == '.' || name[len(name)-1] == '-' {
		return false
	}
	return true
}

func IsBucketExist(directory, bucketName *string) (bool, error) {
	records, err := ReadFile(*directory + "/buckets.csv")
	if err != nil {
		return false, nil
	}
	for _, record := range records {
		if record[0] == *bucketName {
			return true, nil
		}
	}
	return false, nil
}
