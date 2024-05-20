package auth

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/jalong4/stock-service-go/models"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

type AccessClaims struct {
    ID    string `json:"_id"`
    Email string `json:"email"`
    jwt.StandardClaims
}

type RefreshClaims struct {
    ID    string `json:"_id"`
    Email string `json:"email"`
    jwt.StandardClaims
}

// CreateToken generates new access and refresh tokens for the user
func CreateTokens(user *models.User) (string, string, error) {
    accessTokenSigned, refreshTokenSigned, _, err := createTokensWithAtClaim(user)
    if err != nil {
        return "", "", err
    }

    return accessTokenSigned, refreshTokenSigned, nil
}

func GenerateTokens(user *models.User) (string, string, *AccessClaims, error) {
    accessToken, refreshTokeng, accessClaims, err := createTokensWithAtClaim(user)  
    if err != nil {
        return "", "", nil, err
    }
    return accessToken, refreshTokeng, accessClaims, nil
}

// This function returns the middleware for JWT validation
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        const Bearer_schema = "Bearer "
        authHeader := c.GetHeader("Authorization")
        if !strings.HasPrefix(authHeader, Bearer_schema) {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header provided"})
            return
        }

        tokenString := authHeader[len(Bearer_schema):]
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }

            return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
        })

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            // Pass the processing to the next middleware or handler
            c.Set("userID", claims["user_id"])
        } else {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            return
        }
    }
}


// Private functions

func getAccessSecret() ([]byte, error) {
    accessSecret := os.Getenv("ACCESS_TOKEN_SECRET")
    if accessSecret == "" {
        return nil, fmt.Errorf("ACCESS_TOKEN_SECRET not set")
    }
    return []byte(accessSecret), nil
}

func getRefreshSecret() ([]byte, error) {
    refreshSecret := os.Getenv("REFRESH_TOKEN_SECRET")
    if refreshSecret == "" {
        return nil, fmt.Errorf("REFRESH_TOKEN_SECRET not set")
    }
    return []byte(refreshSecret), nil
}

func createTokensWithAtClaim(user *models.User) (string, string, *AccessClaims, error) {
    var err error
    // Creating Access Token
    atClaims := AccessClaims{
        ID:    user.ID.String(),
        Email: user.Email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(), // Token expires after 30 days
        },
    }
    at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

    accessToken, err := getAccessSecret()
    if err != nil {
        return "", "", nil, err
    }

    accessTokenSigned, err := at.SignedString(accessToken)

    if err != nil {
        return "", "", nil, err
    }

    // Creating Refresh Token
    rtClaims := jwt.MapClaims{}
    rtClaims["user_id"] = user.ID
    rtClaims["email"] = user.Email
    rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() // Token expires after 7 days
    rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

    refreshToken, err := getRefreshSecret()
    if err != nil {
        return "", "", nil, err
    }

    refreshTokenSigned, err := rt.SignedString(refreshToken)

    if err != nil {
        log.Println("Refresh Token signing failed:", err.Error())
        return "", "", nil, err
    }

    return accessTokenSigned, refreshTokenSigned, &atClaims, nil
}
