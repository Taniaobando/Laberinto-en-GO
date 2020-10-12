package main

import(
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)
const(
	TAMMURO = 40
)
type muro struct{
	tex *sdl.Texture
	x,y float64
}

func nuevoMuro(renderer *sdl.Renderer,x, y float64, file string) (m muro, err error){
	img, err := sdl.LoadBMP(file)
	if err != nil{
		return muro{}, fmt.Errorf("Cargando bmp: %v", err)
	}
	defer img.Free()
	m.tex, err = renderer.CreateTextureFromSurface(img)
	if err != nil{
		return muro{}, fmt.Errorf("Cargando textura: %v",err)
	}
	m.x = x
	m.y = y

	return m, nil
}

func (m *muro)dibujar(renderer *sdl.Renderer){
	x := m.x
	y := m.y

	renderer.Copy(m.tex,
		&sdl.Rect{X: 0, Y: 0, W: 300, H: 320},
		&sdl.Rect{X: int32(x),Y: int32(y), W: TAMMURO, H:TAMMURO},
	)
}
