package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(960, 540, "Flappy Bird")
	defer rl.CloseWindow()

	var bird_pos = rl.NewVector2(32, (540/2)-16)
	var bird_velocity = rl.NewVector2(0, 0)
	var accel_rate float32 = 25
	var decel_rate = 0.9
	var jump_rate float32 = 200
	var gravity_rate float32 = 0.0005
	var gravity float32 = 0
	var birdImage = rl.LoadTexture("resources/bird.png")

	var pipeScroll float32 = -300
	var pipePositions = [1]rl.Vector2{}
	var pipeTexture = rl.LoadTexture("resources/pipe.png")
	var pipeTopText = rl.LoadTexture("resources/pipe_top.png")

	var font = rl.LoadFont("resources/dos.ttf")

	rl.InitAudioDevice()

	var flapSound = rl.LoadSound("resources/flap.wav")
	var dieSound = rl.LoadSound("resources/death.wav")
	var hsSound = rl.LoadSound("resources/newhs.wav")
	var scoreSound = rl.LoadSound("resources/score.wav")

	var music = rl.LoadSound("resources/song.mp3")
	rl.PlaySound(music)

	var score = 0
	var alive = true

	var checkedForHs = false

	var incrementedScore = false

	var newHS = false
	var playedSfx = false

	for pipeIndex, pipe := range pipePositions {
		pipe.X = 960
		pipe.Y = float32(rand.Intn(412) + 64)

		pipePositions[pipeIndex] = pipe
	}

	rl.SetTargetFPS(240)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		if !rl.IsSoundPlaying(music) {
			rl.PlaySound(music)
		}

		if alive {

			rl.ClearBackground(rl.SkyBlue)

			rl.DrawTexture(birdImage, int32(bird_pos.X), int32(bird_pos.Y), rl.White)

			for pipeIndex, pipe := range pipePositions {
				var top = pipe.Y - 1080 - 100

				if pipe.X < -64 {
					pipe.X = 960
					pipe.Y = float32(rand.Intn(412) + 64)
					incrementedScore = false
				}

				if bird_pos.X > pipe.X && bird_pos.X < pipe.X+32 && !incrementedScore {
					incrementedScore = true

					if bird_pos.Y+16 > (top+1080) && bird_pos.Y < pipe.Y-16 {
						score += 1
						rl.PlaySound(scoreSound)

					} else {
						fmt.Printf("%d, %d, %d\n", bird_pos.Y, pipe.Y, top)
						alive = false
					}
				}

				pipe.X += pipeScroll * rl.GetFrameTime()

				rl.DrawTextureEx(pipeTexture, pipe, 0, 1, rl.White)
				rl.DrawTexture(pipeTopText, int32(pipe.X), int32(top), rl.White)

				pipePositions[pipeIndex] = pipe
			}

			if (bird_pos.Y > 572) || (bird_pos.Y < -32) {
				alive = false
			}

			rl.DrawTextEx(font, (strconv.Itoa(score)), rl.NewVector2(0, 0), 32, 0, rl.Black)

			bird_velocity = rl.Vector2Multiply(bird_velocity, rl.NewVector2(float32(decel_rate), 1))

			gravity += gravity_rate
			bird_velocity.Y += gravity

			if rl.IsKeyDown(rl.KeyRight) {
				bird_velocity.X += accel_rate * rl.GetFrameTime()
			}

			if rl.IsKeyDown(rl.KeyLeft) {
				bird_velocity.X -= accel_rate * rl.GetFrameTime()
			}

			if rl.IsKeyPressed(rl.KeyUp) {
				rl.PlaySound(flapSound)
				gravity = 0
				bird_velocity.Y = -jump_rate * rl.GetFrameTime()
			}

			rl.Vector2Multiply(bird_velocity, rl.NewVector2(rl.GetFrameTime(), rl.GetFrameTime()))

			bird_pos = rl.Vector2Add(bird_pos, bird_velocity)

			bird_pos = rl.Vector2Clamp(bird_pos, rl.NewVector2(0, -10000), rl.NewVector2(960, 10000))

		} else {

			rl.ClearBackground(rl.SkyBlue)

			if !checkedForHs {
				fmt.Println("Checkinf for highscore")
				checkedForHs = true

				var currentHighscore, err = os.ReadFile("highscore.txt")
				var highscore = 0

				if err != nil {
					os.WriteFile("highscore.txt", []byte("0"), 0666)
				}

				var actualHS string = string(currentHighscore)

				highscore, err = strconv.Atoi(actualHS)

				if err != nil {
					panic(err)
				}

				if score > highscore {
					newHS = true
					os.WriteFile("highscore.txt", []byte(strconv.Itoa(score)), 0666)
				} else {
					newHS = false
				}
			}

			if !playedSfx {
				if newHS {
					rl.PlaySound(hsSound)
				} else {
					rl.PlaySound(dieSound)
				}

				playedSfx = true
			}

			rl.DrawTextEx(font, "You Died!", rl.NewVector2(0, 0), 128, 0, rl.Black)
			rl.DrawTextEx(font, "Your score: "+strconv.Itoa(score), rl.NewVector2(0, 128), 64, 0, rl.Black)
			rl.DrawTextEx(font, "Press R to retry", rl.NewVector2(0, 128+64), 64, 0, rl.Black)

			if newHS {
				rl.DrawTextEx(font, "New high score!", rl.NewVector2(0, 256), 64, 0, rl.Black)
			}

			if rl.IsKeyPressed(rl.KeyR) {
				alive = true
				score = 0
				bird_velocity.X = 0
				bird_velocity.Y = 0
				bird_pos.X = 32
				bird_pos.Y = (540 / 2) - 16
				incrementedScore = false
				checkedForHs = false
				playedSfx = false

				for pipeIndex, pipe := range pipePositions {
					pipe.X = 960
					pipe.Y = float32(rand.Intn(412) + 64)

					pipePositions[pipeIndex] = pipe
				}
			}

		}

		rl.EndDrawing()
	}
}
