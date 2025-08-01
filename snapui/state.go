package snapui

import (
	"math"

	"github.com/gabe-lee/OurSweeper/vec2"
)

type (
	Vec2_F32 = vec2.Vec2[float32]
	AABB_F32 = vec2.AABB[float32]
)

type UI_Idx int16

const (
	UI_IDX_NULL int16   = math.MinInt16
	EPSILON     float32 = 0.05
)

type UI_Element[T any] struct {
	OnClick      func(elemList []UI_Element[T], idx UI_Idx)
	OnHover      func(elemList []UI_Element[T], idx UI_Idx)
	OnScroll     func(elemList []UI_Element[T], idx UI_Idx, dx float64, dy float64)
	OnDrag       func(elemList []UI_Element[T], idx UI_Idx, dx float64, dy float64)
	GetSize      func(elemList []UI_Element[T], idx UI_Idx) Vec2_F32
	GetOffset    func(elemList []UI_Element[T], idx UI_Idx) Vec2_F32
	Children     []UI_Idx
	Padding      float32
	Offset       Vec2_F32
	Size         Vec2_F32
	FinalPos     Vec2_F32
	FinalSize    Vec2_F32
	Flags        UI_Flag
	ZIndex       uint16
	AnchorRefIdx UI_Idx
	UserData     T
}

func (elem *UI_Element[T]) Layout(elementList []UI_Element[T], parentIdx UI_Idx, parentZIndex uint16, parentPos Vec2_F32, parentSize Vec2_F32, parentAABB AABB_F32, parentFlags UI_Flag, parentPadding float32, selfIdx UI_Idx, relativeIdxInParent UI_Idx) (finalAABB AABB_F32) {
	elem.ZIndex = parentZIndex + 1
	elem.ZIndex |= elem.Flags.GetPlaneAsZIndex()
	var selfSize Vec2_F32 = elem.Size
	if elem.GetSize != nil {
		selfSize = elem.GetSize(elementList, selfIdx)
	}
	var selfOff Vec2_F32 = elem.Offset
	if elem.GetOffset != nil {
		selfOff = elem.GetOffset(elementList, selfIdx)
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
	var anchorIdx UI_Idx
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
		childAABB := elementList[childIdx].Layout(elementList, selfIdx, elem.ZIndex, truePos, trueSize, finalAABB, elem.Flags, elem.Padding, childIdx, UI_Idx(childRelIdx))
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

func (elem *UI_Element[T]) GetZindexWithinPlane() uint16 {
	return elem.ZIndex & _ZINDEX_WITHOUT_PLANE_MASK
}
func (elem *UI_Element[T]) GetPlaneWithoutZIndex() UI_Plane {
	return UI_Plane((elem.ZIndex & _PLANE_IN_ZINDEX_MASK) >> _PLANE_ZINDEX_SHIFT)
}

func GetAllOnClickCandidates[T any](elementList []UI_Element[T], oldCandidateBuffer []UI_Idx, ignoreList []UI_Idx, cursorPos Vec2_F32) (newCandidateBuffer []UI_Idx) {
	newCandidateBuffer = oldCandidateBuffer[:0]
	var shouldIgnore bool
	for i, element := range elementList {
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
					newCandidateBuffer = append(newCandidateBuffer, UI_Idx(i))
				}
			}
		}
	}
	return newCandidateBuffer
}

func GetAllOnHoverCandidates[T any](elementList []UI_Element[T], oldCandidateBuffer []UI_Idx, ignoreList []UI_Idx, cursorPos Vec2_F32) (newCandidateBuffer []UI_Idx) {
	newCandidateBuffer = oldCandidateBuffer[:0]
	var shouldIgnore bool
	for i, element := range elementList {
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
					newCandidateBuffer = append(newCandidateBuffer, UI_Idx(i))
				}
			}
		}
	}
	return newCandidateBuffer
}

func GetAllOnScrollCandidates[T any](elementList []UI_Element[T], oldCandidateBuffer []UI_Idx, ignoreList []UI_Idx, cursorPos Vec2_F32) (newCandidateBuffer []UI_Idx) {
	newCandidateBuffer = oldCandidateBuffer[:0]
	var shouldIgnore bool
	for i, element := range elementList {
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
					newCandidateBuffer = append(newCandidateBuffer, UI_Idx(i))
				}
			}
		}
	}
	return newCandidateBuffer
}

func GetAllOnDragCandidates[T any](elementList []UI_Element[T], oldCandidateBuffer []UI_Idx, ignoreList []UI_Idx, cursorPos Vec2_F32) (newCandidateBuffer []UI_Idx) {
	newCandidateBuffer = oldCandidateBuffer[:0]
	var shouldIgnore bool
	for i, element := range elementList {
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
					newCandidateBuffer = append(newCandidateBuffer, UI_Idx(i))
				}
			}
		}
	}
	return newCandidateBuffer
}

func GetTopmostCandidate[T any](elementList []UI_Element[T], newCandidateBuffer []UI_Idx) UI_Idx {
	var topmostZIndex uint16
	var topmostElemIdx UI_Idx
	for _, elementIdx := range newCandidateBuffer {
		element := elementList[elementIdx]
		if element.ZIndex > topmostZIndex {
			topmostZIndex = element.ZIndex
			topmostElemIdx = elementIdx
		}
	}
	return topmostElemIdx
}
