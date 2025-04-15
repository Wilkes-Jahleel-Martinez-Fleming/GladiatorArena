package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
	"strconv"
)

type Gladiator struct {
	ID      int `json:"id"`
	Health  int `json:"health"`
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	Speed   int `json:"speed"`
}

func genGladiator(id int) Gladiator {
	return Gladiator{
		ID:      id,
		Health:  rand.Intn(25) + 75,
		Attack:  rand.Intn(25) + 10,
		Defense: rand.Intn(5) + 5,
		Speed:   rand.Intn(10) + 5,
	}
}

type KeyPressData struct {
	Key string `json:"key"`
}

type Player struct {
	ID        int       `json:"id"`
	Nickname  String    `json:"nickname`
	Gladiator Gladiator `json:"gladiator"`
}

type Lobby struct {
	ID           int
	Password     string
	Players      [2]*Player
	PlayerCount  int
	mu           sync.Mutex
	readyCount   int
	cond         *sync.Cond
	GameOver     bool `json:"game_over"`
	Winner       int  `json:"winner"`
	RoundResults map[string]interface{}
}

var lobbies = make(map[int]*Lobby)
var lobbyIDCounter = 1
var mu sync.Mutex

func init() {
	rand.Seed(time.Now().UnixNano())
}

func attack(defender *Gladiator, attacker Gladiator) int {
	dmg := attacker.Attack + rand.Intn(10) - 5 - defender.Defense - rand.Intn(5) + 2
	if dmg < 0 {
		dmg = 0
	}
	defender.Health -= dmg
	return dmg
}

func powAttack(defender *Gladiator, attacker *Gladiator) int {
	attacker.Speed = 0
	dmg := attacker.Attack*2 + rand.Intn(10) - 5 - defender.Defense - rand.Intn(5) + 2
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

	lobbyID := lobbyIDCounter
	lobbyIDCounter++

	lobby := &Lobby{
		ID:          lobbyID,
		Password:    requestData.Password,
		PlayerCount: 0,
		RoundResults: make(map[string]interface{}),
	}
	lobby.cond = sync.NewCond(&lobby.mu)
	lobbies[lobbyID] = lobby

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"lobby_id": lobbyID})
}

func joinLobby(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var requestData struct {
		LobbyID  int    `json:"lobby_id"`
		Password string `json:"password"`
		Nickname  String    `json:"nickname`
	Gladiator Gladiator `json:"gladiator"`
}

type Lobby struct {
	ID           int
	Password     string
	Players      [2]*Player
	PlayerCount  int
	mu           sync.Mutex
	readyCount   int
	cond         *sync.Cond
	GameOver     bool `json:"game_over"`
	Winner       int  `json:"winner"`
	RoundResults map[string]interface{}
}

var lobbies = make(map[int]*Lobby)
var lobbyIDCounter = 1
var mu sync.Mutex

func init() {
	rand.Seed(time.Now().UnixNano())
}

func attack(defender *Gladiator, attacker Gladiator) int {
	dmg := attacker.Attack + rand.Intn(10) - 5 - defender.Defense - rand.Intn(5) + 2
	if dmg < 0 {
		dmg = 0
	}
	defender.Health -= dmg
	return dmg
}

func powAttack(defender *Gladiator, attacker *Gladiator) int {
	attacker.Speed = 0
	dmg := attacker.Attack*2 + rand.Intn(10) - 5 - defender.Defense - rand.Intn(5) + 2
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

	lobbyID := lobbyIDCounter
	lobbyIDCounter++

	lobby := &Lobby{
		ID:          lobbyID,
		Password:    requestData.Password,
		PlayerCount: 0,
		RoundResults: make(map[string]interface{}),
	}
	lobby.cond = sync.NewCond(&lobby.mu)
	lobbies[lobbyID] = lobby

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"lobby_id": lobbyID})
}

