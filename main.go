package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strings"

	dd "github.com/Robert-Duck-by-BB-SR/talking_pond/internal/duck_dom"
	"golang.org/x/term"

	tpc "github.com/Robert-Duck-by-BB-SR/talking_pond/internal/tps_client"
)

// var frame_chars = []byte{' ', '`', '.', ',', '~', '+', '*', '&', '#', '@'}

// type CharMeDaddy struct {
// 	char, count, r, g, b byte
// }

// func encode_frame(img image.Image) []byte {
// 	orig_bounds := img.Bounds().Max
//
// 	scale_x := orig_bounds.X / 80
// 	scale_y := orig_bounds.Y / 40
// 	new_img_x := orig_bounds.X / scale_x
// 	new_img_y := orig_bounds.Y / scale_y
//
// 	encoded_data := []byte{}
// 	all_rle := []CharMeDaddy{}
//
// 	for y := range new_img_y {
// 		for x := range new_img_x {
// 			r, g, b, _ := img.At(x*scale_x, y*scale_y).RGBA()
// 			lum := (19595*r + 38470*g + 7471*b + 1<<15) >> 24
// 			indx := lum * uint32(len(frame_chars)) / 256
// 			// sliding window -> 5 bytes
// 			// 0 - char
// 			// 1 - repeat
// 			// 2 - r
// 			// 3 - g
// 			// 4 - b
// 			// 5 - new line
// 			if x == 0 {
// 				all_rle = append(all_rle, CharMeDaddy{frame_chars[indx], 1, uint8(r), uint8(g), uint8(b)})
// 			} else {
// 				curr_rle := &all_rle[len(all_rle)-1]
// 				if frame_chars[indx] == curr_rle.char &&
// 					uint8(r) == curr_rle.r &&
// 					uint8(g) == curr_rle.g &&
// 					uint8(b) == curr_rle.b {
// 					curr_rle.count += 1
// 				} else {
// 					all_rle = append(all_rle, CharMeDaddy{frame_chars[indx], 1, uint8(r), uint8(g), uint8(b)})
// 				}
// 			}
// 		}
// 		for _, el := range all_rle {
// 			encoded_data = append(encoded_data, el.char, el.count, el.r, el.g, el.b)
// 		}
// 		all_rle = []CharMeDaddy{}
// 		encoded_data = append(encoded_data, '\n')
// 	}
// 	return encoded_data
// }

// func decode_frame(enc_data []byte) {
// 	// sliding window -> 5 bytes
// 	// 0 - char
// 	// 1 - repeat
// 	// 2 - r
// 	// 3 - g
// 	// 4 - b
// 	// 5 - new line
// 	fmt.Print("\033[2J\033[H")
// 	for i := 0; i < len(enc_data); i += 5 {
// 		if enc_data[i] == '\n' {
// 			// or i -= 4
// 			i += 1
// 			fmt.Println()
// 			if i >= len(enc_data) {
// 				break
// 			}
// 		}
//
// 		for reps := 0; reps < int(enc_data[i+1]); reps += 1 {
// 			r := enc_data[i+2]
// 			g := enc_data[i+3]
// 			b := enc_data[i+4]
//
// 			var cell string = fmt.Sprintf("\033[38;2;%d;%d;%dm%c\033[0m", r, g, b, enc_data[i])
// 			fmt.Print(cell)
// 		}
//
// 	}
// }

