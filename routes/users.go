package routes

// Path: routes/users.go
import (
    "fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jalong4/stock-service-go/auth"
	"github.com/jalong4/stock-service-go/config"
	"github.com/jalong4/stock-service-go/models"

    "go.mongodb.org/mongo-driver/bson/primitive"

	"golang.org/x/crypto/bcrypt"
)

// LoginRequest defines the structure of the login request payload
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}


// TokenDetails holds the token information
type TokenDetails struct {
    AccessToken  string `json:"accessToken"`
    RefreshToken string `json:"refreshToken"`
    IssuedAt     int64  `json:"iat"`
    ExpiresAt    int64  `json:"exp"`
}

// AuthResponse defines the structure of the authentication response
type AuthResponse struct {
    Success bool          `json:"success"`
    User    models.User   `json:"user"`
    Auth    TokenDetails  `json:"auth"`
}


// LoginHandler handles the POST request for user login
func LoginHandler(c *gin.Context) {
    var loginReq LoginRequest
    if err := c.BindJSON(&loginReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // Here, simulate database call
    user, err := config.GetUserByEmail(loginReq.Email) // Replace with actual DB call
    if err != nil || user == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    accessToken, refreshToken, err := auth.CreateTokens(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tokens"})
        return
    }

    resp := AuthResponse{
        Success: true,
        User: *user,
        Auth: TokenDetails{
            AccessToken:  accessToken,
            RefreshToken: refreshToken,
            IssuedAt:     time.Now().Unix(),
            ExpiresAt:    time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
        },
    }
    c.JSON(http.StatusOK, gin.H{"response": resp})
}

func RegisterUserHandler(c *gin.Context) {
    log.Println("RegisterUserHandler")
    var req models.RegistrationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
        return
    }

    log.Println("Request:", req)

    // Check required fields and validate passwords
    if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" || req.Password2 == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Please fill in all fields"})
        return
    }
    if req.Password != req.Password2 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
        return
    }

    // Check if the user already exists
    existingUser, err := config.GetUserByEmail(req.Email)
    if err == nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Email %s already exists", req.Email)})
        return
    }

    if existingUser != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("User %s already exists", req.Email)})
        return
    }

    // Hash the password using bcrypt
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
        return
    }

    // Create new user
    newUser := models.User{
        FirstName:       req.FirstName,
        LastName:        req.LastName,
        Email:           req.Email,
        Password:        string(hashedPassword),
        Timezone:        req.Timezone,
        Date:            time.Now(),
        ProfileImageURL: req.ProfileImageURL,
    }

    // Insert the new user into the database
    result, err := config.InsertUser(&newUser)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
        return
    }

    // Retrieve the inserted document ID
    newID := result.InsertedID.(primitive.ObjectID)
    newUser.ID = newID

    // Generate JWT tokens
    accessToken, refreshToken, accessClaims, err := auth.GenerateTokens(&newUser)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tokens"})
        return
    }

    // Return the response
    response := gin.H{
        "response": gin.H{
            "user": gin.H{
                "_id":             newID.Hex(),
                "firstName":       newUser.FirstName,
                "lastName":        newUser.LastName,
                "email":           newUser.Email,
                "password":        newUser.Password,
                "timezone":        newUser.Timezone,
                "date":            newUser.Date.Format(time.RFC3339), // Ensure Date is set to the current time or appropriately
            },
            "auth": gin.H{
                "accessToken":  accessToken,
                "refreshToken": refreshToken,
                "accessTokenProperties": gin.H{
                    "_id":   newID.Hex(),
                    "email": newUser.Email,
                    "iat":   accessClaims.IssuedAt,
                    "exp":   accessClaims.ExpiresAt,
                },
            },
        },
    }

    c.JSON(http.StatusOK, response)
}

func DeleteUserHandler(c *gin.Context) {
    id := c.Param("_id")

    // Delete the user from the database
    user, err := config.DeleteUserByID(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
        return
    }

    // Respond with a success message
    c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s deleted successfully!", user.Email)})
}

func UpdateUserHandler(c *gin.Context) {
    id := c.Param("_id")
    var updatedUser models.User
    if err := c.ShouldBindJSON(&updatedUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
        return
    }
    updatedUser.Password = string(hashedPassword)

    // Update the user in the database
    result, err := config.UpdateUserByID(id, updatedUser)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s updated successfully!\n%v", id, result)})
}


// GetAllUsers retrieves all users
func GetAllUsers(c *gin.Context) {
    users, err := config.GetAllUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    response := gin.H{
        "count": len(users),
        "users": users,
    }

    c.JSON(http.StatusOK, response)
}


// GetUserByID handles the GET request to retrieve a user by their ID
func GetUserByID(c *gin.Context) {
    idParam := c.Param("_id")
    log.Println("ID Param:", idParam)
    user, err := config.GetUserByID(idParam)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, user)
}