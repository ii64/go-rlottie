package main

import (
	"compress/gzip"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"unsafe"

	rlottie "github.com/ii64/go-rlottie"
)

var (
	lottieFile string
)

func usage() {
	fmt.Printf("%s%s%s%s",
		"Usage: \n   go-lottie2gif [lottieFileName] [Resolution] ",
		"[bgColor]\n\nExamples: \n    $ lottie2gif input.json\n   ",
		" $ lottie2gif input.json 200x200\n    $ lottie2gif ",
		"input.json 200x200 ff00ff\n\n",
	)
}

func parseResolution(s string) (w, h uint, ok bool) {
	sp := strings.SplitN(s, "x", 2)
	if len(sp) != 2 {
		ok = false
		return
	}
	var x uint64
	var err error
	if x, err = strconv.ParseUint(sp[0], 10, int(unsafe.Sizeof(w)*8)); err != nil {
		ok = false
		return
	} else {
		w = uint(x)
	}

	if x, err = strconv.ParseUint(sp[0], 10, int(unsafe.Sizeof(h)*8)); err != nil {
		ok = false
		return
	} else {
		h = uint(x)
	}

	ok = true
	return
}

func parseBgColor(s string) bool {
	v, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return false
	}
	setBgColor(uint32(v))
	return true
}

func main() {
	var err error
	var gifFile string
	if len(os.Args) < 2 {
		usage()
		return
	}
	defer func() {
		if err != nil {
			usage()
		}
	}()
	lottieFile = os.Args[1]
	var ow, oh uint
	var odim = false
	if len(os.Args) > 2 {
		if ow, oh, odim = parseResolution(os.Args[2]); !odim {
			usage()
			return
		}
	}
	if len(os.Args) > 3 {
		parseBgColor(os.Args[3])
	}

	//
	var of io.ReadCloser
	if of, err = os.Open(lottieFile); err != nil {
		panic(err)
	}
	defer of.Close()
	if strings.HasSuffix(lottieFile, ".tgs") {
		fmt.Printf("detected telegram sticker :3\n")
		var reader io.Reader
		if reader, err = gzip.NewReader(of); err != nil {
			panic(err)
		}
		of = ioutil.NopCloser(reader)
		defer of.Close()
	}

	var by []byte
	if by, err = io.ReadAll(of); err != nil {
		panic(err)
	}

	var (
		lottieData    = string(by)
		lottieKey     = ""
		lottieResPath = ""
	)

	gifFile = lottieFile + ".gif"
	var ins rlottie.Animation
	ins, err = rlottie.AnimationFromData(lottieData, lottieKey, lottieResPath)
	if err != nil {
		panic(err)
	}
	defer ins.Destroy()

	d := ins.GetSize()
	fmt.Printf("actual dim: %+#v\n", d)
	if odim {
		d.Width = ow
		d.Height = oh
	}
	fmt.Printf("out surface w:%d h:%d\n", d.Width, d.Height)

	dur := ins.GetDuration()
	fmt.Printf("dur: %v\n", dur)

	totalframe := ins.GetTotalFrame()
	fmt.Printf("totalframe: %v\n", totalframe)

	framerate := ins.GetFrameRate()
	fmt.Printf("framerate: %v\n", framerate)
	if framerate > 30 {
		fmt.Printf("framerate: framerate more than 30\n")
	}

	// duration in second * 100
	frameDuration := int(float64(dur*100) / (float64(totalframe)))
	fmt.Printf("frameDuration: %v\n", frameDuration)

	fmt.Printf("creating gif...\n")

	images := []*image.Paletted{}
	var delays []int
	var i uint

	var buffer = make([]uint32, d.Width*d.Height)
	var skipper uint = uint(float64(totalframe) / framerate)
	_ = skipper // frame skipper to support image/gif ?

	for i = 0; i < ins.GetTotalFrame(); {
		buffs := ins.Render(buffer, i, d.Width, d.Height, d.Width*4)

		pallets := argbToRgba(buffs, d.Height*d.Width*4)
		img := image.NewPaletted(
			image.Rect(
				0, 0,
				int(d.Width), int(d.Height),
			),
			pallets)
		setPixel(buffs, img, d.Width, d.Height)

		images = append(images, img)
		delays = append(delays, frameDuration)

		i = i + 1
	}

	var f *os.File
	f, err = os.OpenFile(gifFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
		Config: image.Config{
			ColorModel: nil,
			Width:      int(d.Width),
			Height:     int(d.Height),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("ok\n")
}

var (
	bgColor                      = setBgColor(0xffffffff)
	bgColorR, bgColorG, bgColorB uint8
)

func setBgColor(bgColor uint32) uint32 {
	bgColorR = (uint8)((bgColor & 0xff0000) >> 16)
	bgColorG = (uint8)((bgColor & 0x00ff00) >> 8)
	bgColorB = (uint8)((bgColor & 0x0000ff))
	return bgColor
}

func u8to32(r, g, b, a uint8) uint32 {
	return (((((uint32(r) << 8) ^ uint32(g)) << 8) ^ uint32(b)) << 8) ^ uint32(a)
}

type xcolor struct {
	s uint32
	c color.RGBA
}
type acolor []xcolor

func (a acolor) Len() int {
	return len(a)
}

func (a acolor) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a acolor) Less(i, j int) bool {
	return a[i].s < a[j].s
}

type fcolor struct {
	freq int
	t    *xcolor
}
type hfcolor map[uint32]fcolor
type afcolor []fcolor

func (a afcolor) Len() int {
	return len(a)
}

func (a afcolor) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a afcolor) Less(i, j int) bool {
	return a[i].freq > a[j].freq
}

