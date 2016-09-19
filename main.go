package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/net/context"
)

type Payment struct {
	Name   string
	Amount int64
}

const (
	confirmedKey   = "confirmed.context.key"
	transactionKey = "transaction.context.key"
)

func ProcessPayment(ctx context.Context) {
	confirmed := ctx.Value(confirmedKey).(chan bool)
	payment := ctx.Value(transactionKey).(Payment)

	for {
		select {
		case <-confirmed:
			fmt.Printf("Transaksi sebesar Rp. %d berhasil ", payment.Amount)
			fmt.Printf("a/n %s\n", payment.Name)
			return
		case <-ctx.Done():
			if ctx.Err() == context.Canceled {
				fmt.Printf("Transaksi pembayaran anda dibatalkan! \nUang sebesar Rp. %d dikembalikan.\n", payment.Amount)
				return
			} else if ctx.Err() == context.DeadlineExceeded {
				fmt.Println("\nTransaksi pembayaran anda kadaluarsa. Silahkan kembali lain waktu.")
				os.Exit(0)
			}
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	var (
		cancel        context.CancelFunc
		amount        int64
		name, confirm string
		payment       Payment
	)

	fmt.Print("\nMasukkan nominal pembayaran:\nRp. ")
	fmt.Scan(&amount)

	fmt.Println("\nIsikan nama anda:")
	fmt.Scanln(&name)

	payment.Amount = amount
	payment.Name = name
	if name != "" && amount >= 0 {
		confirmed := make(chan bool)
		ctx := context.WithValue(context.Background(), confirmedKey, confirmed)
		ctx = context.WithValue(ctx, transactionKey, payment)
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)

		go ProcessPayment(ctx)

		fmt.Printf("\nAnda akan melakukan pembayaran sebesar Rp. %d \n", amount)
		fmt.Printf("A/n %s \n", name)

		for {
			fmt.Printf("[K]onfirmasi, (B)atalkan: ")
			fmt.Scanln(&confirm)

			switch confirm {
			case "K":
				confirmed <- true
				return
			case "B":
				cancel()
				return
			default:
				fmt.Printf("\nPilihan anda tidak tersedia: %s. Silahkan coba lagi.\n", confirm)
				fmt.Println("Mohon konfirmasi pembayaran anda:")
			}
		}
	} else if name != "" && amount < 0 {
		fmt.Println("Nilai transaksi tidak boleh kosong")
	} else if name == "" && amount >= 0 {
		fmt.Println("Nama tidak boleh kosong")
	}
}
