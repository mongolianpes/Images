package app

import (
	// "errors"
	pb "image-converter/proto"
	"io"
	"strconv"
)

type ImageServer struct {
	pb.ImageServiceServer
}

func NewImageServer() ImageServer {
	i := ImageServer{}
	return i
}

// создадим структуру OriginalImage которая будет хранить метаданные о наших картинках
type OriginalImage struct {
	Path      string
	Lenght    []int32
	Width     []int32
	Format    string
	Folder    string
	Watermark string
	UUID      string
}

const defaultWatermark = "watermark.png"

func (img ImageServer) DownloadImages(stream pb.ImageService_DownloadImagesServer) error {
	var images []*pb.DownloadImagesRequest
	//принимаем поток от клиента и проверяем, что он дошёл до нас без сбоев
	for {
		image, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil { // возвращаем ошибку в случае неудачного приема
			return stream.SendAndClose(&pb.DownloadImagesResponse{
				Error: err.Error(),
			})
		}
		images = append(images, image)
	}
	/*Проверяем , чтобы клиент при передаче длины и ширины картинок которые
	  необходимо создать передал одинаковое их количество, иначе будет непонятно
	  по каким размерам её нужно будет обрабатывать, и сохраняем сами картинки,
	   также проверяем чтобы в ватермарке присутствовало хотя бы значение
	   по умолчанию, а не пустое */
	var paths []OriginalImage
	for i := range images {
		if len(images[i].Info.Height) != len(images[i].Info.Width) {
			return stream.SendAndClose(&pb.DownloadImagesResponse{
				Error: "different len of lenght and width for picture " + strconv.Itoa(i),
			})
		}
		if images[i].Info.Watermark == "" { // если нам не передали делаем по умолчанию
			images[i].Info.Watermark = defaultWatermark
		}
	}

	/*наше сохранение исходных файлов, полученных от клиента,
	а также помещение всех метаданных в слайс наших структур для
	удобства работы с ними */
	paths = saveSourceFiles(images)
	if len(images) == 0 {
		return stream.SendAndClose(&pb.DownloadImagesResponse{
			Error: "no images in request",
		})
	}

	var err error
	//наложение watermark на наши изображения
	// err := watermark(paths)
	// if err != nil {
	// 	return stream.SendAndClose(&pb.DownloadImagesResponse{
	// 		Error: errors.New("path for watermark is invalid").Error(),
	// 	})
	// }
	//изменение размера картинки и сохранение уже обработанной версии
	uploadPath := resizeAndSave(paths)

	res := &pb.DownloadImagesResponse{
		StoragePath: uploadPath,
	}
	//возврат полученных путей хранения картинок для возможности просмотра
	err = stream.SendAndClose(res)
	if err != nil {
		return stream.SendAndClose(&pb.DownloadImagesResponse{
			Error: "error when receive an responce",
		})
	}

	return nil
}
