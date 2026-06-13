package sync

import (
	"time"

	"github.com/seekky/slinx-node/internal/database"
	"github.com/seekky/slinx-node/internal/sync/slinx"
)

var stopChans = map[uint]chan struct{}{}

func Start() {
	var boards []database.Board
	database.DB.Where("enable = ?", true).Find(&boards)

	for _, b := range boards {
		start(b)
	}
}

func Stop() {
	for _, ch := range stopChans {
		close(ch)
	}
	stopChans = map[uint]chan struct{}{}
}

func Restart() {
	Stop()
	Start()
}

func start(b database.Board) {
	ch := make(chan struct{})
	stopChans[b.ID] = ch

	go func() {
		sync(b)
		ticker := time.NewTicker(time.Duration(b.SyncInterval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				sync(b)
			case <-ch:
				return
			}
		}
	}()
}

func sync(b database.Board) {
	switch b.Type {
	case "SLINX":
		slinx.Sync(b)
	}
}
