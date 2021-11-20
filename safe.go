package safe

import (
	"context"
	"fmt"
	"sync"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

type Group struct {
	cancel  func()
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

func (g *Group) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Infof("Panic Recover(%s)",err)
				fmt.Println(string(debug.Stack()))
			}
		}()
		defer g.wg.Done()
		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}

func Go(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Infof("Panic Recover(%s)",err)
				fmt.Println(string(debug.Stack()))
			}
		}()
		f()
	}()
}
