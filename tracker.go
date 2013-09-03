package main

import (
	"fmt"
	"os"

	"github.com/ddunkin/go-opencv/opencv"
	"github.com/ddunkin/launcher"
)

var faceCascade *opencv.HaarClassifierCascade
var storage *opencv.MemStorage

var launcherDev *launcher.Launcher

func main() {
	faceCascade = opencv.LoadClassifier("haarcascades/haarcascade_frontalface_alt.xml")
	defer faceCascade.Release()
	storage = opencv.CreateMemStorage(0)
	defer storage.Release()

	cap := opencv.NewCameraCapture(opencv.CV_CAP_ANY)
	if cap == nil {
		panic("can not open video")
	}
	defer cap.Release()

	launcherDev = launcher.Create()
	defer launcherDev.Destroy()

	launcherDev.SendCommandDuration(launcher.Up, 100)
	launcherDev.SendCommandDuration(launcher.Down, 100)

	win := opencv.NewWindow("Face Tracker")
	defer win.Destroy()

	fps := 30

	cap.SetProperty(opencv.CV_CAP_PROP_FRAME_WIDTH, 320)
	cap.SetProperty(opencv.CV_CAP_PROP_FRAME_HEIGHT, 240)

	stop := false

	frame := 0

	for {
		if !stop {
			img := cap.QueryFrame()
			if img == nil {
				break
			}

			if frame % 15 == 0 {
				detectAndMove(img)
			}
			frame++

			win.ShowImage(img)
			key := opencv.WaitKey(1000 / fps)
			if key == 27 {
				os.Exit(0)
			}
		} else {
			key := opencv.WaitKey(20)
			if key == 27 {
				os.Exit(0)
			}
		}
	}

	opencv.WaitKey(0)
}

func detectAndMove(image *opencv.IplImage) bool {
	//fmt.Printf("detecting faces\n")
	faceSize := opencv.Size{40, 40}
	imageSize := opencv.Size{image.Width(), image.Height()}
	imageCenter := opencv.Point{X: imageSize.Width / 2, Y: imageSize.Height / 2}
	//faces := faceCascade.DetectObjects(image, storage, 1.1, 2, opencv.CV_HAAR_DO_CANNY_PRUNING, opencv.Size{40, 40}, opencv.Size{0,0})
	faces := faceCascade.DetectObjects(image, storage, 2, 2, 0, faceSize, opencv.Size{0, 0})
	if len(faces) > 1 {
		fmt.Printf("Multiple faces detected\n")
		/*
		for i, face := range faces {
			x := face.X()
			y := face.Y()
			fmt.Printf("  face %d at %d, %d\n", i, x, y)
		}
		*/
	} else if len(faces) == 1 {
		face := faces[0]
		p := opencv.Point{X: face.X() + (faceSize.Width / 2), Y: face.Y() + (faceSize.Height / 2)}
		//fmt.Printf("Face detected at %d, %d\n", face.X(), face.Y())
		fmt.Printf("Face center at %d, %d\n", p.X, p.Y)
		//fmt.Printf("Image center at %d, %d\n", imageCenter.X, imageCenter.Y)

		const distanceThresh = 20
		const moveDur = 100

		if imageCenter.X - p.X > distanceThresh {
			fmt.Printf("Move Left\n")
			launcherDev.SendCommandDuration(launcher.Left, moveDur)
		} else if imageCenter.X - p.X < -distanceThresh {
			fmt.Printf("Move Right\n")
			launcherDev.SendCommandDuration(launcher.Right, moveDur)
		}
		if imageCenter.Y - p.Y > distanceThresh {
			fmt.Printf("Move Up\n")
			launcherDev.SendCommandDuration(launcher.Up, moveDur)
		} else if imageCenter.Y - p.Y < -distanceThresh {
			fmt.Printf("Move Down\n")
			launcherDev.SendCommandDuration(launcher.Down, moveDur)
		}
		return true
	}
	return false
}