func create_main_window(screen *dd.Screen) {

	if screen.Client.Conn == nil {
		conn, err := net.Dial("tcp", screen.Client.ServerAddr)
		if err != nil {
			if !dd.DEBUG_MODE {
				log.Fatalf("Failed to connect: %v", err)
			}
		}
		screen.Client.Conn = conn
	}

	width, height, _ := term.GetSize(int(os.Stdin.Fd()))
	screen.Width = width
	screen.Height = height

	sidebar := dd.CreateWindow(dd.Styles{
		Width:      50,
		Height:     screen.Height - 1,
		Background: dd.PRIMARY_THEME.PrimaryBg,
		Paddding:   1,
		Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.PRIMARY_THEME.SecondaryTextColor},
	})

	sidebar.AddComponent(
		dd.CreateComponent("Deez nuts123123 hello there", dd.Styles{
			MaxWidth:   10,
			MaxHeight:  5,
			TextColor:  dd.PRIMARY_THEME.SecondaryTextColor,
			Background: dd.PRIMARY_THEME.ActiveBg,
			Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.RED_COLOR},
		}),
	)

	sidebar.Components[0].ScrollType = dd.VERTICAL

	sidebar.AddComponent(
		dd.CreateComponent("Deez nuts123123 hello there", dd.Styles{
			MaxWidth:   20,
			MaxHeight:  3,
			TextColor:  dd.PRIMARY_THEME.SecondaryTextColor,
			Background: dd.PRIMARY_THEME.ActiveBg,
			Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.RED_COLOR},
		}),
	)

	sidebar.Components[1].ScrollType = dd.HORIZONTAL

	sidebar.AddComponent(
		dd.CreateComponent("Deez nuts", dd.Styles{
			MaxWidth:   10,
			TextColor:  dd.PRIMARY_THEME.SecondaryTextColor,
			Background: dd.PRIMARY_THEME.ActiveBg,
			Paddding:   1,
			Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.RED_COLOR},
		},
		),
	)
	screen.AddWindow(sidebar)

	content := dd.CreateWindow(dd.Styles{
		Width:      screen.Width - sidebar.Styles.Width - 1,
		Height:     int(float32(screen.Height) * 0.8),
		Background: dd.PRIMARY_THEME.PrimaryBg,
		Paddding:   1,
		Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.PRIMARY_THEME.SecondaryTextColor},
		Direction:  dd.INLINE,
	})

	content.AddComponent(
		dd.CreateComponent(
			"|SIMD|",
			dd.Styles{
				MaxWidth:   10,
				Background: dd.MakeRGBBackground(80, 40, 100),
				Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
			},
		))
	content.AddComponent(
		dd.CreateComponent(
			"LIGMA???",
			dd.Styles{
				MaxWidth:   10,
				Background: dd.MakeRGBBackground(80, 40, 100),
				Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
			},
		))

	screen.AddWindow(content)

	input_bar := dd.CreateWindow(
		dd.Styles{
			Width:      screen.Width - sidebar.Styles.Width - 1,
			Height:     int(float32(screen.Height)*0.2) - 1,
			Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.PRIMARY_THEME.SecondaryTextColor},
			Background: dd.PRIMARY_THEME.PrimaryBg,
		},
	)

	input_bar.AddComponent(
		dd.CreateComponent(
			"",
			dd.Styles{
				MinWidth:   1,
				MaxWidth:   input_bar.Width - 2,
				Background: dd.MakeRGBBackground(100, 40, 100),
			},
		))

	input_bar.Components[0].Inputable = true
	screen.AddWindow(input_bar)

	create_status_bar(screen)

	screen.Activate()
	screen.RenderFull()
}

func create_status_bar(screen *dd.Screen) {
	screen.StatusBar = dd.Window{
		Position: dd.Position{Row: screen.Height, Col: 1},
		Styles: dd.Styles{
			Width:      screen.Width,
			Height:     1,
			Background: dd.MakeRGBBackground(80, 40, 100),
		},
	}
	screen.StatusBar.Oldfart = screen
	screen.StatusBar.Components = []*dd.Component{
		{
			Parent: &screen.StatusBar,
			Buffer: dd.NORMAL,
			Styles: dd.Styles{
				TextColor: dd.PRIMARY_THEME.ActiveTextColor,
				MaxWidth:  screen.Width,
				MaxHeight: 1,
			},
		},
	}

	status_line := screen.StatusBar.Components[0]
	status_line.Action = func() {
		buffer := strings.Split(status_line.Buffer, ":")
		if len(buffer) < 2 {
			dd.DebugMeDaddy(screen, "Brother this is not a command you dumb fuck")
		}
		switch buffer[1] {
		case "q":
			screen.EventLoopIsRunning = false
			return
		case "new":
			create_new_conversation(screen)
		}

	}
}

func create_new_conversation(screen *dd.Screen) {
	width, height, _ := term.GetSize(int(os.Stdin.Fd()))
	screen.Width = width
	screen.Height = height

	modal := dd.CreateWindow(
		dd.Styles{
			Width:      40,
			Height:     40,
			Paddding:   1,
			Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.PRIMARY_THEME.SecondaryTextColor},
			Background: dd.PRIMARY_THEME.PrimaryBg,
		},
	)
	modal.Position = dd.Position{Row: screen.Height/2 - 20, Col: screen.Width/2 - 20}

	modal.AddComponent(dd.CreateComponent("",
		dd.Styles{
			MinWidth:   10,
			MaxWidth:   modal.Width - 2,
			Background: dd.MakeRGBBackground(100, 40, 100),
			Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.PRIMARY_THEME.SecondaryTextColor},
		},
	))

	err, users := tpc.RequestUsers(screen.Client.Config[1], screen.Client.Conn)
	if err != nil {
		users = []string{"bollocks, cannot retreive users at this time"}
	}

	for _, user := range users {
		modal.AddComponent(dd.CreateComponent(user,
			dd.Styles{
				MaxWidth:   modal.Width - 2,
				Background: dd.MakeRGBBackground(100, 40, 100),
				Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.PRIMARY_THEME.SecondaryTextColor},
			},
		))
	}

	screen.AddWindow(modal)

	create_status_bar(screen)

	screen.Activate()
	screen.RenderFull()

}

