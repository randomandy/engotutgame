package scene

import (
	"image/color"
	"log"
	"time"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type Guy struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// MenuScene is responsible for managing the icon
type MenuScene struct{}

func (*MenuScene) Preload() {
	err := engo.Files.Load("icon.png")
	if err != nil {
		log.Println(err)
	}
}

func (*MenuScene) Setup(w *ecs.World) {
	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&SceneSwitcherSystem{
		NextScene: "DefaultScene",
		WaitTime:  time.Second * 3,
	})

	// Retrieve a texture
	texture, err := common.PreloadedSpriteSingle("icon.png")
	if err != nil {
		log.Println(err)
	}

	// Create an entity
	guy := Guy{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 8x
	guy.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{8, 8},
	}
	guy.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    texture.Width() * guy.RenderComponent.Scale.X,
		Height:   texture.Height() * guy.RenderComponent.Scale.Y,
	}

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		}
	}
}

func (*MenuScene) Hide() {
	log.Println("MenuScene is now hidden")
}

func (*MenuScene) Show() {
	log.Println("MenuScene is now shown")
}

func (*MenuScene) Type() string { return "MenuScene" }

// SceneSwitcherSystem is a System that actually calls SetScene
type SceneSwitcherSystem struct {
	NextScene     string
	WaitTime      time.Duration
	secondsWaited float32
}

func (*SceneSwitcherSystem) Priority() int          { return 1 }
func (*SceneSwitcherSystem) Remove(ecs.BasicEntity) {}

func (s *SceneSwitcherSystem) Update(dt float32) {
	s.secondsWaited += dt
	if float64(s.secondsWaited) > s.WaitTime.Seconds() {
		s.secondsWaited = 0

		// Change the world to s.NextScene, and don't override / force World re-creation
		engo.SetScene(&DefaultScene{}, false)

		log.Println("Switched to", s.NextScene)
	}
}
