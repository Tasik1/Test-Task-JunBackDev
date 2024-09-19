package handler

import (
	"TestBackDev/model"
	"TestBackDev/repository"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

type UserHandler interface {
	CreateUser(*gin.Context)
	SignIn(*gin.Context)
	RefreshTokenPair(*gin.Context)
}
type userHandler struct {
	repo repository.UserRepository
}

func NewUserHandler() UserHandler {
	return &userHandler{
		repo: repository.NewUserRepository(),
	}
}

func (h *userHandler) CreateUser(ctx *gin.Context) {
	var input model.User
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.repo.CreateUser(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *userHandler) SignIn(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	dbUser, err := h.repo.GetByLogin(user.Login)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "No Such User Found"})
		return
	}

	tMP, err := GenerateTokenPair(dbUser.ID, ctx.ClientIP())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	token := model.Token{
		GUID:    dbUser.ID,
		Refresh: tMP["refresh_token"],
	}
	err = h.repo.StoreRefreshToken(token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"coudn't store refresh token to DB": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"Access Token": tMP["access_token"], "Refresh Token": tMP["refresh_token"]})
	return
}

func (h *userHandler) RefreshTokenPair(ctx *gin.Context) {
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	tokenReq := tokenReqBody{}
	if err := ctx.ShouldBindJSON(&tokenReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	tokenDec, _ := base64.StdEncoding.DecodeString(tokenReq.RefreshToken)
	token, _ := ParseTokenFromString(string(tokenDec))

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		dbTokenStr, err := h.repo.GetTokenByUserID(uint(claims["GUID"].(float64)))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "No Such User Found"})
			return
		}
		dbTokenDec, _ := base64.StdEncoding.DecodeString(dbTokenStr.Refresh)
		dbToken, _ := ParseTokenFromString(string(dbTokenDec))
		dbClaims := dbToken.Claims.(jwt.MapClaims)

		if ctx.ClientIP() != dbClaims["IP"] {
			err = SendEmail("210103282@stu.sdu.edu.kz", WARNING_SJT, HTML_BODY)
		}

		newTokenPair, err := GenerateTokenPair(dbTokenStr.GUID, ctx.ClientIP())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error:": err})
			return
		}

		dbTokenStr = model.Token{
			Refresh: newTokenPair["refresh_token"],
		}
		err = h.repo.UpdateRefreshToken(dbTokenStr)

		ctx.JSON(http.StatusOK, newTokenPair)
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"error:": "not valid"})
	return
}

func ParseTokenFromString(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}
