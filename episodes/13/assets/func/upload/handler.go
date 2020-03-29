package function

import (
	"fmt"
	"hash/fnv"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/openfaas/openfaas-cloud/sdk"
)

const bucketName = "crossplane-tbs-14"

var (
	bucketRegionKey   = "endpoint"
	bucketUserKey     = "username"
	bucketPasswordKey = "password"

	errMessage = "Unable to service request"
)

var (
	region   string
	user     string
	password string
)

func read(key string) string {
	val, err := sdk.ReadSecret(key)
	if err != nil {
		panic(err)
	}
	return val
}

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprint(h.Sum32())
}

func init() {
	region = read(bucketRegionKey)
	user = read(bucketUserKey)
	password = read(bucketPasswordKey)
}

// Handle a function invocation
func Handle(w http.ResponseWriter, r *http.Request) {
	var err error

	mySession := session.Must(session.NewSession(
		&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(user, password, ""),
		},
	))

	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer file.Close()

	manager := s3manager.NewUploader(mySession)

	imageName := fmt.Sprintf("%s.%s", hash(time.Now().String()), "jpg")

	out, err := manager.Upload(&s3manager.UploadInput{
		ACL:         aws.String(s3.BucketCannedACLPublicRead),
		Bucket:      aws.String(bucketName),
		Key:         aws.String(imageName),
		Body:        file,
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("file uploaded to, %s\n", out.Location)))
}
