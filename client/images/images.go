package images

import (
	"context"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "client/proto"
)

var client pb.ImageServiceClient
var imagesServiceHost = os.Getenv("IMAGES_SERVICE_HOST_GRPC_PORT")

func InitService() error {
	conn, err := grpc.NewClient("localhost:8086", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	client = pb.NewImageServiceClient(conn)
	return nil
}

func SaveImage(imageData []byte) ([]string, error) {
	stream, err := client.DownloadImages(context.Background())
	if err != nil {
		return []string{}, err
	}

	req := &pb.DownloadImagesRequest{
		Info: &pb.ImageInfo{
			Compress: "low",
			Format:   "webp",
			Width:    []int32{200},
			Height:   []int32{200},
		},
		Image: imageData,
	}

	if err := stream.Send(req); err != nil {
		return []string{}, err
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return []string{}, err
	}

	if resp.Error != "" {
		return []string{}, SaveImageError{
			Message:           resp.Error,
			CountStoragePaths: len(resp.StoragePath),
		}
	}

	if len(resp.StoragePath) == 0 {
		return []string{}, SaveImageError{
			Message:           resp.Error,
			CountStoragePaths: len(resp.StoragePath),
		}
	}

	return resp.StoragePath, nil
}

type SaveImageError struct {
	Message           string
	CountStoragePaths int
}

func (saveImage SaveImageError) Error() string {
	if saveImage.CountStoragePaths == 0 {
		return "Не были получены пути к изображениям"
	}

	return saveImage.Message
}
