package main

import (
    "context"
    "fmt"
    "io"
    "time"
    "os"
    "sync"
    pb "github.com/andrejadd/sansibar-fileservice/fileservice"
    "google.golang.org/grpc"
    "log"
    "net"
    "encoding/json"
    "path/filepath"
)

type FileServiceServer struct {
    pb.UnimplementedFileServiceServer
    uploads map[string]int64  // Track file uploads by file ID and offset
    mu      sync.Mutex        // Synchronize access to uploads map
}

type FileMeta struct {
    LastOffset   int64       `json:"lastOffset"`
    UploadStart  time.Time `json:"uploadStart"`
}

func createFileMeta(fileId string, lastOffset int64, dataFolder string) error {
    // Ensure the data folder exists, create it if it doesn't
    if err := os.MkdirAll(dataFolder, os.ModePerm); err != nil {
        return fmt.Errorf("failed to create data folder: %v", err)
    }

    metaFilePath := filepath.Join(dataFolder, fmt.Sprintf("%s_meta.json", fileId))
    fileMeta := FileMeta{
        LastOffset:  lastOffset,
        UploadStart: time.Now(),
    }

    f, err := os.Create(metaFilePath)
    if err != nil {
        return fmt.Errorf("failed to create meta file: %v", err)
    }
    defer f.Close()

    encoder := json.NewEncoder(f)
    encoder.SetIndent("", "  ") // For pretty printing the JSON
    if err := encoder.Encode(fileMeta); err != nil {
        return fmt.Errorf("failed to write to meta file: %v", err)
    }

    fmt.Printf("Meta file created: %s\n", metaFilePath)
    return nil
}

func readFileMeta(fileId string, dataFolder string) (*FileMeta, error) {
    metaFilePath := filepath.Join(dataFolder, fmt.Sprintf("%s_meta.json", fileId))

    file, err := os.Open(metaFilePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open meta file: %v", err)
    }
    defer file.Close()

    var fileMeta FileMeta

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&fileMeta); err != nil {
        return nil, fmt.Errorf("failed to decode meta file: %v", err)
    }

    return &fileMeta, nil
}

func (s *FileServiceServer) Upload(stream pb.FileService_UploadServer) error {

    var local_offset int64 = 0
    var dataFolder string = "./meta_data"

    fmt.Printf("Client connected ..\n ")
    for {
        chunk, err := stream.Recv()
	//fmt.Println(" 1- the chunk offset set was: %d", s.uploads[chunk.FileId])
	fmt.Println(" 1- the chunk offset set was: %d", local_offset)

	if err == io.EOF {
	    // TODO: In case remaining data was sent here (verify), it also needs to be written out before returning.	
	    fmt.Println("Got EOF from io")
            return stream.SendAndClose(&pb.UploadStatus{
                Success: true,
                Message: "File uploaded successfully.",
		//LastOffset: s.uploads[chunk.FileId], // non-valid memory when this code is reached
                //LastOffset: chunk.Offset,  // not in payload, when err is io.EOF
		LastOffset: local_offset,
            })
        }

        if err != nil {
            return err
        }

    	err = createFileMeta(chunk.FileId, chunk.Offset, dataFolder)
    	if err != nil {
        	fmt.Printf("Error: %v\n", err)
    	}

	//s.mu.Lock()
        //if chunk.Offset != s.uploads[chunk.FileId] {
        //    s.mu.Unlock()
        //    return fmt.Errorf("invalid offset, expected: %d, got: %d", s.uploads[chunk.FileId], chunk.Offset)
        //}

	if chunk.Offset != local_offset {
		return fmt.Errorf("invalid offset, expected: %d, got: %d", local_offset, chunk.Offset)
	}

        fmt.Printf("Received chunk for file: %s at offset: %d\n", chunk.FileId, chunk.Offset)
	// TODO: 
	//   - Append chunk to file,
	//   - write chunkOffset to meta data (to disk)
	//   - for resume, read chunk.FileId and chunk.Offset from disk
        //s.uploads[chunk.FileId] = chunk.Offset + int64(len(chunk.Data))  // Update the offset
        local_offset = chunk.Offset + int64(len(chunk.Data))  
	fmt.Println(" 2- the chunk offset set was: %d", local_offset)
        //s.mu.Unlock()
    }
}


func (s *FileServiceServer) ResumeUpload(ctx context.Context, req *pb.ResumeRequest) (*pb.ResumeStatus, error) {
    
    lastOffset, err := readFileMeta(req.FileId, "./meta_data")
    var exists = true
    if err != nil {
    	exists = false
        fmt.Printf("Error reading meta file: %v\n", err)
        return
    }
    
    if !exists {
        return &pb.ResumeStatus{LastOffset: 0}, nil 
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



