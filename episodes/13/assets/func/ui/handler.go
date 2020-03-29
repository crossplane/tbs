package function

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/openfaas/openfaas-cloud/sdk"
)

const bucketName = "crossplane-tbs-14"

var (
	bucketRegionKey   = "endpoint"
	bucketUserKey     = "username"
	bucketPasswordKey = "password"

	errMessage = "Unable to service request"
)

var indexTemplate = template.Must(template.ParseFiles("index.html"))

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

	svc := s3.New(mySession)

	list, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	pictures := []string{}

	for _, p := range list.Contents {
		pictures = append(pictures, fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", bucketName, region, *p.Key))
	}

	w.WriteHeader(http.StatusOK)
	if err := indexTemplate.Execute(w, pictures); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
