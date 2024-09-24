package common

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	game "github.com/z-riley/go-2048-battle"
	"github.com/z-riley/go-2048-battle/backend"
	"github.com/z-riley/go-2048-battle/backend/grid"
	"github.com/z-riley/turdgl"
)

// Tile settings
const (
	tileSizePx       float64 = 60
	tileSpacingPx    float64 = tileSizePx * 1.2
	tileCornerRadius         = 3
	tileFont                 = game.FontPath

	arenaSize = grid.GridLen
)

// tile is a visual representation of a game tile.
type tile struct {
	tb      *turdgl.TextBox
	pos     coord // index of tile on the grid
	destroy bool  // flag for self-destruction
}

// animationData contains animations and the current game state.
type animationState struct {
	animations []animation
	gameState  backend.Game
}

// Arena displays the grid of a game.
type Arena struct {
	pos         turdgl.Vec                               // pixel position of the arena anchor
	tiles       []*tile                                  // every non-zero tile
	bgTiles     [arenaSize][arenaSize]*turdgl.CurvedRect // every grid space
	latestState backend.Game                             // used to detect changes in game state (for animations etc...)
	animationCh chan (animationState)                    // for sending animations to animator goroutine
}

// NewArena constructs a new arena widget.
func NewArena(pos turdgl.Vec) *Arena {
	// Generate background tiles
	bgTiles := [arenaSize][arenaSize]*turdgl.CurvedRect{}
	for i := 0; i < arenaSize; i++ {
		for j := 0; j < arenaSize; j++ {
			bgTiles[j][i] = turdgl.NewCurvedRect(
				tileSizePx, tileSizePx, tileCornerRadius,
				turdgl.Vec{
					X: pos.X + float64(j)*tileSpacingPx,
					Y: pos.Y + float64(i)*tileSpacingPx,
				},
			)
			bgTiles[j][i].SetStyle(turdgl.Style{Colour: turdgl.DarkSlateGrey})
		}
	}

	a := Arena{
		pos:         pos,
		tiles:       make([]*tile, 0, arenaSize*arenaSize),
		bgTiles:     bgTiles,
		latestState: backend.Game{Grid: &grid.Grid{Tiles: [4][4]grid.Tile{}}},
		animationCh: make(chan animationState, 50),
	}

	// Begin listening to animation channel
	go a.handleAnimations()

	return &a
}

// Destroy tears down the arena.
func (a *Arena) Destroy() {
	close(a.animationCh)
}

// Draw draws the arena.
func (a *Arena) Draw(win *turdgl.Window) {
	for i := 0; i < arenaSize; i++ {
		for j := 0; j < arenaSize; j++ {
			win.DrawBackground(a.bgTiles[j][i])
		}
	}
	for _, t := range a.tiles {
		win.DrawForeground(t.tb)
	}
}

// Load updates the arena to match the backend game data.
func (a *Arena) Load(g backend.Game) {
	var newTiles []*tile
	for i := 0; i < arenaSize; i++ {
		for j := 0; j < arenaSize; j++ {
			val := g.Grid.Tiles[i][j].Val
			if val != 0 {
				tb := turdgl.NewTextBox(
					turdgl.NewCurvedRect(
						tileSizePx, tileSizePx, tileCornerRadius,
						turdgl.Vec{
							X: a.pos.X + float64(j)*tileSpacingPx,
							Y: a.pos.Y + float64(i)*tileSpacingPx,
						},
					), tileFont).
					SetTextAlignment(turdgl.AlignTopCentre).
					SetText(fmt.Sprint(val))
				tb.Shape.SetStyle(turdgl.Style{Colour: tileColour(val)})
				newTiles = append(newTiles, &tile{
					tb:  tb,
					pos: coord{j, i},
				})
			}
		}
	}
	a.tiles = newTiles
}

// Reset clears thed current game data from the arena.
func (a *Arena) Reset() {
	a.tiles = make([]*tile, 0, arenaSize*arenaSize)
}

