package snaplay

import (
	"slices"

	"github.com/gabe-lee/OurSweeper/vec2"
)

type (
	Vec2_F32 = vec2.Vec2[float32]
	AABB_F32 = vec2.AABB[float32]
)

type Index interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Param interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

type testUserData struct {
	kind uint32
	p1   uint16
	p2   uint16
	p3   uint16
	p4   uint16
}

type testElem = UI_Element[testUserData, uint16, float32]

type UI_Element[Snap_UserData any, Snap_Idx Index, Snap_Param Param] struct {
	GetSize      func(elemList []UI_Element[Snap_UserData, Snap_Idx, Snap_Param], idx Snap_Idx, userData Snap_UserData) vec2.Vec2[Snap_Param]
	GetOffset    func(elemList []UI_Element[Snap_UserData, Snap_Idx, Snap_Param], idx Snap_Idx, userData Snap_UserData) vec2.Vec2[Snap_Param]
	Padding      Snap_Param
	Offset       vec2.Vec2[Snap_Param]
	Size         vec2.Vec2[Snap_Param]
	FinalPos     vec2.Vec2[Snap_Param]
	FinalSize    vec2.Vec2[Snap_Param]
	Flags        Snap_Flag
	FinalZIndex  Snap_Idx
	PosAnchorIdx Snap_Idx
	WidthRefIdx  Snap_Idx
	HeightRefIdx Snap_Idx
	UserData     Snap_UserData
}

