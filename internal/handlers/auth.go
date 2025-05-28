package handlers

import (
	"datingApp/pkg/model"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler представляет обработчик для авторизации пользователя
type AuthHandler struct {
	DB *gorm.DB // Ссылка на базу данных для выполнения запросов
}

// Секретный ключ для подписи JWT токенов
var JWTSecretKey = []byte(os.Getenv("JWT_SECRET"))

// hashPassword хеширует пароль перед его сохранением в базу данных
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPaswwordHash проверяет, соответствует ли переданный пароль хешированному паролю
func checkPaswwordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken создает JWT-токен для пользователя с указанием срока действия
func GenerateToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Установка времени действия токена на 24 часа
	claims := jwt.MapClaims{
		"userID": user.ID,
		"exp":    expirationTime.Unix(), // Время окончания действия токена
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecretKey)
}

// TokenAuthMiddleware проверяет, является ли переданный токен действительным
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем токен из заголовка Authorization
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is messing"})
			c.Abort()
			return
		}

		// Удаляем префикс "Bearer " из токена
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Проверяем, является ли токен действительным
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("Invalid signing method", jwt.ValidationErrorMalformed)
			}
			return JWTSecretKey, nil
		})

		// Проверка корректности токена
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Извлечение и проверка claims из токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
			c.Abort()
			return
		}

		// Извлечение userID из claims и установка его в контексте запроса
		userID, ok := claims["userID"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid user ID in token"})
			c.Abort()
			return
		}

		c.Set("userID", int(userID))
	}
}

// Register регистрирует нового пользователя
func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	// Проверка входных данных
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка уникальности email
	var user model.User
	if err := h.DB.Where("email = ?", input.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	// Хеширование пароля
	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Создание нового пользователя
	user = model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Создание профиля для нового пользователя
	profile := model.Profile{
		UserID:    user.ID,
		Bio:       "",
		Interests: "",
		Age:       "",
		Gender:    "",
		Location:  "",
	}

	if err := h.DB.Create(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed ti create profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login входит в систему
func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Проверка корректности входных данных
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Поиск пользователя по email
	var user model.User
	if err := h.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Проверка пароля
	if !checkPaswwordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Генерация JWT токена
	token, err := GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "message": "Login successful"})
}

// AuthRoutes задает маршруты для аутентификации
func (h *AuthHandler) AuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register) // Маршрут для регистрации
		auth.POST("/login", h.Login)       // Маршрут для входа
	}
}
