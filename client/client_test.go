package client

import (
	"context"
	// "fmt"
	// "log"
	"testing"
	// "sync"
	"github.com/redis/go-redis/v9"
)

// func Tes_Client(t *testing.T) {
// 	total_client := 10
// 	wg := sync.WaitGroup{}
// 	wg.Add(total_client)

// 	for i := 0; i < total_client; i++ {
// 		go func() {

// 			cl, err := New(":5000")
// 			defer cl.Close()
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			if response, err := cl.Set(context.Background(), fmt.Sprintf("admin_%d", i), [3]int{1,2,3}); err != nil {
// 				log.Fatal(err)
// 			} else {
// 				fmt.Println(response)
// 			}
// 			if response, err := cl.Get(context.Background(), fmt.Sprintf("admin_%d", i)); err != nil {
// 				log.Fatal(err)
// 			} else {
// 				fmt.Printf("%T",response)
// 			}
// 			wg.Done()
// 		}()

// 	}
// 	wg.Wait()
// }


func Test_Client1(t *testing.T) {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:5000",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
	
    err := rdb.Set(ctx, "key", "value", 0).Err()
    if err != nil {
        panic(err)
    }
	// total_client := 10
	// wg := sync.WaitGroup{}
	// wg.Add(total_client)

	// for i := 0; i < total_client; i++ {
	// 	go func() {

	// 		cl, err := New(":5000")
	// 		defer cl.Close()
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		if response, err := cl.Set(context.Background(), fmt.Sprintf("admin_%d", i), [3]int{1,2,3}); err != nil {
	// 			log.Fatal(err)
	// 		} else {
	// 			fmt.Println(response)
	// 		}
	// 		if response, err := cl.Get(context.Background(), fmt.Sprintf("admin_%d", i)); err != nil {
	// 			log.Fatal(err)
	// 		} else {
	// 			fmt.Printf("%T",response)
	// 		}
	// 		wg.Done()
	// 	}()
		
	// }
	// wg.Wait()
}
