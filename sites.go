package uggo

import (
	"fmt"
	pb "github.com/rendicott/uggly"
)

func PageTopMenuFullWidthContent(
	w, h int,
	links []*PageLink,
	current, content string) *pb.PageResponse {
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
		divFillSt := Style("white", "black")
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
			StartX: int32(xPos),
			StartY: int32(menuY),
			Width:  int32(divWidth),
			Height: int32(menuHeight),
			FillSt: divFillSt,
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
	contentDivSt := Style("black", "lightslategray")
	contentTextStyle := Style("black", "lightslategray")
	//contentY := menuY + menuHeight + 1

	contentY := 1
	contentHeight := h - 10
	tpr.DivBoxes.Boxes = append(tpr.DivBoxes.Boxes, &pb.DivBox{
		Name:     "content",
		Border:   false,
		FillChar: ConvertStringCharRune(" "),
		//BorderChar: ConvertStringCharRune("*"),
		StartX: int32(1),
		StartY: int32(contentY),
		Width:  int32(w),
		Height: int32(contentHeight),
		FillSt: contentDivSt,
		//BorderSt: Style("green", "black"),
	})
	tpr.Elements.TextBlobs = append(tpr.Elements.TextBlobs, &pb.TextBlob{
		Content:  content,
		Wrap:     true,
		Style:    contentTextStyle,
		DivNames: []string{"content"},
	})

	// build bottom Menu
	bMenuDivSt := Style("black", "lightslategray")
	bMenuTextStyle := Style("black", "lightslategray")
	bMenuY := contentY + contentHeight + 1
	bMenuHeight := 2
	bMenuContent := "continue (j)"
	tpr.DivBoxes.Boxes = append(tpr.DivBoxes.Boxes, &pb.DivBox{
		Name:     "bMenu",
		Border:   false,
		FillChar: ConvertStringCharRune(" "),
		//BorderChar: ConvertStringCharRune("*"),
		StartX: int32(1),
		StartY: int32(bMenuY),
		Width:  int32(w),
		Height: int32(bMenuHeight),
		FillSt: bMenuDivSt,
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
				Down:    false,
			},
		},
	})
	tpr.KeyStrokes = append(tpr.KeyStrokes, &pb.KeyStroke{
		KeyStroke: "k",
		Action: &pb.KeyStroke_DivScroll{
			&pb.DivScroll{
				DivName: "content",
				Down:    true,
			},
		},
	})
	return &tpr
}
