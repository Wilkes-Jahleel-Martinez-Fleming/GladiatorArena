1.0 VERSION are files: updatedClient.go and updatedServer.go

2.0 VERSION with ATB combat are files: updatedClientATB.go and updatedServerATB.go


This README will be updated with progress.

added first revision of code as separate instances. next step is to implement code together to create simple game simulation between 2 players.

04/10/2025
polish more client and QoL for user experience. look into more robust changes.
implement web interface for the game.
allow nickname instead of player 1, 2, etc

04/11/2025
Andrew's Tasks

Change Defense logic on gameproto.go
add dmg value to gladiator struct
use that dmg value in the attack, powattack and defend functions
remove the health changing inside of the attack moves
change defense to be half the value of dmg in gladiator struct
apply a execute move func that just changes the health value based on the dmg.


Jahleel's Tasks

1)based on the updates that andrew do to the gameproto.go
change the move logic to be more robust and in line with Andrew's Changes.
2)(DONE)Add a display of gladiator values for each player such as their stats.
3)change their nickname on the client. and server so instead of p1 p2 it would be more custumizable.
4)add 2sec delay per message on client console.

extra task for both:

1)Allow each player to pick from 10 random gladiators generated for their 1 vs 1 battle.
