package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Gladiator struct {
	id      int
	health  int
	attack  int
	defense int
	speed   int
	isDead  bool
}

//Move holds the player's chosen action and any calculated damage
type Move struct {
	moveType string // "attack", "powAttack", "defend"
	dmg      int
}

func genGladiator(id int) Gladiator {
	return Gladiator{
		id:      id,
		health:  rand.Intn(25) + 75,
		attack:  rand.Intn(25) + 10,
		defense: rand.Intn(5) + 5,
		speed:   rand.Intn(10) + 5,
	}
}

//generates damage but does NOT apply it
func attack(attacker Gladiator, defender Gladiator) Move {
	dmg := attacker.attack + rand.Intn(10) - 5 - defender.defense - rand.Intn(5) + 2
	if dmg < 0 {
		dmg = 0
	}
	return Move{moveType: "attack", dmg: dmg}
}

//generates higher damage and sets attacker's speed to 0
func powAttack(attacker *Gladiator, defender Gladiator) Move {
	attacker.speed = 0
	dmg := attacker.attack*2 + rand.Intn(10) - 5 - defender.defense - rand.Intn(5) + 2
	if dmg < 0 {
		dmg = 0
	}
	return Move{moveType: "powAttack", dmg: dmg}
}

//applies both players' moves, taking defend into account
func executeMoves(p1 *Gladiator, p2 *Gladiator, m1 Move, m2 Move) {
	//Apply defense buffs before damage
	if m1.moveType == "defend" {
		p1.defense += 5
		fmt.Printf("P1 is defending. Defense temporarily increased.\n")
	}
	if m2.moveType == "defend" {
		p2.defense += 5
		fmt.Printf("P2 is defending. Defense temporarily increased.\n")
	}

	//Apply damage from P1 to P2
	if m1.moveType == "attack" || m1.moveType == "powAttack" {
		if m2.moveType == "defend" {
			m1.dmg /= 2
		}
		p2.health -= m1.dmg
		fmt.Printf("P2 took %d damage from P1.\n", m1.dmg)
	}

	//Apply damage from P2 to P1
	if m2.moveType == "attack" || m2.moveType == "powAttack" {
		if m1.moveType == "defend" {
			m2.dmg /= 2
		}
		p1.health -= m2.dmg
		fmt.Printf("P1 took %d damage from P2.\n", m2.dmg)
	}

	//Reset defense to normal
	if m1.moveType == "defend" {
		p1.defense -= 5
	}
	if m2.moveType == "defend" {
		p2.defense -= 5
	}
}

func main() {
	var p1move, p2move int
	p1glad := genGladiator(1)
	p2glad := genGladiator(2)

	p1baseSpeed := p1glad.speed
	p2baseSpeed := p2glad.speed

	fmt.Printf("P1 Gladiator Stats: %+v\n", p1glad)
	fmt.Printf("P2 Gladiator Stats: %+v\n", p2glad)

	for p1glad.health > 0 && p2glad.health > 0 {
		fmt.Printf("\n--- New Round ---\n")
		fmt.Printf("P1 health: %d | P2 health: %d\n", p1glad.health, p2glad.health)

		fmt.Println("P1 move (1: attack, 2: pow attack, 3: defend):")
		fmt.Scan(&p1move)

		fmt.Println("P2 move (1: attack, 2: pow attack, 3: defend):")
		fmt.Scan(&p2move)

		//Reset moves
		var m1, m2 Move

		//Get P1's move
		switch p1move {
		case 1:
			m1 = attack(p1glad, p2glad)
		case 2:
			m1 = powAttack(&p1glad, p2glad)
		case 3:
			m1 = Move{moveType: "defend"}
		default:
			fmt.Println("Invalid move from P1.")
			m1 = Move{moveType: "none"}
		}

		//Get P2's move
		switch p2move {
		case 1:
			m2 = attack(p2glad, p1glad)
		case 2:
			m2 = powAttack(&p2glad, p1glad)
		case 3:
			m2 = Move{moveType: "defend"}
		default:
			fmt.Println("Invalid move from P2.")
			m2 = Move{moveType: "none"}
		}

		executeMoves(&p1glad, &p2glad, m1, m2)

		//reset speed
		p1glad.speed = p1baseSpeed
		p2glad.speed = p2baseSpeed
	}

	fmt.Println("\n--- Game Over ---")
	if p1glad.health <= 0 && p2glad.health <= 0 {
		fmt.Println("It's a draw!")
	} else if p1glad.health <= 0 {
		fmt.Println("P2 wins!")
	} else {
		fmt.Println("P1 wins!")
	}
}
