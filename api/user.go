package api

import (
	"context"
	shared "expence_management/Shared"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Signup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancel()

	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrBadRequest})
		return
	}

	// User verification information
	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrBadRequest})
		return
	}

	// check password complexity policy
	err = shared.ValidatePassword(user.Password)
	if err != nil {
		// log.Println("Invalid:", err)
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: err.Error()})
		return
	}

	// check email validity
	err = shared.IsValidEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: err.Error()})
		return
	}

	// check  if user role equals to 'user'
	role := shared.Roles[1]

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(passwordHash)

	err = h.Repo.Repo.AddUser(ctx, user.Username, user.Password, user.Email, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrAddingUser})
		return
	}
	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkAddingUser})
}

func (h *Handler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := c.ShouldBindJSON(&credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrWrongInputs})
		return
	}

	if credentials.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrWrongCredentials})
		return
	}

	user, err := h.Repo.Repo.GetUserInfo(ctx, credentials.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{shared.ErrKeyword: shared.ErrWrongCredentials})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrWrongCredentials})
		log.Println(err)
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(TokenLifeTime).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrJWTTokenGenFailed})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.TokenKey: token})
}

func (h *Handler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.UserLogoutMsg})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancel()

	var credentials struct {
		CurrentPassword   string `json:"current_password"`
		NewPassword       string `json:"new_password"`
		NewPasswordRepeat string `json:"new_password_repeat"`
	}

	err := c.ShouldBindJSON(&credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrWrongInputs})
		return
	}

	// check if newpassword and newpasswordrepeate are same
	if credentials.NewPassword != credentials.NewPasswordRepeat {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrWrongInputs})
		return
	}

	// check password complexity policy on new password
	err = shared.ValidatePassword(credentials.NewPassword)
	if err != nil {
		// log.Println("Invalid:", err)
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: err.Error()})
		return
	}

	// check if current password of user is correct
	user := h.GetUserByName(c)
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.CurrentPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrWrongPassword})
		log.Println(err)
		return
	}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(credentials.NewPassword), bcrypt.DefaultCost)

	err = h.Repo.Repo.ChangePassword(ctx, user.ID, string(passwordHash))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrInternalError})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkPassChanged})
}

func (h *Handler) ChangePasswordByAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancel()

	var credentials struct {
		UserID   uint   `json:"user_id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	err := c.ShouldBindJSON(&credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrBadRequest})
		return
	}

	// check password complexity policy on new password
	err = shared.ValidatePassword(credentials.Password)
	if err != nil {
		// log.Println("Invalid:", err)
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: err.Error()})
		return
	}

	// Get bcrypt hash of password
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)

	err = h.Repo.Repo.ChangePassword(ctx, credentials.UserID, string(passwordHash))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrInternalError})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkPassChanged})
}

func (h *Handler) UserInfo(c *gin.Context) {
	user := h.GetUserByName(c)

	type User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}

	u := &User{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: u})
}

func (h *Handler) ListUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancel()

	users, err := h.Repo.Repo.ListUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrListingUsers})
		return
	}
	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: users})
}

func (h *Handler) AddNewUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancel()

	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}

	//var user domain.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrBadRequest})
		return
	}

	// User verification information
	if user.Username == "" || user.Password == "" || user.Email == "" || user.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrBadRequest})
		return
	}
	// check if the selected role is not valid
	if !slices.Contains(shared.Roles, user.Role) {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrWrongRole})
		return
	}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(passwordHash)

	err = h.Repo.Repo.AddUser(ctx, user.Username, user.Password, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrAddingUser})
		return
	}
	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkAddingUser})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancel()

	var user struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		// Password string `json:"password"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{shared.ErrKeyword: shared.ErrBadRequest})
		return
	}

	err = h.Repo.Repo.UpdateUser(ctx, user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrUpdatingUser})
		return
	}
	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkUpdatingUser})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancel()

	intID, _ := strconv.Atoi(c.Query("id"))

	err := h.Repo.Repo.DeleteUser(ctx, uint(intID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{shared.ErrKeyword: shared.ErrDeletingUser})
		return
	}

	c.JSON(http.StatusOK, gin.H{shared.OkKeyword: shared.OkDeletingUser})
}
