package fileupload

import (
	"context"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/majid-cj/go-chat-server/util"
)

// UploadFile ...
type UploadFile struct{}

// UploadFileInterface ...
type UploadFileInterface interface {
	UploadFile(*multipart.FileHeader, multipart.File, string) (string, error)
}

var _ UploadFileInterface = &UploadFile{}

// NewUploadFile ...
func NewUploadFile() *UploadFile {
	return &UploadFile{}
}

// UploadFile ...
func (uf *UploadFile) UploadFile(fileHeader *multipart.FileHeader, file multipart.File, folder string) (string, error) {
	ctx := context.Background()
	cld, _ := cloudinary.NewFromParams(os.Getenv("CLD_NAME"), os.Getenv("CLD_KEY"), os.Getenv("CLD_SECRET"))
	fileHeader.Filename = FormatFile(fileHeader.Filename)
	src, err := fileHeader.Open()

	if err != nil {
		return "", util.GetError("general_error")
	}

	defer src.Close()

	size := fileHeader.Size

	if size > 10<<20 {
		return "", util.GetError("image_size_error")
	}

	buffer := make([]byte, size)
	_, err = file.Read(buffer)
	if err != nil {
		return "", util.GetError("general_error")
	}
	filetype := http.DetectContentType(buffer)

	if filetype != "image/jpeg" {
		return "", util.GetError("not_supported_type")
	}

	resp, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", util.GetError("general_error")
	}
	return resp.SecureURL, nil
}
