package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"
)

type Payment struct {
	Name   string
	Amount int64
}

const (
	confirmedKey = "confirmed.context.key"
)

func ProcessPayment(ctx context.Context, payment *Payment) {
	confirmed := ctx.Value(confirmedKey).(chan bool)

	for {
		select {
		case <-confirmed:
			fmt.Printf("Transaksi sebesar Rp. %d berhasil.\n", payment.Amount)
			return
		case <-ctx.Done():
			if ctx.Err() == context.Canceled {
				fmt.Printf("Transaksi pembayarn anda dibatalkan. Uang sebesar Rp. %d dikembalikan.\n", payment.Amount)
				return
			} else if ctx.Err() == context.DeadlineExceeded {
				fmt.Println("Transaksi pembayaran anda kadaluarsa. Silahkan kembali lain waktu.")
				os.Exit(0)
			}
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	var cancel context.CancelFunc

	confirmed := make(chan bool)
	ctx := context.WithValue(context.Background(), confirmedKey, confirmed)
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)

	go ProcessPayment(ctx, &Payment{
		Name:   "John Doe",
		Amount: 1000000})

	fmt.Print("Transaksi pembayaran anda tertunda. ")
	if deadline, ok := ctx.Deadline(); ok {
		fmt.Printf("Anda mempunyai waktu sekitar %s untuk menyelesaikan pembayaran.\n", deadline.Sub(time.Now()).String())
	}

	fmt.Println()
	fmt.Println("Pilih salah satu pilihan untuk konfirmasi pembyaran:")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("[K]onfirmasi, (B)atalkan: ")
		if line, err := reader.ReadString('\n'); err == nil {
			command := strings.TrimSuffix(line, "\n")
			switch command {
			case "K":
				confirmed <- true
				time.Sleep(500 * time.Millisecond)
				return
			case "B":
				cancel()
				time.Sleep(500 * time.Millisecond)
				return
			default:
				fmt.Printf("\nPilihan anda tidak tersedia: %s. Silahkan coba lagi.\n", command)
				fmt.Println("Mohon konfirmasi pemayaran anda:")
			}
		}
	}
}