func joinLobby(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var requestData struct {
		LobbyID  int    `json:"lobby_id"`
		Password string `json:"password"`
		Nickname  String    `json:"nickname`
	Gladiator Gladiator `json:"gladiator"`
}

type Lobby struct {
	ID           int
	Password     string
	Players      [2]*Player
	PlayerCount  int
	mu           sync.Mutex
	readyCount   int
	cond         *sync.Cond
	GameOver     bool `json:"game_over"`
	Winner       int  `json:"winner"`
	RoundResults map[string]interface{}
}

var lobbies = make(map[int]*Lobby)
var lobbyIDCounter = 1
var mu sync.Mutex

func init() {
	rand.Seed(time.Now().UnixNano())
}

func attack(defender *Gladiator, attacker Gladiator) int {
	dmg := attacker.Attack + rand.Intn(10) - 5 - defender.Defense - rand.Intn(5) + 2
	if dmg < 0 {
		dmg = 0
	}
	defender.Health -= dmg
	return dmg
}

func powAttack(defender *Gladiator, attacker *Gladiator) int {
	attacker.Speed = 0
	dmg := attacker.Attack*2 + rand.Intn(10) - 5 - defender.Defense - rand.Intn(5) + 2
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

	lobbyID := lobbyIDCounter
	lobbyIDCounter++

	lobby := &Lobby{
		ID:          lobbyID,
		Password:    requestData.Password,
		PlayerCount: 0,
		RoundResults: make(map[string]interface{}),
	}
	lobby.cond = sync.NewCond(&lobby.mu)
	lobbies[lobbyID] = lobby

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"lobby_id": lobbyID})
}

func joinLobby(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var requestData struct {
		LobbyID  int    `json:"lobby_id"`
		Password string `json:"password"`
		Nickname  String    `json:"nickname`
	Gladiator Gladiator `json:"gladiator"`
}

type Lobby struct {
	ID           int
	Password     string
	Players      [2]*Player
	PlayerCount  int
	mu           sync.Mutex
	readyCount   int
	cond         *sync.Cond
	GameOver     bool `json:"game_over"`
	Winner       int  `json:"winner"`
	RoundResults map[string]interface{}
}

var lobbies = make(map[int]*Lobby)
var lobbyIDCounter = 1
var mu sync.Mutex

func init() {
	rand.Seed(time.Now().UnixNano())
}

func attack(defender *Gladiator, attacker Gladiator) int {
	dmg := attacker.Attack + rand.Intn(10) - 5 - defender.Defense - rand.Intn(5) + 2
	if dmg < 0 {
		dmg = 0
	}
	defender.Health -= dmg
	return dmg
}

func powAttack(defender *Gladiator, attacker *Gladiator) int {
	attacker.Speed = 0
	dmg := attacker.Attack*2 + rand.Intn(10) - 5 - defender.Defense - rand.Intn(5) + 2
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

	lobbyID := lobbyIDCounter
	lobbyIDCounter++

	lobby := &Lobby{
		ID:          lobbyID,
		Password:    requestData.Password,
		PlayerCount: 0,
		RoundResults: make(map[string]interface{}),
	}
	lobby.cond = sync.NewCond(&lobby.mu)
	lobbies[lobbyID] = lobby

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"lobby_id": lobbyID})
}

func joinLobby(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var requestData struct {
		LobbyID  int    `json:"lobby_id"`
		Password string `json:"password"`
		Nickname  String    `json:"nickname`
	Gladiator Gladiator `json:"gladiator"`
}

type Lobby struct {
	ID           int
	Password     string
	Players      [2]*Player
	PlayerCount  int
	mu           sync.Mutex
	readyCount   int
	cond         *sync.Cond
	GameOver     bool `json:"game_over"`
	Winner       int  `json:"winner"`
	RoundResults map[string]interface{}
}

var lobbies = make(map[int]*Lobby)
var lobbyIDCounter = 1
var mu sync.Mutex

func init() {
	rand.Seed(time.Now().UnixNano())
}

func attack(defender *Gladiator, attacker Gladiator) int {
	dmg := attacker.Attack + rand.Intn(10) - 5 - defender.Defense - rand.Intn(5) + 2
	if dmg < 0 {
		dmg = 0
	}
	defender.Health -= dmg
	return dmg
}

func powAttack(defender *Gladiator, attacker *Gladiator) int {
	attacker.Speed = 0
	dmg := attacker.Attack*2 + rand.Intn(10) - 5 - defender.Defense - rand.Intn(5) + 2
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

	lobbyID := lobbyIDCounter
	lobbyIDCounter++

	lobby := &Lobby{
		ID:          lobbyID,
		Password:    requestData.Password,
		PlayerCount: 0,
		RoundResults: make(map[string]interface{}),
	}
	lobby.cond = sync.NewCond(&lobby.mu)
	lobbies[lobbyID] = lobby

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"lobby_id": lobbyID})
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
	lobby.PlayerCount++
	lobby.Players[playerID] = &Player{
		ID:        playerID + 1,
		Nickname:  requestData.Nickname,
		Gladiator: genGladiator(playerID + 1),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   fmt.Sprintf("Joined lobby %d as player %d", requestData.LobbyID, playerID+1),
		"player_id": playerID + 1,
		"your_stats": lobby.Players[playerID].Gladiator,
	})
}

