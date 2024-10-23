package pkg

type Metadata struct {
	ObjectKey    string `xml:"object-key"`
	Size         int    `xml:"size"`
	ContentType  string `xml:"content_type"`
	LastModified string `xml:"last_modified"`
}

type Bucket struct {
	Name     string            `xml:"name"`
	Data     map[string]string `xml:"data"`
	Metadata Metadata          `xml:"metadata"`
}

type BucketMetadata struct {
	Name         string `xml:"name"`
	CreationTime string `xml:"creation_time"`
	LastModified string `xml:"last_modified"`
	Status       string `xml:"status"`
}

type BucketStore struct {
	Buckets map[string]Bucket `xml:"buckets"`
}

type Response struct {
	Status  int    `xml:"status"`
	Message string `xml:"message"`
}
