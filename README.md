# go-getter-example

## for Git
./go-getter-example "/tmp/go-getter-git" "git@github.com:svennela/azure-storage-example.git"

./go-getter-example "/tmp/go-getter-http" "https://github.com/svennela/azure-storage-example.git"

./go-getter-example "/tmp/go-getter-http" "git@github.com:svennela/azure-storage-example.git"

## GCP - GCS 

./go-getter-example "/tmp/go-getter-gcp" "gs://gs-bucketname/filename"

## AWS - S3 

AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY

./go-getter-example "/tmp/go-getter-s3" "bucketname.s3.amazonaws.com/filename"