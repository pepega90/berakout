package main

import (
	"image/color"
	"log"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	WIDTH  = 800
	HEIGHT = 600
)

var (
	f          font.Face
	menuFont   font.Face
	creditFont font.Face
)

var backgroundImg *ebiten.Image

type player struct {
	img           *ebiten.Image
	width, height int
	x, y          float64
	vx, ax        float64
	darah         int
}

func (p *player) DrawPlayer(screen *ebiten.Image) {
	player_pos := &ebiten.DrawImageOptions{}
	player_pos.GeoM.Translate(p.x, p.y)
	screen.DrawImage(p.img, player_pos)
}

func (p *player) UpdatePlayer() {
	if ebiten.IsKeyPressed(ebiten.KeyRight) && int(p.x) < WIDTH-p.width {
		p.vx += 5
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) && p.x > 0 {
		p.vx -= 5
	}

	p.ax = -p.vx * 0.3
	p.vx += p.ax
	p.x += p.vx
}

type block struct {
	img           *ebiten.Image
	width, height int
	x, y          float64
	collide       bool
}

func (b *block) DrawBlock(screen *ebiten.Image) {

	block_pos := &ebiten.DrawImageOptions{}
	block_pos.GeoM.Translate(b.x, b.y)
	screen.DrawImage(b.img, block_pos)
}

type ball struct {
	radius float64
	x, y   float64
	sx, sy int
}

func (b *ball) DrawBall(screen *ebiten.Image) {
	ebitenutil.DrawCircle(screen, b.x, b.y, b.radius, color.White)
}

func (b *ball) UpdateBall() {
	if b.x > float64(WIDTH)-b.radius || b.x < 0 {
		b.sx *= -1
	}

	if b.y < b.radius {
		b.sy *= -1
	}

	b.x += float64(b.sx)
	b.y += float64(b.sy)
}
func collisionBallToBlock(b ball, p block) bool {
	cx := b.x
	cy := b.y

	if b.x < p.x {
		cx = p.x
	} else if b.x > p.x+float64(p.width) {
		cx = p.x + float64(p.width)
	}

	if b.y < p.y {
		cy = p.y
	} else if b.y > p.y+float64(p.height) {
		cy = p.y + float64(p.height)
	}

	dx := b.x - cx
	dy := b.y - cy

	jarak := math.Sqrt(dx*dx + dy*dy)

	return jarak <= b.radius
}

func collisionBallToPlayer(b ball, p player) bool {
	cx := b.x
	cy := b.y

	if b.x < p.x {
		cx = p.x
	} else if b.x > p.x+float64(p.width) {
		cx = p.x + float64(p.width)
	}

	if b.y < p.y {
		cy = p.y
	} else if b.y > p.y+float64(p.height) {
		cy = p.y + float64(p.height)
	}

	dx := b.x - cx
	dy := b.y - cy

	jarak := math.Sqrt(dx*dx + dy*dy)

	return jarak <= b.radius
}

type Game struct {
	player
	ball
	play              bool
	gameover          bool
	blocks            []block
	blockCollideCount int
}

