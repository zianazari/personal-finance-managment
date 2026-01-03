package api

import (
	"context"
	shared "expence_management/Shared"
	"expence_management/domain"
	"expence_management/repo"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

const TokenLifeTime = time.Minute * 60

// CommonTypeWithStringTime is a common type used for both incomes and expenses in which field of time is in string format
type CommonTypeWithStringTime struct {
	ID          uint    `gorm:"primary_key" json:"id"`
	CreatedAt   string  `json:"created_at"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
}

type Handler struct {
	Repo *repo.DBRepository
}

func NewHandler(repo *repo.DBRepository) *Handler {
	return &Handler{
		Repo: repo,
	}
}

func (h *Handler) OverallSummary(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	from := c.Query("from")
	to := c.Query("to")

	if from == "" || to == "" {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrBadRequest})
		return
	}

	user := h.GetUserByName(c)

	res, err := h.Repo.Repo.OverallSummary(ctx, from, to, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrInternalError})
		return
	}
	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: res})
}

func GetFieldNames(v interface{}) []string {
	t := reflect.TypeOf(v)

	// if pointer, get the element type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// must be a struct
	if t.Kind() != reflect.Struct {
		return nil
	}

	fields := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fields = append(fields, field.Name)
	}

	return fields
}

func (h *Handler) GetUserByName(c *gin.Context) *domain.User {
	u, _ := c.Get("username")
	username, _ := u.(string)
	user, _ := h.Repo.Repo.GetUserInfo(context.Background(), username)

	return user
}
