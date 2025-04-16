package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var nickname string

type KeyPressData struct {
	Key string `json:"key"`
}

type BattleResult struct {
	P1Name    string `json:"p1_name"`
	P2Name    string `json:"p2_name"`
	P1Health  int    `json:"p1_health"`
	P2Health  int    `json:"p2_health"`
	Damage    int    `json:"damage"`
	GameOver  bool   `json:"game_over"`
	Winner    int    `json:"winner"`
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	nickname = readInput("Enter your nickname: ")
	
	for {
		fmt.Println("\n=== Main Menu ===")
		fmt.Println("1. Create Lobby")
		fmt.Println("2. Join Lobby")
		fmt.Println("3. Exit")

		choice := readInput("Select an option: ")
		switch choice {
		case "1":
			createLobby()
		case "2":
			joinLobby()
		case "3":
			fmt.Println("Exiting game.")
			os.Exit(0)
		default:
			fmt.Println("Invalid option, please try again.")
		}
	}
}

func createLobby() {
	password := readInput("Enter a password for the lobby: ")

	data := map[string]string{"password": password}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post("http://146.94.10.168:8080/create", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating lobby:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	lobbyID := int(response["lobby_id"].(float64))
	fmt.Printf("Lobby created! ID: %d\n", lobbyID)
	joinLobbyWithID(lobbyID, password)
}

func joinLobby() {
	lobbyIDStr := readInput("Enter Lobby ID: ")
	lobbyID, err := strconv.Atoi(lobbyIDStr)
	if err != nil {
		fmt.Println("Invalid lobby ID")
		return
	}

	password := readInput("Enter Lobby Password: ")
	joinLobbyWithID(lobbyID, password)
}

func joinLobbyWithID(lobbyID int, password string) {
	data := map[string]interface{}{
		"lobby_id": lobbyID,
		"password": password,
		"nickname": nickname,
	}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post("http://146.94.10.168:8080/join", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error joining lobby:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	if errMsg, exists := response["error"].(string); exists {
		fmt.Println("Error:", errMsg)
		return
	}

	playerID := int(response["player_id"].(float64))
	stats := response["stats"].(map[string]interface{})
	
	fmt.Printf("\nJoined as Player %d\n", playerID)
	fmt.Println("Your Gladiator:")
	fmt.Printf("Health: %v  Attack: %v  Defense: %v  Speed: %v\n",
		stats["health"], stats["attack"], stats["defense"], stats["speed"])

	handleKeyPress(lobbyID, playerID)
}

func handleKeyPress(lobbyID, playerID int) {
	for {
		// Check ATB status
		canMove, atb, err := checkATB(lobbyID, playerID)
		if err != nil {
			fmt.Printf("\nError checking ATB: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if !canMove {
			fmt.Printf("\r[Waiting] ATB: %d%%", atb)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		fmt.Printf("\n[READY] Enter move (1:Attack 2:Power 3:Defend): ")
		move := readInput("")

		if move != "1" && move != "2" && move != "3" {
			fmt.Println("Invalid move! Use 1, 2, or 3.")
			continue
		}

		result, err := submitMove(lobbyID, playerID, move)
		if err != nil {
			fmt.Println("Error submitting move:", err)
			continue
		}

		printBattleResult(result, playerID)
		if result.GameOver {
			if result.Winner == playerID {
				fmt.Println("\nVICTORY! You won the battle!")
			} else {
				fmt.Println("\nDEFEAT! You were defeated!")
			}
			return
		}
	}
}

func checkATB(lobbyID, playerID int) (bool, int, error) {
	url := fmt.Sprintf("http://146.94.10.168:8080/keypress?player_id=%d&lobby_id=%d&check_atb=1", 
		playerID, lobbyID)

	resp, err := http.Get(url)
	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()

	var result struct {
		CanMove  bool `json:"can_move"`
		AtbGauge int  `json:"atb_gauge"`
		Error    string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, 0, err
	}

	if result.Error != "" {
		return false, 0, fmt.Errorf(result.Error)
	}

	return result.CanMove, result.AtbGauge, nil
}

func submitMove(lobbyID, playerID int, move string) (*BattleResult, error) {
	data := KeyPressData{Key: move}
	jsonData, _ := json.Marshal(data)
	url := fmt.Sprintf("http://146.94.10.168:8080/keypress?player_id=%d&lobby_id=%d", 
		playerID, lobbyID)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result BattleResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func printBattleResult(result *BattleResult, playerID int) {
	fmt.Println("\n=== BATTLE ===")
	fmt.Printf("%s: %d HP\n", result.P1Name, result.P1Health)
	fmt.Printf("%s: %d HP\n", result.P2Name, result.P2Health)
	
	if playerID == 1 {
		fmt.Printf("You dealt %d damage!\n", result.Damage)
	} else {
		fmt.Printf("You dealt %d damage!\n", result.Damage)
	}
}
