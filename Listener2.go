package main

import (
    "bufio"
    "fmt"
    "math"
    "math/rand"
    "net"
    "strconv"
    "strings"
    "time"
    //"log"
)

//https://en.wikipedia.org/wiki/Maze_generation_algorithm
const N = 8
const M = 13
const NS = 4 //Numero de servidores
const NC = 3 //Numero de clientes

var servers [NS]string
var allClients map[*Client]string
var userlist [NC]string
var labyrinth [N][M]*Node
var nombreUsuario string
var ReadServer chan string = make(chan string)
var SendServer chan string = make(chan string)
var estado string = "M"
var nServidor int = 0
var mutex_Client bool = true

//------------------------------------------------------------------

//------------------------------------------------------------------
func hablarConEsclavitos(server string, i int){
    conn, err := net.Dial("tcp", server)
    if err != nil {
        fmt.Println("Servidor:", i+1, " está desconectado")
        return
    }
    nWriter := bufio.NewWriter(conn)
    for data := range(SendServer){

        _,err:=nWriter.WriteString(data)
        
        nWriter.Flush()
        if err != nil {
            fmt.Println("Servidor Desconectado")
        }
    }
    conn.Close()

}
func mandarEstadoDeJuego() {
    for i := 0; i < NS; i++ {
        if i != nServidor {
                go hablarConEsclavitos(servers[i],i)
        }
    }
}

type Client struct {
    outgoing chan string
    // reader   *bufio.Reader
    writer *bufio.Writer
    conn   net.Conn
    x      float64
    y      float64
    nJ     float64
}

type Node struct {
    visited    int
    node_up    *Node
    node_down  *Node
    node_left  *Node
    node_right *Node
}

type NodexDir struct {
    node *Node
    dir  string
}

func NewNode() *Node {
    node := &Node{
        visited:    0,
        node_up:    nil,
        node_down:  nil,
        node_left:  nil,
        node_right: nil,
    }
    return node
}
func (nodexdir *NodexDir) getNode() *Node {
    return nodexdir.node
}
func (nodexdir *NodexDir) getDir() string {
    return nodexdir.dir
}
func NewNodexDir(node *Node, dir string) *NodexDir {
    nodexdir := &NodexDir{
        node: node,
        dir:  dir,
    }
    return nodexdir
}
func (node *Node) setNodeUp(conecNode *Node) {
    node.node_up = conecNode
}
func (node *Node) setNodeDown(conecNode *Node) {
    node.node_down = conecNode
}
func (node *Node) setNodeLeft(conecNode *Node) {
    node.node_left = conecNode
}
func (node *Node) setNodeRight(conecNode *Node) {
    node.node_right = conecNode
}

func (node *Node) visit() {
    node.visited = 1
}

func (node *Node) isVisited() int {
    return node.visited
}

func (node *Node) getNodeUp() *Node {
    return node.node_up
}
func (node *Node) getNodeDown() *Node {
    return node.node_down
}
func (node *Node) getNodeLeft() *Node {
    return node.node_left
}
func (node *Node) getNodeRight() *Node {
    return node.node_right
}

func (node *Node) getNeighborsNV() []*NodexDir {
    returnList := []*NodexDir{}
    if node.node_up != nil {
        if node.node_up.isVisited() == 0 {
            returnList = append(returnList, NewNodexDir(node, "up"))
        }
    }
    if node.node_down != nil {
        if node.node_down.isVisited() == 0 {
            returnList = append(returnList, NewNodexDir(node, "down"))
        }
    }
    if node.node_left != nil {
        if node.node_left.isVisited() == 0 {
            returnList = append(returnList, NewNodexDir(node, "left"))
        }
    }
    if node.node_right != nil {
        if node.node_right.isVisited() == 0 {
            returnList = append(returnList, NewNodexDir(node, "right"))
        }
    }
    return returnList

}

func createMatrixOfLaby() [N][M]*Node {
    fmt.Println(N)
    for i := 1; i <= N; i++ {
        for j := 1; j <= M; j++ {
            labyrinth[i-1][j-1] = NewNode()
            //fmt.Print(1)
        }
        fmt.Print("\n")
    }
    return labyrinth
}

func createConectionsOfLaby(labyrinth [N][M]*Node) [N][M]*Node {
    for i := 1; i <= N; i++ {
        for j := 1; j <= M; j++ {
            if i > 1 {
                labyrinth[i-1][j-1].setNodeUp(labyrinth[i-2][j-1])
            }
            if i < N {
                labyrinth[i-1][j-1].setNodeDown(labyrinth[i][j-1])
            }
            if j > 1 {
                labyrinth[i-1][j-1].setNodeLeft(labyrinth[i-1][j-2])
            }
            if j < M {
                labyrinth[i-1][j-1].setNodeRight(labyrinth[i-1][j])
            }

        }
    }
    return labyrinth
}