func (g *Game) Update() error {
	// g.gameover = true
	if !g.play {
		g.player.x = WIDTH/2 - 30
		g.player.y = HEIGHT - 50
		g.ball.x = g.player.x + 50
		g.ball.y = g.player.y - 10
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.play {
		g.play = true
	}

	if g.play {
		g.ball.UpdateBall()
		g.player.UpdatePlayer()
	}

	// check jika bola melebihi batas bawah layar
	if g.ball.y > HEIGHT {
		g.play = false
		g.player.darah -= 1
	}

	// jika darah player == 0, maka game over
	if g.player.darah == 0 {
		g.gameover = true
	}

	// jika game over make blocknya bikin jadi banyak lagi
	if ebiten.IsKeyPressed(ebiten.KeyR) && g.gameover {
		for i := 0; i < len(g.blocks); i++ {
			g.blocks[i].collide = false
		}
		g.gameover = false
		g.player.darah = 3
		g.blockCollideCount = 0
	}

	// collision ke player
	if collisionBallToPlayer(g.ball, g.player) {
		g.ball.sy *= -1
	}

	// collision ke block
	for i := 0; i < len(g.blocks); i++ {
		if collisionBallToBlock(g.ball, g.blocks[i]) && !g.blocks[i].collide {
			g.ball.sy *= -1
			g.blocks[i].collide = true
			g.blockCollideCount += 1
			// g.blocks = append(g.blocks[:i], g.blocks[i+1:]...)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw background
	screen.DrawImage(backgroundImg, &ebiten.DrawImageOptions{})
	g.ball.DrawBall(screen)
	g.player.DrawPlayer(screen)
	for _, b := range g.blocks {
		if !b.collide {
			b.DrawBlock(screen)
		}
	}

	// draw nyawa
	if !g.play && g.player.darah > 0 {
		text.Draw(screen, "Berakout", menuFont, WIDTH/2-60, HEIGHT/2+40, color.White)
		text.Draw(screen, "Tekan \"SPACE\" untuk bermain!", f, WIDTH/2-130, HEIGHT/2+90, color.White)
		text.Draw(screen, "created by aji mustofa @pepega90", creditFont, WIDTH-260, HEIGHT-20, color.White)
	} else if g.gameover {
		text.Draw(screen, "Game Over", menuFont, WIDTH/2-60, HEIGHT/2+40, color.White)
		text.Draw(screen, "Tekan \"R\" untuk restart!", f, WIDTH/2-80, HEIGHT/2+90, color.White)
	}
	if g.blockCollideCount == len(g.blocks) {
		text.Draw(screen, "Thanks for playing!", menuFont, WIDTH/2-150, HEIGHT/2, color.White)
		text.Draw(screen, "created by aji mustofa @pepega90", creditFont, WIDTH-260, HEIGHT-20, color.White)
	}

	text.Draw(screen, "x "+strconv.Itoa(g.player.darah), f, 40, HEIGHT-15, color.White)
	ebitenutil.DrawCircle(screen, 20, HEIGHT-20, g.ball.radius, color.White)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func main() {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Berakout")

	// load font
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	// load assets
	player_img, _, _ := ebitenutil.NewImageFromFile("./assets/paddleRed.png")
	backgroundImg, _, _ = ebitenutil.NewImageFromFile("./assets/bg.png")
	blue, _, _ := ebitenutil.NewImageFromFile("./assets/element_blue_rectangle_glossy.png")
	green, _, _ := ebitenutil.NewImageFromFile("./assets/element_green_rectangle_glossy.png")
	red, _, _ := ebitenutil.NewImageFromFile("./assets/element_red_rectangle_glossy.png")

	g := &Game{}

	// font
	f, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	menuFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    35,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	creditFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    15,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	// player
	g.player.img = player_img
	g.player.width, g.player.height = player_img.Size()
	g.player.x = WIDTH/2 - 30
	g.player.y = HEIGHT - 50
	g.player.darah = 3

	// ball
	g.ball.radius = 10
	g.ball.sx = 4
	g.ball.sy = -4

	// block
	baris := 7
	kolom := 6
	for i := 0; i < baris; i++ {
		for j := 0; j < kolom; j++ {
			if j > 3 {
				g.blocks = append(g.blocks, block{
					img:    green,
					x:      float64(i*100 + 80),
					y:      float64(j*50 + 20),
					width:  64,
					height: 32,
				})
			} else if j > 1 {
				g.blocks = append(g.blocks, block{
					img:    blue,
					x:      float64(i*100 + 80),
					y:      float64(j*50 + 20),
					width:  64,
					height: 32,
				})
			} else {
				g.blocks = append(g.blocks, block{
					img:    red,
					x:      float64(i*100 + 80),
					y:      float64(j*50 + 20),
					width:  64,
					height: 32,
				})
			}
		}
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
