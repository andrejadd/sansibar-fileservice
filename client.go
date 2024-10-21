package main

import (
    "context"
    "io"
    "log"
    "os"
    "time"

    pb "path/to/protos/filetransfer"
    "google.golang.org/grpc"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Did not connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewFileServiceClient(conn)
    
    fileId := "video1"
    filePath := "path/to/video.mp4"
    
    // Try to resume upload if previously interrupted
    resumeOffset := resumeUpload(client, fileId)
    
    // Open the file and seek to the correct position
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

    // Start streaming chunks
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
    }

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