func (elem *UI_Element[T]) Layout(elementList []UI_Element[T], parentIdx Snap_Idx, parentZIndex Snap_Zindex, parentPos Vec2_F32, parentSize Vec2_F32, parentAABB AABB_F32, parentPadding float32, selfIdx Snap_Idx, relativeIdxInParent Snap_Idx) (finalAABB AABB_F32) {
	elem.ZIndex = parentZIndex + 1
	elem.ZIndex |= elem.Flags.GetPlaneAsZIndex()
	var selfSize Vec2_F32 = elem.Size
	if elem.GetSize != nil {
		selfSize = elem.GetSize(elementList, selfIdx, elem.UserData)
	}
	var selfOff Vec2_F32 = elem.Offset
	if elem.GetOffset != nil {
		selfOff = elem.GetOffset(elementList, selfIdx, elem.UserData)
	}
	xFit := elem.Flags.GetXFit()
	yFit := elem.Flags.GetYFit()
	xSize := elem.Flags.GetXSizeMode()
	ySize := elem.Flags.GetYSizeMode()
	xOff := elem.Flags.GetXOffsetMode()
	yOff := elem.Flags.GetYOffsetMode()
	selfAnchor := elem.Flags.GetSelfAnchor()
	externalAnchor := elem.Flags.GetExternalAnchor()
	externalAnchorMode := elem.Flags.GetExternalAnchorMode()
	var anchorAABB AABB_F32
	var anchorSize Vec2_F32
	var anchorIdx Snap_Idx
	switch externalAnchorMode {
	case ANCHOR_TO_REF_IDX_ABSOLUTE:
		refElem := elementList[elem.AnchorRefIdx]
		anchorIdx = elem.AnchorRefIdx
		anchorSize = refElem.FinalSize
		anchorAABB = vec2.AABB_FromPosSize(refElem.FinalPos, refElem.FinalSize)
	case ANCHOR_TO_REF_IDX_RELATIVE:
		refElemIdx := elementList[parentIdx].Children[relativeIdxInParent+elem.AnchorRefIdx]
		refElem := elementList[refElemIdx]
		anchorIdx = refElemIdx
		anchorSize = refElem.FinalSize
		anchorAABB = vec2.AABB_FromPosSize(refElem.FinalPos, refElem.FinalSize)
	default: //ANCHOR_TO_PARENT or invalid
		anchorIdx = parentIdx
		anchorSize = parentSize
		anchorAABB = parentAABB
	}
	var trueSize Vec2_F32
	switch xSize {
	case RELATIVE_ANCHOR_SIZE:
		trueSize.X = anchorSize.X * selfSize.X
	default: // ABSOLUTE or invalid
		trueSize.X = selfSize.X
	}
	switch ySize {
	case RELATIVE_ANCHOR_SIZE:
		trueSize.Y = anchorSize.Y * selfSize.Y
	default: // ABSOLUTE or invalid
		trueSize.Y = selfSize.Y
	}
	trueSize = trueSize.AddScalar(elem.Padding)
	var truePos Vec2_F32
	switch selfAnchor {
	// TOP_LEFT: already at (0, 0)
	case TOP_CENTER:
		truePos.X -= trueSize.X / 2.0
	case TOP_RIGHT:
		truePos.X -= trueSize.X
	case MID_LEFT:
		truePos.Y -= trueSize.Y / 2.0
	case MID_CENTER:
		truePos.X -= trueSize.X / 2.0
		truePos.Y -= trueSize.Y / 2.0
	case MID_RIGHT:
		truePos.X -= trueSize.X
		truePos.Y -= trueSize.Y / 2.0
	case BOT_LEFT:
		truePos.Y -= trueSize.Y
	case BOT_CENTER:
		truePos.X -= trueSize.X / 2.0
		truePos.Y -= trueSize.Y
	case BOT_RIGHT:
		truePos.X -= trueSize.X
		truePos.Y -= trueSize.Y
	}
	switch externalAnchor {
	case TOP_CENTER:
		truePos = truePos.Add(anchorAABB.TopLeft())
	case TOP_RIGHT:
		truePos = truePos.Add(anchorAABB.TopCenter())
	case MID_LEFT:
		truePos = truePos.Add(anchorAABB.MidLeft())
	case MID_CENTER:
		truePos = truePos.Add(anchorAABB.MidCenter())
	case MID_RIGHT:
		truePos = truePos.Add(anchorAABB.MidRight())
	case BOT_LEFT:
		truePos = truePos.Add(anchorAABB.BotLeft())
	case BOT_CENTER:
		truePos = truePos.Add(anchorAABB.BotCenter())
	case BOT_RIGHT:
		truePos = truePos.Add(anchorAABB.BotRight())
	default: //TOP_LEFT or invalid
		truePos = truePos.Add(anchorAABB.TopLeft())
	}
	switch xOff {
	case RELATIVE_ANCHOR_SIZE:
		truePos.X += (anchorSize.X * selfOff.X)
	default: //ABSOLUTE or invalid
		truePos.X += selfOff.X
	}
	switch yOff {
	case RELATIVE_ANCHOR_SIZE:
		truePos.Y += (anchorSize.Y * selfOff.Y)
	default: //ABSOLUTE or invalid
		truePos.Y += selfOff.Y
	}
	if anchorIdx == parentIdx {
		truePos = truePos.AddScalar(parentPadding)
	}
	var initialAABB AABB_F32 = vec2.AABB_FromPosSize(truePos, trueSize)
	var childenCumulativeAABB AABB_F32
	for childRelIdx, childIdx := range elem.Children {
		childAABB := elementList[childIdx].Layout(elementList, selfIdx, elem.ZIndex, truePos, trueSize, finalAABB, elem.Padding, childIdx, Snap_Idx(childRelIdx))
		switch xFit {
		case MIN_GROW, MAX_SHRINK:
			childenCumulativeAABB.XMax = max(childenCumulativeAABB.XMax, childAABB.XMax)
		}
		switch yFit {
		case MIN_GROW, MAX_SHRINK:
			childenCumulativeAABB.YMax = max(childenCumulativeAABB.YMax, childAABB.YMax)
		}
	}
	finalAABB = initialAABB
	switch xFit {
	case MIN_GROW:
		finalAABB.XMax = max(finalAABB.XMax, childenCumulativeAABB.XMax+elem.Padding)
	case MAX_SHRINK:
		finalAABB.XMax = min(finalAABB.XMax, childenCumulativeAABB.XMax+elem.Padding)
	}
	switch yFit {
	case MIN_GROW:
		finalAABB.YMax = max(finalAABB.YMax, childenCumulativeAABB.YMax+elem.Padding)
	case MAX_SHRINK:
		finalAABB.YMax = min(finalAABB.YMax, childenCumulativeAABB.YMax+elem.Padding)
	}
	elem.FinalPos, elem.FinalSize = finalAABB.ToPosSize()
	xEdgeChange := finalAABB.XMax - initialAABB.XMax
	yEdgeChange := finalAABB.YMax - initialAABB.YMax
	needsShift := false
	var shift Vec2_F32
	if xEdgeChange > EPSILON {
		switch selfAnchor {
		case TOP_CENTER, MID_CENTER, BOT_CENTER:
			shift.X = xEdgeChange / -2.0
			needsShift = true
		case TOP_RIGHT, MID_RIGHT, BOT_RIGHT:
			shift.X = -xEdgeChange
			needsShift = true
		}
	} else if xEdgeChange < -EPSILON {
		switch selfAnchor {
		case TOP_CENTER, MID_CENTER, BOT_CENTER:
			shift.X = xEdgeChange / 2.0
			needsShift = true
		case TOP_RIGHT, MID_RIGHT, BOT_RIGHT:
			shift.X = xEdgeChange
			needsShift = true
		}
	}
	if yEdgeChange > EPSILON {
		switch selfAnchor {
		case MID_LEFT, MID_CENTER, MID_RIGHT:
			shift.Y = yEdgeChange / -2.0
			needsShift = true
		case BOT_LEFT, BOT_CENTER, BOT_RIGHT:
			shift.Y = -yEdgeChange
			needsShift = true
		}
	} else if yEdgeChange < -EPSILON {
		switch selfAnchor {
		case MID_LEFT, MID_CENTER, MID_RIGHT:
			shift.Y = yEdgeChange / 2.0
			needsShift = true
		case BOT_LEFT, BOT_CENTER, BOT_RIGHT:
			shift.Y = yEdgeChange
			needsShift = true
		}
	}
	if needsShift {
		elem.Shift(elementList, shift)
		finalAABB = vec2.AABB_FromPosSize(elem.FinalPos, elem.FinalSize)
	}
	return finalAABB
}