func createMaze(labyrinth [N][M]*Node) [N][M]*Node {
    var actualNode *Node = labyrinth[0][0]
    actualNode.visit()
    var listNeighbors []*NodexDir = actualNode.getNeighborsNV()
    fmt.Println(listNeighbors)
    for len(listNeighbors) > 0 {
        rand.Seed(time.Now().UnixNano())
        i := rand.Intn(len(listNeighbors))
        actualNodexDir := listNeighbors[i]
        actualNode := actualNodexDir.getNode()
        actualNode.visit()
        //fmt.Println("Dir :", actualNodexDir.getDir())
        //fmt.Println(actualNode.node_up, actualNode.node_down, actualNode.node_left, actualNode.node_right, "\n")
        listNeighbors = append(listNeighbors[:i], listNeighbors[i+1:]...)
        if actualNodexDir.getDir() == "up" {
            if actualNode.node_up.isVisited() == 0 {
                actualNode.node_up.visit()
                actualNode.node_up.node_down = nil
                listNeighbors = append(listNeighbors, actualNode.node_up.getNeighborsNV()...)
                actualNode.node_up = nil
            }
        } else if actualNodexDir.getDir() == "down" {
            if actualNode.node_down.isVisited() == 0 {
                actualNode.node_down.visit()
                actualNode.node_down.node_up = nil
                listNeighbors = append(listNeighbors, actualNode.node_down.getNeighborsNV()...)
                actualNode.node_down = nil
            }
        } else if actualNodexDir.getDir() == "left" {
            if actualNode.node_left.isVisited() == 0 {
                actualNode.node_left.visit()
                actualNode.node_left.node_right = nil
                listNeighbors = append(listNeighbors, actualNode.node_left.getNeighborsNV()...)
                actualNode.node_left = nil
            }
        } else if actualNodexDir.getDir() == "right" {
            if actualNode.node_right.isVisited() == 0 {
                actualNode.node_right.visit()
                actualNode.node_right.node_left = nil
                listNeighbors = append(listNeighbors, actualNode.node_right.getNeighborsNV()...)
                actualNode.node_right = nil
            }
        }
    }
    return labyrinth
}

func createLabyRaw() [N][M]*Node {
    return createMaze(createConectionsOfLaby(createMatrixOfLaby()))
}

func translationOfLaby(labyrinth [N][M]*Node) string {
    var tamañoI int = 3 + (2 * (N - 1))
    var tamañoJ int = 3 + (2 * (M - 1))
    fmt.Println(tamañoI)
    fmt.Println(tamañoJ)
    var finalString string = ""
    for i := 1; i <= tamañoI; i++ {
        for j := 1; j <= tamañoJ; j++ {
            if (i == 1 || i == tamañoI) && j < tamañoJ {
                finalString = finalString + "1,"
            } else if j == 1 {
                finalString = finalString + "1,"
            } else if j == tamañoJ {
                finalString = finalString + "1*"
            } else if math.Mod(float64(j), 2) == 1 && math.Mod(float64(i), 2) == 1 {
                finalString = finalString + "1,"
            } else if math.Mod(float64(j), 2) == 0 && math.Mod(float64(i), 2) == 0 {
                finalString = finalString + "0,"
            } else {
                if math.Mod(float64(j), 2) == 1 && math.Mod(float64(i), 2) == 0 {
                    if labyrinth[(i/2)-1][((j-1)/2)-1].getNodeRight() != nil {
                        finalString = finalString + "1,"
                    } else {
                        finalString = finalString + "0,"
                    }
                } else if math.Mod(float64(j), 2) == 0 && math.Mod(float64(i), 2) == 1 {
                    if labyrinth[((i-1)/2)-1][(j/2)-1].getNodeDown() != nil {
                        finalString = finalString + "1,"
                    } else {
                        finalString = finalString + "0,"
                    }
                }

            }
        }
    }
    return finalString
}

func (client *Client) Write() {
    for data := range client.outgoing {
        _, err := client.writer.WriteString(data)
        client.writer.Flush()
        if err != nil {
            fmt.Println("El cliente: ", allClients[client], " se ha desconectado")
            delete(allClients, client)
            //Mandar todo a otro servidor
        }
    }
}

func (client *Client) Listen() {
    go client.Write()
}

func NewClient(connection net.Conn) *Client {
    writer := bufio.NewWriter(connection)
    //reader := bufio.NewReader(connection)

    client := &Client{
        // incoming: make(chan string),
        outgoing: make(chan string),
        conn:     connection,
        //reader:   reader,
        writer: writer,
        x:      45,
        y:      45,
        nJ:     0,
    }
    client.Listen()

    return client
}