// Animate animates the arena to match the given game state.
func (a *Arena) Animate(game backend.Game) {
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

		fmt.Print("Animations:")
		for _, a := range animationState.animations {
			fmt.Printf(" %+v,", a.String())
		}
		fmt.Print("\n")

		// Listen to errors being produced by animations
		errCh := make(chan error, arenaSize*arenaSize)

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
			fmt.Printf("Error \"%v\". Resetting to latest game state\n", err)
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
			fmt.Printf("Error \"%v\". Resetting to latest game state\n", err)
			a.Load(animationState.gameState)
		default:
		}
		close(errCh)

		// Check number of tiles
		uiTiles := len(a.tiles)
		backendTiles := animationState.gameState.Grid.NumTiles()
		if uiTiles != backendTiles {
			fmt.Println("Found tile count mismatch. Reloading grid")
			a.Load(animationState.gameState)
		}

	}
}

// animateMove animates a tile moving.
func (a *Arena) animateMove(animation moveAnimation, errCh chan (error), wg *sync.WaitGroup) {
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
	for i := 0; i < steps; i++ {
		tile.tb.Move(moveStep)
		time.Sleep(6 * time.Millisecond)
	}

	// Update local tile state with new position
	tile.pos = dest
}

// animateMoveToCombine animates a tile moving into another like tile, resulting
// in a combine.
func (a *Arena) animateMoveToCombine(animation moveToCombineAnimation, errCh chan (error), wg *sync.WaitGroup) {
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
	for i := 0; i < steps; i++ {
		originTile.tb.Move(moveStep)
		time.Sleep(6 * time.Millisecond)
	}

	// Mark tiles for destruction. The combined tile will be newly spawned seperately
	destTile, err := a.tileAtIdx(dest)
	if err == nil {
		// There may or may not already be a tile at the destination
		destTile.destroy = true
	}
	originTile.destroy = true
}

// animateSpawn animates a tile spawn animation.
func (a *Arena) animateSpawn(animation spawnAnimation, errCh chan (error), wg *sync.WaitGroup) {
	defer wg.Done()

	dest := animation.dest
	newVal := animation.newVal

	t, err := a.tileAtIdx(dest)
	if err == nil {
		errCh <- fmt.Errorf("aninimateSpawn - tile shouldn't already exist at %v", t.pos)
		return
	}

	// Make a small new tile
	const originalSize = tileSizePx / 6
	newTile := &tile{
		tb: turdgl.NewTextBox(
			turdgl.NewCurvedRect(
				originalSize, originalSize, tileCornerRadius,
				turdgl.Vec{
					X: a.pos.X + float64(dest.x)*tileSpacingPx + (tileSizePx-originalSize)/2,
					Y: a.pos.Y + float64(dest.y)*tileSpacingPx + (tileSizePx-originalSize)/2,
				},
			), tileFont).
			SetTextAlignment(turdgl.AlignTopCentre).
			SetText(fmt.Sprint(newVal)),
		pos: dest,
	}
	a.tiles = append(a.tiles, newTile)
	newTile.tb.Shape.SetStyle(turdgl.Style{Colour: tileColour(newVal)})

	// Animate tile growing to normal size
	const (
		growPx   = (tileSizePx - originalSize) / 2
		steps    = 10
		stepSize = growPx / steps
	)
	shape := newTile.tb.Shape
	originalPos := shape.GetPos() // position of shape before animation starts
	for i := float64(0); i <= growPx; i += stepSize {
		shape.SetPos(turdgl.Sub(originalPos, turdgl.Vec{X: i, Y: i}))
		shape.SetHeight(originalSize + i*2)
		shape.SetWidth(originalSize + i*2)
		time.Sleep(25 * time.Millisecond)
	}
}

