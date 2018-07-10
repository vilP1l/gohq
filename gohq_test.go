package gohq

import (
	"testing"
	"fmt"
	"log"
)

func TestHQ(t *testing.T) {
	account, err := New("not sure if this is important")
	if err != nil {
		log.Fatal(err)
	}

	me, err := account.Me()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Lives: " + me.Lives)

	users, err := account.SearchUser("Steve150")
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range users.Data {
		fmt.Println(u.UserID)
	}

	schedule, err := account.Schedule()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(schedule.Upcoming)
}
