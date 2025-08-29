package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"runtime"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/getlantern/systray"
	"github.com/kbinani/screenshot"
)

func GetScreenshot(monitor int) *image.RGBA {
	s, err := screenshot.CaptureDisplay(monitor)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func main() {
	onExit := func() {
		fmt.Println("Success!")
	}
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Screen Effects")
	zoom := systray.AddMenuItem("Zoom", "Zoom the screen")
	mQuit := systray.AddMenuItem("Quit", "Quit the app")

	for {
		select {
			case <- zoom.ClickedCh: {
				raylibWindow()
			}
			case <- mQuit.ClickedCh: {
				os.Exit(0)
				return
			}
		}
	}
}

func raylibWindow() {
	runtime.LockOSThread()
	rl.InitWindow(int32(1920), int32(1080), "raylib [core]")
	monitor := rl.GetCurrentMonitor()
	width := rl.GetMonitorWidth(monitor)
	height := rl.GetMonitorHeight(monitor)
	fmt.Println(height, width)
	rl.ToggleFullscreen()

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	screenshot := GetScreenshot(monitor)
	img := rl.NewImageFromImage(screenshot)
	texture := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)

	camera := rl.Camera2D{}
	camera.Zoom = 1.0

	for !rl.WindowShouldClose() {
		wheel := rl.GetMouseWheelMove()
		mouseWorldPos := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera)
		camera.Offset = rl.GetMousePosition()
		camera.Target = mouseWorldPos
		scale := 0.2 * float64(wheel)
		x := float32(math.Exp(math.Log(float64(camera.Zoom)) + scale))
		camera.Zoom = rl.Clamp(x, 0.125, 64.0)

		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)
		rl.BeginMode2D(camera)

		rl.DrawTexture(texture, 0, 0, rl.White)

		rl.EndMode2D()
		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)
}
