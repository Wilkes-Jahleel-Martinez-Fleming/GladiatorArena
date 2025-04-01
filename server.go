package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type KeyPressData struct {
	Key string `json:"key"`
}

type Player struct {
	ID  int
	Key string
}

type Lobby struct {
	ID         int
	Password   string
	Players    [2]*Player
	Keys       [2]string
	PlayerCount int
	mu         sync.Mutex
	readyCount int
	cond       *sync.Cond
}

var lobbies = make(map[int]*Lobby)
var lobbyIDCounter = 1
var mu sync.Mutex

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
		ID:         lobbyID,
		Password:   requestData.Password,
		PlayerCount: 0,
	}
	lobby.cond = sync.NewCond(&lobby.mu)
	lobbies[lobbyID] = lobby

	log.Printf("Lobby created! Lobby ID: %d", lobbyID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"lobby_id": lobbyID})
}

func joinLobby(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var requestData struct {
		LobbyID  int    `json:"lobby_id"`
		Password string `json:"password"`
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
	lobby.Players[playerID] = &Player{ID: playerID}

	log.Printf("Player %d joined lobby %d", playerID+1, requestData.LobbyID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   fmt.Sprintf("Joined lobby %d as player %d", requestData.LobbyID, playerID+1),
		"player_id": playerID + 1,
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

	lobby, exists := lobbies[toInt(lobbyID)]
	if !exists {
		http.Error(w, `{"error": "Lobby not found"}`, http.StatusNotFound)
		return
	}

	id := toInt(playerID) - 1
	if id < 0 || id > 1 {
		http.Error(w, `{"error": "Invalid player ID"}`, http.StatusBadRequest)
		return
	}

	lobby.mu.Lock()
	defer lobby.mu.Unlock()

	lobby.Keys[id] = keyPress.Key
	lobby.readyCount++

	if lobby.readyCount < 2 {
		lobby.cond.Wait()
	} else {
		lobby.cond.Broadcast()
	}

	log.Printf("Lobby %d: Player 1 (%s) vs Player 2 (%s)", lobby.ID, lobby.Keys[0], lobby.Keys[1])
	lobby.readyCount = 0

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "player_key": keyPress.Key})
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

func toInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func setupRoutes() {
	http.HandleFunc("/create", createLobby)
	http.HandleFunc("/join", joinLobby)
	http.HandleFunc("/keypress", handleKeyPress)
}

func main() {
	go consoleListener() // Start listening for console input
	setupRoutes()
	fmt.Println("Server is running on http://146.94.10.168:8080")
	log.Fatal(http.ListenAndServe("146.94.10.168:8080", nil))
}