func setNJClient() {
    var cont float64 = 1
    for cli := range allClients {
        cli.nJ = cont
        cont++
    }
}

/*func ServerWrite(s string) {
    for {
        var b []byte = make([]byte, 1)
        os.Stdin.Read(b)
        for client, i := range allClients {
            client.outgoing <- s + ":"
            fmt.Println(i)
        }
    }
}*/
func ServerWriteIA() {
    for {
        for client := range allClients {
            client.outgoing <- "ALIVE:\n"
            time.Sleep(500 * time.Millisecond)
        }
    }
}
func ServerWriteAMsg(s string) {
    for client, i := range allClients {
        client.outgoing <- s
        fmt.Println("Menssage Sent to Client:", i)
    }
}
func GetAllUsers() string {
    s := ""
    for cli, i := range allClients {
        s = s + i + ";" + strconv.FormatFloat(cli.x, 'g', -1, 64) + ";" + strconv.FormatFloat(cli.y, 'g', -1, 64) + ";" + strconv.FormatFloat(cli.nJ, 'g', -1, 64) + ","
    }
    if s != "" {
        s = s[:len(s)-1]
    }
    fmt.Println(s)
    return s
}

func ServerRead(nReader *bufio.Reader) {
    for {
        //fmt.Print(1)
        line, err := nReader.ReadString('\n')
        if err == nil {
            //fmt.Print(line)
            //fmt.Print("antes")
            //go func() {
            ReadServer <- line
            //}()
            //fmt.Print("despues")
        }
    }
}

func escucharServers() {
    listener2, _ := net.Listen("tcp", servers[nServidor])
    for {
        conn, err := listener2.Accept()
        if err != nil {
            fmt.Println(err.Error())
        }
        nReader := bufio.NewReader(conn)
        go ServerRead(nReader)
    }
}
var mutex_Server bool = true

func Decode() {
    var x, y, nombre string
    for {
        for i := range ReadServer {
            //fmt.Println("i:", i)
            chain := strings.Split(i, ":")
            //fmt.Println(chain)
            if chain[0] == "USERNAME" {
                //fmt.Println("Prueba1", chain)
                fmt.Println("true")
                nombreUsuario = chain[1]
            }
            if chain[0] == "COORDS" {
                fmt.Println(chain)
                x = chain[1]
                y = chain[2]
                nombre = chain[3]
                fmt.Println("Nombre:", nombre)
                for !mutex_Client{}
                mutex_Client=false
                for cli, n := range allClients {
                    fmt.Println("Name:", n)
                    if n == nombre {
                        cli.x, _ = strconv.ParseFloat(x, 64)
                        cli.y, _ = strconv.ParseFloat(y, 64)
                    }

                }
                LU := "USERS:" + GetAllUsers() + ":\n"
                ServerWriteAMsg(LU)
                LUL := "USERLIST:" + GetAllUsers() + ":\n"
                
                go func () {
                    for !mutex_Server{}
                    mutex_Server = false
                    SendServer <- LUL
                    fmt.Println("SendServer: ",SendServer)
                    mutex_Server = true
                } ()  
                go mandarEstadoDeJuego()
            }
            if chain[0] == "USERLIST" {
                chain2 := strings.Split(chain[1], ",")
                for i := 0; i < NC; i++ {
                    userlist[i] = chain2[i]
                }
            }
        }
    }
}

func reconectarCliente() {
    for i := 0; i < NC; i++ {

    }
}

func main() {
    servers[0] = "192.168.121.26:8081"
    servers[1] = "192.168.121.24:8082"
    servers[2] = "ip:8083"
    servers[3] = "ip:8084"
    //exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    allClients = make(map[*Client]string)
    go escucharServers()
    listener, _ := net.Listen("tcp", "192.168.121.26:8080")
    var Laby [N][M]*Node = createLabyRaw()
    var matrix string = "MATRIX:" + translationOfLaby(Laby) + ":\n"
    fmt.Println(matrix)
    go Decode()
    go ServerWriteIA()
    //go ServerWriteAMsg(s)
    for {
        for estado == "M" {
            conn, err := listener.Accept()
            if err != nil {
                fmt.Println(err.Error())
            }
            nReader := bufio.NewReader(conn)
            go ServerRead(nReader)
            listOfUsers := "USERS:" + GetAllUsers() + ":\n"
            client := NewClient(conn)
            allClients[client] = "temp"
            setNJClient()
            ServerWriteAMsg(matrix)
            ServerWriteAMsg(listOfUsers)
            for nombreUsuario == "" {
            }
            allClients[client] = nombreUsuario
            nombreUsuario = ""
            listOfUsers = "USERS:" + GetAllUsers() + ":\n"
            ServerWriteAMsg(listOfUsers)
            //go handleConection(nReader,allClients[client])
        }
    }
}
