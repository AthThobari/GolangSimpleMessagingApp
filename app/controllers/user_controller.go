package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/app/repository"
	"github.com/kooroshh/fiber-boostrap/pkg/jwt_token"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) error {
	user := new(models.User)

	err := ctx.BodyParser(user)
	if err != nil {
		errResponse := fmt.Errorf("failed to parse request: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}
	err = user.Validate()
	if err != nil {
		errResponse := fmt.Errorf("failed to validate request: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		errResponse := fmt.Errorf("failed to encrypt the password: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}
	user.Password = string(hashPassword)

	err = repository.InsertNewUser(ctx.Context(), user)
	if err != nil {
		errResponse := fmt.Errorf("failed to insert new user: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	resp := user
	resp.Password = ""
	return response.SendsuccessResponse(ctx, resp)
}

func Login(ctx *fiber.Ctx) error {
	loginReq := new(models.LoginRequest)
	resp := models.LoginResponse{}
	now := time.Now()

	err := ctx.BodyParser(loginReq)
	if err != nil {
		errResponse := fmt.Errorf("failed to parse request: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}
	err = loginReq.Validate()
	if err != nil {
		errResponse := fmt.Errorf("failed to validate request: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}

	user, err := repository.GetUserByUsername(ctx.Context(), loginReq.Username)
	if err != nil {
		errResponse := fmt.Errorf("failed to get username: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "email/password incorrect", nil)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		errResponse := fmt.Errorf("failed to get username: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "email/password incorrect", nil)
	}

	token, err := jwt_token.GenerateToken(ctx.Context(), user.Username, user.Fullname, "token")
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "terjadi kesalahan pada sistem", nil)
	}
	refreshToken, err := jwt_token.GenerateToken(ctx.Context(), user.Username, user.Fullname, "refresh_token")
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "terjadi kesalahan pada sistem", nil)
	}
	log.Println("Token duration:", jwt_token.MapTypeToken["token"])
	log.Println("Refresh Token duration:", jwt_token.MapTypeToken["refresh_token"])

	userSession := models.UserSession{
		UserID:              user.ID,
		Token:               token,
		RefreshToken:        refreshToken,
		TokenExpired:        now.Add(jwt_token.MapTypeToken["token"]),
		RefreshTokenExpired: now.Add(jwt_token.MapTypeToken["refresh_token"]),
	}

	err = repository.InsertNewUserSession(ctx.Context(), &userSession)

	if err != nil {
		errResponse := fmt.Errorf("failed insert user session: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "terjadi kesalahan pada sistem", nil)
	}

	resp.Username = user.Username
	resp.Fullname = user.Fullname
	resp.Token = token
	resp.RefreshToken = refreshToken
	return response.SendsuccessResponse(ctx, resp)
}

func Logout(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")
	err := repository.DeleteUserSessionByToken(ctx.Context(), token)
	if err != nil {
		errResponse := fmt.Errorf("failed delete user session: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada sistem", nil)
	}
	return response.SendsuccessResponse(ctx, nil)
}

func RefreshToken(ctx *fiber.Ctx) error {
	now := time.Now()
	RefreshToken := ctx.Get("Authorization")
	username := ctx.Locals("username").(string)
	fullName := ctx.Locals("full_name").(string)

	token, err := jwt_token.GenerateToken(ctx.Context(), username, fullName, "token")
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada sistem", nil)
	}

	err = repository.UpdateUserSessionToken(ctx.Context(), token, now.Add(jwt_token.MapTypeToken["token"]), RefreshToken)
	if err != nil {
		errResponse := fmt.Errorf("failed to update token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada sistem", nil)
	}
	
	return response.SendsuccessResponse(ctx, fiber.Map{
		"token": token,
	})
}
