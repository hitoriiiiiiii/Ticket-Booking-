//Worker 
package notifications

import "log"

func StartWorker(repo *Repository) {

	go func() {
		log.Println("🔔 Notification Worker Started...")

		for job := range NotificationQueue {

			log.Println("📩 Processing Notification:", job.Message)

			err := repo.Save(job)
			if err != nil {
				log.Println("❌ Failed to save notification:", err)
			} else {
				log.Println("✅ Notification saved successfully")
			}
		}
	}()
}