func setPixel(buff []uint8, img *image.Paletted, width, height uint) {
	var y, x uint
	pos := 0
	for y = 0; y < height; y++ {
		for x = 0; x < width; x++ {
			c := color.RGBA{
				R: buff[pos],
				G: buff[pos+1],
				B: buff[pos+2],
				// A: buff[pos+3], //interlace problem
				A: 255,
			}
			img.Set(int(x), int(y), c)
			pos = pos + 4
		}
	}
}
func argbToRgba(buff []uint8, totalbytes uint) color.Palette {
	var cpallete = acolor{}
	var cfreq = hfcolor{}

	var i uint
	for i = 0; i < totalbytes; i = i + 4 {
		a := buff[i+3]
		if a != 0 {
			var (
				r = buff[i+2]
				g = buff[i+1]
				b = buff[i]
			)
			if a != 255 {
				var (
					r2 = (uint8)((float64)(bgColorR) * ((float64)(255-a) / 255))
					g2 = (uint8)((float64)(bgColorG) * ((float64)(255-a) / 255))
					b2 = (uint8)((float64)(bgColorB) * ((float64)(255-a) / 255))
				)
				buff[i] = r + r2
				buff[i+1] = g + g2
				buff[i+2] = b + b2
			} else {
				// only sizzle r and b
				buff[i] = r
				buff[i+2] = b
			}
		} else {
			buff[i+2] = bgColorB
			buff[i+1] = bgColorG
			buff[i] = bgColorR
		}

		var (
			r = buff[i]
			g = buff[i+1]
			b = buff[i+2]
			// a = buff[i+3]
		)
		c := color.RGBA{r, g, b, a}
		ck := u8to32(r, g, b, a)
		if i, ok := cfreq[ck]; !ok {
			xc := xcolor{
				s: ck, c: c,
			}
			cfreq[ck] = fcolor{
				freq: 1,
				t:    &xc,
			}
			cpallete = append(cpallete, xc)
		} else {
			i.freq = i.freq + 1
			cfreq[ck] = i
		}
	}
	sort.Sort(cpallete)
	return cfreq.GetPalette()
}

func (f hfcolor) GetPalette() color.Palette { // get most palettes colour
	var carrr = afcolor{}
	for _, v := range f {
		carrr = append(carrr, v)
	}
	sort.Sort(carrr)
	cp := color.Palette{}
	for i, c := range carrr[:] {
		if i >= 256 {
			break
		}
		cp = append(cp, c.t.c)
	}
	return cp
}
