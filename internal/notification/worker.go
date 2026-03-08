// Worker 
package notification

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hitorii/ticket-booking/internal/config"
	"github.com/hitorii/ticket-booking/internal/queue"
)

func StartWorker(repo *Repository) {
	// Initialize Redis connection
	cfg := config.Load()
	
	log.Println("🔌 Connecting to Redis...")
	if err := queue.InitRedis(cfg.RedisURL); err != nil {
		log.Printf("❌ Failed to connect to Redis: %v", err)
		log.Println("⚠️  Worker cannot start without Redis - job queue will be unavailable")
		return
	}
	
	log.Println("✅ Connected to Redis successfully")
	
	// Get hostname for consumer identification
	hostname, _ := os.Hostname()
	consumerName := fmt.Sprintf("worker-%s-%d", hostname, time.Now().UnixNano())
	
	log.Printf("👤 Starting Redis consumer: %s", consumerName)
	
	// Start consuming from Redis stream
	go func() {
		err := queue.ConsumeJobs(context.Background(), consumerName, func(job queue.JobPayload) error {
			log.Printf("📩 Processing Job: Type=%s, UserID=%s, Message=%s", job.Type, job.UserID, job.Message)

			switch job.Type {
			case JobTypeNotification:
				err := repo.Save(Job{
					Type:    job.Type,
					UserID:  job.UserID,
					Message: job.Message,
				})
				if err != nil {
					log.Println("❌ Failed to save notification:", err)
					return err
				}
				log.Println("✅ Notification saved successfully")

			case JobTypeEmail:
				// Process email job (in a real app, this would send an email)
				log.Printf("📧 Sending email to user %s: %s", job.UserID, job.Message)
				log.Println("✅ Email job processed")

			case JobTypePayment:
				// Process payment notification
				log.Printf("💳 Processing payment notification for user %s: %s", job.UserID, job.Message)
				err := repo.Save(Job{
					Type:    job.Type,
					UserID:  job.UserID,
					Message: job.Message,
					Data:    job.Data,
				})
				if err != nil {
					log.Println("❌ Failed to save payment notification:", err)
					return err
				}
				log.Println("✅ Payment notification saved successfully")

			case JobTypeBooking:
				// Process booking notification
				log.Printf("🎫 Processing booking notification for user %s: %s", job.UserID, job.Message)
				err := repo.Save(Job{
					Type:    job.Type,
					UserID:  job.UserID,
					Message: job.Message,
					Data:    job.Data,
				})
				if err != nil {
					log.Println("❌ Failed to save booking notification:", err)
					return err
				}
				log.Println("✅ Booking notification saved successfully")

			default:
				log.Printf("⚠️ Unknown job type: %s", job.Type)
			}
			
			return nil
		})
		
		if err != nil && err != context.Canceled {
			log.Printf("❌ Error in Redis consumer: %v", err)
			log.Println("⚠️  Restarting consumer in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}()
}
