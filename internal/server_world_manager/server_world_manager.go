package server_world_manager

import "github.com/gabe-lee/OurSweeper/internal/common"

type (
	ServerWorld = common.ServerWorld
)

const (
	MaxWorlds = 3
)

type WorldManager struct {
	activeWorlds [MaxWorlds]ServerWorld
}

func (w *WorldManager) CanLoginToWorldIndex(idx int) bool {
	return idx < MaxWorlds
}
