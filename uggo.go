/* Package uggo provides helper utilities and page templates
for building uggly servers in golang. */
package uggo

import (
	"fmt"
	pb "github.com/rendicott/uggly"
	"github.com/google/uuid"
)

func init() {
	ThemeDefault = &Theme{}
	ThemeGreen = &Theme{
		StyleTextBoxDescription: Style("yellowgreen", "darkgreen"),
		StyleTextBoxCursor:      Style("aquamarine", "darkseagreen"),
		StyleTextBoxText:        Style("black", "green"),
		StyleTextBoxFill:        Style("black", "green"),
		StyleDivFill:            Style("darkkhaki", "black"),
		StyleDivBorder:          Style("yellowgreen", "darkgreen"),
		StyleTextBlob:           Style("yellowgreen", "black"),
		DivBorderWidth:          int32(1),
		DivBorderChar:           ConvertStringCharRune("="),
		DivFillChar:             ConvertStringCharRune("."),
	}
}


var ThemeDefault *Theme
var ThemeGreen *Theme

type PageLink struct {
	Page      string
	KeyStroke string
}

// ConvertStringCharRune takes a string and converts it to a rune slice
// then grabs the rune at index 0 in the slice so that it can return
// an int32 to satisfy the Uggly protobuf struct for border and fill chars
// and such. If the input string is less than zero length then it will just
// rune out a space char and return that int32.
func ConvertStringCharRune(s string) int32 {
	if len(s) == 0 {
		s = " "
	}
	runes := []rune(s)
	return runes[0]
}

// Style takes a foreground and background color
// string and makes assumptions about the attr
// style to return a uggly Style struct. Easier
// to use than building struct pointers for common
// use cases.
func Style(fg, bg string) *pb.Style {
	return &pb.Style{
		Fg:   fg,
		Bg:   bg,
		Attr: "4",
	}
}

// StrokeMap provides a list of unique keyboard shortcuts that require no
// modifies and are generally easy for the user to hit. Useful for iterating
// through when dynamically assigning links to lists of items
var StrokeMap = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
	"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

// Percent takes a percentage as an integer (e.g., 20) and then calculates
// that percentage of the second argument then returns as an int32. Useful
// for designing DivBox sizes and positioning components on screen.
func Percent(perc int, outof int) int32 {
	p := float64(float64(perc) / float64(100))
	o := float64(outof)
	return int32(p*o)
}

// GrowBox takes a page and a box name then locates it and modifies it's size
func MoveBox(pr *pb.PageResponse, boxName string, x, y int) *pb.PageResponse {
	for _, box := range pr.DivBoxes.Boxes {
		if box.Name == boxName {
			box.StartX += int32(x)
			box.StartY += int32(y)
		}
	}
	return pr
}

// GrowBox takes a page and a box name then locates it and modifies it's size
func GrowBox(pr *pb.PageResponse, boxName string, width, height int) *pb.PageResponse {
	for _, box := range pr.DivBoxes.Boxes {
		if box.Name == boxName {
			box.Width += int32(width)
			box.Height += int32(height)
		}
	}
	return pr
}

// GenPageLittleBox generates a new page with a small box at the given coordinates
func GenPageLittleBox(x, y int) *pb.PageResponse {
	localPage := pb.PageResponse{
		Name:          "generated",
		DivBoxes:      &pb.DivBoxes{},
		Elements:      &pb.Elements{},
		StreamDelayMs: int32(5),
	}
	localPage.DivBoxes.Boxes = append(localPage.DivBoxes.Boxes,
		ThemeDefault.StylizeDivBox(&pb.DivBox{
			Name:   "generated",
			Border: true,
			StartX: int32(x),
			StartY: int32(y),
			Width:  4,
			Height: 4,
		}))
	return &localPage
}