func (elem *UI_Element[T]) Shift(elementList []UI_Element[T], delta Vec2_F32) {
	elem.FinalPos = elem.FinalPos.Add(delta)
	for _, child := range elem.Children {
		elementList[child].Shift(elementList, delta)
	}
}

func (elem *UI_Element[T]) GetZindexWithinPlane() Snap_Zindex {
	return elem.ZIndex & _ZINDEX_WITHOUT_PLANE_MASK
}
func (elem *UI_Element[T]) GetPlaneWithoutZIndex() UI_Plane {
	return UI_Plane((uint16(elem.ZIndex) & _PLANE_IN_ZINDEX_MASK) >> _PLANE_ZINDEX_SHIFT)
}

type UI_Plane[T any] struct {
	Height int
	List   []UI_Element[T]
}

func (p *UI_Plane[T]) Layout(width float32, height float32) (totalAABB AABB_F32) {
	pos := vec2.New[float32](0, 0)
	size := vec2.New(width, height)
	totalAABB = p.List[0].Layout(p.List, UI_IDX_NULL, 0, pos, size, vec2.AABB_FromPosSize(pos, size), 0, 0, 0)
	return
}

type UI_Manager[T any] struct {
	Planes []UI_Plane[T]
}

func NewUiManager[T any]() UI_Manager[T] {
	return UI_Manager[T]{
		Planes: make([]UI_Plane[T], 0),
	}
}

func (m *UI_Manager[T]) AddPlane(height int, initCap int) *UI_Plane[T] {
	idx := len(m.Planes)
	m.Planes = append(m.Planes, UI_Plane[T]{
		Height: height,
		List:   make([]UI_Element[T], 0, initCap),
	})
	return &m.Planes[idx]
}

func (m *UI_Manager[T]) Layout(width float32, height float32) (totalAABB AABB_F32) {
	first := true
	for i := range m.Planes {
		planeAABB := m.Planes[i].Layout(width, height)
		if first {
			totalAABB = planeAABB
			first = false
		} else {
			totalAABB = totalAABB.Combine(planeAABB)
		}
	}
	return
}

