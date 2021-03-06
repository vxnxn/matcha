package view

import (
	"fmt"
	"time"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/animate"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pointer"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewScrollView", func() view.View {
		return NewScrollView()
	})
}

type ScrollView struct {
	view.Embed
	scrollPosition *view.ScrollPosition
}

func NewScrollView() *ScrollView {
	return &ScrollView{
		scrollPosition: &view.ScrollPosition{},
	}
}

func (v *ScrollView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.Subscribe(v.scrollPosition)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(v.scrollPosition)
	}
}

func (v *ScrollView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	childLayouter := &table.Layouter{}
	for i := 0; i < 5; i++ {
		childLayouter.Add(NewTableCell(), nil)
	}

	scrollview := view.NewScrollView()
	scrollview.ScrollPosition = v.scrollPosition
	scrollview.PaintStyle = &paint.Style{BackgroundColor: colornames.Blue}
	scrollview.ContentLayouter = childLayouter
	scrollview.ContentChildren = childLayouter.Views()
	g1 := l.Add(scrollview, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(0)
		s.Width(200)
		s.HeightEqual(l.Height())
	})

	textView := view.NewTextView()
	textView.PaintStyle = &paint.Style{BackgroundColor: colornames.Red}
	textView.String = fmt.Sprintln("Position:", v.scrollPosition.X.Value(), v.scrollPosition.Y.Value())
	textView.MaxLines = 2
	g2 := l.Add(textView, func(s *constraint.Solver) {
		s.Top(50)
		s.LeftEqual(g1.Right())
		s.RightEqual(l.Right())
		s.Height(100)
	})

	button := view.NewButton()
	button.String = "Scroll"
	button.PaintStyle = &paint.Style{BackgroundColor: colornames.White}
	button.OnPress = func() {
		a := &animate.Basic{
			Start: v.scrollPosition.Y.Value(),
			End:   200,
			Dur:   time.Second / 5,
		}
		v.scrollPosition.Y.Run(a)
	}
	_ = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g2.Bottom())
		s.LeftEqual(g1.Right())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Green},
	}
}

type TableCell struct {
	view.Embed
}

func NewTableCell() *TableCell {
	return &TableCell{}
}

func (v *TableCell) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(200)
	})

	chl := NewTableButton()
	l.Add(chl, func(s *constraint.Solver) {
		s.LeftEqual(l.Left().Add(10))
		s.RightEqual(l.Right().Add(-10))
		s.TopEqual(l.Top().Add(50))
		s.BottomEqual(l.Bottom().Add(-50))
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

type TableButton struct {
	view.Embed
}

func NewTableButton() *TableButton {
	return &TableButton{}
}

func (v *TableButton) Build(ctx view.Context) view.Model {
	return view.Model{
		Painter: &paint.Style{BackgroundColor: colornames.Blue},
		Options: []view.Option{
			pointer.GestureList{
				&pointer.ButtonGesture{
					OnEvent: func(e *pointer.ButtonEvent) {
						v.Signal()
					},
				},
			},
		},
	}
}
