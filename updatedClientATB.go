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

func clearLine() {
    fmt.Print("\033[2K\r") // ANSI escape sequence to clear line
}

func handleKeyPress(lobbyID, playerID int) {
    gameOverChan := make(chan bool, 1)
    moveChan := make(chan struct{}) // Empty struct for synchronization
    
    // Start the dedicated polling goroutine
    go pollGameState(lobbyID, playerID, gameOverChan, moveChan)

    for {
        select {
        case <-gameOverChan:
            time.Sleep(2 * time.Second)
            return

        default:
            // Check ATB
            canMove, atb, err := checkATB(lobbyID, playerID)
            if err != nil || !canMove {
                fmt.Printf("\rWaiting for ATB: %d%%", atb)
                time.Sleep(300 * time.Millisecond)
                continue
            }

            // Signal poller to pause updates
            moveChan <- struct{}{}
            
            // Get player input
            fmt.Printf("\n[READY] Enter move (1:Attack 2:Power 3:Defend): ")
            move := readInput("")
            
            // Resume updates
            moveChan <- struct{}{}

            // Process move
            if move != "1" && move != "2" && move != "3" {
                fmt.Println("Invalid move! Use 1, 2, or 3.")
                continue
            }

            result, err := submitMove(lobbyID, playerID, move)
            if err != nil {
                fmt.Println("Error submitting move:", err)
                continue
            }

	
		
	 if !result.GameOver {
	    printBattleResult(result, playerID)
	 	}
	 
            // Check if game ended
            if result.GameOver {
             time.Sleep(2200 * time.Millisecond)
                return
            }
           
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

func pollGameState(lobbyID, playerID int, gameOverChan chan bool, moveChan chan struct{}) {
    ticker := time.NewTicker(300 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // Check if we should pause updates (during player input)
            select {
            case <-moveChan:
                // If we get a signal, pause updates until we get another
                <-moveChan
                continue
            default:
                // Proceed with normal update
            }

            // Fetch game state
            url := fmt.Sprintf("http://146.94.10.168:8080/gamestate?lobby_id=%d", lobbyID)
            resp, err := http.Get(url)
            if err != nil {
                continue
            }

            var state struct {
                P1Name    string `json:"p1_name"`
                P2Name    string `json:"p2_name"`
                P1Health  int    `json:"p1_health"`
                P2Health  int    `json:"p2_health"`
                GameOver  bool   `json:"game_over"`
                Winner    int    `json:"winner"`
            }

            if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
                resp.Body.Close()
                continue
            }
            resp.Body.Close()

            // Handle game over
            if state.GameOver {
                clearLine()
                if state.Winner == playerID {
                fmt.Println("\nYou struck first! Victory! You were The Winner!")
                time.Sleep(1800 * time.Millisecond)
                    fmt.Println("\nTaking you back to Main Menu")
                    
                } else {
                fmt.Println("\nYour Opponent was Faster! DEFEAT! You were defeated!")
                time.Sleep(1800 * time.Millisecond)
                fmt.Println("\nTaking you back to Main Menu")
                }
                gameOverChan <- true
                return
            }

            // Update display
            clearLine()
            if playerID == 1 {
                fmt.Printf("%s: %d HP | %s: %d HP", state.P1Name, state.P1Health, state.P2Name, state.P2Health)
            } else {
                fmt.Printf("%s: %d HP | %s: %d HP", state.P2Name, state.P2Health, state.P1Name, state.P1Health)
            }

        case <-gameOverChan:
            // Emergency exit
            return
        }
    }
}


func playerSelect(p1, p2, playerID int) int {
    if playerID == 1 {
        return p1
    }
    return p2
}

func printBattleResult(result *BattleResult, playerID int) {
    clearLine()
    fmt.Println("\n=== BATTLE RESULT ===")
    
    // Show health from player's perspective
    if playerID == 1 {
        fmt.Printf("%s: %d HP\n", result.P1Name, result.P1Health)
        fmt.Printf("%s: %d HP\n", result.P2Name, result.P2Health)
    } else {
        fmt.Printf("%s: %d HP\n", result.P2Name, result.P2Health)
        fmt.Printf("%s: %d HP\n", result.P1Name, result.P1Health)
    }

    // Show damage dealt or defend action
    if result.Damage > 0 {
        opponent := result.P2Name
        if playerID == 2 {
            opponent = result.P1Name
        }
        fmt.Printf("You dealt %d damage to %s!\n", result.Damage, opponent)
    } else {
        fmt.Println("You defended!")
    }
    time.Sleep(800 * time.Millisecond) // Pause for readability
}