// AddTextBoxToPage adds a textbox and puts the given text in it
func AddTextBoxToPage(pr *pb.PageResponse, text string) *pb.PageResponse {
	// grab last divbox
	width := len(text) + 3
	if len(pr.DivBoxes.Boxes) > 0 {
		box := pr.DivBoxes.Boxes[len(pr.DivBoxes.Boxes)-1]
		startY := box.StartY + 2
		startX := box.StartX
		pr.DivBoxes.Boxes = append(pr.DivBoxes.Boxes,
			ThemeDefault.StylizeDivBox(&pb.DivBox{
				Name:   "added",
				Border: true,
				StartX: startX,
				StartY: startY,
				Width:  int32(width),
				Height: 3,
			}))
		pr.Elements.TextBlobs = append(pr.Elements.TextBlobs,
			ThemeDefault.StylizeTextBlob(&pb.TextBlob{
				Content:  text,
				Wrap:     true,
				DivNames: []string{"added"},
			}))
	}
	return pr
}

// GenPageSimple generates a page that fits within the given w, h and adds
// the provided text
func GenPageSimple(w, h int, text string) *pb.PageResponse {
	newPage := pb.PageResponse{
		Name:     "generated",
		DivBoxes: &pb.DivBoxes{},
		Elements: &pb.Elements{},
	}
	newPage.DivBoxes.Boxes = append(newPage.DivBoxes.Boxes, &pb.DivBox{
		Name:   "generated",
		Border: true,
		StartX: 5,
		StartY: 5,
		Width:  int32(w - 10),
		Height: int32(h - 15),
	})
	newPage.Elements.TextBlobs = append(newPage.Elements.TextBlobs, &pb.TextBlob{
		Content:  text,
		Wrap:     true,
		DivNames: []string{"generated"},
	})
	return ThemeDefault.StylizePage(&newPage)
}

// FindLastBoxEndY takes a page and finds its last div box then returns the extreme
// lower coordinate as an int32. If it encounters any errors it just returns 0
func FindLastBoxEndY(presp *pb.PageResponse) int32 {
	if len(presp.DivBoxes.Boxes) > 0 {
		lastBox := presp.DivBoxes.Boxes[len(presp.DivBoxes.Boxes)-1]
		return lastBox.StartY + lastBox.Height
	}
	return int32(0)
}

// FindLastBoxEndX takes a page and finds its last div box then returns the extreme
// righthand coordinate as an int32. If it encounters any errors it just returns 0
func FindLastBoxEndX(presp *pb.PageResponse) int32 {
	if len(presp.DivBoxes.Boxes) > 0 {
		lastBox := presp.DivBoxes.Boxes[len(presp.DivBoxes.Boxes)-1]
		return lastBox.StartX + lastBox.Width
	}
	return int32(0)
}

// FindLastBoxStartX takes a page and finds its last div box then returns the extreme
// lefthand coordinate as an int32. If it encounters any errors it just returns 0
func FindLastBoxStartX(presp *pb.PageResponse) int32 {
	if len(presp.DivBoxes.Boxes) > 0 {
		lastBox := presp.DivBoxes.Boxes[len(presp.DivBoxes.Boxes)-1]
		return lastBox.StartX
	}
	return int32(0)
}

// FindLastBoxStartY takes a page and finds its last div box then returns the extreme
// top coordinate as an int32. If it encounters any errors it just returns 0
func FindLastBoxStartY(presp *pb.PageResponse) int32 {
	if len(presp.DivBoxes.Boxes) > 0 {
		lastBox := presp.DivBoxes.Boxes[len(presp.DivBoxes.Boxes)-1]
		return lastBox.StartY
	}
	return int32(0)
}

func FindFirstBoxName(presp *pb.PageResponse) string {
	if len(presp.DivBoxes.Boxes) > 0 {
		firstBox := presp.DivBoxes.Boxes[0]
		return firstBox.Name
	}
	return ""
}

