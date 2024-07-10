package utils

import (
	"gotransact/apps/transaction/utils"
	"log"

	"github.com/robfig/cron"
)

func Cron() {
	c := cron.New()
	c.AddFunc("@every 1h", func() {
		transactions := utils.FetchTransactionsLast24Hours()
		filePath, err := utils.GenerateExcel(transactions)
		if err != nil {
			log.Fatalf("failed to generate excel: %v", err)
		}
		utils.SendMailWithAttachment("sangadevshatwar143@gmail.com", filePath)
	})
	c.Start()
}