func handleKeyPress(w http.ResponseWriter, r *http.Request) {
	var keyPress KeyPressData
	if err := json.NewDecoder(r.Body).Decode(&keyPress); err != nil {
		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	playerID := r.URL.Query().Get("player_id")
	lobbyID := r.URL.Query().Get("lobby_id")

	if playerID == "" || lobbyID == "" {
		http.Error(w, `{"error": "Missing parameters"}`, http.StatusBadRequest)
		return
	}

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

	lobby, exists := lobbies[lid]
	if !exists {
		http.Error(w, `{"error": "Lobby not found"}`, http.StatusNotFound)
		return
	}

	lobby.mu.Lock()
	defer lobby.mu.Unlock()

	if lobby.GameOver {
		http.Error(w, `{"error": "Game is already over"}`, http.StatusConflict)
		return
	}

	// Store player's move
	lobby.RoundResults[fmt.Sprintf("p%d_move", pid)] = keyPress.Key
	lobby.readyCount++

	// Wait for both players
	if lobby.readyCount < 2 {
		lobby.cond.Wait()
		
		// After waking up, send the results
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lobby.RoundResults)
		return
	}

	// Both players ready - process moves
	p1 := lobby.Players[0]
	p2 := lobby.Players[1]
	var p1Dmg, p2Dmg int

	// Get moves
	p1Move := lobby.RoundResults["p1_move"].(string)
	p2Move := lobby.RoundResults["p2_move"].(string)

	// Resolve moves based on speed
	if p1.Gladiator.Speed >= p2.Gladiator.Speed {
		p1Dmg = resolveMove(p1, p2, p1Move)
		p2Dmg = resolveMove(p2, p1, p2Move)
	} else {
		p2Dmg = resolveMove(p2, p1, p2Move)
		p1Dmg = resolveMove(p1, p2, p1Move)
	}

	// Check for game over
	if p1.Gladiator.Health <= 0 {
		lobby.GameOver = true
		lobby.Winner = 2
	} else if p2.Gladiator.Health <= 0 {
		lobby.GameOver = true
		lobby.Winner = 1
	}

	// Prepare response
	lobby.RoundResults = map[string]interface{}{
		"status":    "success",
		"game_over": lobby.GameOver,
		"winner":    lobby.Winner,
		"p1_health": p1.Gladiator.Health,
		"p2_health": p2.Gladiator.Health,
		"p1_dmg":    p1Dmg,
		"p2_dmg":    p2Dmg,
	}

	// Broadcast to both players
	lobby.cond.Broadcast()
	lobby.readyCount = 0

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lobby.RoundResults)
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

func deleteLobby(lobbyID int, password string) {
	mu.Lock()
	defer mu.Unlock()

	lobby, exists := lobbies[lobbyID]
	if !exists {
		fmt.Println("Error: Lobby not found")
		return
	}

	if lobby.Password != password {
		fmt.Println("Error: Incorrect password")
		return
	}

	delete(lobbies, lobbyID)
	fmt.Printf("Lobby %d deleted successfully!\n", lobbyID)
}

func consoleListener() {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		if input == "~\n" {
			fmt.Println("Enter Lobby ID to delete:")
			var lobbyID int
			fmt.Scanf("%d", &lobbyID)

			fmt.Println("Enter Lobby Password:")
			var password string
			fmt.Scanf("%s", &password)

			deleteLobby(lobbyID, password)
		}
	}
}

func setupRoutes() {
	http.HandleFunc("/create", createLobby)
	http.HandleFunc("/join", joinLobby)
	http.HandleFunc("/keypress", handleKeyPress)
}

func main() {
	go consoleListener()
	setupRoutes()
	fmt.Println("Server is running on http://146.94.10.168:8080")
	log.Fatal(http.ListenAndServe("146.94.10.168:8080", nil))
}
