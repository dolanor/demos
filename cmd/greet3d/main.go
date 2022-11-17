package main

import (
	_ "embed"
	"image"
	"log"
	"time"

	"github.com/dolanor/demos/greeting"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/window"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var (
	title    *gui.Label // title with wins
	numbWins = 0        //number of wins
	//go:embed earth.vert
	shaderEarthVertex string
	//go:embed earth.frag
	shaderEarthFrag string
	firstSrv        = "http://192.168.20.95:8080"
	secondSrv       = "http://localhost:8080"
)

func main() {
	a := app.App()
	scene := core.NewNode()
	e := Earth{
		app:   a,
		scene: scene,
	}
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	e.setupGUI()

	e.start()
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
		e.mainPanel.SetSize(float32(width), float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
		e.Update(deltaTime)
	})
}

func (e *Earth) setupGUI() {
	dl := gui.NewDockLayout()
	width, height := e.app.GetSize()
	e.mainPanel = gui.NewPanel(float32(width), float32(height))
	e.mainPanel.SetRenderable(false)
	e.mainPanel.SetEnabled(false)
	e.mainPanel.SetLayout(dl)
	e.scene.Add(e.mainPanel)
	gui.Manager().Set(e.mainPanel)

	headerColor := math32.Color4{R: 13.0 / 256.0, G: 41.0 / 256.0, B: 62.0 / 256.0, A: 1}
	lightTextColor := math32.Color4{R: 0.8, G: 0.8, B: 0.8, A: 1}
	header := gui.NewPanel(600, 40)
	header.SetBorders(0, 0, 1, 0)
	header.SetPaddings(4, 4, 4, 4)
	header.SetColor4(&headerColor)
	header.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockTop})

	// Horizontal box layout for the header
	hbox := gui.NewHBoxLayout()
	header.SetLayout(hbox)
	e.mainPanel.Add(header)

	// Header title
	const fontSize = 50
	title = gui.NewLabel(" ")
	title.SetFontSize(fontSize)
	title.SetLayoutParams(&gui.HBoxLayoutParams{AlignV: gui.AlignCenter})
	title.SetText(greeting.Greet("Dagger friends!"))
	title.SetColor4(&lightTextColor)
	header.Add(title)

}

type Earth struct {
	app *app.Application

	mainPanel *gui.Panel

	sphere *graphic.Mesh
	scene  *core.Node
}

func (e *Earth) start() {
	// Create Skybox
	skybox, err := graphic.NewSkybox(graphic.SkyboxData{
		DirAndPrefix: "./images/space/dark-s_", Extension: "jpg",
		Suffixes: [6]string{"px", "nx", "py", "ny", "pz", "nz"}})
	if err != nil {
		panic(err)
	}
	e.scene.Add(skybox)

	// Adds directional front light
	dir1 := light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, 0.9)
	dir1.SetPosition(0, 0, 100)
	e.scene.Add(dir1)

	// Create custom shader
	e.app.Renderer().AddShader("shaderEarthVertex", shaderEarthVertex)
	e.app.Renderer().AddShader("shaderEarthFrag", shaderEarthFrag)
	e.app.Renderer().AddProgram("shaderEarth", "shaderEarthVertex", "shaderEarthFrag")

	// Helper function to load a texture and handle errors
	newTexture := func(path string) *texture.Texture2D {
		tex, err := texture.NewTexture2DFromImage(path)
		if err != nil {
			log.Fatalf("Error loading texture: %s", err)
		}
		tex.SetFlipY(false)
		return tex
	}

	// Create earth textures
	texDay := newTexture("./images/earth_clouds_big.jpg")
	texSpecular := newTexture("./images/earth_spec_big.jpg")
	texNight := newTexture("./images/earth_night_big.jpg")
	//texBump, err := newTexture("./images/earth_bump_big.jpg")
	texGreeting := text2Tex(greeting.Greet("Dagger friends"))

	// Create custom material using the custom shader
	matEarth := NewEarthMaterial(&math32.Color{R: 1, G: 1, B: 1})
	matEarth.SetShininess(20)
	//matEarth.SetSpecularColor(&math32.Color{0., 1, 1})
	//matEarth.SetColor(&math32.Color{0.8, 0.8, 0.8})

	matEarth.AddTexture(texDay)
	matEarth.AddTexture(texSpecular)
	matEarth.AddTexture(texNight)
	matEarth.AddTexture(texGreeting)

	// Create sphere
	geom := geometry.NewSphere(1, 32, 16)
	e.sphere = graphic.NewMesh(geom, matEarth)
	e.scene.Add(e.sphere)

	// Create sun sprite
	texSun, err := texture.NewTexture2DFromImage("./images/lensflare0_alpha.png")
	if err != nil {
		log.Fatalf("Error loading texture: %s", err)
	}
	sunMat := material.NewStandard(&math32.Color{R: 1, G: 1, B: 1})
	sunMat.AddTexture(texSun)
	sunMat.SetTransparent(true)
	sun := graphic.NewSprite(10, 10, sunMat)
	sun.SetPositionZ(20)
	e.scene.Add(sun)

	// Add axes helper
	//axes := helper.NewAxes(5)
	//e.scene.Add(axes)
}

// Update is called every frame.
func (t *Earth) Update(deltaTime time.Duration) {
	t.sphere.RotateY(-0.5 * float32(deltaTime.Seconds()))
}

type EarthMaterial struct {
	material.Standard // Embedded standard material
}

// NewEarthMaterial creates and returns a pointer to a new earth material
func NewEarthMaterial(color *math32.Color) *EarthMaterial {

	pm := new(EarthMaterial)
	pm.Standard.Init("shaderEarth", color)
	return pm
}

func text2Tex(text string) *texture.Texture2D {
	const (
		width        = 350
		height       = 100
		startingDotX = 6
		startingDotY = 60
	)

	f, err := opentype.Parse(goitalic.TTF)
	if err != nil {
		log.Fatalf("Parse: %v", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("NewFace: %v", err)
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	d := font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(startingDotX, startingDotY),
	}
	_ = d
	d.DrawString(text)

	tex := texture.NewTexture2DFromRGBA(dst)
	return tex
}
