package gohq

import (
	"testing"
	"fmt"
	"log"
)

func TestHQ(t *testing.T) {
	account, err := New("ac4321NGx06pCFVHZEfSmD4k5caYE3NbR8utLrvduGJPYGTpkoctVdMGukC5VMFF")
	if err != nil {
		log.Fatal(err)
	}

	me, err := account.Me()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Lives: " + me.Lives)

	users, err := account.SearchUser("RedSinclair")
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range users.Data {
		fmt.Println(u.UserID)
		fmt.Println(account.AddFriend("21591399"))
	}

	schedule, err := account.Schedule()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(schedule.Upcoming)
}
