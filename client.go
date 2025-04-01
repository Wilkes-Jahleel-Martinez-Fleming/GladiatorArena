package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type KeyPressData struct {
	Key string `json:"key"`
}

func main() {
	for {
		fmt.Println("\n=== Main Menu ===")
		fmt.Println("1. Create Lobby")
		fmt.Println("2. Join Lobby")
		fmt.Println("3. Exit")
		fmt.Print("Select an option: ")

		var choice int
		fmt.Scanf("%d", &choice)

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
	fmt.Print("Enter a password for the lobby: ")
	var password string
	fmt.Scanf("%s", &password)

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
	fmt.Print("Enter Lobby ID: ")
	var lobbyID int
	fmt.Scanf("%d", &lobbyID)

	fmt.Print("Enter Lobby Password: ")
	var password string
	fmt.Scanf("%s", &password)

	joinLobbyWithID(lobbyID, password)
}

func joinLobbyWithID(lobbyID int, password string) {
	data := map[string]interface{}{
		"lobby_id":  lobbyID,
		"password":  password,
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

	// Start sending keypresses
	handleKeyPress(lobbyID, playerID)
}

func handleKeyPress(lobbyID, playerID int) {
	for {
		var key string

		for {
		fmt.Print("Press a key to send to the server: ")
		fmt.Scanf("%s", &key)
			
			if len(key) == 1 {
				if key == "1" || key == "2" || key == "3" {
					break 
				}
				fmt.Println("Invalid input! Please enter either 1, 2, or 3.")
			} else {
				fmt.Println("Invalid input! Please enter only one character.")
			}
		}
		
		data := KeyPressData{Key: key}
		jsonData, _ := json.Marshal(data)
		url := fmt.Sprintf("http://146.94.10.168:8080/keypress?player_id=%d&lobby_id=%d", playerID, lobbyID)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error: Unable to connect to server.")
			fmt.Println("Returning to main menu...")
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		var respMessage map[string]string
		if err := json.Unmarshal(body, &respMessage); err != nil {
			fmt.Println("Error decoding response from server.")
			fmt.Println("Returning to main menu...")
			return
		}

		// Handle errors
		if errMsg, exists := respMessage["error"]; exists {
			fmt.Println("Error:", errMsg)
			fmt.Println("Returning to main menu...")
			return
		}

		fmt.Println("Keypress sent successfully! Waiting for next round...")
	}
}
