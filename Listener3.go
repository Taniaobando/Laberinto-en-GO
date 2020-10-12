package main

import (
    "bufio"
    "fmt"
    "net"
    "strconv"
    "strings"
    "time"
)

//https://en.wikipedia.org/wiki/Maze_generation_algorithm
const N = 8
const M = 13
const NS = 4 //Numero de servidores
const NC = 3 //Numero de clientes

var servers [NS]string
var allClients map[*Client]string
var userlist [NC]string
var empty [NC]string
var nombreUsuario string
var ReadServer chan string = make(chan string)
var estado string = "E"
var nServidor int = 1
var mutex_varServer bool = true

//------------------------------------------------------------------

//------------------------------------------------------------------

func mandarEstadoDeJuego() {
    for i := 0; i < NS; i++ {
        if i != nServidor {
            conn, err := net.Dial("tcp", servers[i])
            if err != nil {
                fmt.Println("Servidor:", i+1, " está desconectado", err)
                continue
            }
            nWriter := bufio.NewWriter(conn)
            listOfUsers := "USERLIST:" + GetAllUsers() + ":\n"
            nWriter.WriteString(listOfUsers)
            conn.Close()
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
            client.outgoing <- "ALIVE" + ":\n"
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
            fmt.Println(line, "line")
            for !mutex_varServer {
            }
            mutex_varServer = false
            ReadServer <- line
            mutex_varServer = true
            fmt.Println("despues")
        } else {
            fmt.Println(err)
        }
    }
}

func escucharServers(listener2 net.Listener) {
    for {
        conn, err := listener2.Accept()
        if err != nil {
            fmt.Println(err.Error())
        }
        nReader2 := bufio.NewReader(conn)
        go ServerRead(nReader2)
    }
}

func Decode() {
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
                x := chain[1]
                y := chain[2]
                nombre := chain[3]
                for cli, n := range allClients {
                    if n == nombre {
                        cli.x, _ = strconv.ParseFloat(x, 64)
                        cli.y, _ = strconv.ParseFloat(y, 64)
                    }

                }
                LU := "USERS:" + GetAllUsers() + ":\n"
                ServerWriteAMsg(LU)
                go mandarEstadoDeJuego()
            }
            if chain[0] == "USERLIST" {
                chain2 := strings.Split(chain[1], ",")
                for i := 0; i < NC; i++ {
                    fmt.Println("llegó a userlist=chain2")
                    userlist[i] = chain2[i]
                }
            }
            if chain[0] == "USER_RECONECT" {
                //fmt.Println("Prueba1", chain)
                fmt.Println("true")
                nombreUsuario = chain[1]
                estado = "M"
            }
        }
    }
}
func (client *Client) reconectarCliente() {
    for i := 0; i < NC; i++ {
        Jugador := strings.Split(userlist[i], ";")
        if Jugador[0] == nombreUsuario {
            client.x, _ = strconv.ParseFloat(Jugador[1], 64)
            client.y, _ = strconv.ParseFloat(Jugador[2], 64)
            client.nJ, _ = strconv.ParseFloat(Jugador[3], 64)
        }
    }
}

func main() {
    servers[0] = "192.168.121.26:8081"
    servers[1] = "192.168.121.24:8082"
    servers[2] = "ip:8083"
    servers[3] = "ip:8084"
    //exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    allClients = make(map[*Client]string)
    listener2, _ := net.Listen("tcp", servers[nServidor])
    go escucharServers(listener2)
    listener, _ := net.Listen("tcp", "192.168.121.24:8080")
    go Decode()
    go ServerWriteIA()
    for {
        for estado == "M" {
            conn, err := listener.Accept()
            if err != nil {
                fmt.Println(err.Error())
            }
            fmt.Println("jelou")
            nReader := bufio.NewReader(conn)
            go ServerRead(nReader)
            for userlist == empty {
            }
            client := NewClient(conn)
            client.reconectarCliente()
            allClients[client] = nombreUsuario
            nombreUsuario = ""
        }
    }
}
