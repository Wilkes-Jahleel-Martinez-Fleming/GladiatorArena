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
)

type KeyPressData struct {
	Key string `json:"key"`
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	for {
		fmt.Println("\n=== Main Menu ===")
		fmt.Println("1. Create Lobby")
		fmt.Println("2. Join Lobby")
		fmt.Println("3. Exit")

		choiceStr := readInput("Select an option: ")
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Invalid option, please try again.")
			continue
		}

		switch choice {
		case 1:
			createLobby()
		case 2:
			joinLobby()
		case 3:
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
	var response map[string]int
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error decoding response.")
		return
	}

	lobbyID, exists := response["lobby_id"]
	if !exists {
		fmt.Println("Error: Could not retrieve lobby ID.")
		return
	}

	fmt.Printf("Lobby created! ID: %d\n", lobbyID)
	fmt.Println("Waiting for another player to join...")
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
		fmt.Println("Error decoding response.")
		return
	}

	if errMsg, exists := response["error"].(string); exists {
		fmt.Println("Error:", errMsg)
		return
	}

	playerID := int(response["player_id"].(float64))
	fmt.Printf("Joined lobby %d as Player %d\n", lobbyID, playerID)
	
	yourStatsJSON, _ := json.Marshal(response["your_stats"])
	var yourStats map[string]interface{}
	json.Unmarshal(yourStatsJSON, &yourStats)
	
	fmt.Println("\n=== Gladiator Stats ===")
	fmt.Println("Your Gladiator:")
	fmt.Printf("Health: %v | Attack: %v | Defense: %v | Speed: %v\n",
	yourStats["health"], yourStats["attack"], yourStats["defense"], yourStats["speed"])
		
	handleKeyPress(lobbyID, playerID)
}

func handleKeyPress(lobbyID, playerID int) {
    for {
        key := readInput("\nEnter move (1: Attack, 2: Power Attack, 3: Defend): ")

        if key != "1" && key != "2" && key != "3" {
            fmt.Println("Invalid move! Use 1, 2, or 3.")
            continue
        }

        data := KeyPressData{Key: key}
        jsonData, _ := json.Marshal(data)
        url := fmt.Sprintf("http://146.94.10.168:8080/keypress?player_id=%d&lobby_id=%d", playerID, lobbyID)

        resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            fmt.Println("Error: Unable to connect to server.")
            return
        }
        defer resp.Body.Close()

        var response map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
            fmt.Println("Error decoding response:", err)
            return
        }

        if errMsg, exists := response["error"].(string); exists {
            fmt.Println("Error:", errMsg)
            return
        }

        // Safe type assertions
        p1Health, ok1 := response["p1_health"].(float64)
        p2Health, ok2 := response["p2_health"].(float64)
        p1Dmg, ok3 := response[fmt.Sprintf("p%d_dmg", playerID)].(float64)
        p2Dmg, ok4 := response[fmt.Sprintf("p%d_dmg", 3-playerID)].(float64)
        gameOver, ok5 := response["game_over"].(bool)

        if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
            fmt.Println("Error: Invalid server response format")
            return
        }

        fmt.Printf("\n=== Round Results ===\n")
        fmt.Printf("Player 1 Health: %d\n", int(p1Health))
        fmt.Printf("Player 2 Health: %d\n", int(p2Health))
        fmt.Printf("Damage Dealt: You=%d, Opponent=%d\n", int(p1Dmg), int(p2Dmg))

        if gameOver {
            if winner, ok := response["winner"].(float64); ok && int(winner) == playerID {
                fmt.Println("\nVICTORY! You win the battle!")
            } else {
                fmt.Println("\nDEFEAT! You have been slain!")
            }
            fmt.Println("\nReturning to main menu...")
            return
        }
    }
}