func (a *Arena) animateNewFromCombine(animation newFromCombineAnimation, errCh chan (error), wg *sync.WaitGroup) {
	defer wg.Done()

	dest := animation.dest
	newVal := animation.newVal

	t, err := a.tileAtIdx(dest)
	if err == nil {
		errCh <- fmt.Errorf("animateNewFromCombine - tile shouldn't already exist at %v", t.pos)
		return
	}

	// Make a new tile
	newTile := &tile{
		tb: turdgl.NewTextBox(
			turdgl.NewCurvedRect(
				tileSizePx, tileSizePx, 3,
				turdgl.Vec{
					X: a.pos.X + float64(dest.x)*tileSpacingPx,
					Y: a.pos.Y + float64(dest.y)*tileSpacingPx,
				},
			), tileFont).
			SetTextAlignment(turdgl.AlignTopCentre).
			SetText(fmt.Sprint(newVal)),
		pos: dest,
	}
	a.tiles = append(a.tiles, newTile)
	newTile.tb.Shape.SetStyle(turdgl.Style{Colour: tileColour(newVal)})

	// Animate tile growing and shrinking back to normal size
	const expandPx = 5
	shape := newTile.tb.Shape
	originalPos := shape.GetPos() // position of shape before animation starts
	for i := float64(1); i <= expandPx; i++ {
		shape.SetPos(turdgl.Sub(originalPos, turdgl.Vec{X: i, Y: i}))
		shape.SetHeight(tileSizePx + i*2)
		shape.SetWidth(tileSizePx + i*2)
		time.Sleep(30 * time.Millisecond)
	}
	for i := float64(expandPx) - 1; i > 0; i-- {
		shape.SetPos(turdgl.Sub(originalPos, turdgl.Vec{X: i, Y: i}))
		shape.SetHeight(tileSizePx + i*2)
		shape.SetWidth(tileSizePx + i*2)
		time.Sleep(30 * time.Millisecond)
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

// tilePos generates the pixel position of a tile on the grid based on
// its x and y index.
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

// coord contains Cartesian coordinates.
type coord struct{ x, y int }

// equals returns whether c1 is equal to c2.
func (c1 *coord) equals(c2 coord) bool {
	return c1.x == c2.x && c1.y == c2.y
}

var errFieldDoesNotExist error = errors.New("field does not exist for this animation type")

// animation provides data needed for animating tiles.
type animation interface {
	// Origin returns the start position of the tile. Error if doesn't exist.
	Origin() (coord, error)
	// Dest returns the final position of the tile.
	Dest() coord
	// NewVal returns the new value resulting from a the animation. Error if doesn't exist.
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

func (a moveAnimation) Origin() (coord, error) {
	return a.origin, nil
}

func (a moveAnimation) Dest() coord {
	return a.dest
}

func (a moveAnimation) NewVal() (int, error) {
	return 0, errFieldDoesNotExist
}

func (a moveAnimation) String() string {
	return fmt.Sprint("move from ", a.origin, " to ", a.dest)
}

// spawnAnimation represents a new tile spawning. Satisfies the animation interface.
type spawnAnimation struct {
	dest   coord // tile index
	newVal int   // value of a newly spawned tile. 0 if N/A
}

func (a spawnAnimation) Origin() (coord, error) {
	return coord{-1, -1}, errFieldDoesNotExist
}

func (a spawnAnimation) Dest() coord {
	return a.dest
}

func (a spawnAnimation) NewVal() (int, error) {
	return a.newVal, nil
}

func (a spawnAnimation) String() string {
	return fmt.Sprint("spawn at ", a.dest)
}

// moveToCombineAnimation represents the movement of a tile into another, to
// combine it. Satisfies the animation interface.
type moveToCombineAnimation struct {
	origin coord // tile index
	dest   coord // tile index
}

func (a moveToCombineAnimation) Origin() (coord, error) {
	return a.origin, nil
}

func (a moveToCombineAnimation) Dest() coord {
	return a.dest
}

func (a moveToCombineAnimation) NewVal() (int, error) {
	return -1, errFieldDoesNotExist
}

func (a moveToCombineAnimation) String() string {
	return fmt.Sprint("move-to-combine from ", a.origin, " to ", a.dest)
}

// newFromCombineAnimation represents a new tile being created from a combination.
// Satisfies the animation interface.
type newFromCombineAnimation struct {
	dest   coord // tile index
	newVal int   // the value of the newly made tile
}

func (a newFromCombineAnimation) Origin() (coord, error) {
	return coord{-1, -1}, errFieldDoesNotExist
}

func (a newFromCombineAnimation) Dest() coord {
	return a.dest
}

func (a newFromCombineAnimation) NewVal() (int, error) {
	return a.newVal, nil
}

func (a newFromCombineAnimation) String() string {
	return fmt.Sprint("new-from-combine at ", a.dest)
}

// generateAnimations generates animation data for transitioning between grid states.
func generateAnimations(before, after [arenaSize][arenaSize]grid.Tile, dir grid.Direction) []animation {
	var moves []animation

	// Index "before" tiles by UUID
	beforeUUIDs := make(map[uuid.UUID]coord, arenaSize*arenaSize)
	for y := 0; y < len(before); y++ {
		for x := 0; x < len(before[y]); x++ {
			beforeUUIDs[before[y][x].UUID] = coord{x, y}
		}
	}

	// Evaluate "after" tiles
	for y := 0; y < len(after); y++ {
		for x := 0; x < len(after[y]); x++ {
			// Tiles with the same UUIDs have moved
			beforeCoord, ok := beforeUUIDs[after[y][x].UUID]
			if ok && !beforeCoord.equals(coord{x, y}) {
				moves = append(moves, moveAnimation{
					origin: beforeCoord,
					dest:   coord{x, y},
				})
			}
			// Tiles with the Cmb set are from combinations
			if after[y][x].Cmb {
				// Newly formed tile from combination
				moves = append(moves, newFromCombineAnimation{
					dest:   coord{x, y},
					newVal: after[y][x].Val,
				})

				var beforeRow [arenaSize]grid.Tile
				var afterRow [arenaSize]grid.Tile
				if dir == grid.DirLeft || dir == grid.DirRight {
					// Horizontal move
					beforeRow = before[y]
					afterRow = after[y]
				} else {
					// Vertical move
					beforeRow = [arenaSize]grid.Tile{before[0][x], before[1][x], before[2][x], before[3][x]}
					afterRow = [arenaSize]grid.Tile{after[0][x], after[1][x], after[2][x], after[3][x]}
				}

				// The tiles that combined in the last turn can be ascertained from which tiles are
				// missing, and the direction of movement

				// Get the other tiles in the row which went missing after the move
				missingTiles := make(map[uuid.UUID]grid.Tile)
				for _, tile := range beforeRow {
					if tile.Val != 0 {
						missingTiles[tile.UUID] = tile
					}
				}
				for _, tile := range afterRow {
					delete(missingTiles, tile.UUID)
				}

				// Get the positions of the tiles which went missing after the move
				// TODO: need to behave differently for when origin is 2,2,2,2
				for uuid := range missingTiles {
					origin := beforeUUIDs[uuid]
					dest := coord{x, y}
					if !origin.equals(dest) {
						moves = append(moves, moveToCombineAnimation{
							origin: origin,
							dest:   dest,
						})
					}
				}
			}
		}
	}

	// Non-combined tiles with unique UUIDs are newly spawned
	for y := 0; y < len(after); y++ {
		for x := 0; x < len(after[y]); x++ {
			_, ok := beforeUUIDs[after[y][x].UUID]
			if !ok && after[y][x].Val != 0 && !after[y][x].Cmb {
				moves = append(moves, spawnAnimation{
					dest:   coord{x, y},
					newVal: after[y][x].Val,
				})
			}
		}
	}

	return moves
}