func AddFormSimpleTextBox(
		presp *pb.PageResponse,
		x, y int,
		keyStroke, submitPage, divName string) *pb.PageResponse {

	formName := "simple"
	presp.Elements.Forms = append(presp.Elements.Forms, &pb.Form{
		Name:    formName,
		DivName: divName,
		SubmitLink: &pb.Link{
			PageName: submitPage,
		},
		TextBoxes: []*pb.TextBox{
			ThemeDefault.StylizeTextBox(&pb.TextBox{
				Name:            "generated",
				TabOrder:        1,
				DefaultValue:    "hi",
				PositionX:       int32(x),
				PositionY:       int32(y),
				Height:          1,
				Width:           30,
				ShowDescription: false,
			}),
		},
	})
	presp.KeyStrokes = append(presp.KeyStrokes, &pb.KeyStroke{
		KeyStroke: keyStroke,
		Action: &pb.KeyStroke_FormActivation{
			FormActivation: &pb.FormActivation{
				FormName: formName,
			}}})
	return presp
}

// AddFormLogin takes a page and adds a login div, form, helptext, and activation link to it
// it starts the formDiv a few lines below the last divbox on the page
func AddFormLogin(presp *pb.PageResponse, keyStroke, submitPage string) *pb.PageResponse {
	tbPositionX := FindLastBoxStartX(presp) + 10
	divStartX := FindLastBoxStartX(presp) + 2
	divStartY := FindLastBoxStartY(presp) + 4
	presp.DivBoxes.Boxes = append(presp.DivBoxes.Boxes,
		ThemeDefault.StylizeDivBox(&pb.DivBox{
			Name:   "formDiv",
			Border: true,
			StartX: divStartX,
			StartY: divStartY,
			Width:  75,
			Height: 7,
		}))
	loginFormName := "loginForm"
	presp.Elements.Forms = append(presp.Elements.Forms, &pb.Form{
		Name:    loginFormName,
		DivName: "formDiv",
		SubmitLink: &pb.Link{
			PageName: submitPage,
		},
		TextBoxes: []*pb.TextBox{
			ThemeDefault.StylizeTextBox(&pb.TextBox{
				Name:            "password",
				TabOrder:        2,
				DefaultValue:    "",
				Description:     "Password: ",
				PositionX:       tbPositionX,
				PositionY:       4,
				Height:          1,
				Width:           30,
				ShowDescription: true,
				Password:        true,
			}),
			ThemeDefault.StylizeTextBox(&pb.TextBox{
				Name:            "username",
				TabOrder:        1,
				DefaultValue:    "",
				Description:     "Username: ",
				PositionX:       tbPositionX,
				PositionY:       2,
				Height:          1,
				Width:           30,
				ShowDescription: true,
			}),
		},
	})
	presp.KeyStrokes = append(presp.KeyStrokes, &pb.KeyStroke{
		KeyStroke: keyStroke,
		Action: &pb.KeyStroke_FormActivation{
			FormActivation: &pb.FormActivation{
				FormName: loginFormName,
			}}})
	msg := fmt.Sprintf("Hit (%s) to activate form", keyStroke)
	presp.Elements.TextBlobs = append(presp.Elements.TextBlobs,
		ThemeDefault.StylizeTextBlob(&pb.TextBlob{
			Content:  msg,
			Wrap:     true,
			DivNames: []string{"formDiv"},
		}))
	return presp

}

// AddLink adds a link to the given same-server destination on the provided page 
// using the provided keystroke with an indication as to whether or not 
// it's a link to a stream
func AddLink(presp *pb.PageResponse, keyStroke, pageName string, stream bool) *pb.PageResponse {
	presp.KeyStrokes = append(presp.KeyStrokes, &pb.KeyStroke{
		KeyStroke: keyStroke,
		Action: &pb.KeyStroke_Link{
			Link: &pb.Link{
				PageName: pageName,
				Stream:   stream,
			}}})
	return presp
}

