package handler

import (
	"flash_cards/internal"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

type handler struct {
	handler *gin.Engine
	service service
}

var (
	PORT = os.Getenv("FLASH_CARDS_GENERATOR_PORT") // порт на котором работает http сервер
)

type service interface {
	SetTask(internal.Task) (internal.Task, error)
	GetTask(id int64) (internal.Task, error)
	SearchCards(params map[string]interface{}) ([]internal.FlashCard, error)
	UpdateCards(cards []internal.FlashCard) ([]internal.FlashCard, error)
	ExelExport(params map[string]interface{}) (*xlsx.File, error)
	ExelImport(file multipart.File) ([]internal.FlashCard, error)
}

func New(service service) *handler {
	return &handler{handler: gin.Default(), service: service}
}

func (h *handler) Run() {
	h.handler.POST("/task", h.setTask)
	h.handler.GET("/task/:task_id", h.getTask)
	h.handler.GET("/cards", h.searchCards)
	h.handler.POST("/cards", h.updateCards)

	h.handler.GET("/excel/export/cards", h.excelExport)
	h.handler.POST("/excel/import/cards", h.excelImport)

	if PORT == "" {
		PORT = "8000"
	}
	httpServer := &http.Server{
		Addr:    ":" + PORT,
		Handler: h.handler,
	}
	httpServer.ListenAndServe()
}

func (h *handler) setTask(c *gin.Context) {
	task := internal.Task{}
	if err := c.BindJSON(&task); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			map[string]string{
				"details": err.Error(),
				"status":  "ERROR"},
		)
		return
	}

	task, err := h.service.SetTask(task)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			map[string]string{
				"details": err.Error(),
				"status":  "ERROR"},
		)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *handler) getTask(c *gin.Context) {
	taskId, err := strconv.ParseInt(c.Param("task_id"), 0, 64)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			map[string]string{
				"details": err.Error(),
				"status":  "ERROR"},
		)
		return
	}

	task, err := h.service.GetTask(taskId)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			map[string]string{
				"details": err.Error(),
				"status":  "ERROR"},
		)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *handler) searchCards(c *gin.Context) {
	params := make(map[string]interface{})
	queryParams := c.Request.URL.Query()
	for key, values := range queryParams {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	cards, err := h.service.SearchCards(params)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			map[string]string{
				"details": err.Error(),
				"status":  "ERROR"},
		)
		return
	}

	c.JSON(http.StatusOK, cards)
}

func (h *handler) updateCards(c *gin.Context) {
	var cards []internal.FlashCard
	if err := c.BindJSON(&cards); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			map[string]string{
				"details": err.Error(),
				"status":  "ERROR"},
		)
		return
	}

	cards, err := h.service.UpdateCards(cards)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			map[string]string{
				"details": err.Error(),
				"status":  "ERROR"},
		)
		return
	}

	c.JSON(http.StatusOK, cards)
}

func (h *handler) excelExport(c *gin.Context) {
	params := make(map[string]interface{})
	queryParams := c.Request.URL.Query()
	for key, values := range queryParams {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	// Создаем новый файл Excel
	file, err := h.service.ExelExport(params)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			map[string]string{
				"details": err.Error(),
				"status":  "ERROR"},
		)
		return
	}

	// Устанавливаем заголовок Content-Disposition для указания имени файла
	c.Header("Content-Disposition", "attachment; filename=output.xlsx")

	// Устанавливаем Content-Type для указания типа файла
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Сохраняем файл в HTTP response body
	err = file.Write(c.Writer)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			map[string]string{
				"details": err.Error(),
				"status":  "ERROR"},
		)
		return
	}
}

func (h *handler) excelImport(c *gin.Context) {
	file, _, err := c.Request.FormFile("cards_excel_import") // Имя поля, из которого ожидаем файл
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			map[string]string{
				"details": err.Error(),
				"status":  "ERROR"},
		)
	}
	defer file.Close()

	cards, err := h.service.ExelImport(file)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			map[string]string{
				"details": err.Error(),
				"status":  "ERROR"},
		)
	}

	// Преобразование данных в JSON и отправка ответа
	c.JSON(http.StatusOK, cards)
}
