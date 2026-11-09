package cron

import (
	"github.com/samber/do"
)

func Init(injector *do.Injector) {
	// Register the reminder cron service
	do.Provide(injector, NewReminderCronService)

	// Start the centralized reminder cron service
	reminderCronService := do.MustInvoke[ReminderCronService](injector)
	reminderCronService.Start()
}
