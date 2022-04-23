# uggo
helper utilities for uggly ecosystem

```
package uggo // import "github.com/rendicott/uggo"


    Package uggo provides helper utilities and page templates
for building uggly servers in golang.

FUNCTIONS

func AddFormLogin(presp *pb.PageResponse, keyStroke, submitPage string) *pb.PageResponse
    AddFormLogin takes a page and adds a login div, form, helptext, and
    activation link to it it starts the formDiv a few lines below the last
    divbox on the page

func AddFormNewUser(presp *pb.PageResponse, keyStroke, submitPage string) *pb.PageResponse
    AddFormNewUser takes a page and adds a new user div, form, helptext, and
    activation link to it it starts the formDiv a few lines below the last
    divbox on the page

func AddLink(presp *pb.PageResponse, keyStroke, pageName string, stream bool) *pb.PageResponse
    AddLink adds a link to the given destination on the provided page using the
    provided keystroke with an indication as to whether or not it's a link to a
    stream

func AddTextBoxToPage(pr *pb.PageResponse, text string) *pb.PageResponse
    AddTextBoxToPage adds a textbox and puts the given text in it

func ConvertStringCharRune(s string) int32
    ConvertStringCharRune takes a string and converts it to a rune slice then
    grabs the rune at index 0 in the slice so that it can return an int32 to
    satisfy the Uggly protobuf struct for border and fill chars and such. If the
    input string is less than zero length then it will just rune out a space
    char and return that int32.

func FindLastBoxEndX(presp *pb.PageResponse) int32
    FindLastBoxEndX takes a page and finds its last div box then returns the
    extreme righthand coordinate as an int32. If it encounters any errors it
    just returns 0

func FindLastBoxEndY(presp *pb.PageResponse) int32
    FindLastBoxEndY takes a page and finds its last div box then returns the
    extreme lower coordinate as an int32. If it encounters any errors it just
    returns 0

func FindLastBoxStartX(presp *pb.PageResponse) int32
    FindLastBoxStartX takes a page and finds its last div box then returns the
    extreme lefthand coordinate as an int32. If it encounters any errors it just
    returns 0

func FindLastBoxStartY(presp *pb.PageResponse) int32
    FindLastBoxStartY takes a page and finds its last div box then returns the
    extreme top coordinate as an int32. If it encounters any errors it just
    returns 0

func GenPageLittleBox(x, y int) *pb.PageResponse
    GenPageLittleBox generates a new page with a small box at the given
    coordinates

func GenPageSimple(w, h int, text string) *pb.PageResponse
    GenPageSimple generates a page that fits within the given w, h and adds the
    provided text

func GrowBox(pr *pb.PageResponse, boxName string, width, height int) *pb.PageResponse
    GrowBox takes a page and a box name then locates it and modifies it's size

func MoveBox(pr *pb.PageResponse, boxName string, x, y int) *pb.PageResponse
    GrowBox takes a page and a box name then locates it and modifies it's size

func PageTopMenuFullWidthContent(
	w, h int,
	links []*PageLink,
	current, content string) *pb.PageResponse
func Style(fg, bg string) *pb.Style
    Style takes a foreground and background color string and makes assumptions
    about the attr style to return a uggly Style struct. Easier to use than
    building struct pointers for common use cases.


TYPES

type PageLink struct {
	Page      string
	KeyStroke string
}

```
