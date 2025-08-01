package server_world_manager

import (
	"github.com/gabe-lee/OurSweeper/internal/active_worlds_response"
	"github.com/gabe-lee/OurSweeper/internal/common"
	C "github.com/gabe-lee/OurSweeper/internal/consts"
)

type (
	ServerWorld          = common.ServerWorld
	ActiveWorldsResponse = active_worlds_response.ActiveWorldsReport
	ActiveWorldStats     = active_worlds_response.ActiveWorldStats
)

type WorldManager struct {
	activeWorlds [C.MAX_ACTIVE_WORLDS]ServerWorld
	activeLen    int
}

func (w *WorldManager) CanLoginToWorldIndex(idx int) bool {
	return idx < w.activeLen
}

func (w *WorldManager) GetActiveWorldsResponse() ActiveWorldsResponse {
	resp := ActiveWorldsResponse{
		WorldsLen: byte(w.activeLen),
	}
	for i := range w.activeLen {
		resp.Worlds[i] = ActiveWorldStats{
			ID:              w.activeWorlds[i].Id.Load(),
			Expires:         w.activeWorlds[i].Expires,
			TotalMines:      uint32(w.activeWorlds[i].TotalMines),
			RemainingMines:  uint32(w.activeWorlds[i].RemainingMines.Load()),
			RemainingSpaces: uint32(w.activeWorlds[i].RemainingTiles.Load()),
			CurrentUsers:    uint16(w.activeWorlds[i].CurrentUsers.Load()),
			Difficulty:      w.activeWorlds[i].Difficulty,
		}
	}
	return resp
}
