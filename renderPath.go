package main

import (
	"os"
	"sort"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

type pathConfig struct {
	drawFontSize       float64
	drawColumnWidth    float64
	sheetFontSize      float64
	sheetCircleRadius  float64
	sheetCircleSpacing float64
}

var defaultPathConfig = pathConfig{
	drawFontSize:       60,
	drawColumnWidth:    75,
	sheetFontSize:      100,
	sheetCircleRadius:  20,
	sheetCircleSpacing: 20,
}

func renderPath(config pathConfig, startGame *game, draws []draw, outFileName string) {
	// Create new canvases of dimension big enough for the columns and rows of  gears in mm

	columns := make([]drawColumn, len(draws))
	for i, d := range draws {
		columns[i] = drawColumn{
			centerX: (float64(i) + 0.5) * config.drawColumnWidth,
			heading: d.name,
			nodes:   make(map[string]*pathNode),
		}
	}
	colmap := make(map[string]*drawColumn)
	for i := range columns {
		c := columns[i]
		colmap[c.heading] = &c
	}
	nodeByHeights := make(map[string]*pathNode)
	heighths := make([]string, 0)

	var addGame func(g *game, wl string) *pathNode
	addGame = func(g *game, wl string) *pathNode {
		c := colmap[g.drawName]
		node, isNew := c.addNode(g)
		if isNew {
			nodeByHeights[wl] = node
			heighths = append(heighths, wl)
		}

		if g.winnerTo != nil && node.winNode == nil {
			node.winNode = addGame(g.winnerTo, wl+"w")
		}
		if g.loserTo != nil && node.loseNode == nil {
			node.loseNode = addGame(g.loserTo, wl+"l")
		}

		return node
	}
	addGame(startGame, "")

	sort.Slice(heighths, func(l, r int) bool {
		left := heighths[l]
		right := heighths[r]
		for i := 0; ; i++ {
			if i < len(left) && i < len(right) {
				if left[i] == right[i] {
					continue
				} else {
					return left[i] < right[i]
				}
			} else {
				if i < len(left) {
					return left[i] == 'l'
				} else {
					return right[i] == 'w'
				}
			}
		}
	})
	for i, h := range heighths {
		node := nodeByHeights[h]
		node.heightY = float64(i)*config.sheetCircleSpacing + config.sheetCircleSpacing*1.5
	}

	canvasWidth := float64(len(draws)) * config.drawColumnWidth
	canvasHeight := config.sheetCircleSpacing*float64(len(heighths)+3) + config.drawFontSize/4
	drawCanvas := canvas.New(canvasWidth, canvasHeight)
	// Create a canvas contexts used to keep drawing state
	ctx := canvas.NewContext(drawCanvas)
	sheetFace := loadFontFace(config.sheetFontSize)
	drawFace := loadFontFace(config.drawFontSize)

	renderPathNode := func(node pathNode) {

		ctx.Style = defaultStyle
		if node.winNode != nil {
			ctx.MoveTo(node.col.centerX, node.heightY)
			ctx.LineTo(node.winNode.col.centerX, node.winNode.heightY)
			ctx.Stroke()
		}
		if node.loseNode != nil {
			ctx.MoveTo(node.col.centerX, node.heightY)
			ctx.LineTo(node.loseNode.col.centerX, node.loseNode.heightY)
			ctx.Stroke()
		}
		ctx.MoveTo(node.col.centerX+config.sheetCircleRadius, node.heightY)
		ctx.Arc(config.sheetCircleRadius, config.sheetCircleRadius, 0, 0, 360)
		ctx.FillStroke()

		txt := canvas.NewTextLine(sheetFace, node.sheetLetter, canvas.Center)
		ctx.DrawText(node.col.centerX, node.heightY-config.sheetFontSize/8, txt)
	}

	renderColumn := func(col drawColumn) {
		txt := canvas.NewTextLine(drawFace, col.heading, canvas.Center)
		ctx.DrawText(col.centerX, canvasHeight-config.drawFontSize/3, txt)
		for _, n := range col.nodes {
			renderPathNode(*n)
		}
	}

	// nodeA := pathNode{
	// 	sheetLetter: "A",
	// 	heightY:     config.drawColumnWidth * 5 / 2,
	// 	col:         &columns[0],
	// }
	// columns[0].nodes[nodeA.sheetLetter] = &nodeA
	// nodeB := pathNode{
	// 	sheetLetter: "B",
	// 	heightY:     config.drawColumnWidth*5/2 + 30,
	// 	col:         &columns[1],
	// }
	// columns[1].nodes[nodeB.sheetLetter] = &nodeB
	// nodeC := pathNode{
	// 	sheetLetter: "C",
	// 	heightY:     config.drawColumnWidth*5/2 - 30,
	// 	col:         &columns[2],
	// }
	// columns[2].nodes[nodeC.sheetLetter] = &nodeC
	// nodeA.winNode = &nodeB
	// nodeA.loseNode = &nodeC

	// renderPathNode(nodeA)
	// renderPathNode(nodeB)
	// renderPathNode(nodeC)
	ctx.Style = defaultStyle
	ctx.MoveTo(0, 0)
	ctx.LineTo(canvasWidth, 0)
	ctx.LineTo(canvasWidth, canvasHeight)
	ctx.LineTo(0, canvasHeight)
	ctx.LineTo(0, 0)
	ctx.FillStroke()

	dashed := defaultStyle
	dashed.Dashes = []float64{5, 5}
	ctx.Style = dashed
	for i := 1; i < len(columns); i++ {
		ctx.MoveTo(float64(i)*config.drawColumnWidth, 0)
		ctx.LineTo(float64(i)*config.drawColumnWidth, canvasHeight)
		ctx.Stroke()
	}
	for _, c := range columns {
		renderColumn(c)
	}

	f, err := os.Create(outFileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := renderers.PNG()
	if err := w(f, drawCanvas); err != nil {
		panic(err)
	}

}

type drawColumn struct {
	centerX  float64
	currentY float64
	heading  string
	nodes    map[string]*pathNode
}

func (dc *drawColumn) addNode(game *game) (*pathNode, bool) {
	letter := game.sheetName[4:]
	if dc.nodes[letter] != nil {
		return dc.nodes[letter], false
	} else {
		n := pathNode{
			sheetLetter: letter,
			heightY:     dc.currentY,
			col:         dc,
		}
		dc.nodes[letter] = &n
		dc.currentY -= 60
		return &n, true
	}
}

type pathNode struct {
	sheetLetter string
	heightY     float64
	col         *drawColumn
	winNode     *pathNode
	loseNode    *pathNode
}

var defaultStyle = canvas.Style{
	FillColor:    canvas.White,
	StrokeColor:  canvas.Black,
	StrokeWidth:  1.0,
	StrokeCapper: canvas.ButtCap,
	StrokeJoiner: canvas.MiterJoin,
	DashOffset:   0.0,
	Dashes:       []float64{},
	FillRule:     canvas.NonZero,
}

func loadFontFace(fontSize float64) *canvas.FontFace {
	fontFamily := canvas.NewFontFamily("Arial")
	style := canvas.FontRegular
	if err := fontFamily.LoadLocalFont("Arial", style); err != nil {
		panic(err)
	}
	return fontFamily.Face(fontSize, canvas.Black, style, canvas.FontNormal)
}
