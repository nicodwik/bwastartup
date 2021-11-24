package handler

import (
	"bwastartup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) Index(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": err})
		return
	}

	c.HTML(http.StatusOK, "user_index.html", gin.H{"users": users})
}

func (h *userHandler) Create(c *gin.Context) {
	c.HTML(http.StatusOK, "user_create.html", nil)
}

func (h *userHandler) Store(c *gin.Context) {
	var input user.FormCreateUserInput

	err := c.ShouldBind(&input)
	if err != nil {
		input.Error = err
		c.HTML(http.StatusOK, "user_create.html", input)
		return
	}

	registerInput := user.RegisterUserInput{}
	registerInput.Name = input.Name
	registerInput.Occupation = input.Occupation
	registerInput.Email = input.Email
	registerInput.Password = input.Password

	_, err = h.userService.RegisterUser(registerInput)
	if err != nil {
		input.Error = err
		c.HTML(http.StatusOK, "user_create.html", input)
		return
	}

	c.Redirect(http.StatusFound, "/users")
}

func (h *userHandler) Edit(c *gin.Context) {
	id := c.Param("id")
	paramsId, _ := strconv.Atoi(id)

	user, err := h.userService.GetUserByID(paramsId)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err})
		return
	}

	c.HTML(http.StatusOK, "user_edit.html", gin.H{"user": user})
}

func (h *userHandler) Update(c *gin.Context) {
	id := c.Param("id")
	paramsId, _ := strconv.Atoi(id)

	var input user.FormUpdateUserInput
	input.Id = paramsId
	err := c.ShouldBind(&input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "user_edit.html", gin.H{"error": err, "user": input})
		return
	}

	_, err = h.userService.Updateuser(input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "user_edit.html", gin.H{"error": err})
		return
	}

	c.Redirect(http.StatusFound, "/users")
}

func (h *userHandler) EditAvatar(c *gin.Context) {
	id := c.Param("id")

	c.HTML(http.StatusOK, "user_avatar.html", gin.H{"Id": id})
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	id := c.Param("id")
	idParam, _ := strconv.Atoi(id)

	file, err := c.FormFile("avatar")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err})
		return
	}

	path := fmt.Sprintf("images/%d-%s", idParam, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err})
		return
	}

	_, err = h.userService.SaveAvatar(idParam, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err})
		return
	}

	c.Redirect(http.StatusFound, "/users")
}
