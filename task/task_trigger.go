package task

import "time"

//startAt,每天的这个时间点运行任务
func StartTask(startAt int, f func()) {
	go func() {
		for {
			f()

			//寻找下一个20点
			now := time.Now()
			next := now
			if now.Hour() >= startAt {
				next = now.Add(time.Hour * 24)
			}
			next = time.Date(next.Year(), next.Month(), next.Day(), startAt, 0, 0, 0, next.Location())

			t := time.NewTicker(next.Sub(now))
			<-t.C
		}
	}()
}
