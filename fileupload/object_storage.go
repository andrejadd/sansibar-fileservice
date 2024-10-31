package main 

import (
   "fmt"
    "context"
    "log"
    "path"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

func PutFile(file_path string) (status bool) {

    endpoint := "minio:9000"  
    accessKeyID := "QYyRTG5zKyWZdL82QTcB"   
    secretAccessKey := "pdGByveGweKMZf31zZg5qeWsizg68F7t5nqbGEac"
    useSSL := false

    minioClient, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
        Secure: useSSL,
    })
    if err != nil {
        log.Fatalf("Error creating MinIO client: %v", err)
	return false
    }

    bucketName := "incoming"
    objectName := path.Base(file_path)
    contentType := "application/octet-stream"

    _, err = minioClient.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})
    if err == nil {
        fmt.Printf("Object %s already exists. Skipping upload.\n", objectName)
        return false
    }
    info, err := minioClient.FPutObject(context.Background(), bucketName, objectName, file_path, minio.PutObjectOptions{
        ContentType: contentType,
    })
    if err != nil {
        log.Fatalf("Error uploading file: %v", err)
	return false
    }

    log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
    return true
}

