package common

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/z-riley/go-2048-battle/backend"
	"github.com/z-riley/go-2048-battle/backend/grid"
	"github.com/z-riley/go-2048-battle/log"
	"github.com/z-riley/turdgl"
	"golang.org/x/exp/constraints"
)

// Adjustable settings.
const (
	TileSizePx        float64 = 72   // the width and height of a tile, in pixels
	TileCornerRadius  float64 = 3    // the radius, in pixels, of the rounded corners of the tiles
	TileBoundryFactor float64 = 0.15 // the gap between tiles as a proportion of the tile size
)

// Derived constants.
const (
	ArenaSizePx   = tileSpacingPx*4 + TileSizePx*TileBoundryFactor // the width and height of arena, in pixels
	tileSpacingPx = TileSizePx * (1 + TileBoundryFactor)
	tileFont      = FontPathBold
	numTiles      = grid.GridSize
)

// tile is a visual representation of a game tile.
type tile struct {
	tb      *turdgl.TextBox
	pos     coord // index of tile on the grid
	destroy bool  // flag for self-destruction
}

// newTile constructs a new tile with the correct style.
func newTile(sizePx float64, pos turdgl.Vec, val int, posIdx coord) *tile {
	return &tile{
		tb: turdgl.NewTextBox(turdgl.NewCurvedRect(
			sizePx, sizePx, TileCornerRadius, pos,
		).SetStyle(turdgl.Style{Colour: tileColour(val)}), strconv.Itoa(val), tileFont).
			SetTextSize(tileFontSize(val)).
			SetTextColour(tileTextColour(val)),
		pos: posIdx,
	}
}

// animationData contains animations and the current game state.
type animationState struct {
	animations []animation
	gameState  backend.Game
}

// Arena displays the grid of a game.
type Arena struct {
	pos         turdgl.Vec                             // pixel position of the arena anchor
	tiles       []*tile                                // every non-zero tile
	bgTiles     [numTiles][numTiles]*turdgl.CurvedRect // every grid space
	background  *turdgl.CurvedRect                     // the background of the arena
	latestState backend.Game                           // used to detect changes in game state (for animations etc...)
	animationCh chan animationState                    // for sending animations to animator goroutine
}

// NewArena constructs a new arena widget. pos is the top-left pixel of the
// top-left tile (excluding the arena background).
func NewArena(pos turdgl.Vec) *Arena {
	// Generate background tiles
	bgTiles := [numTiles][numTiles]*turdgl.CurvedRect{}
	for i := range numTiles {
		for j := range numTiles {
			bgTiles[j][i] = turdgl.NewCurvedRect(
				TileSizePx, TileSizePx, TileCornerRadius,
				turdgl.Vec{
					X: pos.X + float64(j)*tileSpacingPx,
					Y: pos.Y + float64(i)*tileSpacingPx,
				},
			)
			bgTiles[j][i].SetStyle(turdgl.Style{Colour: TileBackgroundColour})
		}
	}

	arenaBG := turdgl.NewCurvedRect(
		ArenaSizePx, ArenaSizePx,
		TileCornerRadius,
		turdgl.Vec{
			X: pos.X - TileSizePx*TileBoundryFactor,
			Y: pos.Y - TileSizePx*TileBoundryFactor,
		},
	)
	arenaBG.SetStyle(turdgl.Style{Colour: ArenaBackgroundColour})

	a := Arena{
		pos:         pos,
		tiles:       make([]*tile, 0, numTiles*numTiles),
		bgTiles:     bgTiles,
		background:  arenaBG,
		latestState: backend.Game{Grid: &grid.Grid{Tiles: [4][4]grid.Tile{}}},
		animationCh: make(chan animationState, 50),
	}

	// Begin listening to animation channel
	go a.handleAnimations()

	return &a
}

// Destroy tears down the arena.
func (a *Arena) Destroy() {}

// Draw draws the arena.
func (a *Arena) Draw(buf *turdgl.FrameBuffer) {
	a.background.Draw(buf)

	for i := range numTiles {
		for j := range numTiles {
			a.bgTiles[j][i].Draw(buf)
		}
	}

	for _, t := range a.tiles {
		t.tb.Draw(buf)
	}
}

