/* Package uggo provides helper utilities and page templates 
for building uggly servers in golang. */
package uggo

import (
	"fmt"
	pb "github.com/rendicott/uggly"
)

type PageLink struct {
	Page string
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



func PageTopMenuFullWidthContent(
	w,h int,
	links []*PageLink,
	current, content string) (*pb.PageResponse) {
	tpr := pb.PageResponse{
		Name:     "generic",
		DivBoxes: &pb.DivBoxes{},
		Elements: &pb.Elements{},
	}
	// build menu
	xPos := 0
	menuY := 0
	menuHeight := 1
	for i, link := range links {
		text := fmt.Sprintf("%s (%s)", link.Page, link.KeyStroke)
		divWidth := len(text) + 3
		if i != 0 {
			xPos = xPos + divWidth
		}
		divName := fmt.Sprintf("menu-%d", i)
		divFillSt := Style("white","black")
		textStyle := Style("white", "black")
		if link.Page == current {
			divFillSt = Style("white", "dimgray")
			textStyle = Style("white", "dimgray")
		}
		tpr.DivBoxes.Boxes = append(tpr.DivBoxes.Boxes, &pb.DivBox{
			Name:     divName,
			Border:   false,
			FillChar: ConvertStringCharRune(" "),
			//BorderChar: ConvertStringCharRune("*"),
			StartX:   int32(xPos),
			StartY:   int32(menuY),
			Width:    int32(divWidth),
			Height:   int32(menuHeight),
			FillSt:   divFillSt,
			//BorderSt: Style("green", "black"),
		})
		tpr.Elements.TextBlobs = append(tpr.Elements.TextBlobs, &pb.TextBlob{
			Content:  text,
			Wrap:     false,
			Style:    textStyle,
			DivNames: []string{divName},
		})
		tpr.KeyStrokes = append(tpr.KeyStrokes, &pb.KeyStroke{
			KeyStroke: link.KeyStroke,
			Action: &pb.KeyStroke_Link{
				&pb.Link{
					PageName: link.Page,
				},
			},
		})
	}

	// build content
	contentDivSt := Style("black","lightslategray")
	contentTextStyle := Style("black","lightslategray")
	//contentY := menuY + menuHeight + 1

	contentY := 1
	contentHeight := h-10
	tpr.DivBoxes.Boxes = append(tpr.DivBoxes.Boxes, &pb.DivBox{
		Name:     "content",
		Border:   false,
		FillChar: ConvertStringCharRune(" "),
		//BorderChar: ConvertStringCharRune("*"),
		StartX:   int32(1),
		StartY:   int32(contentY),
		Width:    int32(w),
		Height:   int32(contentHeight),
		FillSt:   contentDivSt,
		//BorderSt: Style("green", "black"),
	})
	tpr.Elements.TextBlobs = append(tpr.Elements.TextBlobs, &pb.TextBlob{
		Content:  content,
		Wrap:     true,
		Style:    contentTextStyle,
		DivNames: []string{"content"},
	})

	// build bottom Menu
	bMenuDivSt:= Style("black","lightslategray")
	bMenuTextStyle := Style("black","lightslategray")
	bMenuY := contentY + contentHeight + 1
	bMenuHeight := 2
	bMenuContent := "continue (j)"
	tpr.DivBoxes.Boxes = append(tpr.DivBoxes.Boxes, &pb.DivBox{
		Name:     "bMenu",
		Border:   false,
		FillChar: ConvertStringCharRune(" "),
		//BorderChar: ConvertStringCharRune("*"),
		StartX:   int32(1),
		StartY:   int32(bMenuY),
		Width:    int32(w),
		Height:   int32(bMenuHeight),
		FillSt:   bMenuDivSt,
		//BorderSt: Style("green", "black"),
	})
	tpr.Elements.TextBlobs = append(tpr.Elements.TextBlobs, &pb.TextBlob{
		Content:  bMenuContent,
		Wrap:     true,
		Style:    bMenuTextStyle,
		DivNames: []string{"bMenu"},
	})
	tpr.KeyStrokes = append(tpr.KeyStrokes, &pb.KeyStroke{
		KeyStroke: "j",
		Action: &pb.KeyStroke_DivScroll{
			&pb.DivScroll{
				DivName: "content",
				Down: false,
			},
		},
	})
	tpr.KeyStrokes = append(tpr.KeyStrokes, &pb.KeyStroke{
		KeyStroke: "k",
		Action: &pb.KeyStroke_DivScroll{
			&pb.DivScroll{
				DivName: "content",
				Down: true,
			},
		},
	})
	return &tpr
}

// GrowBox takes a page and a box name then locates it and modifies it's size
func MoveBox(pr *pb.PageResponse, boxName string, x, y int) (*pb.PageResponse) {
	for _, box := range pr.DivBoxes.Boxes {
		if box.Name == boxName {
			box.StartX+=int32(x)
			box.StartY+=int32(y)
		}
	}
	return pr
}

// GrowBox takes a page and a box name then locates it and modifies it's size
func GrowBox(pr *pb.PageResponse, boxName string, width, height int) (*pb.PageResponse) {
	for _, box := range pr.DivBoxes.Boxes {
		if box.Name == boxName {
			box.Width+=int32(width)
			box.Height+=int32(height)
		}
	}
	return pr
}

// GenPageLittleBox generates a new page with a small box at the given coordinates
func GenPageLittleBox(x, y int) (*pb.PageResponse) {
	localPage := pb.PageResponse{
		Name: "generated",
		DivBoxes: &pb.DivBoxes{},
		Elements: &pb.Elements{},
		StreamDelayMs: int32(5),
	}
	localPage.DivBoxes.Boxes = append(localPage.DivBoxes.Boxes, &pb.DivBox{
		Name:     "generated",
		Border:   true,
		BorderW:  int32(1),
		BorderChar: ConvertStringCharRune("^"),
		FillChar: ConvertStringCharRune(""),
		StartX:   int32(x),
		StartY:   int32(y),
		Width:    4,
		Height:   4,
		BorderSt: Style("orange", "black"),
		FillSt: Style("white","black"),
	})
	//localPage.StreamDelayMs = int32(5)
	return &localPage
}

// AddTextBoxToPage adds a textbox and puts the given text in it
func AddTextBoxToPage(pr *pb.PageResponse, text string) (*pb.PageResponse) {
	// grab last divbox
	width := len(text) + 3
	if len(pr.DivBoxes.Boxes) > 0 {
		box := pr.DivBoxes.Boxes[len(pr.DivBoxes.Boxes) -1]
		startY := box.StartY + 2
		startX := box.StartX
		pr.DivBoxes.Boxes = append(pr.DivBoxes.Boxes, &pb.DivBox{
			Name:     "added",
			Border:   true,
			BorderW:  int32(1),
			BorderChar: ConvertStringCharRune("^"),
			FillChar: ConvertStringCharRune(""),
			StartX:   startX,
			StartY:   startY,
			Width:    int32(width),
			Height:   3,
			BorderSt: Style("orange", "black"),
			FillSt: Style("white","black"),
		})
		pr.Elements.TextBlobs = append(pr.Elements.TextBlobs, &pb.TextBlob{
			Content: text,
			Wrap:    true,
			Style: Style("white","black"),
			DivNames: []string{"added"},
		})
	}
	return pr
}


// GenPageSimple generates a page that fits within the given w, h and adds
// the provided text
func GenPageSimple(w, h int, text string) (*pb.PageResponse) {
	newPage := pb.PageResponse{
		Name: "generated",
		DivBoxes: &pb.DivBoxes{},
		Elements: &pb.Elements{},
	}
	newPage.DivBoxes.Boxes = append(newPage.DivBoxes.Boxes, &pb.DivBox{
		Name:     "generated",
		Border:   true,
		BorderW:  int32(1),
		BorderChar: ConvertStringCharRune("^"),
		FillChar: ConvertStringCharRune(""),
		StartX:   5,
		StartY:   5,
		Width:    int32(w-10),
		Height:   int32(h-15),
		BorderSt: Style("orange", "black"),
		FillSt: Style("white","black"),
	})
	newPage.Elements.TextBlobs = append(newPage.Elements.TextBlobs, &pb.TextBlob{
		Content: text,
		Wrap:    true,
		Style: Style("white","black"),
		DivNames: []string{"generated"},
	})
	return &newPage
}

// FindLastBoxEndY takes a page and finds its last div box then returns the extreme
// lower coordinate as an int32. If it encounters any errors it just returns 0
func FindLastBoxEndY(presp *pb.PageResponse) (int32) {
	if len(presp.DivBoxes.Boxes) > 0 {
		lastBox := presp.DivBoxes.Boxes[len(presp.DivBoxes.Boxes) - 1]
		return lastBox.StartY + lastBox.Height
	}
	return int32(0)
}

// FindLastBoxEndX takes a page and finds its last div box then returns the extreme
// righthand coordinate as an int32. If it encounters any errors it just returns 0
func FindLastBoxEndX(presp *pb.PageResponse) (int32) {
	if len(presp.DivBoxes.Boxes) > 0 {
		lastBox := presp.DivBoxes.Boxes[len(presp.DivBoxes.Boxes) - 1]
		return lastBox.StartX + lastBox.Width
	}
	return int32(0)
}

// FindLastBoxStartX takes a page and finds its last div box then returns the extreme
// lefthand coordinate as an int32. If it encounters any errors it just returns 0
func FindLastBoxStartX(presp *pb.PageResponse) (int32) {
	if len(presp.DivBoxes.Boxes) > 0 {
		lastBox := presp.DivBoxes.Boxes[len(presp.DivBoxes.Boxes) - 1]
		return lastBox.StartX
	}
	return int32(0)
}

// FindLastBoxStartY takes a page and finds its last div box then returns the extreme
// top coordinate as an int32. If it encounters any errors it just returns 0
func FindLastBoxStartY(presp *pb.PageResponse) (int32) {
	if len(presp.DivBoxes.Boxes) > 0 {
		lastBox := presp.DivBoxes.Boxes[len(presp.DivBoxes.Boxes) - 1]
		return lastBox.StartY
	}
	return int32(0)
}



// AddFormLogin takes a page and adds a login div, form, helptext, and activation link to it
// it starts the formDiv a few lines below the last divbox on the page
func AddFormLogin(presp *pb.PageResponse, keyStroke, submitPage string) (*pb.PageResponse) {
	tbPositionX := FindLastBoxStartX(presp) + 10
	divStartX := FindLastBoxStartX(presp) + 2
	divStartY := FindLastBoxStartY(presp) + 4
	presp.DivBoxes.Boxes = append(presp.DivBoxes.Boxes, &pb.DivBox{
		Name:     "formDiv",
		Border:   true,
		BorderW:  int32(1),
		BorderChar: ConvertStringCharRune("#"),
		FillChar: ConvertStringCharRune(""),
		StartX:   divStartX,
		StartY:   divStartY,
		Width:    75,
		Height:   7,
		BorderSt: Style("blue", "black"),
		FillSt: Style("white","black"),
	})
	loginFormName := "loginForm"
	presp.Elements.Forms = append(presp.Elements.Forms, &pb.Form{
		Name: loginFormName,
		DivName: "formDiv",
		SubmitLink: &pb.Link{
			PageName: submitPage,
		},
		TextBoxes: []*pb.TextBox{
			&pb.TextBox{
				Name: "password",
				TabOrder: 2,
				DefaultValue: "",
				Description: "Password: ",
				PositionX: tbPositionX,
				PositionY: 4,
				Height: 1,
				Width: 30,
				StyleCursor: Style("black", "gray"),
				StyleFill: Style("black", "blue"),
				StyleText: Style("white", "blue"),
				StyleDescription: Style("red", "black"),
				ShowDescription: true,
				Password: true,
			},
			&pb.TextBox{
				Name: "username",
				TabOrder: 1,
				DefaultValue: "",
				Description: "Username: ",
				PositionX: tbPositionX,
				PositionY: 2,
				Height: 1,
				Width: 30,
				StyleCursor: Style("black", "gray"),
				StyleFill: Style("black", "blue"),
				StyleText: Style("white", "blue"),
				StyleDescription: Style("red", "black"),
				ShowDescription: true,
			},
		},
	})
	presp.KeyStrokes = append(presp.KeyStrokes, &pb.KeyStroke{
		KeyStroke: keyStroke,
		Action: &pb.KeyStroke_FormActivation{
			FormActivation: &pb.FormActivation{
				FormName: loginFormName,
	}}})
	msg := fmt.Sprintf("Hit (%s) to activate form", keyStroke)
	presp.Elements.TextBlobs = append(presp.Elements.TextBlobs, &pb.TextBlob{
		Content: msg,
		Wrap:    true,
		Style: Style("white", "black"),
		DivNames: []string{"formDiv"},
	})
	return presp

}


// AddLink adds a link to the given destination on the provided page using the provided keystroke
// with an indication as to whether or not it's a link to a stream
func AddLink(presp *pb.PageResponse, keyStroke, pageName string, stream bool) (*pb.PageResponse) {
	presp.KeyStrokes = append(presp.KeyStrokes, &pb.KeyStroke{
		KeyStroke: keyStroke,
		Action: &pb.KeyStroke_Link{
			Link: &pb.Link{
				PageName: pageName,
				Stream: stream,
	}}})
	return presp
}

// AddFormNewUser takes a page and adds a new user div, form, helptext, and activation link to it
// it starts the formDiv a few lines below the last divbox on the page
func AddFormNewUser(presp *pb.PageResponse, keyStroke, submitPage string) (*pb.PageResponse) {
	tbPositionX := FindLastBoxStartX(presp) + 18
	divStartX := FindLastBoxStartX(presp) + 2
	divStartY := FindLastBoxStartY(presp) + 4
	presp.DivBoxes.Boxes = append(presp.DivBoxes.Boxes, &pb.DivBox{
		Name:     "formDiv",
		Border:   true,
		BorderW: int32(1),
		BorderChar: ConvertStringCharRune("#"),
		FillChar: ConvertStringCharRune(""),
		StartX:   divStartX,
		StartY:   divStartY,
		Width:    75,
		Height:   9,
		BorderSt: Style("blue", "black"),
		FillSt: Style("white","black"),
	})
	loginFormName := "newUserForm"
	presp.Elements.Forms = append(presp.Elements.Forms, &pb.Form{
		Name: loginFormName,
		DivName: "formDiv",
		SubmitLink: &pb.Link{
			PageName: submitPage,
		},
		TextBoxes: []*pb.TextBox{
			&pb.TextBox{
				Name: "password1",
				TabOrder: 2,
				DefaultValue: "",
				Description: "Password: ",
				PositionX: tbPositionX,
				PositionY: 4,
				Height: 1,
				Width: 30,
				StyleCursor: Style("black", "gray"),
				StyleFill: Style("black", "blue"),
				StyleText: Style("white", "blue"),
				StyleDescription: Style("red", "black"),
				ShowDescription: true,
				Password: true,
			},
			&pb.TextBox{
				Name: "password2",
				TabOrder: 3,
				DefaultValue: "",
				Description: "Re-Type Password: ",
				PositionX: tbPositionX,
				PositionY: 6,
				Height: 1,
				Width: 30,
				StyleCursor: Style("black", "gray"),
				StyleFill: Style("black", "blue"),
				StyleText: Style("white", "blue"),
				StyleDescription: Style("red", "black"),
				ShowDescription: true,
				Password: true,
			},
			&pb.TextBox{
				Name: "username",
				TabOrder: 1,
				DefaultValue: "",
				Description: "Username: ",
				PositionX: tbPositionX,
				PositionY: 2,
				Height: 1,
				Width: 30,
				StyleCursor: Style("black", "gray"),
				StyleFill: Style("black", "blue"),
				StyleText: Style("white", "blue"),
				StyleDescription: Style("red", "black"),
				ShowDescription: true,
			},
		},
	})
	presp.KeyStrokes = append(presp.KeyStrokes, &pb.KeyStroke{
		KeyStroke: keyStroke,
		Action: &pb.KeyStroke_FormActivation{
			FormActivation: &pb.FormActivation{
				FormName: loginFormName,
	}}})
	msg := fmt.Sprintf("Hit (%s) to activate form", keyStroke)
	presp.Elements.TextBlobs = append(presp.Elements.TextBlobs, &pb.TextBlob{
		Content: msg,
		Wrap:    true,
		Style: Style("white", "black"),
		DivNames: []string{"formDiv"},
	})
	return presp
}
