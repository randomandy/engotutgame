package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

var (
	WalkUpAction    *common.Animation
	WalkDownAction  *common.Animation
	WalkLeftAction  *common.Animation
	WalkRightAction *common.Animation
	StopAction      *common.Animation
	SkillAction     *common.Animation
	actions         []*common.Animation

	upButton    = "up"
	downButton  = "down"
	leftButton  = "left"
	rightButton = "right"
	model       = "motw.png"
	width       = 52
	height      = 72
)

type DefaultScene struct{}

type Animation struct {
	ecs.BasicEntity
	common.AnimationComponent
	common.RenderComponent
	common.SpaceComponent
	ControlComponent
}

type ControlComponent struct {
	SchemeVert  string
	SchemeHoriz string

	// oldY is (optionally) the old Y-location of the mouse / touch - used to determine drag direction
	oldY float32
}

func (*DefaultScene) Preload() {

	engo.Files.Load(model)

	StopAction = &common.Animation{
		Name:   "stop",
		Frames: []int{1},
	}

	WalkUpAction = &common.Animation{
		Name:   "up",
		Frames: []int{36, 37, 38},
		Loop:   true,
	}

	WalkDownAction = &common.Animation{
		Name:   "down",
		Frames: []int{0, 1, 2},
		Loop:   true,
	}

	WalkLeftAction = &common.Animation{
		Name:   "left",
		Frames: []int{12, 13, 14},
		Loop:   true,
	}

	WalkRightAction = &common.Animation{
		Name:   "right",
		Frames: []int{24, 25, 26},
		Loop:   true,
	}

	actions = []*common.Animation{
		StopAction,
		WalkUpAction,
		WalkDownAction,
		WalkLeftAction,
		WalkRightAction,
	}

	engo.Input.RegisterButton(upButton, engo.W, engo.ArrowUp)
	engo.Input.RegisterButton(leftButton, engo.A, engo.ArrowLeft)
	engo.Input.RegisterButton(rightButton, engo.D, engo.ArrowRight)
	engo.Input.RegisterButton(downButton, engo.S, engo.ArrowDown)
}

func (scene *DefaultScene) Setup(w *ecs.World) {
	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.AnimationSystem{})
	w.AddSystem(&ControlSystem{})

	engo.Input.RegisterAxis("vertical", engo.AxisKeyPair{engo.ArrowUp, engo.ArrowDown})
	engo.Input.RegisterAxis("horizontal", engo.AxisKeyPair{engo.ArrowLeft, engo.ArrowRight})

	spriteSheet := common.NewSpritesheetFromFile(model, width, height)

	hero := scene.CreateEntity(engo.Point{0, 0}, spriteSheet)

	hero.ControlComponent = ControlComponent{
		SchemeHoriz: "horizontal",
		SchemeVert:  "vertical",
	}

	// Add our hero to the appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&hero.BasicEntity, &hero.RenderComponent, &hero.SpaceComponent)
		case *common.AnimationSystem:
			sys.Add(&hero.BasicEntity, &hero.AnimationComponent, &hero.RenderComponent)
		case *ControlSystem:
			sys.Add(&hero.BasicEntity, &hero.AnimationComponent, &hero.ControlComponent, &hero.SpaceComponent)
		}
	}
}

func (*DefaultScene) Type() string { return "GameWorld" }

func (*DefaultScene) CreateEntity(point engo.Point, spriteSheet *common.Spritesheet) *Animation {
	entity := &Animation{BasicEntity: ecs.NewBasic()}

	entity.SpaceComponent = common.SpaceComponent{
		Position: point,
		Width:    float32(width),
		Height:   float32(height),
	}
	entity.RenderComponent = common.RenderComponent{
		Drawable: spriteSheet.Cell(0),
		Scale:    engo.Point{1, 1},
	}
	entity.AnimationComponent = common.NewAnimationComponent(spriteSheet.Drawables(), 0.1)

	entity.AnimationComponent.AddAnimations(actions)
	entity.AnimationComponent.AddDefaultAnimation(StopAction)

	return entity
}

type controlEntity struct {
	*ecs.BasicEntity
	*common.AnimationComponent
	*ControlComponent
	*common.SpaceComponent
}

type ControlSystem struct {
	entities []controlEntity
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, anim *common.AnimationComponent, control *ControlComponent, space *common.SpaceComponent) {
	c.entities = append(c.entities, controlEntity{basic, anim, control, space})
}

func (c *ControlSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range c.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		c.entities = append(c.entities[:delete], c.entities[delete+1:]...)
	}
}

func (c *ControlSystem) Update(dt float32) {
	for _, e := range c.entities {

		if engo.Input.Button(upButton).JustPressed() {
			e.AnimationComponent.SelectAnimationByAction(WalkUpAction)
		} else if engo.Input.Button(downButton).JustPressed() {
			e.AnimationComponent.SelectAnimationByAction(WalkDownAction)
		} else if engo.Input.Button(leftButton).JustPressed() {
			e.AnimationComponent.SelectAnimationByAction(WalkLeftAction)
		} else if engo.Input.Button(rightButton).JustPressed() {
			e.AnimationComponent.SelectAnimationByAction(WalkRightAction)
		}

		speed := engo.GameWidth()*dt - 15

		vert := engo.Input.Axis(e.ControlComponent.SchemeVert)
		e.SpaceComponent.Position.Y += speed * vert.Value()

		horiz := engo.Input.Axis(e.ControlComponent.SchemeHoriz)
		e.SpaceComponent.Position.X += speed * horiz.Value()

		if (e.SpaceComponent.Height + e.SpaceComponent.Position.Y) > engo.GameHeight() {
			e.SpaceComponent.Position.Y = engo.GameHeight() - e.SpaceComponent.Height
		} else if e.SpaceComponent.Position.Y < 0 {
			e.SpaceComponent.Position.Y = 0
		}

		if (e.SpaceComponent.Width + e.SpaceComponent.Position.X) > engo.GameWidth() {
			e.SpaceComponent.Position.X = engo.GameWidth() - e.SpaceComponent.Width
		} else if e.SpaceComponent.Position.X < 0 {
			e.SpaceComponent.Position.X = 0
		}

	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Animation Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &DefaultScene{})
}