// Pos returns the top left pixel coordinate of the whole arena.
func (a *Arena) Pos() turdgl.Vec {
	return a.background.GetPos()
}

// Width returns the total width of the arena.
func (a *Arena) Width() float64 {
	return a.background.Width()
}

// Height returns the total width of the arena.
func (a *Arena) Height() float64 {
	return a.background.Height()
}

// Load updates the arena to match the backend game data.
func (a *Arena) Load(g backend.Game) {
	var newTiles []*tile
	for i := range numTiles {
		for j := range numTiles {
			val := g.Grid.Tiles[i][j].Val
			if val != 0 {
				newTiles = append(newTiles,
					newTile(
						TileSizePx,
						turdgl.Vec{
							X: a.pos.X + float64(j)*tileSpacingPx,
							Y: a.pos.Y + float64(i)*tileSpacingPx,
						},
						val,
						coord{j, i},
					))
			}
		}
	}
	a.tiles = newTiles
}

// Reset clears the current game data from the arena.
func (a *Arena) Reset() {
	a.tiles = make([]*tile, 0, numTiles*numTiles)
	a.SetNormal()
}

// SetNormal makes the arena show its losing state.
func (a *Arena) SetNormal() {
	a.background.SetStyle(turdgl.Style{Colour: ArenaBackgroundColour})
}

// SetLose makes the arena show its losing state.
func (a *Arena) SetLose() {
	a.background.SetStyle(turdgl.Style{
		Colour: turdgl.DarkRed,
		Bloom:  15,
	})
}

// SetWin makes the arena show its losing state.
func (a *Arena) SetWin() {
	a.background.SetStyle(turdgl.Style{
		Colour: turdgl.LimeGreen,
		Bloom:  15,
	})
}

// Update animates the arena to match the given game state.
func (a *Arena) Update(game backend.Game) {
	defer func() {
		// Update the local state upon exit
		a.latestState.Grid.Tiles = game.Grid.Tiles
	}()

	// Return early if the grid hasn't changed
	if grid.EqualGrid(a.latestState.Grid.Tiles, game.Grid.Tiles) {
		return
	}

	// Calculate the movement of each tile
	tileAnimations := generateAnimations(a.latestState.Grid.Tiles, game.Grid.Tiles, game.Grid.LastMove)
	if len(tileAnimations) == 0 {
		return
	}

	// Send the animations for the turn down the animation channel
	a.animationCh <- animationState{tileAnimations, game}
}

// handleAnimations executes animations from the animation channel.
func (a *Arena) handleAnimations() {
	for animationState := range a.animationCh {
		// Listen to errors being produced by animations
		errCh := make(chan error, numTiles*numTiles)

		// Animate stage 1: tiles moving and combining
		var wg sync.WaitGroup
		for _, animation := range animationState.animations {
			switch animation := animation.(type) {
			case moveAnimation:
				wg.Add(1)
				go a.animateMove(animation, errCh, &wg)
			case moveToCombineAnimation:
				wg.Add(1)
				go a.animateMoveToCombine(animation, errCh, &wg)
			}
		}
		wg.Wait()

		// Handle errors from stage 1
		select {
		case err := <-errCh:
			log.Printf("Error \"%v\". Resetting to latest game state\n", err)
			a.Load(animationState.gameState)
		default:
		}

		// Remove tiles that have been marked for destruction
		a.trimTiles()

		// Animate stage 2: spawn new tiles
		for _, animation := range animationState.animations {
			switch animation := animation.(type) {
			case spawnAnimation:
				wg.Add(1)
				go a.animateSpawn(animation, errCh, &wg)
			case newFromCombineAnimation:
				wg.Add(1)
				go a.animateNewFromCombine(animation, errCh, &wg)
			}
		}
		wg.Wait()

		// Handle errors from stage 2
		select {
		case err := <-errCh:
			log.Printf("Error \"%v\". Resetting to latest game state\n", err)
			a.Load(animationState.gameState)
		default:
		}
		close(errCh)

		// Check number of tiles
		uiTiles := len(a.tiles)
		backendTiles := animationState.gameState.Grid.NumTiles()
		if uiTiles != backendTiles {
			log.Println("Found tile count mismatch. Reloading grid")
			a.Load(animationState.gameState)
		}
	}
}

