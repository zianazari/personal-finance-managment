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

func (h *Handler) AddExpense(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	var expense struct {
		Time        uint64
		Amount      float64
		Description string
		Category    string
	}
	err := c.ShouldBind(&expense)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrAddingExpense})
		return
	}

	user := h.GetUserByName(c)

	err = h.Repo.Repo.AddExpense(ctx, expense.Time, expense.Amount, expense.Description, expense.Category, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrAddingExpense})
		return
	}

	// add category if it is new
	_ = h.Repo.Repo.AddCategory(ctx, "expense", expense.Category, user.ID)

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkAddingExpense})
}

func (h *Handler) UpdateExpense(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	var expense struct {
		ID          uint
		Time        uint64
		Amount      float64
		Description string
		Category    string
	}

	err := c.ShouldBind(&expense)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrUpdatingExpense})
		return
	}

	user := h.GetUserByName(c)

	err = h.Repo.Repo.UpdateExpense(ctx, expense.ID, expense.Time, expense.Amount, expense.Description, expense.Category, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrAddingExpense})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkUpdatingExpense})
}

func (h *Handler) DeleteExpense(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	user := h.GetUserByName(c)

	ExpInt, _ := strconv.Atoi(c.Query("id"))
	err := h.Repo.Repo.DeleteExpense(ctx, uint(ExpInt), user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrDeletingExpense})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkDeletingExpense})
}

func (h *Handler) Expenses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	from := c.Query("from")
	to := c.Query("to")

	user := h.GetUserByName(c)

	expenses, err := h.Repo.Repo.ListExpenses(ctx, from, to, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrListingExpenses})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: expenses})
}

func (h *Handler) ReportExpenses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	from := c.Query("from")
	to := c.Query("to")

	user := h.GetUserByName(c)

	expenses, err := h.Repo.Repo.ReportExpenses(ctx, from, to, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrListingExpenses)
		return
	}

	// var newExpenses []CommonTypeWithStringTime

	// for i, _ := range expenses {
	// 	ni := CommonTypeWithStringTime{
	// 		ID:          expenses[i].ID,
	// 		Amount:      expenses[i].Amount,
	// 		CreatedAt:   shared.UnixTimeToRFC339(expenses[i].CreatedAt),
	// 		Description: expenses[i].Description,
	// 		Category:    expenses[i].Category,
	// 	}

	// 	newExpenses = append(newExpenses, ni)
	// }

	c.JSON(http.StatusOK, expenses)
}

func (h *Handler) ExportExpenses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	from := c.Query("from")
	to := c.Query("to")

	user := h.GetUserByName(c)

	expenses, err := h.Repo.Repo.ReportExpenses(ctx, from, to, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrListingIncomes})
		return
	}

	var newExpenses []CommonTypeWithStringTime

	for i, _ := range expenses {
		ni := CommonTypeWithStringTime{
			ID:          expenses[i].ID,
			Amount:      expenses[i].Amount,
			CreatedAt:   shared.UnixTimeToRFC339(expenses[i].CreatedAt),
			Description: expenses[i].Description,
			Category:    expenses[i].Category,
		}

		newExpenses = append(newExpenses, ni)
	}

	// Create an in-memory buffer for the CSV
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// Write header
	writer.Write(GetFieldNames(domain.Income{}))

	// Write data
	for _, i := range newExpenses {
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
	c.Header("Content-Disposition", "attachment; filename=expenses.csv")
	c.Data(http.StatusOK, "text/csv", buf.Bytes())
}

func (h *Handler) ExpenseSummary(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	from := c.Query("from")
	to := c.Query("to")

	user := h.GetUserByName(c)

	res, err := h.Repo.Repo.ExpensesSummary(ctx, from, to, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrInternalError})
		return
	}
	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: res})
}
