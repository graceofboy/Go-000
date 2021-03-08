package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	svr := http.NewServeMux()

	g.Go(func() error {
		fmt.Println("http")
		go func() {
			<-ctx.Done() //?
			fmt.Println("http ctx done")
			svr.Shutdown(context.TODO()) //?
		}()
		return svt.Start()
	})

	g.Go(func() error {
		exitSignal := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT} // ?
		sig := make(chan os.Signal, len(exitSignal))                                              //?
		signal.Notify(sig, exitSignal...)
		for {
			fmt.Println("signal")
			select {
			case <-ctx.Done():
				fmt.Println("signal ctx done")
				return ctx.Err()
			case <-sig:
				return nil
			}
		}
	})

	g.Go(func() error {
		fmt.Println("sleep")
		time.Sleep(time.Second)
		return errors.New("sleep error")
	})
	err := g.Wait()
	fmt.Println(err)
}