// animateMove animates a tile moving.
func (a *Arena) animateMove(animation moveAnimation, errCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	origin, dest := animation.origin, animation.dest

	// Move origin tile to destination
	tile, err := a.tileAtIdx(origin)
	if err != nil {
		errCh <- fmt.Errorf("animateMove could not find origin tile at %v", origin)
		return
	}
	moveVec := turdgl.Sub(a.tilePos(dest), a.tilePos(origin))
	const steps = 20
	moveStep := moveVec.SetMag(moveVec.Mag() / steps)
	for range steps {
		tile.tb.Move(moveStep)
		time.Sleep(5 * time.Millisecond)
	}

	// Update local tile state with new position
	tile.pos = dest
}

// animateMoveToCombine animates a tile moving into another like tile, resulting
// in a combine.
func (a *Arena) animateMoveToCombine(animation moveToCombineAnimation, errCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	origin, dest := animation.origin, animation.dest

	// Move origin tile to destination
	originTile, err := a.tileAtIdx(origin)
	if err != nil {
		errCh <- fmt.Errorf("animateMoveToCombine could not find origin tile at %v", origin)
		return
	}

	moveVec := turdgl.Sub(a.tilePos(dest), a.tilePos(origin))
	const steps = 20
	moveStep := moveVec.SetMag(moveVec.Mag() / steps)
	for range steps {
		originTile.tb.Move(moveStep)
		time.Sleep(5 * time.Millisecond)
	}

	// Mark tiles for destruction. The combined tile will be newly spawned separately
	destTile, err := a.tileAtIdx(dest)
	if err == nil {
		// Mark tile if it exists
		destTile.destroy = true
	}
	originTile.destroy = true
}

// animateSpawn animates a tile spawn animation.
func (a *Arena) animateSpawn(animation spawnAnimation, errCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	dest := animation.dest
	newVal := animation.newVal

	t, err := a.tileAtIdx(dest)
	if err == nil {
		errCh <- fmt.Errorf("aninimateSpawn - tile shouldn't already exist at %v", t.pos)
		return
	}

	// Make a small new tile
	const originalSize = TileSizePx / 6
	newTile := newTile(
		originalSize,
		turdgl.Vec{
			X: a.pos.X + float64(dest.x)*tileSpacingPx + (TileSizePx-originalSize)/2,
			Y: a.pos.Y + float64(dest.y)*tileSpacingPx + (TileSizePx-originalSize)/2,
		},
		newVal,
		dest,
	)
	a.tiles = append(a.tiles, newTile)

	// Animate tile growing to normal size
	const (
		growPx   = (TileSizePx - originalSize) / 2
		steps    = 10
		stepSize = growPx / steps
	)
	shape := newTile.tb.Shape.(*turdgl.CurvedRect)
	originalPos := shape.GetPos() // position of shape before animation starts
	for i := float64(0); i <= growPx; i += stepSize {
		shape.SetPos(turdgl.Sub(originalPos, turdgl.Vec{X: i, Y: i}))
		shape.SetHeight(originalSize + i*2)
		shape.SetWidth(originalSize + i*2)
		time.Sleep(10 * time.Millisecond)
	}
}

