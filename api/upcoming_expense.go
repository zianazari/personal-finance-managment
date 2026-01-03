package api

import (
	"context"
	shared "expence_management/Shared"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddUpcomingExpense(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	var expense struct {
		Time        uint64
		Amount      float64
		Description string
		Category    string
		Repeat      string
	}
	err := c.ShouldBind(&expense)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrAddingExpense})
		return
	}

	user := h.GetUserByName(c)

	if expense.Repeat == "once" {
		res := h.addUpcomingExpenseMuliple(ctx, expense.Time, expense.Amount, expense.Description, expense.Category, user.ID)
		if res != nil {
			c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrAddingExpense})
			return
		}
	} else if expense.Repeat == "monthly" {

		beginingOfUserSelectedDay := time.Unix(int64(expense.Time), 0)

		// year, month, day := time.Now().Date()
		// if month == time.November && day == 10 {
		// 	fmt.Println("Happy Go day!")
		// }

		// beginingOfToday := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

		var times [12]time.Time

		times[0] = beginingOfUserSelectedDay

		var i = 0
		for i = 1; i <= 11; i++ {
			oneMonthLater := beginingOfUserSelectedDay.AddDate(0, i, 0)
			times[i] = oneMonthLater
		}

		for _, t := range times {
			res := h.addUpcomingExpenseMuliple(ctx, uint64(t.Unix()), expense.Amount, expense.Description, expense.Category, user.ID)
			if res != nil {
				c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrAddingExpense})
				return
			}
		}

	}

	// err = h.Repo.Repo.AddUpcomingExpense(ctx, expense.Time, expense.Amount, expense.Description, expense.Category, user.ID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrAddingExpense})
	// 	return
	// }

	// // add category if it is new
	// _ = h.Repo.Repo.AddCategory(ctx, "expense", expense.Category, user.ID)

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkAddingExpense})
}

func (h *Handler) addUpcomingExpenseMuliple(ctx context.Context, time uint64, amount float64, description, category string, userID uint) error {
	err := h.Repo.Repo.AddUpcomingExpense(ctx, time, amount, description, category, userID)
	if err != nil {
		return err
	}

	// add category if it is new
	_ = h.Repo.Repo.AddCategory(ctx, "expense", category, userID)

	return nil
}

func (h *Handler) UpdateUpcomingExpense(c *gin.Context) {
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

	err = h.Repo.Repo.UpdateUpcomingExpense(ctx, expense.ID, expense.Time, expense.Amount, expense.Description, expense.Category, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrAddingExpense})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkUpdatingExpense})
}

func (h *Handler) DeleteUpcomingExpense(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	user := h.GetUserByName(c)

	ExpInt, _ := strconv.Atoi(c.Query("id"))
	err := h.Repo.Repo.DeleteUpcomingExpense(ctx, uint(ExpInt), user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrDeletingExpense})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkDeletingExpense})
}

func (h *Handler) PayUpcomingExpense(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	user := h.GetUserByName(c)

	ExpInt, _ := strconv.Atoi(c.Query("id"))
	err := h.Repo.Repo.PayUpcomingExpense(ctx, uint(ExpInt), user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrPayingExpense})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkPayingExpense})
}

func (h *Handler) UpcomingExpenses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	from := c.Query("from")
	to := c.Query("to")

	user := h.GetUserByName(c)

	expenses, err := h.Repo.Repo.ListUpcomingExpenses(ctx, from, to, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrListingExpenses})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: expenses})
}

func (h *Handler) UpcomingExpenseSummary(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	from := c.Query("from")
	to := c.Query("to")

	user := h.GetUserByName(c)

	res, err := h.Repo.Repo.UpcomingExpensesSummary(ctx, from, to, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrInternalError})
		return
	}
	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: res})
}