func create_login_screen(screen *dd.Screen) {
	width, height, _ := term.GetSize(int(os.Stdin.Fd()))
	screen.Width = width
	screen.Height = height

	login := dd.CreateWindow(
		dd.Styles{
			Width:      40,
			Height:     10,
			Paddding:   1,
			Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.PRIMARY_THEME.SecondaryTextColor},
			Background: dd.PRIMARY_THEME.PrimaryBg,
		},
	)
	login.Position = dd.Position{Row: screen.Height/2 - 5, Col: screen.Width/2 - 20}

	login.AddComponent(
		dd.CreateComponent(
			"",
			dd.Styles{
				MinWidth:   10,
				MaxWidth:   login.Width - 3,
				MaxHeight:  1,
				Background: dd.MakeRGBBackground(80, 40, 100),
			},
		))
	login.AddComponent(
		dd.CreateComponent(
			"",
			dd.Styles{
				MinWidth:   10,
				MaxWidth:   login.Width - 3,
				MaxHeight:  3,
				Background: dd.MakeRGBBackground(80, 40, 100),
				Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
			},
		))
	login.AddComponent(
		dd.CreateComponent(
			"connect",
			dd.Styles{
				MaxWidth:   10,
				Background: dd.MakeRGBBackground(80, 40, 100),
				Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
			},
		))

	screen.AddWindow(login)

	create_status_bar(screen)

	ip := login.Components[0]
	key := login.Components[1]

	ip.Inputable = true
	key.Inputable = true

	login_button := login.Components[2]
	login_button.Action = func() {
		os.Create(".secrets")

		f, err := os.Create(".secrets")
		if err != nil {
			panic(err)
		}
		_, err = f.Write([]byte(ip.Buffer + "\n" + key.Buffer + "\n"))
		if err != nil {
			panic(err)
		}
		screen.Windows = []*dd.Window{}
		create_main_window(screen)
	}

	screen.Activate()
	screen.RenderFull()
}

func enable_tracing() {
	trace_file, err := os.Create("trace.out")
	if err != nil {
		log.Fatalln(err)
	}
	if err := trace.Start(trace_file); err != nil {
		log.Fatalln(err)
	}
}

func enable_mem() {
	runtime.MemProfileRate = 1
	mem_prof_file, err := os.Create("mem.prof")
	if err != nil {
		log.Fatalln(err)
	}
	defer mem_prof_file.Close()
	if err := pprof.WriteHeapProfile(mem_prof_file); err != nil {
		log.Fatalln(err)
	}

	runtime.MemProfileRate = 512
}

func enable_cpu() {
	cpu_prof_file, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatalln(err)
	}
	if err := pprof.StartCPUProfile(cpu_prof_file); err != nil {
		log.Fatalln(err)
	}
}

func enable_dev() {
	dd.DEBUG_MODE = true
}

func main() {

	var cpu_prof_file, trace_file *os.File

	for _, arg := range os.Args {
		switch arg {
		case "trace":
			enable_tracing()
		case "cpu":
			enable_cpu()
		case "dev":
			enable_dev()
		}
	}

	defer func() {
		for _, arg := range os.Args {
			switch arg {
			case "trace":
				defer trace.Stop()
				defer trace_file.Close()
			case "cpu":
				defer cpu_prof_file.Close()
				defer pprof.StopCPUProfile()
			case "mem":
				enable_mem()
			}
		}
	}()
	old_state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error enabling raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), old_state)

	dd.ClearScreen()
	screen := dd.Screen{State: &dd.Normal, EventLoopIsRunning: true}

	screen.Client = tpc.Client{}

	if !screen.Client.LoadClient() {
		create_login_screen(&screen)
	} else {
		create_main_window(&screen)
	}

	stdin_buffer := make([]byte, 1)
	for screen.EventLoopIsRunning {
		if screen.RenderQueue.Len() != 0 {
			screen.Render()
		}

		_, err := os.Stdin.Read(stdin_buffer)
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		screen.State.HandleKeypress(&screen, stdin_buffer)
	}
	// restart to default settings
	fmt.Print(dd.SHOW_CURSOR)
	// TODO: any assert should have show cursor
	//screen.Client.Conn.Close()
}
