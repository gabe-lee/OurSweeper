package common



// type DebugRelCoords struct {
// 	Table [11][11]byte
// }

// func (d *DebugRelCoords) Init() {
// 	for y := range 11 {
// 		for x := range 11 {
// 			d.Table[y][x] = byte('_')
// 		}
// 	}
// }

func DoActionOn8NearbyCoords(center Coord, action func(nearPos Coord)) {
	queue := NewCascadeQueue(center)
	next, hasMoreNext := queue.NextToCheck()
	for hasMoreNext {
		action(next.Pos)
		next, hasMoreNext = queue.NextToCheck()
	}
}

func DoActionOn8NearbyCoordsInRange(center Coord, minX, maxX, minY, maxY int, action func(nearPos Coord, nearBit uint64)) {
	queue := NewCascadeQueue(center)
	next, hasMoreNext := queue.NextToCheck()
	for hasMoreNext {
		if next.Pos.IsInRange(minX, maxX, minY, maxY) {
			action(next.Pos, uint64(1)<<next.RelativeIdx)
		}
		next, hasMoreNext = queue.NextToCheck()
	}
}
