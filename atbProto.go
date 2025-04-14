package main

import (
        "fmt"
        "math/rand"
        "time"
)

type Gladiator struct {
        id      int
        name    string
        health  int
        attack  int
        defense int
        speed   int
}

//Command represents move
type Command struct {
        playerID int
        move     string
}

//Generate random gladiator
func genGladiator(id int) Gladiator {
        names := []string{"g1", "g2", "g3", "g4", "g5", "g6", "g7", "g8", "g9", "g10"}
        return Gladiator{
                id:      id,
                name:    names[rand.Intn(len(names))],
                health:  rand.Intn(30) + 70,
                attack:  rand.Intn(15) + 10,
                defense: rand.Intn(6) + 5,
                speed:   rand.Intn(5) + 5,
        }
}

//Handles each gladiator's independent timer and input loop
func actionTimer(glad Gladiator, commandChan chan<- Command) {
        for glad.health > 0 {
                wait := time.Duration(5000-glad.speed*5) * time.Millisecond
                time.Sleep(wait)

                fmt.Printf("\n[%s READY] Enter move (1: attack, 2: defend): ", glad.name)
                var input int
                fmt.Scan(&input)

                move := "wait"
                switch input {
                case 1:
                        move = "attack"
                case 2:
                        move = "defend"
                default:
                        fmt.Println("Invalid input. Skipping turn.")
                }

                commandChan <- Command{playerID: glad.id, move: move}
        }
}

//Applies command from player
func applyMove(cmd Command, p1, p2 *Gladiator) {
        var attacker, defender *Gladiator
        if cmd.playerID == p1.id {
                attacker = p1
                defender = p2
        } else {
                attacker = p2
                defender = p1
        }

        switch cmd.move {
        case "attack":
                dmg := attacker.attack + rand.Intn(6) - defender.defense
                if dmg < 0 {
                        dmg = 0
                }
                defender.health -= dmg
                fmt.Printf("\n %s attacks %s for %d damage! (%d HP left)\n", attacker.name, defender.name, dmg, defender.health)
        case "defend":
                attacker.defense += 3
                fmt.Printf("\n %s defends! (+3 defense this turn)\n", attacker.name)
                time.AfterFunc(2*time.Second, func() {
                        attacker.defense -= 3
                        fmt.Printf("\n %s's defense boost wore off.\n", attacker.name)
                })
        default:
                fmt.Println("No valid move received.")
        }
}

func main() {
        rand.Seed(time.Now().UnixNano())

        p1 := genGladiator(1)
        p2 := genGladiator(2)

        fmt.Printf("\n Player 1 Gladiator: %+v\n", p1)
        fmt.Printf(" Player 2 Gladiator: %+v\n", p2)

        //shared channel
        commandChan := make(chan Command)

        //independent action timers
        go actionTimer(p1, commandChan)
        go actionTimer(p2, commandChan)

        //listen for and apply commands
        for {
                if p1.health <= 0 || p2.health <= 0 {
                        break
                }
                cmd := <-commandChan
                applyMove(cmd, &p1, &p2)
        }

        fmt.Println("\n Game Over")
        if p1.health <= 0 && p2.health <= 0 {
                fmt.Println("It's a draw!")
        } else if p1.health <= 0 {
                fmt.Printf(" %s wins!\n", p2.name)
        } else {
                fmt.Printf(" %s wins!\n", p1.name)
        }
}