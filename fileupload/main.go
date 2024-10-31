package main

import (
    "context"
    "fmt"
    "io"
    "time"
    "os"
    "sync"
    pb "github.com/andrejadd/sansibar-fileservice/fileupload/fileservice"
    "google.golang.org/grpc"
    "log"
    "net"
    "encoding/json"
    "path/filepath"
    //"github.com/google/uuid"
)

// FileInfo holds information about the data and meta files
type FileInfo struct {
    DataFilePath string
    MetaFilePath string
}
var dataFolder string = "/incoming"
var uploadQueue = make(chan FileInfo, 100) // Buffered channel with a capacity of 10

type FileServiceServer struct {
    pb.UnimplementedFileServiceServer
    uploads map[string]int64  // Track file uploads by file ID and offset
    mu      sync.Mutex        // Synchronize access to uploads map
}

type FileMeta struct {
    LastOffset   int64       `json:"lastOffset"`
    FileExt      string      `json:"fileExt"`
    FileHash     string      `json:"fileHash"`
    FileName     string      `json:"fileName"`
    UploadStart  time.Time `json:"uploadStart"`
}

func createFileMeta(file_hash string, file_ext string, lastOffset int64) error {
    if err := os.MkdirAll(dataFolder, os.ModePerm); err != nil {
        return fmt.Errorf("failed to create data folder: %v", err)
    }

    var file_name string = fmt.Sprintf("%s/%s%s", dataFolder, file_hash, file_ext)
    metaFilePath := filepath.Join(dataFolder, fmt.Sprintf("%s_meta.json", file_hash))
    fileMeta := FileMeta{
        LastOffset:  lastOffset,
	FileExt: file_ext,
	FileHash: file_hash,
	FileName: file_name,
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
    return nil
}

func updateOffset(file_hash string, newOffset int64) error {
    filename := filepath.Join(dataFolder, fmt.Sprintf("%s_meta.json", file_hash))
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open file: %v", err)
    }
    defer file.Close()

    // Read the JSON file into a struct
    var meta FileMeta
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&meta); err != nil {
        return fmt.Errorf("failed to decode JSON: %v", err)
    }

    // Update the LastOffset
    meta.LastOffset = newOffset

    // Write the updated struct back to the file
    // To overwrite, we need to open the file with write permissions
    file, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
        return fmt.Errorf("failed to open file for writing: %v", err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ") // Optional: for pretty printing
    if err := encoder.Encode(meta); err != nil {
        return fmt.Errorf("failed to encode JSON: %v", err)
    }

    return nil
}

func new_data_arrived() {
    for fileInfo := range uploadQueue {
	fmt.Printf("Uploading data file to Minio: %s\n", fileInfo.DataFilePath)
	result := PutFile(fileInfo.DataFilePath)
	os.Remove(fileInfo.DataFilePath)
	if !result {
       	   fmt.Printf("Failed to upload data file %s to object storage --- needs handling \n", fileInfo.DataFilePath)
    	}

    	fmt.Printf("Uploading meta file to Minio: %s\n", fileInfo.MetaFilePath)
 	result = PutFile(fileInfo.MetaFilePath)
        if !result {
       	   fmt.Printf("Failed to upload meta file %s to object storage --- needs handling \n", fileInfo.MetaFilePath)
        }
    }
}

func (s *FileServiceServer) Upload(stream pb.FileService_UploadServer) error {

    fmt.Printf("Client connected ..\n ")
    var offset int64 = -1 
    var file_ext string = "unknown"
    var file_name string = "unknown"
    var file_out_fd *os.File
    var file_hash string = "unknown"

    for {
        chunk, err := stream.Recv()
	if err == io.EOF {
	    fmt.Printf("[Upload()] Adding upload job to queue ..\n")
	    ready_for_storage := FileInfo{
    		DataFilePath: file_name,
    		MetaFilePath: filepath.Join(dataFolder, fmt.Sprintf("%s_meta.json", file_hash)),
	    }
	    uploadQueue <- ready_for_storage

            return stream.SendAndClose(&pb.UploadStatus{
                Success: true,
                Message: "File uploaded successfully.",
		LastOffset: offset,
            })
        }

        if err != nil {
            return err
        }
        	
	if offset == -1 {  // init the upload config on first iteration
	    file_meta, exists := readFileMeta(chunk.FileId)
	    if !exists {
		fmt.Printf("meta file does not exist")
		return nil // should return error
    	    }
	    offset = file_meta.LastOffset 
	    file_ext = file_meta.FileExt
	    file_name = file_meta.FileName
	    file_hash = chunk.FileId
	    fmt.Printf("Initialized params for file %s - offset: %d, file_ext: %s\n", file_name, offset, file_ext)

    	    fmt.Printf("Creating data file %s\n", file_name)
    	    file_out_fd, err = os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    	    if err != nil {
                return err
    	    }
    	}

	//fmt.Printf("Got chunk offset %d vs expected offset %d\n", chunk.Offset, offset)
	if chunk.Offset != offset {
	     fmt.Printf("Local and client offsets do not match!\n")
	     return fmt.Errorf("invalid offset, expected: %d, got: %d", offset, chunk.Offset)
	}

        //fmt.Printf("Writing to file .. %d\n", chunk.Offset)
	//   - for resume, read chunk.FileId and chunk.Offset from disk
	if file_out_fd == nil {
        	fmt.Printf("file_out_fd is nil\n")
	}
        _, err = file_out_fd.Write(chunk.Data)
        if err != nil {
	    fmt.Printf("Failed writing to file with error: %v", err)
            return err
        } 

	offset = chunk.Offset + int64(len(chunk.Data))  
        //fmt.Printf("Updating offset record .. %d\n", offset)
	err = updateOffset(chunk.FileId, offset)
    	if err != nil {
        	fmt.Printf("Error: %v\n", err)
    	}
    }
    defer file_out_fd.Close()
    return nil
}

func readFileMeta(file_hash string) (*FileMeta, bool) {
    metaFilePath := filepath.Join(dataFolder, fmt.Sprintf("%s_meta.json", file_hash))
    fmt.Printf("reading from meta file %s\n", metaFilePath)
    
    file, err := os.Open(metaFilePath)
    if err != nil {
	fmt.Errorf("failed to open meta file: %v", err)
        return nil, false 
    }
    defer file.Close()

    var fileMeta FileMeta

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&fileMeta); err != nil {
	fmt.Errorf("failed to decode meta file: %v", err)
        return nil, false 
    }

    return &fileMeta, true
}

func get_offset(file_hash string)(int64) {

    fileMeta, exists := readFileMeta(file_hash)
    if !exists {
	fmt.Printf("meta file does not exist")
	return 0
    }
    return fileMeta.LastOffset
}

func (s *FileServiceServer) InitiateUpload(ctx context.Context, req *pb.ResumeRequest) (*pb.ResumeStatus, error) {

    last_offset := get_offset(req.FileId)
    if last_offset == 0 {
	fmt.Printf("creating new meta file for file %s\n", req.FileId)
    	createFileMeta(req.FileId, req.FileExt, 0)
    }
    fmt.Printf("\nInitiateUpload() - meta says file hash %s with offset %d\n", req.FileId, last_offset)
    return &pb.ResumeStatus{LastOffset: last_offset}, nil 
}

func main() {

    go new_data_arrived()

    if _, err := os.Stat(dataFolder); os.IsNotExist(err) {
        err := os.MkdirAll(dataFolder, 0755) // 0755 is the permission mode
        if err != nil {
            log.Fatalf("Failed to create folder: %v", err)
	    return 
        }
    } 
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