func (m *UI_Manager[T]) DoActionOnAllVisibleElementsBottomToTop(action func(element UI_Element[T], elementList []UI_Manager[T], elemIdx Snap_Idx)) {
	sortedPlanes := make([]UI_Plane[T], len(m.Planes))
	copy(sortedPlanes, m.Planes)
	slices.SortFunc(sortedPlanes, func(a, b UI_Plane[T]) int {
		if a.Height < b.Height {
			return -1
		} else if a.Height > b.Height {
			return 1
		}
		return 0
	})
	for _, plane := range sortedPlanes {

	}
}

func (m *UI_Manager[T]) GetAllOnClickCandidates(cursorPos Vec2_F32, oldCandidateBuffer []Snap_Idx, ignoreList []Snap_Idx) (newCandidateBuffer []Snap_Idx) {
	if oldCandidateBuffer != nil {
		newCandidateBuffer = oldCandidateBuffer[:0]
	} else {
		newCandidateBuffer = make([]Snap_Idx, 0, 16)
	}
	var shouldIgnore bool
	for i, element := range m.List {
		shouldIgnore = element.OnClick == nil
		if !shouldIgnore {
			for _, ignoreIdx := range ignoreList {
				if i == int(ignoreIdx) {
					shouldIgnore = true
				}
			}
			if !shouldIgnore {
				aabb := vec2.AABB_FromPosSize(element.FinalPos, element.FinalSize)
				if aabb.PointIsWithin(cursorPos) {
					newCandidateBuffer = append(newCandidateBuffer, Snap_Idx(i))
				}
			}
		}
	}
	return newCandidateBuffer
}

func (m *UI_Manager[T]) GetTopmostOnClick(cursorPos Vec2_F32, ignoreList []Snap_Idx) (elemIdx Snap_Idx, found bool) {
	var topmostZIndex Snap_Zindex
	var shouldIgnore bool
	for i, element := range m.List {
		shouldIgnore = element.OnClick == nil || element.ZIndex <= topmostZIndex
		if !shouldIgnore {
			for _, ignoreIdx := range ignoreList {
				if i == int(ignoreIdx) {
					shouldIgnore = true
				}
			}
			if !shouldIgnore {
				aabb := vec2.AABB_FromPosSize(element.FinalPos, element.FinalSize)
				if aabb.PointIsWithin(cursorPos) {
					topmostZIndex = element.ZIndex
					elemIdx = Snap_Idx(i)
					found = true
				}
			}
		}
	}
	return
}

func (m *UI_Manager[T]) GetAllOnHoverCandidates(cursorPos Vec2_F32, oldCandidateBuffer []Snap_Idx, ignoreList []Snap_Idx) (newCandidateBuffer []Snap_Idx) {
	if oldCandidateBuffer != nil {
		newCandidateBuffer = oldCandidateBuffer[:0]
	} else {
		newCandidateBuffer = make([]Snap_Idx, 0, 16)
	}
	var shouldIgnore bool
	for i, element := range m.List {
		shouldIgnore = element.OnHover == nil
		if !shouldIgnore {
			for _, ignoreIdx := range ignoreList {
				if i == int(ignoreIdx) {
					shouldIgnore = true
				}
			}
			if !shouldIgnore {
				aabb := vec2.AABB_FromPosSize(element.FinalPos, element.FinalSize)
				if aabb.PointIsWithin(cursorPos) {
					newCandidateBuffer = append(newCandidateBuffer, Snap_Idx(i))
				}
			}
		}
	}
	return newCandidateBuffer
}

func (m *UI_Manager[T]) GetTopmostOnHover(cursorPos Vec2_F32, ignoreList []Snap_Idx) (elemIdx Snap_Idx, found bool) {
	var topmostZIndex Snap_Zindex
	var shouldIgnore bool
	for i, element := range m.List {
		shouldIgnore = element.OnHover == nil || element.ZIndex <= topmostZIndex
		if !shouldIgnore {
			for _, ignoreIdx := range ignoreList {
				if i == int(ignoreIdx) {
					shouldIgnore = true
				}
			}
			if !shouldIgnore {
				aabb := vec2.AABB_FromPosSize(element.FinalPos, element.FinalSize)
				if aabb.PointIsWithin(cursorPos) {
					topmostZIndex = element.ZIndex
					elemIdx = Snap_Idx(i)
					found = true
				}
			}
		}
	}
	return
}

