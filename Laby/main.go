package main

import (
"bufio"
"fmt"
"github.com/veandco/go-sdl2/sdl"
"strings"
//"io/ioutil"
"net"
"os"
"strconv"
"time"
)

const (
ancho       = 1080
alto        = 680
tpsobjetivo = 60
)
const NS = 4

var servers [NS]string
var jugadores []*jugador
var userList []string
var str string = ""
var Read chan string = make(chan string)
var mutex_Sincronizar = true
var mutex_Escritura = true
var jugadoract string
var conectado bool = false
var nReader *bufio.Reader
var nWriter *bufio.Writer

func ClientRead() {
for {
line, err := nReader.ReadString('\n')
if err == nil {
//fmt.Println(line)
//go func() {
Read <- line
//}()
} else {
break
}
}
}
func ClientWriteIA() {
	for {
		if mutex_Escritura {
			mutex_Escritura = false
			_, err := nWriter.WriteString("ALIVE:\n")
			if err != nil {
				fmt.Println("Servidor desconectado")
				mutex_Escritura = true
				conectado = false
				break
			}
			nWriter.Flush()
			mutex_Escritura = true
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func ClientWriteAMsg(s string) {
	for !mutex_Escritura {}
	mutex_Escritura = false
	_, err := nWriter.WriteString(s)
	if err != nil {
		fmt.Println("Servidor desconectado", err)
		mutex_Escritura= true
		return
		}
	fmt.Println("Mensaje enviado")
	nWriter.Flush()
	mutex_Escritura = true
}

func Decode() {
	for i := range Read {
		//fmt.Println("i:", i)
		chain := strings.Split(i, ":")
		//fmt.Println(chain[0])
		switch chain[0] {
			case "MATRIX":
				//fmt.Println("Prueba1", chain)
				str = chain[1]
			case "USERS":
				mutex_Sincronizar = false
				if chain[1] != "" {
					userList = strings.Split(chain[1], ",")
					//fmt.Println(userList)
					sincronizar()
					mutex_Sincronizar = true
				}
			case "WINNER":
				for{
					fmt.Println("GANADOR: ", chain[1])
				}
			default:
		}
	}
}
func ExistingUser(name string) bool {
for _, j := range userList {
//fmt.Println(j)
k := strings.Split(j, ";")
i := k[0]
//fmt.Println(i)
if i == name {
return true
}
}
return false
}

func UserName() string {
inputUserName := bufio.NewReader(os.Stdin)
fmt.Print("Inserte su nombre de usuario: ")
name, _ := inputUserName.ReadString('\n')
name = name[:len(name)-1]
var in bool = true
for in {
in = ExistingUser(name)
if in == true {
fmt.Print("El nombre de usuario ya ha sido utilizado escoja uno nuevo : ")
name, _ = inputUserName.ReadString('\n')
name = name[:len(name)-1]
}
}
return name

}

func matriz(cad string, m int, n int) [][]string {
mtr := make([][]string, m)
cadena := strings.Split(cad, "*")
for i := 0; i < m; i++ {
cad2 := strings.Split(cadena[i], ",")
mtr[i] = make([]string, n)
for j := 0; j < n; j++ {
mtr[i][j] = cad2[j]
}
}
return mtr
}

func sincronizar() {

for _, j := range userList {
split := strings.Split(j, ";")
xnew, _ := strconv.ParseFloat(split[1], 64)
ynew, _ := strconv.ParseFloat(split[2], 64)
t := split[0]
for _, jug := range jugadores {
if t == jug.tag && t != jugadoract {
jug.x = xnew
jug.y = ynew
}
}
}
}

func intentarReconectar() {
	cont:=1
	for _, i := range servers {
		fmt.Println("Intentando Reconectar en el servidor:",cont)
	conn, err := net.Dial("tcp", i)
	if err != nil {
		cont ++
		continue
	}
	nReader = bufio.NewReader(conn)
	nWriter = bufio.NewWriter(conn)
	go ClientWriteIA()
	go ClientRead()
	fmt.Println("antes")
	ClientWriteAMsg("USER_RECONECT:" + jugadoract + ":\n")
	fmt.Println("despues")
	conectado = true
	fmt.Println("Conexión exitosa al servidor:",cont)
	return
	}
}

var muros []muro
var infin []muro

func main() {
	servers[0] = "192.168.121.11:8080"
	servers[1] = "192.168.121.5:8080"
	servers[2] = "192.168.121.31:8080"
	servers[3] = "192.168.121.19:8080"
	conn, err := net.Dial("tcp", "192.168.121.11:8080")
	if err != nil {
	panic(err)
	}
	conectado = true
	nReader = bufio.NewReader(conn)
	nWriter = bufio.NewWriter(conn)
	go ClientWriteIA()
	go ClientRead()
	go Decode()
	jugadoract = UserName()
	name := "USERNAME:" + jugadoract + ":\n"
	//fmt.Println(name)
	go ClientWriteAMsg(name)
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
	fmt.Println("Iniciando SDL:", err)
	return
	}
	var m, n int
	m = 17
	n = 27

	//str="1,0,1,1,0,0,0,1,0,1,0,1,1,1,0,1,1,0,1*1,0,1,1,0,0,0,0,0,1,0,1,0,1,0,1,1,0,1*1,0,1,1,0,0,0,1,0,1,0,1,1,0,0,1,1,1,1*1,0,1,1,0,0,0,1,0,1,0,1,0,1,0,1,1,0,1*1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1*1,0,1,1,0,0,0,0,0,1,0,1,0,0,1,1,1,0,1*1,0,1,1,0,1,1,0,0,1,0,1,0,1,0,1,1,0,1*1,1,1,1,0,0,0,1,0,1,0,1,0,1,0,1,1,0,1*1,0,1,1,0,0,0,0,0,1,0,1,0,0,0,1,1,0,1*1,0,1,1,0,1,0,1,0,1,0,1,0,0,0,1,1,0,1*1,0,1,1,0,0,0,0,0,1,0,1,0,1,0,1,1,0,1*1,0,1,1,0,0,1,0,0,1,0,1,0,1,0,1,1,0,1*1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1"
	//str = "1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1*1,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,1,0,1,0,0,0,0,0,1*1,0,1,0,1,0,1,1,1,1,1,1,1,0,1,1,1,0,1,0,1,0,1,0,1,0,1*1,0,1,0,1,0,1,0,0,0,0,0,1,0,1,0,1,0,1,0,0,0,1,0,1,0,1*1,0,1,0,1,1,1,0,1,1,1,1,1,0,1,0,1,1,1,0,1,1,1,1,1,0,1*1,0,1,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,1,0,0,0,1,0,0,0,1*1,1,1,0,1,1,1,0,1,1,1,0,1,1,1,0,1,1,1,0,1,1,1,1,1,0,1*1,0,1,0,1,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,1,0,1*1,0,1,0,1,1,1,0,1,1,1,0,1,1,1,0,1,1,1,0,1,1,1,1,1,0,1*1,0,0,0,0,0,1,0,0,0,1,0,0,0,1,0,1,0,0,0,0,0,0,0,1,0,1*1,1,1,1,1,0,1,0,1,0,1,1,1,1,1,1,1,1,1,0,1,0,1,0,1,1,1*1,0,0,0,0,0,1,0,1,0,0,0,0,0,1,0,0,0,1,0,1,0,1,0,1,0,1*1,0,1,0,1,1,1,1,1,1,1,0,1,1,1,0,1,0,1,0,1,0,1,0,1,0,1*1,0,1,0,0,0,1,0,0,0,0,0,0,0,0,0,1,0,1,0,1,0,1,0,0,0,1*1,0,1,0,1,1,1,1,1,0,1,1,1,0,1,0,1,1,1,1,1,1,1,1,1,0,1*1,0,1,0,1,0,0,0,0,0,0,0,1,0,1,0,0,0,0,0,0,0,0,0,1,0,1*1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1*"
	for str == "" {
	//fmt.Println("str=", str)
	}
	tablero := matriz(str, m, n)
	window, err := sdl.CreateWindow(
	"Laby",
	50, 50,
	ancho, alto,
	sdl.WINDOW_OPENGL)
	if err != nil {
	fmt.Println("Iniciando SDL:", err)
	return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
	fmt.Println("Iniciando SDL:", err)
	return
	}
	defer renderer.Destroy()

	var incx, incy float64
	//fmt.Println(userList)
	for {
	if len(userList) == 2 {
	break
	}
	fmt.Println("Esperando a otro jugador.")
	time.Sleep(1000 * time.Millisecond)
	}

	for _, j := range userList {
	split := strings.Split(j, ";")
	x, _ := strconv.ParseFloat(split[1], 64)
	y, _ := strconv.ParseFloat(split[2], 64)
	num := split[3]
	t := split[0]
	jug, _ := nuevoJugador(renderer, x, y, num, t)
	jugadores = append(jugadores, &jug)
	}
	if err != nil {
	fmt.Println("Creando jugador:", err)
	return
	}
	for i := 0; i < m; i++ {
	incx = 0
	for j := 0; j < n; j++ {
	if i == 1 && j == 1 {
	x := incx
	y := incy
	ob, err := nuevoMuro(renderer, x, y, "sprites/inicio.bmp")
	if err != nil {
	fmt.Println("Creando muro: ", err)
	}
	infin = append(infin, ob)
	} else if tablero[i][j] == "1" {
	x := incx
	y := incy
	ob, err := nuevoMuro(renderer, x, y, "sprites/muro.bmp")
	if err != nil {
	fmt.Println("Creando muro: ", err)
	return
	}
	muros = append(muros, ob)
	} else if i == m-2 && j == n-2 {
	x := incx
	y := incy
	ob, err := nuevoMuro(renderer, x, y, "sprites/final.bmp")
	if err != nil {
	fmt.Println("Creando muro: ", err)
	}
	infin = append(infin, ob)
	}
	incx += TAMMURO
	}
	incy += TAMMURO
	}
	for {
		for conectado {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
			return
			}
			}

			renderer.SetDrawColor(129, 228, 198, 1)
			renderer.Clear()

			for _, mr := range muros {
			mr.dibujar(renderer)
			}

			for _, inf := range infin {
			inf.dibujar(renderer)
			}

			for _, player := range jugadores {
				player.dibujar(renderer)
				if player.tag == jugadoract {
					if mutex_Sincronizar {
						if player.actualizar() {
							fmt.Println(jugadoract)
							cadena := "COORDS:" + strconv.FormatFloat(player.x, 'g', -1, 64) + ":" + strconv.FormatFloat(player.y, 'g', -1, 64) + ":" + player.tag + ":\n"
							ClientWriteAMsg(cadena)
							fmt.Println(cadena)
							if player.ganar() {
								wincad := "WIN:"+player.tag+":\n"
								ClientWriteAMsg(wincad)
							}
						}
					}
				}
			}
			renderer.Present()
		}
		intentarReconectar()
		if !conectado {
			window.Destroy()
			break
		}
	}
}

