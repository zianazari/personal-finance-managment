package api

import (
	"bytes"
	"context"
	"encoding/csv"
	shared "expence_management/Shared"
	"expence_management/domain"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddIncome(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	var income struct {
		Time        uint64
		Amount      float64
		Description string
		Category    string
	}
	err := c.ShouldBind(&income)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrAddingIncome})
		return
	}

	user := h.GetUserByName(c)

	err = h.Repo.Repo.AddIncome(ctx, income.Time, income.Amount, income.Description, income.Category, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrAddingIncome})
		return
	}

	// add category if it is new
	// keep categories uniq for (name,userid) -> then insert in db
	_ = h.Repo.Repo.AddCategory(ctx, "income", income.Category, user.ID)

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkAddingIncome})
}

func (h *Handler) UpdateIncome(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	var income struct {
		ID          uint
		Time        uint64
		Amount      float64
		Description string
		Category    string
	}

	err := c.ShouldBind(&income)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrUpdatingIncome})
		return
	}

	user := h.GetUserByName(c)

	err = h.Repo.Repo.UpdateIncome(ctx, income.ID, income.Time, income.Amount, income.Description, income.Category, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrAddingIncome})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkUpdatingIncome})
}

func (h *Handler) DeleteIncome(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	user := h.GetUserByName(c)

	IncID, _ := strconv.Atoi(c.Query("id"))

	err := h.Repo.Repo.DeleteIncome(ctx, uint(IncID), user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrDeleteIncome})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkDeletingIncome})
}

func (h *Handler) Incomes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	from := c.Query("from")
	to := c.Query("to")

	user := h.GetUserByName(c)

	incomes, err := h.Repo.Repo.ListIncomes(ctx, from, to, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrListingIncomes})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: incomes})
}

func (h *Handler) ReportIncomes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	from := c.Query("from")
	to := c.Query("to")

	user := h.GetUserByName(c)

	incomes, err := h.Repo.Repo.ReportIncomes(ctx, from, to, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrListingIncomes)
		return
	}

	c.JSON(http.StatusOK, incomes)
}

func (h *Handler) ExportIncomes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	from := c.Query("from")
	to := c.Query("to")

	user := h.GetUserByName(c)

	incomes, err := h.Repo.Repo.ReportIncomes(ctx, from, to, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrListingIncomes})
		return
	}

	var newIncomes []CommonTypeWithStringTime

	for i, _ := range incomes {
		ni := CommonTypeWithStringTime{
			ID:          incomes[i].ID,
			Amount:      incomes[i].Amount,
			CreatedAt:   shared.UnixTimeToRFC339(incomes[i].CreatedAt),
			Description: incomes[i].Description,
			Category:    incomes[i].Category,
		}

		newIncomes = append(newIncomes, ni)
	}

	// Create an in-memory buffer for the CSV
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// Write header
	writer.Write(GetFieldNames(domain.Income{}))

	// Write data
	for _, i := range newIncomes {
		record := []string{
			fmt.Sprintf("%d", i.ID),
			fmt.Sprintf("%d", i.CreatedAt),
			strconv.FormatFloat(i.Amount, 'f', 2, 64),
			i.Description,
			i.Category,
		}
		writer.Write(record)
	}
	writer.Flush()

	// Set headers so browser downloads file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=incomes.csv")
	c.Data(http.StatusOK, "text/csv", buf.Bytes())
}

func (h *Handler) IncomeSummary(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	from := c.Query("from")
	to := c.Query("to")

	user := h.GetUserByName(c)

	res, err := h.Repo.Repo.IncomesSummary(ctx, from, to, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrInternalError})
		return
	}
	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: res})
}
