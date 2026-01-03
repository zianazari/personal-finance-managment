package api

import (
	"context"
	shared "expence_management/Shared"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	var category struct {
		Type string
		Name string
	}

	err := c.ShouldBind(&category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrAddingCategory})
		return
	}

	user := h.GetUserByName(c)

	err = h.Repo.Repo.AddCategory(ctx, category.Type, category.Name, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrAddingCategory})
		return
	}
	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkAddingCategory})
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	var category struct {
		ID   uint
		Type string
		Name string
	}

	err := c.ShouldBind(&category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrUpdatingCategory})
		return
	}

	user := h.GetUserByName(c)

	err = h.Repo.Repo.UpdateCategory(ctx, category.ID, category.Type, category.Name, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrUpdatingCategory})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkUpdatingCategory})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	user := h.GetUserByName(c)

	IncID, _ := strconv.Atoi(c.Query("id"))

	err := h.Repo.Repo.DeleteCategory(ctx, uint(IncID), user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrDeleteCategory})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkDeletingCategory})
}

func (h *Handler) Categories(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	user := h.GetUserByName(c)
	_type, ok := c.GetQuery("type")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrBadRequest})
		return
	}

	categories, err := h.Repo.Repo.ListCategories(ctx, _type, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrListingCategories})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: categories})
}