func (m *UI_Manager[T]) GetAllOnScrollCandidates(cursorPos Vec2_F32, oldCandidateBuffer []Snap_Idx, ignoreList []Snap_Idx) (newCandidateBuffer []Snap_Idx) {
	if oldCandidateBuffer != nil {
		newCandidateBuffer = oldCandidateBuffer[:0]
	} else {
		newCandidateBuffer = make([]Snap_Idx, 0, 16)
	}
	var shouldIgnore bool
	for i, element := range m.List {
		shouldIgnore = element.OnScroll == nil
		if !shouldIgnore {
			for _, ignoreIdx := range ignoreList {
				if i == int(ignoreIdx) {
					shouldIgnore = true
				}
			}
			if !shouldIgnore {
				aabb := vec2.AABB_FromPosSize(element.FinalPos, element.FinalSize)
				if aabb.PointIsWithin(cursorPos) {
					newCandidateBuffer = append(newCandidateBuffer, Snap_Idx(i))
				}
			}
		}
	}
	return newCandidateBuffer
}

func (m *UI_Manager[T]) GetTopmostOnScroll(cursorPos Vec2_F32, ignoreList []Snap_Idx) (elemIdx Snap_Idx, found bool) {
	var topmostZIndex Snap_Zindex
	var shouldIgnore bool
	for i, element := range m.List {
		shouldIgnore = element.OnScroll == nil || element.ZIndex <= topmostZIndex
		if !shouldIgnore {
			for _, ignoreIdx := range ignoreList {
				if i == int(ignoreIdx) {
					shouldIgnore = true
				}
			}
			if !shouldIgnore {
				aabb := vec2.AABB_FromPosSize(element.FinalPos, element.FinalSize)
				if aabb.PointIsWithin(cursorPos) {
					topmostZIndex = element.ZIndex
					elemIdx = Snap_Idx(i)
					found = true
				}
			}
		}
	}
	return
}

func (m *UI_Manager[T]) GetAllOnDragCandidates(cursorPos Vec2_F32, oldCandidateBuffer []Snap_Idx, ignoreList []Snap_Idx) (newCandidateBuffer []Snap_Idx) {
	if oldCandidateBuffer != nil {
		newCandidateBuffer = oldCandidateBuffer[:0]
	} else {
		newCandidateBuffer = make([]Snap_Idx, 0, 16)
	}
	var shouldIgnore bool
	for i, element := range m.List {
		shouldIgnore = element.OnDrag == nil
		if !shouldIgnore {
			for _, ignoreIdx := range ignoreList {
				if i == int(ignoreIdx) {
					shouldIgnore = true
				}
			}
			if !shouldIgnore {
				aabb := vec2.AABB_FromPosSize(element.FinalPos, element.FinalSize)
				if aabb.PointIsWithin(cursorPos) {
					newCandidateBuffer = append(newCandidateBuffer, Snap_Idx(i))
				}
			}
		}
	}
	return newCandidateBuffer
}

func (m *UI_Manager[T]) GetTopmostOnDrag(cursorPos Vec2_F32, ignoreList []Snap_Idx) (elemIdx Snap_Idx, found bool) {
	var topmostZIndex Snap_Zindex
	var shouldIgnore bool
	for i, element := range m.List {
		shouldIgnore = element.OnDrag == nil || element.ZIndex <= topmostZIndex
		if !shouldIgnore {
			for _, ignoreIdx := range ignoreList {
				if i == int(ignoreIdx) {
					shouldIgnore = true
				}
			}
			if !shouldIgnore {
				aabb := vec2.AABB_FromPosSize(element.FinalPos, element.FinalSize)
				if aabb.PointIsWithin(cursorPos) {
					topmostZIndex = element.ZIndex
					elemIdx = Snap_Idx(i)
					found = true
				}
			}
		}
	}
	return
}
