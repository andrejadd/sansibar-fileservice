package main

import (
    "context"
    "fmt"
    "io"
    "sync"
    "time"

    pb "fileservice"
    "google.golang.org/grpc"
    "log"
    "net"
)

type FileServiceServer struct {
    pb.UnimplementedFileServiceServer
    uploads map[string]int64  // Track file uploads by file ID and offset
    mu      sync.Mutex        // Synchronize access to uploads map
}

func (s *FileServiceServer) Upload(stream pb.FileService_UploadServer) error {
    for {
        chunk, err := stream.Recv()
        if err == io.EOF {
            // File transfer completed
            return stream.SendAndClose(&pb.UploadStatus{
                Success: true,
                Message: "File uploaded successfully.",
                LastOffset: chunk.Offset,
            })
        }
        if err != nil {
            return err
        }

        s.mu.Lock()
        if chunk.Offset != s.uploads[chunk.FileId] {
            s.mu.Unlock()
            return fmt.Errorf("invalid offset, expected: %d, got: %d", s.uploads[chunk.FileId], chunk.Offset)
        }
        // Simulate saving chunk to disk or storage
        fmt.Printf("Received chunk for file: %s at offset: %d\n", chunk.FileId, chunk.Offset)
        s.uploads[chunk.FileId] = chunk.Offset + int64(len(chunk.Data))  // Update the offset
        s.mu.Unlock()
    }
}

func (s *FileServiceServer) ResumeUpload(ctx context.Context, req *pb.ResumeRequest) (*pb.ResumeStatus, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    lastOffset, exists := s.uploads[req.FileId]
    if !exists {
        return &pb.ResumeStatus{LastOffset: 0}, nil  // Start from beginning
    }

    return &pb.ResumeStatus{LastOffset: lastOffset}, nil
}

func main() {
    server := grpc.NewServer()
    pb.RegisterFileServiceServer(server, &FileServiceServer{
        uploads: make(map[string]int64),
    })

    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    log.Println("Server is running...")
    if err := server.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}



