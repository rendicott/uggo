package uggo

import (
	pb "github.com/rendicott/uggly"
)

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


