package api

import (
	"context"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
	"zhaisan-medods/db"
	"zhaisan-medods/utils"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type AuthAPI struct {
	db *mongo.Client
	config *utils.Config
}

func NewAuthAPI(client *mongo.Client, config *utils.Config) *AuthAPI {
	return &AuthAPI{
		db: client,
		config: config,
	}
}

func (api *AuthAPI) GenerateToken(c *gin.Context) {
	userIDHex := c.Query("user_id")
	if userIDHex == "" {
		c.JSON(400, gin.H{"error": "User id must be provided"})
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		UserID: userIDHex,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(api.config.JWTSecretKey))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create access token"})
		return
	}

	refreshToken := utils.GenerateRefreshToken(userIDHex)

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash refresh token"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID format"})
		return
	}

	refreshTokenModel := db.Token{
		UserID:    userID,
		TokenHash: string(hashedRefreshToken),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	ctx := context.Background()
	err = db.SaveRefreshToken(ctx, api.db, refreshTokenModel)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save refresh token"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  tokenString,
		"refresh_token": refreshToken,
	})
}

func (api *AuthAPI) RefreshToken(c *gin.Context) {
	var requestData struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	refreshTokenEncoded := requestData.RefreshToken
	if refreshTokenEncoded == "" {
		c.JSON(400, gin.H{"error": "Refresh token must be provided"})
		return
	}

	refreshTokenBytes, err := base64.StdEncoding.DecodeString(refreshTokenEncoded)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid refresh token"})
		return
	}

	refreshToken := string(refreshTokenBytes)
	refreshTokenParts := strings.Split(refreshToken, ".")
	if len(refreshTokenParts) != 2 {
		c.JSON(400, gin.H{"error": "Invalid refresh token format"})
		return
	}
	userID := refreshTokenParts[0]

	ctx := context.Background()
	token, err := db.FindRefreshTokenByUserID(ctx, api.db, userID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Refresh token not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(token.TokenHash), []byte(refreshTokenEncoded))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid refresh token"})
		return
	}

	if time.Now().After(token.ExpiresAt) {
		c.JSON(400, gin.H{"error": "Refresh token expired"})
		return
	}

	err = db.DeleteRefreshToken(ctx, api.db, token.TokenHash)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete old refresh token"})
		return
	}

	api.GenerateTokenForUserID(c, userID)
}


func (api *AuthAPI) GenerateTokenForUserID(c *gin.Context, userIDHex string) {
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID format"})
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		UserID: userID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	accessToken, err := token.SignedString([]byte(api.config.JWTSecretKey))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create access token"})
		return
	}

	refreshToken := utils.GenerateRefreshToken(userID.Hex())

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash refresh token"})
		return
	}

	refreshTokenModel := db.Token{
		UserID:    userID,
		TokenHash: string(hashedRefreshToken),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	ctx := context.Background()
	err = db.SaveRefreshToken(ctx, api.db, refreshTokenModel)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save refresh token"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

