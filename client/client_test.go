package client

import (
	"context"
	"fmt"
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
	defer rdb.Close()
    err := rdb.Set(ctx, "key", "value", 0).Err()
    if err != nil {
        panic(err)
    }
	val, err := rdb.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("key", val)

    val2, err := rdb.Get(ctx, "key2").Result()
    if err == redis.Nil {
        fmt.Println("key2 does not exist")
    } else if err != nil {
        panic(err)
    } else {
        fmt.Println("key2", val2)
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
