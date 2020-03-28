 package main

import (
	"os"
	"regexp"
	"fmt"
	"strings"
	"strconv"
	"log"
	"time"

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
	watch chan time.Time
	apiKey = "CC238ZlLq4J0m-JTvrKBlmx5XNA"
	re = regexp.MustCompile("[A-Z]+:")
	re2 = regexp.MustCompile("\\.\\!\\?")
)

var session = cleverbot.New(apiKey)


func main() {
	c = bot.NewClient()

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


func bed(msg string) error {
	log.Println("Bed requested")
	x,y,z := c.Player.GetBlockPos()
	success := 0
	// look for a bed nearby
	var xb, yb, zb int
	for xb = x-3 ; xb < x+3; xb++ {
		for yb = y-3 ; yb < y+3; yb++ {
			for zb = z-3 ; zb < z+3; zb++ {
				block := c.Wd.GetBlock(xb,yb,zb)
				if block.String() == "minecraft:white_bed" {
					log.Println(fmt.Sprintf("Bed found at %d,%d,%d\n",xb,yb,zb))
					success = 1
					break
				}
			}
			if success == 1 {
				break
			}
		}
		if success == 1 {
			break
		}
	}
	if success == 0 {
		log.Println("No bed found.")
		leave()
	} else {
		err := c.UseBlock(0,xb,yb,zb,1,0.5,1,0.5,false)
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
		c.Chat("In bed")
	}
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
	success := 0
	var xb, yb, zb int
	for xb = x-30 ; xb < x+30; xb++ {
		for yb = Max(1,y-30) ; yb < Min(255,y+30); yb++ {
			for zb = z-30 ; zb < z+30; zb++ {
				block := c.Wd.GetBlock(xb,yb,zb)
				if block.String() == blockname {
					c.Chat(fmt.Sprintf("%s found at %d,%d,%d\n",blockname,xb,yb,zb))
					success = 1
					break
				}
			}
			if success == 1 {
				break
			}
		}
		if success == 1 {
			break
		}
	}
	if success == 0 {
		c.Chat(fmt.Sprintf("%s not found.",blockname))
	}
	return nil
}

func moveShip(mspl []string) error {
	xold,yold,zold := c.Player.GetBlockPos()

	x, errx := strconv.Atoi(mspl[2])
	y, erry := strconv.Atoi(mspl[3])
	z, errz := strconv.Atoi(mspl[4])
	if errx != nil || erry != nil ||errz != nil {
		c.Chat("invalid coordinate.")
		return nil
	}

	xnew, errx := strconv.Atoi(mspl[5])
	ynew, erry := strconv.Atoi(mspl[6])
	znew, errz := strconv.Atoi(mspl[7])
	if errx != nil || erry != nil ||errz != nil {
		c.Chat("invalid coordinate.")
		return nil
	}

	c.Chat("Ship warp requested. Materializing on your bridge now, captain.")
	c.Chat("/kill Telleilogical")
	time.Sleep(2 * time.Second)
	c.Chat(fmt.Sprintf("/teleport Telleilogical %d %d %d",x,y,z))
	c.Chat("Hello, captain. I require a crystaline structure to align the phase. One diamond, please.")
	time.Sleep(5 * time.Second)

	// Check if given a diamond.
	// If not, leave.
	if false {
		c.Chat("I cannot warp the ship without a diamond.")
		c.Chat(fmt.Sprintf("/teleport Telleilogical %d %d %d", xold,yold,zold))
		return nil
	}
	// If so, continue.
	c.Chat("Thank you, captain.")

	var xl,yl,zl,xu,yu,zu int
	var xs,ys,zs int
	// first we find a gold block.
	var xb, yb, zb int
	success := 0
	for xb = x-30 ; xb < x+30; xb++ {
		for yb = Max(1,y-30) ; yb < Min(255,y+30); yb++ {
			for zb = z-30 ; zb < z+30; zb++ {
				block := c.Wd.GetBlock(xb,yb,zb)
				if block.String() == "minecraft:gold_block" {
					success = 1
					xs = xb
					ys = yb
					zs = zb
					break
				}
			}
			if success == 1 {
				break
			}
		}
		if success == 1 {
			break
		}
	}
	if success == 0 {
		c.Chat(fmt.Sprintf("Please move your captain closer to a ship corner."))
		return nil
	} else {
		c.Chat(fmt.Sprintf("Found seed block at %d,%d,%d.",xs, ys, zs))
	}

	// Now we seek the other gold blocks.
	success = 0
	d := 0
	for success == 0 && d < 100 {
		d = d+1
		block := c.Wd.GetBlock(xs-d,ys,zs)
		if block.String() == "minecraft:gold_block" {
			success = 1
			xl = xs-d
			xu = xs
			c.Chat(fmt.Sprintf("Found new x block at %d,%d,%d.",xl, ys, zs))

		}
		block = c.Wd.GetBlock(xs+d,ys,zs)
		if block.String() == "minecraft:gold_block" {
			success = 1
			xl = xs
			xu = xs+d
			c.Chat(fmt.Sprintf("Found new x block at %d,%d,%d.",xu, ys, zs))
		}
	}
	if success == 0 {
		c.Chat("your ship is too big in x.")
		c.Chat(fmt.Sprintf("seed block: %d,%d,%d. Min: %d,%d,%d. Max: %d,%d,%d.",xs,ys,zs,xl,yl,zl,xu,yu,zu))
		return nil
	}
	success = 0
	d = 0
	for success == 0 && d < 100 {
		d = d+1
		block := c.Wd.GetBlock(xs,Max(ys-d,1),zs)
		if block.String() == "minecraft:gold_block" {
			success = 1
			yl = Max(ys-d,1)
			yu = ys
			c.Chat(fmt.Sprintf("Found new y block at %d,%d,%d.",xs, yl, zs))
		}
		block = c.Wd.GetBlock(xs,Min(ys+d,255),zs)
		if block.String() == "minecraft:gold_block" {
			success = 1
			yl = ys
			yu = Min(ys+d,255)
			c.Chat(fmt.Sprintf("Found new y block at %d,%d,%d.",xs, yu, zs))
		}
	}
	if success == 0 {
		c.Chat("your ship is too big in y.")
		c.Chat(fmt.Sprintf("seed block: %d,%d,%d. Min: %d,%d,%d. Max: %d,%d,%d.",xs,ys,zs,xl,yl,zl,xu,yu,zu))
		return nil
	}
	success = 0
	d = 0
	for success == 0 && d < 100 {
		d = d+1
		block := c.Wd.GetBlock(xs,ys,zs-d)
		c.Chat(fmt.Sprintf("%d,%d,%d: %s",xs,ys,zs-d,block.String()))
		if block.String() == "minecraft:gold_block" {
			success = 1
			zl = zs-d
			zu = zs
			c.Chat(fmt.Sprintf("Found new z block at %d,%d,%d.",xs, ys, zl))
		}
		block = c.Wd.GetBlock(xs,ys,zs+d)
		c.Chat(fmt.Sprintf("%d,%d,%d: %s",xs,ys,zs+d,block.String()))
		if block.String() == "minecraft:gold_block" {
			success = 1
			zl = zs
			zu = zs+d
			c.Chat(fmt.Sprintf("Found new z block at %d,%d,%d.",xs, ys, zu))
		}
	}
	if success == 0 {
		c.Chat("your ship is too big in z.")
		c.Chat(fmt.Sprintf("seed block: %d,%d,%d. Min: %d,%d,%d. Max: %d,%d,%d.",xs,ys,zs,xl,yl,zl,xu,yu,zu))
		return nil
	}
	c.Chat("I've computed the boundary of your ship. Warping now...")

	// Begin warp.
	xdest := xl + (xnew-x)
	ydest := Max( 1, Min(yl + (ynew-y), 255-(yu-yl)) )
	zdest := zl + (znew-z)

	c.Chat(fmt.Sprintf("/clone %d %d %d %d %d %d %d %d %d replace move",xl,yl,zl,xu,yu,zu,xdest,ydest,zdest))
	c.Chat(fmt.Sprintf("/teleport Telleilogical %d %d %d",xnew,ynew,znew))
	c.Chat("Warp complete, captain.")
	time.Sleep(1 * time.Second)
	c.Chat("Returning to base now.")
	c.Chat(fmt.Sprintf("/teleport Telleilogical %d %d %d", xold, yold, zold))
	return nil
}

func onChatMsg(cm chat.Message, pos byte) error {
	log.Println("Chat:", cm)
	spl := strings.Split(cm.String(), "> ")
	if len(spl) <= 1 {
		return nil
	}

	msg := spl[1]
	if len(msg) > 2 && strings.ToLower(msg[:3]) == "bed" {
		err := bed(msg)
		if err != nil {
			log.Fatal(err)
		}
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
		} else if len(pmsg) > 5 && strings.ToLower(pmsg[:5]) == "drive" {
			err := moveShip(mspl)
			if err != nil {
				log.Fatal(err)
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

	return nil
}

func onDisconnect(c chat.Message) error {
	log.Println("Disconnect:", c)
	return nil
}

