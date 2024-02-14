package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

var users = make(map[string]string)

func main() {
    router := gin.Default()

    router.POST("/register", registerHandler)
    router.POST("/login", loginHandler)
    router.GET("/profile", profileHandler)

    router.Run(":8080")
}

func registerHandler(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
        return
    }

    users[user.Email] = string(hashedPassword)
    c.Status(http.StatusCreated)
}

func loginHandler(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, ok := users[user.Email]
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    c.Status(http.StatusOK)
}

func profileHandler(c *gin.Context) {
    // Получить информацию о текущем пользователе из сессии или токена
    // В этом примере просто отправляем OK, так как это демонстрационный пример
    c.Status(http.StatusOK)
}
