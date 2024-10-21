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
    if err := os.MkdirAll(dataFolder, os.ModePerm); err != nil {
        return fmt.Errorf("failed to create data folder: %v", err)
    }

    metaFilePath := filepath.Join(dataFolder, fmt.Sprintf("%s_meta.json", fileId))
    fileMeta := FileMeta{
        LastOffset:  lastOffset,
	UploadStart: time.Now(),  // TODO: needs to be retained, consider separate timestamp like lastUpload
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


func (s *FileServiceServer) Upload(stream pb.FileService_UploadServer) error {

    var offset int64 = 0
    var dataFolder string = "./meta_data"
    //var fileId string = "None"

    fmt.Printf("Client connected ..\n ")
    for {
        chunk, err := stream.Recv()

	//if fileId == "None" {
	//    fileId = chunk.FileId
	//}

	if err == io.EOF {
	    // TODO: In case remaining data was sent here (verify), it also needs to be written out before returning.	
	    //fmt.Println("Got EOF from io")
	    //err = createFileMeta(fileId, offset, dataFolder)
    	    //if err != nil {
            //	fmt.Printf("Error: %v\n", err)
    	    //}
            return stream.SendAndClose(&pb.UploadStatus{
                Success: true,
                Message: "File uploaded successfully.",
		LastOffset: offset,
            })
        }

        if err != nil {
            return err
        }

    	
	if chunk.Offset != offset {
		return fmt.Errorf("invalid offset, expected: %d, got: %d", offset, chunk.Offset)
	}

        fmt.Printf("Received chunk for file: %s at offset: %d\n", chunk.FileId, chunk.Offset)
	// TODO: 
	//   - Append chunk to file,
	//   - for resume, read chunk.FileId and chunk.Offset from disk
        
	offset = chunk.Offset + int64(len(chunk.Data))  
	err = createFileMeta(chunk.FileId, offset, dataFolder)
    	if err != nil {
        	fmt.Printf("Error: %v\n", err)
    	}

    }
}

func readFileMeta(fileId string, dataFolder string) (*FileMeta, bool) {
    metaFilePath := filepath.Join(dataFolder, fmt.Sprintf("%s_meta.json", fileId))

    file, err := os.Open(metaFilePath)
    if err != nil {
        return nil, false //fmt.Errorf("failed to open meta file: %v", err)
    }
    defer file.Close()

    var fileMeta FileMeta

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&fileMeta); err != nil {
        return nil, false //fmt.Errorf("failed to decode meta file: %v", err)
    }

    return &fileMeta, true
}

func (s *FileServiceServer) ResumeUpload(ctx context.Context, req *pb.ResumeRequest) (*pb.ResumeStatus, error) {
    
    fileMeta, exists := readFileMeta(req.FileId, "./meta_data")
    if !exists {
        fmt.Printf("readFileMeta() returned non-existent meta file for file ID: %s\n", req.FileId)
        return &pb.ResumeStatus{LastOffset: 0}, nil 
    }

    fmt.Printf("lastOffset info read for file ID %s with offset: %d\n", req.FileId, fileMeta.LastOffset)
    return &pb.ResumeStatus{LastOffset: fileMeta.LastOffset}, nil
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