func (a *Arena) animateNewFromCombine(animation newFromCombineAnimation, errCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	dest := animation.dest
	newVal := animation.newVal

	t, err := a.tileAtIdx(dest)
	if err == nil {
		errCh <- fmt.Errorf("animateNewFromCombine - tile shouldn't already exist at %v", t.pos)
		return
	}

	// Make a new tile
	newTile := newTile(
		TileSizePx,
		turdgl.Vec{
			X: a.pos.X + float64(dest.x)*tileSpacingPx,
			Y: a.pos.Y + float64(dest.y)*tileSpacingPx,
		},
		newVal,
		dest,
	)
	a.tiles = append(a.tiles, newTile)

	// Animate tile growing and shrinking back to normal size
	const expandPx = 5
	shape := newTile.tb.Shape.(*turdgl.CurvedRect)
	originalPos := shape.GetPos() // position of shape before animation starts
	for i := float64(1); i <= expandPx; i++ {
		shape.SetPos(turdgl.Sub(originalPos, turdgl.Vec{X: i, Y: i}))
		shape.SetHeight(TileSizePx + i*2)
		shape.SetWidth(TileSizePx + i*2)
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(30 * time.Millisecond)
	for i := float64(expandPx) - 1; i > 0; i-- {
		shape.SetPos(turdgl.Sub(originalPos, turdgl.Vec{X: i, Y: i}))
		shape.SetHeight(TileSizePx + i*2)
		shape.SetWidth(TileSizePx + i*2)
		time.Sleep(10 * time.Millisecond)
	}
}

// tileAtIdx returns a reference to the first found tile at a given position on the grid.
// If the tile doesn't exist, an error is returned.
func (a *Arena) tileAtIdx(pos coord) (*tile, error) {
	for i := range a.tiles {
		if a.tiles[i].pos.equals(pos) {
			return a.tiles[i], nil
		}
	}
	return nil, fmt.Errorf("could not find tile at pos: %v", pos)
}

// tilePos generates the pixel position of a tile on the grid based on its x and y index.
func (a *Arena) tilePos(pos coord) turdgl.Vec {
	return turdgl.Vec{
		X: a.pos.X + float64(pos.x)*tileSpacingPx,
		Y: a.pos.Y + float64(pos.y)*tileSpacingPx,
	}
}

// trimTiles removes tiles that have been marked for destruction.
func (a *Arena) trimTiles() {
	var remainingTiles []*tile
	for _, t := range a.tiles {
		if !t.destroy {
			remainingTiles = append(remainingTiles, t)
		}
	}
	a.tiles = remainingTiles
}

// tileFontSize returns the font size for a tile of a given value.
func tileFontSize(val int) float64 {
	chars := len(strconv.Itoa(val))
	switch {
	case chars < 3:
		return 36
	case chars < 4:
		return 30
	case chars < 5:
		return 26
	default:
		return 18
	}
}

// coord contains Cartesian coordinates.
type coord struct{ x, y int }

// equals returns whether c1 is equal to c2.
func (c1 *coord) equals(c2 coord) bool {
	return c1.x == c2.x && c1.y == c2.y
}

var errFieldDoesNotExist = errors.New("field does not exist for this animation type")

// animation provides data needed for animating tiles.
type animation interface {
	// Origin returns the start position of the tile. Error if not applicable.
	Origin() (coord, error)
	// Dest returns the final position of the tile.
	Dest() coord
	// NewVal returns the new value resulting from a the animation. Error if not applicable.
	NewVal() (int, error)
	// String returns a human readible string of the animation.
	String() string
}

// moveAnimation represents the movement of a tile from one position to another, without
// combining. Satisfies the animation interface.
type moveAnimation struct {
	origin coord // tile index
	dest   coord // tile index
}

// Origin satisfies the animation interface.
func (a moveAnimation) Origin() (coord, error) {
	return a.origin, nil
}

// Dest satisfies the animation interface.
func (a moveAnimation) Dest() coord {
	return a.dest
}

// NewVal satisfies the animation interface.
func (a moveAnimation) NewVal() (int, error) {
	return 0, errFieldDoesNotExist
}

// String satisfies the animation interface.
func (a moveAnimation) String() string {
	return fmt.Sprint("move from ", a.origin, " to ", a.dest)
}

// spawnAnimation represents a new tile spawning. Satisfies the animation interface.
type spawnAnimation struct {
	dest   coord // tile index
	newVal int   // value of a newly spawned tile. 0 if N/A
}

// Origin satisfies the animation interface.
func (a spawnAnimation) Origin() (coord, error) {
	return coord{-1, -1}, errFieldDoesNotExist
}

// Dest satisfies the animation interface.
func (a spawnAnimation) Dest() coord {
	return a.dest
}

// NewVal satisfies the animation interface.
func (a spawnAnimation) NewVal() (int, error) {
	return a.newVal, nil
}

// String satisfies the animation interface.
func (a spawnAnimation) String() string {
	return fmt.Sprint("spawn at ", a.dest)
}

// moveToCombineAnimation represents the movement of a tile into another, to
// combine it. Satisfies the animation interface.
type moveToCombineAnimation struct {
	origin coord // tile index
	dest   coord // tile index
}

// Origin satisfies the animation interface.
func (a moveToCombineAnimation) Origin() (coord, error) {
	return a.origin, nil
}

// Dest satisfies the animation interface.
func (a moveToCombineAnimation) Dest() coord {
	return a.dest
}

// NewVal satisfies the animation interface.
func (a moveToCombineAnimation) NewVal() (int, error) {
	return -1, errFieldDoesNotExist
}

// String satisfies the animation interface.
func (a moveToCombineAnimation) String() string {
	return fmt.Sprint("move-to-combine from ", a.origin, " to ", a.dest)
}

// newFromCombineAnimation represents a new tile being created from a combination.
// Satisfies the animation interface.
type newFromCombineAnimation struct {
	dest   coord // tile index
	newVal int   // the value of the newly made tile
}

// Origin satisfies the animation interface.
func (a newFromCombineAnimation) Origin() (coord, error) {
	return coord{-1, -1}, errFieldDoesNotExist
}

// Dest satisfies the animation interface.
func (a newFromCombineAnimation) Dest() coord {
	return a.dest
}

// NewVal satisfies the animation interface.
func (a newFromCombineAnimation) NewVal() (int, error) {
	return a.newVal, nil
}

// String satisfies the animation interface.
func (a newFromCombineAnimation) String() string {
	return fmt.Sprint("new-from-combine at ", a.dest)
}

// generateAnimations generates animation data for transitioning between grid states.
func generateAnimations(before, after [numTiles][numTiles]grid.Tile, dir grid.Direction) []animation {
	var animations []animation

	if dir == grid.DirLeft || dir == grid.DirRight {
		// Horizontal move; evaluate row-by-row
		for i := range before {
			rowAnimations := generateRowAnimations(before[i], after[i], dir)

			for _, rowAnimation := range rowAnimations {
				switch a := rowAnimation.(type) {
				case moveRowAnimation:
					animations = append(animations, moveAnimation{
						origin: coord{must(a.Origin()), i},
						dest:   coord{a.Dest(), i},
					})
				case spawnRowAnimation:
					animations = append(animations, spawnAnimation{
						dest:   coord{a.Dest(), i},
						newVal: must(a.NewVal()),
					})
				case moveToCombineRowAnimation:
					animations = append(animations, moveToCombineAnimation{
						origin: coord{must(a.Origin()), i},
						dest:   coord{a.Dest(), i},
					})
				case newFromCombineRowAnimation:
					animations = append(animations, newFromCombineAnimation{
						dest:   coord{a.Dest(), i},
						newVal: must(a.NewVal()),
					})
				}
			}
		}
	} else {
		// Vertical move; evaluate column-by-column
		for i := range before {
			rowAnimations := generateRowAnimations(
				[numTiles]grid.Tile{before[0][i], before[1][i], before[2][i], before[3][i]},
				[numTiles]grid.Tile{after[0][i], after[1][i], after[2][i], after[3][i]},
				dir,
			)

			for _, rowAnimation := range rowAnimations {
				switch a := rowAnimation.(type) {
				case moveRowAnimation:
					animations = append(animations, moveAnimation{
						origin: coord{i, must(a.Origin())},
						dest:   coord{i, a.Dest()},
					})
				case spawnRowAnimation:
					animations = append(animations, spawnAnimation{
						dest:   coord{i, a.Dest()},
						newVal: must(a.NewVal()),
					})
				case moveToCombineRowAnimation:
					animations = append(animations, moveToCombineAnimation{
						origin: coord{i, must(a.Origin())},
						dest:   coord{i, a.Dest()},
					})
				case newFromCombineRowAnimation:
					animations = append(animations, newFromCombineAnimation{
						dest:   coord{i, a.Dest()},
						newVal: must(a.NewVal()),
					})
				}
			}
		}
	}

	return animations
}

// rowAnimation provides data for animating a row of tiles.
type rowAnimation interface {
	// Origin returns the start index of the tile. Error if not applicable.
	Origin() (int, error)
	// Dest returns the final index of the tile.
	Dest() int
	// NewVal returns the new value resulting from a the animation. Error if not applicable.
	NewVal() (int, error)
}

// moveAnimation represents the movement of a tile from one position to another, without
// combining. Satisfies the animation interface.
type moveRowAnimation struct {
	origin int // tile index
	dest   int // tile index
}

// Origin satisfies rowAnimation.
func (a moveRowAnimation) Origin() (int, error) {
	return a.origin, nil
}

// Dest satisfies rowAnimation.
func (a moveRowAnimation) Dest() int {
	return a.dest
}

// NewVal satisfies rowAnimation.
func (a moveRowAnimation) NewVal() (int, error) {
	return 0, errFieldDoesNotExist
}

// String satisfies rowAnimation.
func (a moveRowAnimation) String() string {
	return fmt.Sprint("move from ", a.origin, " to ", a.dest)
}

// spawnRowAnimation represents a new tile spawning. Satisfies the animation interface.
type spawnRowAnimation struct {
	dest   int // tile index
	newVal int // value of a newly spawned tile. 0 if N/A
}

// Origin satisfies rowAnimation.
func (a spawnRowAnimation) Origin() (int, error) {
	return -1, errFieldDoesNotExist
}

// Dest satisfies rowAnimation.
func (a spawnRowAnimation) Dest() int {
	return a.dest
}

// NewVal satisfies rowAnimation.
func (a spawnRowAnimation) NewVal() (int, error) {
	return a.newVal, nil
}

// String satisfies rowAnimation.
func (a spawnRowAnimation) String() string {
	return fmt.Sprint("spawn at ", a.dest)
}

// moveToCombineRowAnimation represents the movement of a tile into another, to
// combine it. Satisfies the rowAnimation interface.
type moveToCombineRowAnimation struct {
	origin int // tile index
	dest   int // tile index
}

// Origin satisfies rowAnimation.
func (a moveToCombineRowAnimation) Origin() (int, error) {
	return a.origin, nil
}

// Dest satisfies rowAnimation.
func (a moveToCombineRowAnimation) Dest() int {
	return a.dest
}

// NewVal satisfies rowAnimation.
func (a moveToCombineRowAnimation) NewVal() (int, error) {
	return -1, errFieldDoesNotExist
}

// String satisfies rowAnimation.
func (a moveToCombineRowAnimation) String() string {
	return fmt.Sprint("move-to-combine from ", a.origin, " to ", a.dest)
}

// newFromCombineRowAnimation represents a new tile being created from a combination.
// Satisfies the rowAnimation interface.
type newFromCombineRowAnimation struct {
	dest   int // tile index
	newVal int // the value of the newly made tile
}

// Origin satisfies rowAnimation.
func (a newFromCombineRowAnimation) Origin() (int, error) {
	return -1, errFieldDoesNotExist
}

// Dest satisfies rowAnimation.
func (a newFromCombineRowAnimation) Dest() int {
	return a.dest
}

// NewVal satisfies rowAnimation.
func (a newFromCombineRowAnimation) NewVal() (int, error) {
	return a.newVal, nil
}

// String satisfies rowAnimation.
func (a newFromCombineRowAnimation) String() string {
	return fmt.Sprint("new-from-combine at ", a.dest)
}

// generateAnimations generates animation data for a row of tiles.
func generateRowAnimations(before, after [numTiles]grid.Tile, dir grid.Direction) []rowAnimation {
	var rowAnimations []rowAnimation

	// Build a map of each tile's "before" position, indexed by their UUID
	beforeUUIDs := make(map[uuid.UUID]int, numTiles)
	for x := range before {
		beforeUUIDs[before[x].UUID] = x
	}

	for x := range after {
		// Like UUIDs indicates that a tile has moved
		beforePos, ok := beforeUUIDs[after[x].UUID]
		if ok && !(beforePos == x) {
			rowAnimations = append(rowAnimations, moveRowAnimation{
				origin: beforePos,
				dest:   x,
			})
		}

		// Tiles with the Cmb flag set are from combinations
		if after[x].Cmb {
			// Newly formed tile from combination
			rowAnimations = append(rowAnimations, newFromCombineRowAnimation{
				dest:   x,
				newVal: after[x].Val,
			})

			// The tiles that combined in the last turn can be ascertained from which tiles are
			// missing, and the direction of movement

			// Get the other tiles in the row that went missing after the move
			missingTiles := make(map[uuid.UUID]grid.Tile)
			for _, tile := range before {
				if tile.Val != 0 {
					missingTiles[tile.UUID] = tile
				}
			}
			for _, tile := range after {
				delete(missingTiles, tile.UUID)
			}

			// Count the number of combinations as a result of the move
			var numCombines int
			for _, tile := range after {
				if tile.Cmb {
					numCombines++
				}
			}

			// Get the positions of the tiles that went missing after the move
			for uuid := range missingTiles {
				origin := beforeUUIDs[uuid]
				dest := x

				isLegalMove := func(origin, dest int, dir grid.Direction) bool {
					if origin == dest {
						return false
					}

					// If there 4 tiles before, the maximum travel distance cannot be
					// more than 2 in the direction of travel
					const maxTravelDist = 2
					travelDist := abs(dest - origin)

					isFourLikeTiles := numCombines == 2 &&
						before[0].Val == before[1].Val &&
						before[0].Val == before[2].Val &&
						before[0].Val == before[3].Val

					isTwoPairsOfTiles := !isFourLikeTiles && numCombines == 2 &&
						before[0].Val == before[1].Val &&
						before[2].Val == before[3].Val

					switch dir {
					case grid.DirLeft, grid.DirUp:
						// Index should get smaller
						if origin < dest {
							return false
						}
						// Manual exclusion for 2,2,2,2 situations
						if isFourLikeTiles {
							if travelDist > maxTravelDist {
								return false
							}
							if origin == 2 && dest == 0 {
								return false
							}
						}
						// Manual exclusion for 2,2,4,4 situations
						if isTwoPairsOfTiles {
							if origin == 2 && dest == 0 {
								return false
							}
							if origin == 3 && dest == 0 {
								return false
							}
						}

					case grid.DirRight, grid.DirDown:
						// Index value should get larger
						if origin > dest {
							return false
						}
						// Manual exclusion for 2,2,2,2 and 2,2,4,4 situations
						if isFourLikeTiles || isTwoPairsOfTiles {
							if travelDist > maxTravelDist {
								return false
							}
							if origin == 1 && dest == 3 {
								return false
							}
						}
					}

					return true
				}(origin, dest, dir)

				if isLegalMove {
					rowAnimations = append(rowAnimations, moveToCombineRowAnimation{
						origin: origin,
						dest:   dest,
					})
				}
			}
		}
	}

	// Non-combined tiles with unique UUIDs are newly spawned
	for x := range after {
		_, ok := beforeUUIDs[after[x].UUID]
		if !ok && after[x].Val != 0 && !after[x].Cmb {
			rowAnimations = append(rowAnimations, spawnRowAnimation{
				dest:   x,
				newVal: after[x].Val,
			})
		}
	}

	return rowAnimations
}

// must panics if err is not nil.
func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// abs returns the absolute value of a signed integer.
func abs[T constraints.Signed](a T) T {
	if a >= 0 {
		return a
	}
	return -a
}
