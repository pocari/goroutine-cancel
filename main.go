package main

import (
	"context"
	"fmt"
	"time"
)

func generator(_ctx context.Context) <-chan int {
	// バッファありチャンネルにしておかないと、 親でキャンセルされたタイミングで、
	// default側のチャンネルへの送信まちになってしまってブロックする
	c := make(chan int, 1)
	go func() {
		defer close(c)

		i := 0
		for {
			select {
			case <-_ctx.Done():
				fmt.Println("canceled")
				return
			default:
				c <- i
				i++
			}
		}
	}()

	return c
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	gen := generator(ctx)

	for i := 0; i < 5; i++ {
		fmt.Printf("i ... %d\n", <-gen)
	}

	cancel()
	time.Sleep(time.Second * 1)
}
