package createImage

import (
	"myapp/internal/pkg/s3uploader"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateImage struct {
}

func New() *CreateImage {
	return &CreateImage{}
}

func (e *CreateImage) Status(c echo.Context) error {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Файл не найден"})
	}

	// Открываем файл
	file, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка открытия файла"})
	}
	defer file.Close()

	// Загружаем файл в S3
	url, err := s3uploader.UploadToS3(file, fileHeader.Filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Возвращаем URL загруженного файла
	return c.JSON(http.StatusOK, map[string]string{"url": url})
}
