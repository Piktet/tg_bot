package speachservice

import (
	"context"

	"github.com/Piktet/tg_bot/internal/model"

	"golang.org/x/sync/errgroup"
)

type SpeachService struct {
	*speachOption
	chTaskRequest  chan *model.SpeachTaskData
	chTaskResponse chan *model.SpeachTaskResponse
}

func New(opts ...Option) *SpeachService {
	return &SpeachService{
		speachOption: newSpeachOption(opts...),
	}
}

func (p *SpeachService) AddTask(x *model.SpeachTaskData) {
	p.chTaskRequest <- x
}

func (p *SpeachService) GetTaskResponse() chan *model.SpeachTaskResponse {
	return p.chTaskResponse
}

func (p *SpeachService) Start(ctx context.Context, cnt, size int) error {
	p.chTaskRequest = make(chan *model.SpeachTaskData, size)
	p.chTaskResponse = make(chan *model.SpeachTaskResponse, size)
	defer close(p.chTaskRequest)
	defer close(p.chTaskResponse)
	var wg errgroup.Group
	for range cnt {
		wg.Go(func() error {
			return p.Worker(ctx)
		})
	}
	return wg.Wait()
}

func (p *SpeachService) Worker(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		case taskData := <-p.chTaskRequest:
			task := NewTask(p.speachOption)
			if x, err := task.Process(ctx, taskData); err == nil {
				p.chTaskResponse <- x
			}
		}
	}

}
