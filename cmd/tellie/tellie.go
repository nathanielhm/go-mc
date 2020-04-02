 package main

import (
	"os"
	"regexp"
	"fmt"
	"strings"
	"strconv"
	"log"
	"time"
	"math"

	"github.com/Tnze/go-mc/yggdrasil"
	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/realms"
	_ "github.com/Tnze/go-mc/data/lang/en-us"
	"github.com/ugjka/cleverbot-go"
)

const timeout = 45
const username string = "lanelawley@gmail.com"
const password string = "MbuRobots2"
const realm_name string = "Strawberry City"

var (
	r     *realms.Realms
	c     *bot.Client
	realm_address = ""
	realm_port = 0

	warping = false
	xbase,ybase,zbase int
	ship_xl,ship_yl,ship_zl,ship_xu,ship_yu,ship_zu int
	ship_x,ship_y,ship_z int
	crew []string

	watch chan time.Time
	apiKey = "CC238ZlLq4J0m-JTvrKBlmx5XNA"
	re = regexp.MustCompile("[A-Z]+:")
	re2 = regexp.MustCompile("\\.\\!\\?")
)

var session = cleverbot.New(apiKey)


func main() {
	c = bot.NewClient()
	xbase = 90
	ybase = 66
	zbase = -247

	// log in
	auth, err := yggdrasil.Authenticate(username,password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c.Auth.UUID, c.Name = auth.SelectedProfile()
	c.AsTk = auth.AccessToken()

	fmt.Println("user:", c.Name)
	fmt.Println("uuid:", c.Auth.UUID)
	fmt.Println("astk:", c.AsTk)

	// parse realms
	r = realms.New("1.14.2", c.Name, c.AsTk, c.Auth.UUID)
	servers,err := r.Worlds()

	if err != nil {
		panic(err)
	}

	for _,v := range servers {
		if v.Name == realm_name {
			fmt.Println("Found Realm", realm_name)
			fmt.Printf("v is %s\n", v)
			address, err := r.Address(v)
			if err != nil {
				panic(err)
			}
			rholder := strings.SplitN(address,":",2)
			realm_address = rholder[0]
			realm_port,err = strconv.Atoi(rholder[1])
			fmt.Println(realm_address, realm_port)
		}
	}
	if realm_address == "" {
		panic("Realm not found!")
	}

	// join server
	if err := c.JoinServer(realm_address, realm_port); err != nil {
		log.Fatal(err)
	}
	log.Println("Login success")

	//Register event handlers
	c.Events.GameStart = onGameStart
	c.Events.ChatMsg = onChatMsg
	c.Events.Disconnect = onDisconnect
	c.Events.SoundPlay = onSound
	c.Events.Die = onDeath

	//JoinGame
	err = c.HandleGame()
	if err != nil {
		log.Fatal(err)
	}
}
func onDeath() error {
	log.Println("Death")

	c.Chat("Respawning...")
	c.Respawn()

	if warping == false {
		c.Chat(fmt.Sprintf("/teleport Telleilogical %d %d %d",xbase,ybase,zbase))
	}
	return nil
}

func onGameStart() error {
	log.Println("Game start")

	c.Chat("hello")

	watch = make(chan time.Time)
	return nil
}

func onSound(name string, category int, x, y, z float64, volume, pitch float32) error {
	return nil
}

func leave() int {
	// Sign out
	err := yggdrasil.SignOut(username, password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
	return 0
}


func Max(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}
func Min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}


func bed(xb,yb,zb int) error {
	log.Println("Bed requested")
	// look for a bed nearby
	err := c.UseBlock(0,xb,yb,zb,1,0.5,1,0.5,false)
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	c.Chat("In bed")
	return nil
}

func find(mspl []string) error {
	var blockname string
	if len(mspl) < 5 {
		c.Chat("I can only find stuff if you format > tellie find [block] [x] [y] [z].")
		return nil
	}
	if strings.ToLower(mspl[2]) == "diamonds" {
		blockname = "minecraft:diamond_ore"
	} else {
		blockname = strings.ToLower(mspl[2])
	}
	x, errx := strconv.Atoi(mspl[3])
	y, erry := strconv.Atoi(mspl[4])
	z, errz := strconv.Atoi(mspl[5])
	if errx != nil || erry != nil ||errz != nil {
		c.Chat("invalid coordinate.")
		return nil
	}
	bfx := []int{}
	bfy := []int{}
	bfz := []int{}
	var xb, yb, zb int
	for xb = x-30 ; xb < x+30; xb++ {
		for yb = Max(1,y-30) ; yb < Min(255,y+30); yb++ {
			for zb = z-30 ; zb < z+30; zb++ {
				block := c.Wd.GetBlock(xb,yb,zb)
				if block.String() == blockname {
					bfx = append(bfx, xb)
					bfy = append(bfy, yb)
					bfz = append(bfz, zb)
				}
			}
		}
	}
	if len(bfx) == 0 {
		c.Chat(fmt.Sprintf("%s not found.",blockname))
	} else {
		minj := 0
		mindist := bfx[0]-x + bfy[0]-y + bfz[0]-z
		for j := 0; j < len(bfx); j++ {
			if bfx[j]-x + bfy[j]-y + bfz[j]-z < mindist {
				mindist = bfx[j]-x + bfy[j]-y + bfz[j]-z
				minj = j
			}
		}
		c.Chat(fmt.Sprintf("%d instances of %s found. The nearest is #%d at %d,%d,%d\n",len(bfx),blockname,minj,bfx[minj],bfy[minj],bfz[minj]))

	}
	return nil
}

func learnShip(spl []string) error {
	c.Chat("Ship warp requested. Materializing on your bridge now, captain.")

	xf,yf,zf := c.Player.GetPosition()
	x := int(math.Floor(xf))
	y := int(math.Floor(yf))
	z := int(math.Floor(zf))
	xl,errxl := strconv.Atoi(spl[1])
	yl,erryl := strconv.Atoi(spl[2])
	zl,errzl := strconv.Atoi(spl[3])
	xu,errxu := strconv.Atoi(spl[4])
	yu,erryu := strconv.Atoi(spl[5])
	zu,errzu := strconv.Atoi(spl[6])

	if errxl != nil || erryl != nil || errzl != nil {
		c.Chat("There's something wrong with your ship coding.")
		return nil
	}
	if errxu != nil || erryu != nil || errzu != nil {
		c.Chat("There's something wrong with your ship coding.")
		return nil
	}
	ship_x = x
	ship_y = y
	ship_z = z
	ship_xl = x + xl
	ship_yl = y + yl
	ship_zl = z + zl
	ship_xu = x + xu
	ship_yu = y + yu
	ship_zu = z + zu

	c.Chat("Captain, I've familiarized myself with your ship. You can warp whenever you're ready.")
	c.Chat(fmt.Sprintf("/teleport Telleilogical %d %d %d",xbase,ybase,zbase))

	return nil
}
func moveShip(mspl []string, captain string) error {
	if warping == false {
		c.Chat("I'm unfamiliar with any ships right now, I'm afraid.")
		return nil
	}

	xnew, errx := strconv.Atoi(mspl[2])
	ynew, erry := strconv.Atoi(mspl[3])
	znew, errz := strconv.Atoi(mspl[4])
	if errx != nil || erry != nil ||errz != nil {
		c.Chat("invalid coordinate.")
		return nil
	}
	c.Chat(fmt.Sprintf("/teleport Telleilogical %d %d %d",ship_x,ship_y,ship_z))

	x := ship_x
	y := ship_y
	z := ship_z
	xl := ship_xl
	yl := ship_yl
	zl := ship_zl
	xu := ship_xu
	yu := ship_yu
	zu := ship_zu

	c.Chat(fmt.Sprintf("/tell %s Very well. I've familiarized myself with the boundaries of your ship and we're ready to warp.",captain))
	c.Chat(fmt.Sprintf("/tell %s I require a crystaline structure to align the phase. One diamond, please.",captain))

	time.Sleep(5 * time.Second)

	// Check if given a diamond.
	// If not, leave.
	if false {
		c.Chat(fmt.Sprintf("/tell %s I cannot warp the ship without a diamond.",captain))
		c.Chat(fmt.Sprintf("/teleport Telleilogical %d %d %d", xbase, ybase, zbase))
		return nil
	}
	// If so, continue.
	c.Chat("Thank you, captain.")

	c.Chat("Warping now...")

	// Begin warp.

	xdest := xl + (xnew-x)
	ydest := Max( 1, Min(yl + (ynew-y), 255-(yu-yl)) )
	zdest := zl + (znew-z)

	c.Chat(fmt.Sprintf("/forceload add %d %d %d %d", xdest-(xu-xl), zdest-(zu-zl), xdest, zdest))

	c.Chat(fmt.Sprintf("/clone %d %d %d %d %d %d %d %d %d replace move",xl,yl,zl,xu,yu,zu,xdest,ydest,zdest))
//	c.Chat(fmt.Sprintf("/teleport %s %d %d %d",captain, xnew,ynew+1,znew))
	// c.Chat(fmt.Sprintf("/teleport %s %d %d %d","scefing", xnew,ynew+1,znew))
	// c.Chat(fmt.Sprintf("/teleport %s %d %d %d","CowSnail", xnew,ynew+1,znew))
	for _, member := range crew {
		c.Chat(fmt.Sprintf("/teleport %s %d %d %d", member, xnew, ynew+1, znew))
	}
	c.Chat(fmt.Sprintf("/teleport Telleilogical %d %d %d",xnew,ynew+1,znew))
	c.Chat(fmt.Sprintf("/tell %s Warp complete, captain.",captain))
	time.Sleep(1 * time.Second)
	c.Chat(fmt.Sprintf("/tell %s Returning to base now.",captain))
	c.Chat(fmt.Sprintf("/teleport Telleilogical %d %d %d", xbase, ybase, zbase))
	warping = false
	c.Chat(fmt.Sprintf("/forceload remove %d %d %d %d", xdest-(xu-xl), zdest-(zu-zl), xdest, zdest))

	c.Chat(fmt.Sprintf("x: %d, y: %d, z: %d", x,y,z))
	c.Chat(fmt.Sprintf("xnew: %d, ynew: %d, znew: %d", xnew,ynew,znew))
	c.Chat(fmt.Sprintf("xl: %d, yl: %d, zl: %d", xl,yl,zl))
	c.Chat(fmt.Sprintf("xu: %d, yu: %d, zu: %d", xu,yu,zu))
	c.Chat(fmt.Sprintf("xdest: %d, ydest: %d, zdest: %d", xdest,ydest,zdest))
	return nil
}

func onChatMsg(cm chat.Message, pos byte) error {
	log.Println("Chat:", cm)

	cmstr := cm.String()
	if false == true {
		// this is just here for now.
	} else {
		// it's a standard message.
		var spl, spl2 []string
		if cmstr[0] == '[' {
			spl = strings.Split(cmstr, "] ")
			spl2 = strings.Split(spl[0],"[")
		} else if cmstr[0] == '<' {
			spl = strings.Split(cmstr, "> ")
			spl2 = strings.Split(spl[0],"<")
		} else {
			return nil
		}
		if len(spl) <= 1 {
			return nil
		}

		msg := spl[1]
		requester := spl2[1]
		if len(msg) > 2 && strings.ToLower(msg[:3]) == "bed" {
			err := bed(89,66,-249)
			if err != nil {
				log.Fatal(err)
			}
		} else if len(msg) > 6 && msg[:6] == "BSSSRT" {
			warping = true
			mspl := strings.Split(msg, " ")
			learnShip(mspl)
			return nil
		} else if len(msg) > 6 && strings.ToLower(msg[:6]) == "tellie" {
			mspl := strings.Split(msg, " ")
			pmsg := msg
			if len(mspl) > 1 {
				pmsg = strings.Join(mspl[1:], " ")
			}

			if pmsg == "leave" {
				log.Println("Requested to leave")
				leave()
			} else if pmsg == "see" {
				x,y,z := c.Player.GetBlockPos()
				block := c.Wd.GetBlock(x,y-1,z)
				c.Chat(block.String())
			} else if len(pmsg) > 6 && strings.ToLower(pmsg[:6]) == "select" {
				j, err := strconv.Atoi(mspl[2])
				if err != nil {
					c.Chat("I don't understand that slot.")
					return nil
				} else if j > 8 || j < 0 {
					c.Chat("That slot isn't valid.")
				}
				c.SelectItem(j)
			} else if pmsg == "what are you holding" {
				c.Chat(fmt.Sprintf("%d", c.Player.HeldItem))
			} else if len(pmsg)>4 && strings.ToLower(pmsg[:4]) == "find" {
				err := find(mspl)
				if err != nil {
					log.Fatal(err)
				}
			} else if len(pmsg) > 4 && strings.ToLower(pmsg[:4]) == "warp" {
				err := moveShip(mspl,requester)
				if err != nil {
					log.Fatal(err)
				}
			} else if len(pmsg) > 4 && strings.ToLower(pmsg[:4]) == "crew" {
				for _, member := range crew {
					c.Chat(member)
				}
				cspl := strings.Split(pmsg, " ")
				c.Chat("Warp crew set to:")
				crew = cspl[1:]
				for _, member := range crew {
					c.Chat(fmt.Sprintf("   %s", member))
				}
			} else {
				resp, err := session.Ask(pmsg)
				if err != nil {
					fmt.Printf("Cleverbot error: %v\n", err)
				} else {
					c.Chat(resp)
				}
				/*
				inp := fmt.Sprintf("MAN: %s WOMAN: ", pmsg)
				out, err := exec.Command("/bin/bash", "./cmd.sh", inp).Output()
				if err != nil {
					// log.Fatal(err)
					fmt.Printf("GPT2 error: %v\n", err)
				}
				proc := re.Split(string(out), -1)
				tellieResp := strings.Split(strings.Trim(proc[2], " 	\n"), "\n")[0]
				proc2 := re2.Split(tellieResp, -1)
				if len(proc2) > 1 {
					tellieResp = strings.Join(proc2[:len(proc2)-1], " ")
				}
				c.Chat(tellieResp)
				*/
			}
		}
	}

	return nil
}

func onDisconnect(c chat.Message) error {
	log.Println("Disconnect:", c)
	return nil
}

