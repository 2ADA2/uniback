package deleteImage

import (
	"myapp/internal/pkg/s3DeleteImage"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DeleteImage struct {
}

func New() *DeleteImage {
	return &DeleteImage{}
}

func (e *DeleteImage) Status(c echo.Context) error {
	fileURL := c.QueryParam("url") // Получаем URL из запроса

	if fileURL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL файла не указан"})
	}

	err := s3DeleteImage.DeleteFromS3(fileURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Файл удален"})
}
