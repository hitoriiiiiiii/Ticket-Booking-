//Worker 
package notifications

import "log"

func StartWorker(repo *Repository) {

	go func() {
		log.Println("🔔 Notification Worker Started...")

		for job := range NotificationQueue {
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
				} else {
					log.Println("✅ Notification saved successfully")
				}

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
				} else {
					log.Println("✅ Payment notification saved successfully")
				}

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
				} else {
					log.Println("✅ Booking notification saved successfully")
				}

			default:
				log.Printf("⚠️ Unknown job type: %s", job.Type)
			}
		}
	}()
}
