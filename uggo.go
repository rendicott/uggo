/* Package uggo provides helper utilities and page templates
for building uggly servers in golang. */
package uggo

import (
	"fmt"
	pb "github.com/rendicott/uggly"
)

func init() {
	ThemeDefault = &Theme{}
	ThemeGreen = &Theme{
		StyleTextBoxDescription: Style("yellowgreen", "darkgreen"),
		StyleTextBoxCursor:      Style("aquamarine", "darkseagreen"),
		StyleTextBoxText:        Style("springgreen", "green"),
		StyleTextBoxFill:        Style("springgreen", "green"),
		StyleDivFill:            Style("darkkhaki", "darkgreen"),
		StyleDivBorder:          Style("yellowgreen", "darkgreen"),
		StyleTextBlob:           Style("yellowgreen", "darkgreen"),
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

// Theme holds visual settings for components and provides methods
// for applying those settings to pages and components of pages
type Theme struct {
	StyleTextBoxDescription,
	StyleTextBoxCursor,
	StyleTextBoxText,
	StyleTextBoxFill,
	StyleDivFill,
	StyleDivBorder,
	StyleTextBlob *pb.Style
	DivBorderWidth int32
	DivBorderChar,
	DivFillChar rune
}

// Init instantiates any nil Theme settings with default
// values
func (t *Theme) Init() {
	if t.StyleTextBoxDescription == nil {
		t.StyleTextBoxDescription = Style("white", "black")
	}
	if t.StyleTextBoxCursor == nil {
		t.StyleTextBoxCursor = Style("black", "white")
	}
	if t.StyleTextBoxText == nil {
		t.StyleTextBoxText = Style("white", "blue")
	}
	if t.StyleTextBoxFill == nil {
		t.StyleTextBoxFill = Style("white", "blue")
	}
	if t.StyleDivFill == nil {
		t.StyleDivFill = Style("white", "black")
	}
	if t.StyleDivBorder == nil {
		t.StyleDivBorder = Style("white", "black")
	}
	if t.StyleTextBlob == nil {
		t.StyleTextBlob = Style("white", "black")
	}
	if t.DivFillChar == 0 {
		t.DivFillChar = ConvertStringCharRune(" ")
	}
	if t.DivBorderChar == 0 {
		t.DivBorderChar = ConvertStringCharRune("*")
	}
	if t.DivBorderWidth == 0 {
		t.DivBorderWidth = int32(1)
	}
}

// StylizePage takes a page and applies this theme to the divBoxes,
// textBlobs, and textBoxes within it and returns the modified page
func (t *Theme) StylizePage(page *pb.PageResponse) *pb.PageResponse {
	for i, divBox := range page.DivBoxes.Boxes {
		page.DivBoxes.Boxes[i] = t.StylizeDivBox(divBox)
	}
	for i, textBlob := range page.Elements.TextBlobs {
		page.Elements.TextBlobs[i] = t.StylizeTextBlob(textBlob)
	}
	for i, form := range page.Elements.Forms {
		for j, textBox := range form.TextBoxes {
			page.Elements.Forms[i].TextBoxes[j] = t.StylizeTextBox(textBox)
		}
	}
	return page
}

// StylizeTextBox applies this theme to the provided textBox and returns it
func (t *Theme) StylizeTextBox(textBox *pb.TextBox) *pb.TextBox {
	t.Init()
	textBox.StyleCursor = t.StyleTextBoxCursor
	textBox.StyleFill = t.StyleTextBoxFill
	textBox.StyleDescription = t.StyleTextBoxDescription
	textBox.StyleText = t.StyleTextBoxText
	return textBox
}

// StylizeDivBox applies this theme to the provided divBox and returns it
func (t *Theme) StylizeDivBox(divBox *pb.DivBox) *pb.DivBox {
	t.Init()
	divBox.FillSt = t.StyleDivFill
	divBox.BorderSt = t.StyleDivBorder
	divBox.BorderW = t.DivBorderWidth
	divBox.FillChar = t.DivFillChar
	divBox.BorderChar = t.DivBorderChar
	return divBox
}

// StylizeTextBlob applies this theme to the provided textBlob and returns it
func (t *Theme) StylizeTextBlob(textBlob *pb.TextBlob) *pb.TextBlob {
	t.Init()
	textBlob.Style = t.StyleTextBlob
	return textBlob
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

// AddLink adds a link to the given destination on the provided page using the provided keystroke
// with an indication as to whether or not it's a link to a stream
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
