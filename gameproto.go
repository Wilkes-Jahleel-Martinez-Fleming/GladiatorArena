package main

import (
	"fmt"
	"time"
	"math/rand"
)

/*
Use goroutines to manage the players' moves. Maybe use a switch statement (select statement) to differentiate between the cases of gladiators' attack speeds.
We need a playerCount that stores the current amount of players in the server. Maybe we could use concurrency to decide what players get paired together. Maybe each fight
would be a goroutine.

*/

func init() {
	rand.Seed(time.Now().UnixNano())
}
type Gladiator struct {
	id int
	health int
	attack int
	defense int
	speed int
	isDead bool
}

func genGladiator(id int) Gladiator {

	return Gladiator{
		id: id,
		health: rand.Intn(25) + 75,
		attack: rand.Intn(25) + 10,
		defense: rand.Intn(5) + 5,
		speed: rand.Intn(10) + 5,
	}
}

/*
func gladProducer(ch chan Gladiator, numGlad int) {
	for i := 1; i <= numGlad; i++ {
		
		glad := genGladiator(i)

		ch <- glad
		fmt.Printf("Produced Gladiator %d\n", glad.id)

		time.Sleep(time.Millisecond * 300)
	}
	close(ch)
}
*/

func attack(rglad *Gladiator, gglad Gladiator) {
	dmg := gglad.attack + rand.Intn(10) - 5 - rglad.defense - rand.Intn(5) + 2
	rglad.health = rglad.health - dmg

	fmt.Printf("\nPlayer %d gladiator took %d damage.\n\n", rglad.id, dmg)
	
}

func powAttack(rglad *Gladiator, gglad *Gladiator) {
	gglad.speed = 0
	dmg := gglad.attack * 2 + rand.Intn(10) - 5 - rglad.defense - rand.Intn(5) + 2
	rglad.health = rglad.health - dmg

	fmt.Printf("\nPlayer %d gladiator took %d damage.\n\n", rglad.id, dmg)
}

func defend(glad *Gladiator){
glad.defense += 5
fmt.Printf("\nPlayer %d gladiator is defending. Defense increased.\n\n", glad.id)
}


func main() {
//	rand.Seed(time.Now().UnixNano())
	var p1move int
	var p2move int
	p1glad := genGladiator(1)
	p2glad := genGladiator(2)

	p1baseSpeed := p1glad.speed
	p2baseSpeed := p2glad.speed

	fmt.Printf("P1 Gladiator Stats: ")
	
	for p1glad.health > 0 && p2glad.health > 0 {
		fmt.Printf("P1 gladiator health = %d\n", p1glad.health)
		fmt.Println("P1 move (1 for attack 2 for defend):")
		fmt.Scan(&p1move)
		fmt.Printf("P2 gladiator health = %d\n", p2glad.health)
		fmt.Println("P2 move (1 for attack 2 for defend):")
		fmt.Scan(&p2move)
		
		// Adjust speed after moves are decided
		if p1move == 2 {
			p1glad.speed = 0 // Powerful attack reduces speed
		}
		if p2move == 2 {
			p2glad.speed = 0 // Powerful attack reduces speed
		}

		if p1glad.speed > p2glad.speed {

   			switch p1move {
  			case 1:
        			attack(&p2glad, p1glad)
			case 2:
        			powAttack(&p2glad, &p1glad)
			case 3:
				defend(&p1glad)
    			default:
				fmt.Println("Invalid move entered. Turn skipped.")
    			}

    			switch p2move {
    			case 1:
        			attack(&p1glad, p2glad)
    			case 2:
        			powAttack(&p1glad, &p2glad)
			case 3:
				defend(&p2glad)
    			default:
        			fmt.Println("Invalid move entered. Turn skipped.")
    			}
		} else {
    			switch p2move {
    			case 1:
        			attack(&p1glad, p2glad)
    			case 2:
        			powAttack(&p1glad, &p2glad)
			case 3:
				defend(&p2glad)
    			default:
        			fmt.Println("Invalid move entered. Turn skipped.")
    			}
    			
    			switch p1move {
    			case 1:
        			attack(&p2glad, p1glad)
    			case 2:
        			powAttack(&p2glad, &p1glad)
			case 3:
				defend(&p1glad)
    			default:
        			fmt.Println("Invalid move entered. Turn skipped.")
    			}
		}
		
		// Reset speeds after the round
		p1glad.speed = p1baseSpeed
		p2glad.speed = p2baseSpeed

	}

	if p1glad.health <= 0 {
		fmt.Println("P2 wins!")
	} else {
		fmt.Println("P1 wins!")
	}

	/*
	rand.Seed(time.Now().UnixNano())

	gladChan := make(chan Gladiator)

	go gladProducer(gladChan, 5)

	for glad := range gladChan {

		fmt.Printf("Gladiator %d info: Health = %d, Attack = %d, Defense = %d\n", glad.id, glad.health, glad.attack, glad.defense)
	}
*/
}
