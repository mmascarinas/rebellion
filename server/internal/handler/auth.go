package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"rebellion/internal/config"
	"rebellion/internal/response"
	"rebellion/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type AuthHandler struct {
	rdb         *redis.Client
	oauthCfg    *oauth2.Config
	userInfoURL string
}

func NewAuthHandler(rdb *redis.Client, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		rdb:         rdb,
		userInfoURL: cfg.AuthUserOauthURL,
		oauthCfg: &oauth2.Config{
			ClientID:     cfg.AuthClientID,
			ClientSecret: cfg.AuthClientSecret,
			RedirectURL:  cfg.AuthRedirectURL,
			Scopes:       []string{cfg.AuthUserScopeURL + ".email", cfg.AuthUserScopeURL + ".profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx := context.Background()

	sessionID, _ := c.Cookie("Authorization")
	if sessionID != "" {
		userData, err := h.rdb.Get(ctx, sessionID).Result()
		if err == nil {
			var user UserInfo
			err := json.Unmarshal([]byte(userData), &user)
			if err == nil {
				response.Success(c, 0, "already authenticated", user)
				return
			} else {
				response.Error(c, http.StatusInternalServerError, "failed to parse session data", err.Error())
				return
			}
		}
	}

	code := c.Query("code")
	if code == "" {
		response.Error(c, http.StatusBadRequest, "missing code parameter", nil)
		return
	}

	token, err := h.oauthCfg.Exchange(ctx, code)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid code", err.Error())
		return
	}

	client := h.oauthCfg.Client(ctx, token)
	resp, err := client.Get(h.userInfoURL)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch user info", err.Error())
		return
	}

	user, err := utils.DecodeResponse[UserInfo](resp)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to parse user info", err.Error())
		return
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to encode session", err.Error())
		return
	}

	newSessionID := uuid.New().String()
	err = h.rdb.Set(ctx, newSessionID, string(userJSON), 24*time.Hour).Err()

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to save session", err.Error())
		return
	}

	c.SetCookie("Authorization", newSessionID, 86400, "/", "", true, true)

	response.Success(c, 0, "login successful", user)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID, _ := c.Cookie("Authorization")
	if sessionID != "" {
		h.rdb.Del(context.Background(), sessionID)
	}

	c.SetCookie("Authorization", "", -1, "/", "", true, true)

	response.Success(c, 0, "logged out", nil)
}
