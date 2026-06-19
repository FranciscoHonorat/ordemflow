package main

import (
	"fmt"

	"github.com/FranciscoHonorat/ordemflow/shared/events"
)

func main() {
	var id events.EventID = "example-event-id"
	fmt.Println(id)
}
