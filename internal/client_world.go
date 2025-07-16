package internal

// const (
// 	ICON_0      = " "
// 	ICON_1      = ansi.FG_BLU + "1" + ansi.CLEAR
// 	ICON_2      = ansi.FG_CYA + "2" + ansi.CLEAR
// 	ICON_3      = ansi.FG_GRN + "3" + ansi.CLEAR
// 	ICON_4      = ansi.FG_YEL + "4" + ansi.CLEAR
// 	ICON_5      = ansi.FG_RED + "5" + ansi.CLEAR
// 	ICON_6      = ansi.FG_MAG + "6" + ansi.CLEAR
// 	ICON_7      = ansi.FG_WHT + "7" + ansi.CLEAR
// 	ICON_8      = ansi.FG_BLK + "8" + ansi.CLEAR
// 	ICON_BOMB   = ansi.INV_RED + "X" + ansi.CLEAR
// 	ICON_SKULL  = ansi.FG_RED + "💀" + ansi.CLEAR
// 	ICON_FLAG   = "█"
// 	ICON_OPAQUE = "▒"
// )

type ClientWorld struct {
	Id            uint32
	Score         uint32
	Tiles         [WORLD_TILE_COUNT]byte
	TotalMines    uint32
	ExplodedMines uint32
	SweptTiles    uint32
	Ended         bool
	Expires       int64
}

// // This is NOT safe for concurrent use when world is being played on,
// // this is only for testing/debugging purposes
// func (w *ClientWorld) DrawState(wr io.Writer) {
// 	var i uint = 0
// 	var buf [4]byte
// 	capFill := strings.Repeat("═", int(WORLD_TILE_WIDTH))
// 	var cap string = "╔" + capFill + "╗\n"
// 	wr.Write([]byte(cap))
// 	for range WORLD_TILE_HEIGHT {
// 		n := utf8.EncodeRune(buf[:], '║')
// 		wr.Write(buf[:n])
// 		for range WORLD_HALF_WIDTH {
// 			iconCode := w.Tiles[i].GetIconServer()
// 			switch iconCode {
// 			case ICON_CODE_BOMB:
// 				wr.Write([]byte(ICON_SKULL))
// 			case ICON_CODE_OPAQUE:
// 				wr.Write([]byte(ICON_OPAQUE))
// 			case ICON_CODE_FLAG:
// 				wr.Write([]byte(ICON_FLAG))
// 			case ICON_CODE_0:
// 				wr.Write([]byte(ICON_0))
// 			case ICON_CODE_1:
// 				wr.Write([]byte(ICON_1))
// 			case ICON_CODE_2:
// 				wr.Write([]byte(ICON_2))
// 			case ICON_CODE_3:
// 				wr.Write([]byte(ICON_3))
// 			case ICON_CODE_4:
// 				wr.Write([]byte(ICON_4))
// 			case ICON_CODE_5:
// 				wr.Write([]byte(ICON_5))
// 			case ICON_CODE_6:
// 				wr.Write([]byte(ICON_6))
// 			case ICON_CODE_7:
// 				wr.Write([]byte(ICON_7))
// 			case ICON_CODE_8:
// 				wr.Write([]byte(ICON_8))
// 			default:
// 				wr.Write([]byte(ICON_0))
// 			}
// 			i++
// 		}
// 		n = utf8.EncodeRune(buf[:], '║')
// 		buf[n] = 0x0A
// 		wr.Write(buf[:n+1])
// 	}
// 	cap = "╚" + capFill + "╝\n"
// 	wr.Write([]byte(cap))
// }

// // This is NOT safe for cuncurrent use when world is being played on,
// // this is only for testing/debugging purposes
// func (w *ClientWorld) DrawMines(wr io.Writer) {
// 	var i uint = 0
// 	var buf [4]byte
// 	capFill := strings.Repeat("═", int(WORLD_TILE_WIDTH))
// 	var cap string = "╔" + capFill + "╗\n"
// 	wr.Write([]byte(cap))
// 	for range WORLD_TILE_HEIGHT {
// 		n := utf8.EncodeRune(buf[:], '║')
// 		wr.Write(buf[:n])
// 		for range WORLD_TILE_WIDTH {
// 			isMine := w.Tiles[i].IsMine()
// 			if isMine {
// 				wr.Write([]byte(ICON_BOMB))
// 			} else {
// 				wr.Write([]byte(ICON_0))
// 			}
// 			i++
// 		}
// 		n = utf8.EncodeRune(buf[:], '║')
// 		buf[n] = 0x0A
// 		wr.Write(buf[:n+1])
// 	}
// 	cap = "╚" + capFill + "╝\n"
// 	wr.Write([]byte(cap))
// }

// // This is NOT safe for cuncurrent use when world is being played on,
// // this is only for testing/debugging purposes
// func (w *ClientWorld) DrawNearby(wr io.Writer) {
// 	var i uint = 0
// 	var buf [4]byte
// 	capFill := strings.Repeat("═", int(WORLD_TILE_WIDTH))
// 	var cap string = "╔" + capFill + "╗\n"
// 	wr.Write([]byte(cap))
// 	for range WORLD_TILE_HEIGHT {
// 		n := utf8.EncodeRune(buf[:], '║')
// 		wr.Write(buf[:n])
// 		for range WORLD_TILE_WIDTH {
// 			isMine := w.Tiles[i].IsMine()
// 			if isMine {
// 				wr.Write([]byte(ICON_BOMB))
// 			} else {
// 				nearby := w.Tiles[i].GetNearby()
// 				switch nearby {
// 				case 0:
// 					wr.Write([]byte(ICON_0))
// 				case 1:
// 					wr.Write([]byte(ICON_1))
// 				case 2:
// 					wr.Write([]byte(ICON_2))
// 				case 3:
// 					wr.Write([]byte(ICON_3))
// 				case 4:
// 					wr.Write([]byte(ICON_4))
// 				case 5:
// 					wr.Write([]byte(ICON_5))
// 				case 6:
// 					wr.Write([]byte(ICON_6))
// 				case 7:
// 					wr.Write([]byte(ICON_7))
// 				case 8:
// 					wr.Write([]byte(ICON_8))
// 				default:
// 					wr.Write([]byte(ICON_0))
// 				}
// 			}
// 			i++
// 		}
// 		n = utf8.EncodeRune(buf[:], '║')
// 		buf[n] = 0x0A
// 		wr.Write(buf[:n+1])
// 	}
// 	cap = "╚" + capFill + "╝\n"
// 	wr.Write([]byte(cap))
// }
