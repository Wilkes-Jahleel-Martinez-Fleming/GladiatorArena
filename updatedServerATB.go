package main

import (
	//"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	//"os"
	"sync"
	"time"
	"strconv"
)

type Gladiator struct {
	ID      int `json:"id"`
	Health  int `json:"health"`
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	Speed   int `json:"speed"` // Lower speed = faster turns
}

func genGladiator(id int) Gladiator {
	return Gladiator{
		ID:      id,
		Health:  rand.Intn(25) + 75,
		Attack:  rand.Intn(25) + 10,
		Defense: rand.Intn(5) + 5,
		Speed:   rand.Intn(10) + 5, // Speed range: 5-14
	}
}

type KeyPressData struct {
	Key string `json:"key"`
}

type Player struct {
	ID        int       `json:"id"`
	Nickname  string    `json:"nickname"`
	Gladiator Gladiator `json:"gladiator"`
}

type Lobby struct {
	ID           int
	Password     string
	Players      [2]*Player
	PlayerCount  int
	mu           sync.Mutex
	AtbGauges    [2]int       // 0-100 for each player
	LastTickTime time.Time
	GameOver     bool `json:"game_over"`
	Winner       int  `json:"winner"`
}

var lobbies = make(map[int]*Lobby)
var lobbyIDCounter = 1
var mu sync.Mutex

func init() {
	rand.Seed(time.Now().UnixNano())
}

func attack(defender *Gladiator, attacker Gladiator) int {
	dmg := attacker.Attack + rand.Intn(10) - defender.Defense
	if dmg < 0 {
		dmg = 0
	}
	defender.Health -= dmg
	return dmg
}

func powAttack(defender *Gladiator, attacker *Gladiator) int {
	dmg := attacker.Attack*2 + rand.Intn(10) - defender.Defense
	attacker.Speed = 0 // Power attack exhausts the attacker
	if dmg < 0 {
		dmg = 0
	}
	defender.Health -= dmg
	return dmg
}

func defend(glad *Gladiator) {
	glad.Defense += 5
}

func createLobby(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var requestData struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	lobby := &Lobby{
		ID:          lobbyIDCounter,
		Password:    requestData.Password,
		PlayerCount: 0,
		AtbGauges:    [2]int{0, 0},
		LastTickTime: time.Now(),
	}
	lobbies[lobbyIDCounter] = lobby
	lobbyIDCounter++

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"lobby_id": lobby.ID,
		"message":  "Lobby created",
	})
}

func joinLobby(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var requestData struct {
		LobbyID  int    `json:"lobby_id"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	lobby, exists := lobbies[requestData.LobbyID]
	if !exists {
		http.Error(w, `{"error": "Lobby not found"}`, http.StatusNotFound)
		return
	}

	if lobby.Password != requestData.Password {
		http.Error(w, `{"error": "Incorrect password"}`, http.StatusUnauthorized)
		return
	}

	if lobby.PlayerCount >= 2 {
		http.Error(w, `{"error": "Lobby is full"}`, http.StatusConflict)
		return
	}

	playerID := lobby.PlayerCount
	lobby.Players[playerID] = &Player{
		ID:        playerID + 1,
		Nickname:  requestData.Nickname,
		Gladiator: genGladiator(playerID + 1),
	}
	lobby.PlayerCount++

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"player_id": playerID + 1,
		"stats":     lobby.Players[playerID].Gladiator,
	})
}

func handleKeyPress(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("player_id")
	lobbyID := r.URL.Query().Get("lobby_id")
	checkATB := r.URL.Query().Get("check_atb") == "1"

	pid, err := strconv.Atoi(playerID)
	if err != nil || pid < 1 || pid > 2 {
		http.Error(w, `{"error": "Invalid player ID"}`, http.StatusBadRequest)
		return
	}

	lid, err := strconv.Atoi(lobbyID)
	if err != nil {
		http.Error(w, `{"error": "Invalid lobby ID"}`, http.StatusBadRequest)
		return
	}

	mu.Lock()
	lobby, exists := lobbies[lid]
	mu.Unlock()
	if !exists {
		http.Error(w, `{"error": "Lobby not found"}`, http.StatusNotFound)
		return
	}

	lobby.mu.Lock()
	defer lobby.mu.Unlock()

	if lobby.PlayerCount < 2 {
		http.Error(w, `{"error": "Waiting for opponent"}`, http.StatusPreconditionFailed)
		return
	}

	// Update ATB gauges
	now := time.Now()
	elapsed := now.Sub(lobby.LastTickTime).Seconds()
	lobby.LastTickTime = now

	for i := 0; i < 2; i++ {
		lobby.AtbGauges[i] += int(float64(20-lobby.Players[i].Gladiator.Speed) * elapsed * 10)
		if lobby.AtbGauges[i] > 100 {
			lobby.AtbGauges[i] = 100
		}
	}

	if checkATB {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"can_move":  lobby.AtbGauges[pid-1] >= 100,
			"atb_gauge": lobby.AtbGauges[pid-1],
		})
		return
	}

	// Process move
	if lobby.AtbGauges[pid-1] < 100 {
		http.Error(w, `{"error": "ATB gauge not full"}`, http.StatusForbidden)
		return
	}

	var keyPress KeyPressData
	if err := json.NewDecoder(r.Body).Decode(&keyPress); err != nil {
		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	lobby.AtbGauges[pid-1] = 0
	p1 := lobby.Players[0]
	p2 := lobby.Players[1]

	var dmg int
	if pid == 1 {
		dmg = resolveMove(p1, p2, keyPress.Key)
	} else {
		dmg = resolveMove(p2, p1, keyPress.Key)
	}

	if p1.Gladiator.Health <= 0 {
		lobby.GameOver = true
		lobby.Winner = 2
	} else if p2.Gladiator.Health <= 0 {
		lobby.GameOver = true
		lobby.Winner = 1
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"game_over": lobby.GameOver,
		"p1_name":   p1.Nickname,
		"p2_name":   p2.Nickname,
		"p1_health": p1.Gladiator.Health,
		"p2_health": p2.Gladiator.Health,
		"p1_atb":    lobby.AtbGauges[0],
		"p2_atb":    lobby.AtbGauges[1],
		"damage":    dmg,
		"winner":    lobby.Winner,
	})
}

func resolveMove(attacker, defender *Player, move string) int {
	switch move {
	case "1":
		return attack(&defender.Gladiator, attacker.Gladiator)
	case "2":
		return powAttack(&defender.Gladiator, &attacker.Gladiator)
	case "3":
		defend(&attacker.Gladiator)
		return 0
	default:
		return 0
	}
}

func main() {
	http.HandleFunc("/create", createLobby)
	http.HandleFunc("/join", joinLobby)
	http.HandleFunc("/keypress", handleKeyPress)

	fmt.Println("Server running on http://146.94.10.168:8080")
	log.Fatal(http.ListenAndServe("146.94.10.168:8080", nil))
}
