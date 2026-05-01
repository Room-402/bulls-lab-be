package jobs

import (
	"context"
	"log"
	"bulls-lab-be/internal/core/ports"
)

type CancelExpiredOrdersJob struct {
	orderService ports.OrderService
}

func NewCancelExpiredOrdersJob(orderService ports.OrderService) *CancelExpiredOrdersJob {
	return &CancelExpiredOrdersJob{
		orderService: orderService,
	}
}

func (j *CancelExpiredOrdersJob) Name() string {
	return "CancelExpiredOrders"
}

func (j *CancelExpiredOrdersJob) Run() {
	ctx := context.Background()
	log.Println("[CancelExpiredOrdersJob] Checking for expired orders...")
	count, err := j.orderService.CancelExpiredOrders(ctx)
	if err != nil {
		log.Printf("[CancelExpiredOrdersJob] Error cancelling expired orders: %v", err)
		return
	}
	log.Printf("[CancelExpiredOrdersJob] Finished. Cancelled %d expired orders", count)
}
