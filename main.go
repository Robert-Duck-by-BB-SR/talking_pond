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

func debug_sidebar(sidebar *dd.Window) {
	sidebar.AddComponent(
		dd.CreateComponent("Deez nuts 123 456 789 100 110", dd.Styles{
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
}

func debug_content(content *dd.Window) {
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
}

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

	if !dd.DEBUG_MODE {
		// err, _ := tpc.RequestToConnect(&screen.Client)
		// if err != nil {
		// 	log.Println("boiii you're not allowed here")
		// 	screen.EventLoopIsRunning = false
		// 	return
		// }
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

	if dd.DEBUG_MODE {
		debug_sidebar(sidebar)
	} else {
		// sidebar.OnRender = func() {
		// NOTE: quick fix to not add more components on rerender
		// sidebar.Components = []*dd.Component{}
		// err, conversations := tpc.RequestConversations(&screen.Client)
		// if err != nil {
		// 	log.Println("can't fetch conversations")
		// 	screen.EventLoopIsRunning = false
		// 	return
		// }
		// for _, con := range conversations {
		// 	data := strings.Split(con, string([]byte{255}))
		//
		// 	//FIXME: there is something funky happening here
		// 	if len(data) == 2 {
		// 		chat := dd.CreateComponent(
		// 			fmt.Sprint(data[1], "|", data[0]),
		// 			dd.Styles{
		// 				MaxWidth:   sidebar.Width - 4,
		// 				MaxHeight:  5,
		// 				TextColor:  dd.PRIMARY_THEME.SecondaryTextColor,
		// 				Background: dd.PRIMARY_THEME.ActiveBg,
		// 				Paddding:   1,
		// 				Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.RED_COLOR},
		// 			},
		// 		)
		// 		sidebar.AddComponent(chat)
		// 		chat.Action = func() {
		// 			screen.Client.Conversation = data[0]
		// 			// TODO: should not ignore error
		// 			_, messages := tpc.RequestMessages(&screen.Client)
		// 			content := screen.Windows[1]
		// 			dd.CreateMessages(content, data[0], messages)
		// 			screen.Activate(1)
		// 		}
		// 	}
		// 	}
		// }
	}

	screen.AddWindow(sidebar)

	content := dd.CreateWindow(dd.Styles{
		Width:      screen.Width - sidebar.Styles.Width - 1,
		Height:     int(float32(screen.Height) * 0.8),
		Background: dd.PRIMARY_THEME.PrimaryBg,
		Paddding:   1,
		Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.PRIMARY_THEME.SecondaryTextColor},
		Direction:  dd.INLINE,
	})

	if dd.DEBUG_MODE {
		debug_content(content)
	}

	screen.AddWindow(content)

	input_bar := dd.CreateWindow(
		dd.Styles{
			Width:      screen.Width - sidebar.Styles.Width - 1,
			Height:     int(float32(screen.Height)*0.2) - 1,
			Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.PRIMARY_THEME.SecondaryTextColor},
			Background: dd.PRIMARY_THEME.PrimaryBg,
		},
	)

	input := dd.CreateComponent(
		"",
		dd.Styles{
			MinWidth:   1,
			Width:      input_bar.Width - 2,
			Background: dd.MakeRGBBackground(200, 40, 100),
			Height:     input_bar.Height - 2,
		},
	)

	input_bar.AddComponent(input)
	input.Inputable = true
	input.ScrollType = dd.VERTICAL
	input.Action = func() {
		// if len(input.Buffer) != 0 {
		// 	tpc.SendMessage(&screen.Client, input.Buffer)
		// }
		// input.Buffer = ""
		// // maybe it should be render_content
		// input.Render(&screen.RenderQueue)
	}
	screen.AddWindow(input_bar)

	create_status_bar(screen)

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
			Parent:    &screen.StatusBar,
			Buffer:    dd.NORMAL,
			Inputable: true,
			Styles: dd.Styles{
				TextColor: dd.PRIMARY_THEME.ActiveTextColor,
				Width:     screen.Width,
				Height:    1,
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

	if !dd.DEBUG_MODE {
		// err, users := tpc.RequestUsers(screen.Client.Config[1], screen.Client.Conn)
		// if err != nil {
		// 	users = []string{"bollocks, cannot retreive users at this time"}
		// }
		//
		// for _, user := range users {
		// 	available_user := dd.CreateComponent(user,
		// 		dd.Styles{
		// 			MaxWidth:   modal.Width - 2,
		// 			Background: dd.MakeRGBBackground(100, 40, 100),
		// 			Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.PRIMARY_THEME.SecondaryTextColor},
		// 		},
		// 	)
		//
		// 	available_user.Action = func() {
		// 		tpc.CreateConversation(screen.Client.Config[1], user, screen.Client.Conn)
		// 		screen.CloseModal()
		// 	}
		//
		// 	modal.AddComponent(available_user)
		// }

	}
	screen.AddWindow(modal)

	create_status_bar(screen)

	screen.ActivateModal()
	screen.Render()
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
	go screen.RenderQueueStart()

	if !screen.Client.LoadClient() {
		create_login_screen(&screen)
	} else {
		create_main_window(&screen)
	}

	stdin_buffer := make([]byte, 1)
	for screen.EventLoopIsRunning {
		screen.Render()

		_, err := os.Stdin.Read(stdin_buffer)
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		screen.State.HandleKeypress(&screen, stdin_buffer)
	}
	// restart to default settings
	fmt.Print(dd.VISIBLE_CURSOR)
	// TODO: any assert should have show cursor
	if screen.Client.Conn != nil {
		screen.Client.Conn.Close()
	}
}
