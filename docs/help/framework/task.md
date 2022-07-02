定时任务配置
```cassandraql
// c.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
// c.AddFunc("@hourly",      func() { fmt.Println("Every hour") })
// c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty") })
// c.AddFunc("@every 5s", func() { fmt.Println("Every hour thirty") })

// 每5秒执行一次
// c.AddFunc("@every 2m", r.syncAnnouncement)
// c.AddFunc("0 0 0 * * *", r.syncUser)
// c.AddFunc("@every 3s", r.syncAnnouncement)
```