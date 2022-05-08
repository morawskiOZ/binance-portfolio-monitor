package task

import (
	"fmt"
	"sync"
	"time"

	"github.com/morawskioz/binance-monitor/interal/binance"
	"github.com/morawskioz/binance-monitor/pkg/mail"
)

// BinanceTask is a task that implements the tasker interface so it can be scheduled
type BinanceTask struct {
	BinanceClient *binance.Client
	MailClient    *mail.Client
	Recipient     string
}

// RunTask runs the task
func (bt *BinanceTask) RunTask(wg *sync.WaitGroup) error {
	defer wg.Done()
	portfolioValue, err := bt.BinanceClient.GetPortfolioTotalValue()
	if err != nil {
		return err
	}

	if portfolioValue > 27000 {
		msg := fmt.Sprintf("Time to sell, portfolio value is: %v", portfolioValue)
		err := bt.sendNotification(msg, bt.Recipient, "Portfolio alert")
		if err != nil {
			return err
		}
	}

	return nil
}

// SetupTicker returns a ticker that will run the task every 60 minutes
func (bt *BinanceTask) SetupTicker() *time.Ticker {
	return time.NewTicker(time.Second * 10)
}

func (bt *BinanceTask) sendNotification(msg string, recipient string, subject string) error {
	if err := bt.MailClient.Send(recipient, subject, msg); err != nil {
		return err
	}

	return nil
}