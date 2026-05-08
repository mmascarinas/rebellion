package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"rebellion/internal/game"
	"rebellion/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type GameHandler struct {
	manager *game.Manager
}

func NewGameHandler() *GameHandler {
	return &GameHandler{
		manager: game.NewManager(),
	}
}

func (h *GameHandler) CreateGame(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		response.Error(c, http.StatusBadRequest, "missing name parameter", nil)
		return
	}

	playersStr := c.Query("players")
	maxPlayers, err := strconv.Atoi(playersStr)
	if err != nil || maxPlayers < 2 || maxPlayers > 6 {
		response.Error(c, http.StatusBadRequest, "players must be between 2 and 6", nil)
		return
	}

	gameID := uuid.New().String()
	h.manager.Create(gameID, maxPlayers)

	session, _ := h.manager.Get(gameID)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to upgrade to websocket", err.Error())
		return
	}

	defer conn.Close()

	playerID := uuid.New().String()
	player := game.NewPlayer(playerID, name, conn)
	session.AddPlayer(player)
	defer session.RemovePlayer(playerID)

	msg, _ := json.Marshal(gin.H{
		"event":   "game_created",
		"game_id": gameID,
		"name":    name,
		"players": session.PlayerCount(),
		"max":     maxPlayers,
	})
	conn.WriteMessage(websocket.TextMessage, msg)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var evt struct {
			Event string `json:"event"`
			Data  string `json:"data"`
		}
		if json.Unmarshal(message, &evt) == nil {
			if evt.Event == "leave_game" {
				break
			}
			if evt.Event == "response" {
				player.ReceiveResponse(evt.Data)
				continue
			}
			if evt.Event == "start_game" {
				go session.StartGame()
				continue
			}
		}

		session.Broadcast(message)
	}

	msg, _ = json.Marshal(gin.H{
		"event":   "player_left",
		"player":  playerID,
		"players": session.PlayerCount(),
	})
	session.Broadcast(msg)
}

func (h *GameHandler) JoinGame(c *gin.Context) {
	gameID := c.Param("id")

	session, ok := h.manager.Get(gameID)
	if !ok {
		response.Error(c, http.StatusNotFound, "game not found", nil)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	playerID := uuid.New().String()
	playerName := c.Query("name")
	player := game.NewPlayer(playerID, playerName, conn)

	if !session.AddPlayer(player) {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"game is full"}`))
		return
	}
	defer session.RemovePlayer(playerID)

	msg, _ := json.Marshal(gin.H{
		"event":   "player_joined",
		"player":  playerID,
		"players": session.PlayerCount(),
	})
	session.Broadcast(msg)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var evt struct {
			Event string `json:"event"`
			Data  string `json:"data"`
		}
		if json.Unmarshal(message, &evt) == nil {
			if evt.Event == "leave_game" {
				break
			}
			if evt.Event == "response" {
				player.ReceiveResponse(evt.Data)
				continue
			}
		}

		session.Broadcast(message)
	}

	msg, _ = json.Marshal(gin.H{
		"event":   "player_left",
		"player":  playerID,
		"players": session.PlayerCount(),
	})
	session.Broadcast(msg)
}
