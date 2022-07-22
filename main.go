package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	rxgo "github.com/reactivex/rxgo/v2"
)

type Customer struct {
	ID             int
	Name, LastName string
	Age            int
	TaxNumber      string
}

func produce(ch chan rxgo.Item) {
	// to handle pathing, map
	var dat []interface{}
	ok := json.Unmarshal([]byte("[{\"x\": {\"y\": 4}}]"), &dat)
	if ok != nil {
		not_ok := ok
		panic(not_ok)
	}
	for _, e := range dat {
		ch <- rxgo.Item{V: e}
	}
}

func main() {
	ch := make(chan rxgo.Item)
	go produce(ch)
	config := "path"
	path := []string{"x", "y"}
	observable := rxgo.FromChannel(ch).
		Map(func(_ context.Context, e interface{}) (interface{}, error) {
			if config == "path" {
				var curs any = e
				for _, frag := range path {
					as_map, ok := curs.(map[string]interface{})
					if !ok {
						panic("todo")
					}
					curs = as_map[frag]
				}
				return curs, nil
			}
			return nil, errors.New("idfk what you gave me dude")
		})
	out_ch := observable.Observe()
	x := <-out_ch
	fmt.Printf("It worked! %s", x.V)
}
