package main

import (
    "context"
    "io"
    "log"
    "os"
    "fmt"
    "time"
    "crypto/sha256"
    //"github.com/google/uuid"
    pb "github.com/andrejadd/sansibar-fileservice/fileservice"
    "google.golang.org/grpc"
)

func generateFileID(filePath string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", fmt.Errorf("failed to open file: %v", err)
    }
    defer file.Close()

    hash := sha256.New()

    if _, err := io.Copy(hash, file); err != nil {
        return "", fmt.Errorf("failed to hash file: %v", err)
    }

    fileHash := fmt.Sprintf("%x", hash.Sum(nil))

    return fileHash, nil
}

func main() {

    if len(os.Args) < 2 {
        fmt.Println("Usage: go run client.go <filename>")
        return
    }

    filePath := os.Args[1]
    fileId, err := generateFileID(filePath)
    if err != nil {
        fmt.Printf("Error generating file ID: %v\n", err)
        return
    }
    fmt.Printf("Generated File ID (SHA-256): %s\n", fileId)

    
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Did not connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewFileServiceClient(conn)
    resumeOffset := resumeUpload(client, fileId)
    fmt.Println("resumeOffset is %d", resumeOffset)

    file, err := os.Open(filePath)
    if err != nil {
        log.Fatalf("Could not open file: %v", err)
    }
    defer file.Close()

    if resumeOffset > 0 {
        _, err = file.Seek(resumeOffset, io.SeekStart)
        if err != nil {
            log.Fatalf("Could not seek to offset: %v", err)
        }
    }

    stream, err := client.Upload(context.Background())
    if err != nil {
        log.Fatalf("Could not initiate upload: %v", err)
    }

    const chunkSize = 1024 * 1024  // 1 MB chunks
    buffer := make([]byte, chunkSize)
    offset := resumeOffset

    for {
        n, err := file.Read(buffer)
        if err == io.EOF {
	    // TODO: if there is actually data left, need to send it before leaving this loop.
	    fmt.Println("EOF reached when reading from file")
            break
        }
        if err != nil {
            log.Fatalf("Error reading file: %v", err)
        }

        chunk := &pb.Chunk{
            FileId: fileId,
            Offset: offset,
            Data:   buffer[:n],
        }

        err = stream.Send(chunk)
        if err != nil {
            log.Fatalf("Error sending chunk: %v", err)
        }

        offset += int64(n)
	fmt.Println("uploaded chunk offset %d", offset)
    }
    fmt.Println("finished upload for loop!")
    
    status, err := stream.CloseAndRecv()
    if err != nil {
        log.Fatalf("Error receiving response: %v", err)
    }

    log.Printf("Upload status: %v", status.Message)
}

func resumeUpload(client pb.FileServiceClient, fileId string) int64 {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
    defer cancel()

    resp, err := client.ResumeUpload(ctx, &pb.ResumeRequest{FileId: fileId})
    if err != nil {
        log.Fatalf("Failed to resume upload: %v", err)
    }

    return resp.LastOffset
}