// AddFormNewUser takes a page and adds a new user div, form, helptext, and activation link to it
// it starts the formDiv a few lines below the last divbox on the page
func AddFormNewUser(presp *pb.PageResponse, keyStroke, submitPage string) *pb.PageResponse {
	tbPositionX := FindLastBoxStartX(presp) + 18
	divStartX := FindLastBoxStartX(presp) + 2
	divStartY := FindLastBoxStartY(presp) + 4
	presp.DivBoxes.Boxes = append(presp.DivBoxes.Boxes,
		ThemeDefault.StylizeDivBox(&pb.DivBox{
			Name:   "formDiv",
			Border: true,
			StartX: divStartX,
			StartY: divStartY,
			Width:  75,
			Height: 9,
		}))
	loginFormName := "newUserForm"
	presp.Elements.Forms = append(presp.Elements.Forms, &pb.Form{
		Name:    loginFormName,
		DivName: "formDiv",
		SubmitLink: &pb.Link{
			PageName: submitPage,
		},
		TextBoxes: []*pb.TextBox{
			ThemeDefault.StylizeTextBox(&pb.TextBox{
				Name:            "password1",
				TabOrder:        2,
				DefaultValue:    "",
				Description:     "Password: ",
				PositionX:       tbPositionX,
				PositionY:       4,
				Height:          1,
				Width:           30,
				ShowDescription: true,
				Password:        true,
			}),
			ThemeDefault.StylizeTextBox(&pb.TextBox{
				Name:            "password2",
				TabOrder:        3,
				DefaultValue:    "",
				Description:     "Re-Type Password: ",
				PositionX:       tbPositionX,
				PositionY:       6,
				Height:          1,
				Width:           30,
				ShowDescription: true,
				Password:        true,
			}),
			ThemeDefault.StylizeTextBox(&pb.TextBox{
				Name:            "username",
				TabOrder:        1,
				DefaultValue:    "",
				Description:     "Username: ",
				PositionX:       tbPositionX,
				PositionY:       2,
				Height:          1,
				Width:           30,
				ShowDescription: true,
			}),
		},
	})
	presp.KeyStrokes = append(presp.KeyStrokes, &pb.KeyStroke{
		KeyStroke: keyStroke,
		Action: &pb.KeyStroke_FormActivation{
			FormActivation: &pb.FormActivation{
				FormName: loginFormName,
			}}})
	msg := fmt.Sprintf("Hit (%s) to activate form", keyStroke)
	presp.Elements.TextBlobs = append(presp.Elements.TextBlobs,
		ThemeDefault.StylizeTextBlob(&pb.TextBlob{
			Content:  msg,
			Wrap:     true,
			DivNames: []string{"formDiv"},
		}))
	return presp
}
// NewUuid returns a Uuid string. Helpful for guaranteeing unique component names
// or authorzing messages within packages.
func NewUuid() string {
	return uuid.New().String()
}

// AddTextAtGeo takes a page, string and coordinates and then modifies the page with
// a new DivBox at the given coordinates and a TextBlob with the provided string
// Useful when you want to add text at specific locations on a page without the hassle
// DivBox border is forcibly removed despite what the Theme specifies. 
func AddTextAt(presp *pb.PageResponse, x, y, w, h int, text string) *pb.PageResponse {
	divName := fmt.Sprintf("generated-%s", uuid.New())
	divBox := &pb.DivBox{
		Name:   divName,
		Border: false,
		StartX: int32(x),
		StartY: int32(y),
		Width:  int32(w),
		Height: int32(h),
	}
	textBlob := &pb.TextBlob{
		Content:  text,
		Wrap:     true,
		DivNames: []string{divName},
	}
	// stylize so we can clear some things
	divBox = ThemeDefault.StylizeDivBox(divBox)
	textBlob = ThemeDefault.StylizeTextBlob(textBlob)
	divBox.Border = false
	presp.DivBoxes.Boxes = append(presp.DivBoxes.Boxes, divBox)
	presp.Elements.TextBlobs = append(presp.Elements.TextBlobs, textBlob)
	return presp
}
