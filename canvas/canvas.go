package canvas

import "strings"

type BaseCanvas struct {
	RowData []string
}


type Canvas struct {
	Canvas BaseCanvas
	Width int
	Height int
}

func (c *Canvas) InitializeCanvas() {
	c.Height = 40
	c.Width = 150
	for i := 0; i < c.Height; i++ {
		c.Canvas.RowData = append(c.Canvas.RowData, strings.Repeat("+", c.Width))
	}
}