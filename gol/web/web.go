package web

import (
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

var indexTemplate = template.Must(template.ParseFiles("gol/web/static/index.html"))
var upgrader = websocket.Upgrader{}

const ROWS int64 = 100
const COLS int64 = 100

type Grid [ROWS][COLS]int64

func mod(a, b int64) (c int64) {
	c = ((a % b) + b) % b
	return
}

func (g *Grid) randInit(threshold float64) {
	for y := range ROWS {
		for x := range COLS {
			if rand.Float64() < threshold {
				g[y][x] = 1
			} else {
				g[y][x] = 0
			}
		}
	}
}

func (g *Grid) gliderInit(offset int) {
	g[0+offset][1+offset] = 1
	g[1+offset][2+offset] = 1
	g[2+offset][0+offset] = 1
	g[2+offset][1+offset] = 1
	g[2+offset][2+offset] = 1
}

func (g Grid) update() Grid {
	var newGrid Grid
	for y := range ROWS {
		for x := range COLS {
			var neighbors int64

			neighbors += g[mod(y-1, ROWS)][mod(x-1, COLS)]
			neighbors += g[mod(y-1, ROWS)][mod(x, COLS)]
			neighbors += g[mod(y-1, ROWS)][mod(x+1, COLS)]

			neighbors += g[mod(y, ROWS)][mod(x-1, COLS)]
			neighbors += g[mod(y, ROWS)][mod(x+1, COLS)]

			neighbors += g[mod(y+1, ROWS)][mod(x-1, COLS)]
			neighbors += g[mod(y+1, ROWS)][mod(x, COLS)]
			neighbors += g[mod(y+1, ROWS)][mod(x+1, COLS)]

			if neighbors < 2 || neighbors > 3 {
				newGrid[y][x] = 0
			}
			if (neighbors == 2 || neighbors == 3) && g[y][x] == 1 {
				newGrid[y][x] = 1
			}
			if neighbors == 3 && g[y][x] == 0 {
				newGrid[y][x] = 1
			}
		}
	}
	return newGrid
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		data := map[string]string{
			"Width":         strconv.FormatInt(COLS*10, 10),
			"Height":        strconv.FormatInt(ROWS*10, 10),
			"WebSocketPath": "ws://" + r.Host + "/grid",
		}
		indexTemplate.Execute(w, data)
		return
	}
}

func grid(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	log.Print("User connected")

	var gridBoard Grid
	gridBoard.randInit(0.4)
	// gridBoard.gliderInit(5)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			log.Print("User disconnected")
			return
		case <-ticker.C:
			err = c.WriteJSON(gridBoard)
			if err != nil {
				log.Println("write:", err)
				return
			}
			gridBoard = gridBoard.update()
		}
	}
}

func Run() {
	// log.SetFlags(0)

	_, file, _, _ := runtime.Caller(0)

	staticPath := filepath.Join(filepath.Dir(file), "static")
	fs := http.FileServer(http.Dir(staticPath))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/grid", grid)
	http.HandleFunc("/", index)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
