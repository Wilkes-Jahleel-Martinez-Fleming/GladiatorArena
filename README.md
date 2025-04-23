How to Run

either Run UpdatedClient.go with UpdatedServer.go 
OR
updatedClientATB.go with updatedServerATB.go

Prerequisites
    • - Go 1.18+ installed

Server
1. Open a terminal.
2. Navigate to the directory containing `updatedServerATB.go`.
-Modify file to use local IP address. (optional) 
3. Run the server with:
   go run updatedServerATB.go

Client
1. Open a new terminal window.
2. Navigate to the directory containing `updatedClientATB.go`.
-Modify file to use local IP address. (optional) 
3. Run the client with:
   go run updatedClientATB.go

How It Works

Server
    • - Manages game states, players, and lobbies.
    • - Uses goroutines to handle each HTTP request concurrently.
    • - Updates the ATB bars on a timed loop in the background.

Client
    • - Continuously polls the server to get updated game state using a goroutine.
    • - Displays the player's current status and accepts input during ATB activation.

Concurrency Details
    • - Goroutines manage polling and game ticking independently.
    • - Channels are used to pause/resume polling during player input.
    • - Shared data is protected with `sync.Mutex` to prevent race conditions.
