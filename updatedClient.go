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
var nickname string

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

		nickname = readInput("Enter your nickname: ") 	
	
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
	joinLobbyWithID(lobbyID, password, nickname)

}

func joinLobby() {
	lobbyIDStr := readInput("Enter Lobby ID: ")
	lobbyID, err := strconv.Atoi(lobbyIDStr)
	if err != nil {
		fmt.Println("Invalid lobby ID")
		return
	}

	password := readInput("Enter Lobby Password: ")
	joinLobbyWithID(lobbyID, password, nickname)

}

func joinLobbyWithID(lobbyID int, password string, nickname string) {
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

        // Get names
        p1Name := response["p1_name"].(string)
        p2Name := response["p2_name"].(string)
        yourName := p1Name
        opponentName := p2Name
        if playerID == 2 {
            yourName = p2Name
            opponentName = p1Name
        }

        // Get health and damage
        p1Health := int(response["p1_health"].(float64))
        p2Health := int(response["p2_health"].(float64))
        yourDmg := int(response[fmt.Sprintf("p%d_dmg", playerID)].(float64))
        opponentDmg := int(response[fmt.Sprintf("p%d_dmg", 3-playerID)].(float64))
        gameOver := response["game_over"].(bool)

        // Display results
        fmt.Printf("\n=== Round Results ===\n")
        fmt.Printf("%s's Health: %d\n", p1Name, p1Health)
        fmt.Printf("%s's Health: %d\n", p2Name, p2Health)
        fmt.Printf("Damage Dealt: %s=%d, %s=%d\n", yourName, yourDmg, opponentName, opponentDmg)

        // Handle game over
        if gameOver {
            winnerName := response["winner_name"].(string)
            if winnerName == yourName {
                fmt.Printf("\nVICTORY! %s wins the battle!\n", yourName)
            } else {
                fmt.Printf("\nDEFEAT! %s has been slain by %s!\n", yourName, winnerName)
            }
            fmt.Println("\nReturning to main menu...")
            return
        }
    }

	
	
}
