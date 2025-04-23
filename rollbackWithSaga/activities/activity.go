package activities

import (
	"context"
	"fmt"
	"time"
)

type OrderDetails struct {
	OrderID string
	Amount  float64
	Address string
}

func ProcessOrder(ctx context.Context) error {
	fmt.Println("Processing order")
	// Simulate processing
	time.Sleep(1 * time.Second)
	return nil
}

func CancelOrder(ctx context.Context) error {
	fmt.Println("Cancelling order")
	// Simulate cancellation
	time.Sleep(1 * time.Second)
	return nil
}

func ChargePayment(ctx context.Context) error {
	fmt.Println("Charging payment")
	// Simulate payment charge
	time.Sleep(1 * time.Second)
	return nil
}

func RefundPayment(ctx context.Context) error {
	fmt.Println("Refunding payment")
	// Simulate payment refund
	time.Sleep(1 * time.Second)
	return nil
}

func ShipOrder(ctx context.Context) error {
	fmt.Println("Shipping order")
	// Simulate order shipment
	time.Sleep(1 * time.Second)
	return fmt.Errorf("failed to ship order")
}

func CancelShipment(ctx context.Context) error {
	fmt.Println("Cancelling shipment")
	// Simulate shipment cancellation
	time.Sleep(1 * time.Second)
	return nil
}

// @@@SNIPEND
