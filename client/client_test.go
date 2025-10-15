package client

import (
	"context"
	"fmt"
	"log"
	"testing"
	"sync"
)

func Test_Client(t *testing.T) {
	total_client := 10
	wg := sync.WaitGroup{}
	wg.Add(total_client)

	for i := 0; i < total_client; i++ {
		go func() {

			cl, err := New(":5000")
			defer cl.Close()
			if err != nil {
				log.Fatal(err)
			}
			if response, err := cl.Set(context.Background(), fmt.Sprintf("admin_%d", i), fmt.Sprintf("Apurv_%d", i)); err != nil {
				log.Fatal(err)
			} else {
				fmt.Println(response)
			}
			if response, err := cl.Get(context.Background(), fmt.Sprintf("admin_%d", i)); err != nil {
				log.Fatal(err)
			} else {
				fmt.Println(response)
			}
			wg.Done()
		}()
		
	}
	wg.Wait()
}
