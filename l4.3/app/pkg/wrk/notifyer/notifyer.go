package notifyer

import (
	"app/internal/models"
	"app/pkg/logger"
	"app/pkg/sender"
	"container/heap"
	"context"
	"fmt"
	"sync"
	"time"
)

const BufSize = 50

type PriorityQueue []models.EventTask

func (pq PriorityQueue) Len() int {
	return len(pq)
}
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Event.Date.Before(pq[j].Event.Date)
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(models.EventTask))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

type Notifyer struct {
	pq          PriorityQueue
	timer       *time.Timer
	emailSender *sender.Sender
}

func New(snd *sender.Sender) *Notifyer {
	s := &Notifyer{
		pq:          make(PriorityQueue, 0),
		emailSender: snd,
	}
	heap.Init(&s.pq)
	return s
}

func (s *Notifyer) Start(ctx context.Context, wg *sync.WaitGroup, taskChan chan models.EventTask) {
	for {
		if s.pq.Len() > 0 {
			nextTask := s.pq[0]
			now := time.Now()

			if nextTask.Event.Date.Before(now) {
				s.sendNotification(ctx, nextTask)
				heap.Pop(&s.pq)
				continue
			}

			if s.timer == nil {
				s.timer = time.NewTimer(nextTask.Event.Date.Sub(now))
			} else {
				s.timer.Reset(nextTask.Event.Date.Sub(now))
			}
		}
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case newTask := <-taskChan:
			heap.Push(&s.pq, newTask)
			if s.timer != nil {
				s.timer.Stop()
			}
		case <-func() <-chan time.Time {
			if s.timer != nil {
				return s.timer.C
			}
			return nil
		}():
			if s.pq.Len() > 0 {
				task := heap.Pop(&s.pq).(models.EventTask)
				s.sendNotification(ctx, task)
			}
		}
	}
}

func (s *Notifyer) sendNotification(ctx context.Context, task models.EventTask) {
	lg := logger.LoggerFromCtx(ctx).Lg
	err := s.emailSender.SendReminder(task.Email, task.Event)
	if err != nil {
		lg.Error().Str("worker", "notifyer").Err(err).Msg(fmt.Sprintf("Failed to send email for event %s", task.Event.EventId))
	} else {
		lg.Info().Str("worker", "notifyer").Msg(fmt.Sprintf("Reminder sent for event %s to %s", task.Event.EventId, task.Email))
	}
}
