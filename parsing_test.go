package list

import (
	"log/slog"
	"math"
	"testing"
)

func getElements() []*Element {
	return []*Element{
		{Name: "01 Hello, World"},
		{Name: "010 Goodbye, World"},
		{Name: "02 ZZZ, World"},
		{Name: "03 Sleepy, World"},
	}
}

type expect struct {
	pat string
	exp []int
}

func TestParser(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	cases := []expect{
		{pat: "[0]", exp: []int{0}},
		{pat: SortS + "r", exp: []int{3, 2, 1, 0}},
		{pat: SortS + "n", exp: []int{0, 2, 3, 1}},
		{pat: SortS + "rn", exp: []int{1, 3, 2, 0}},
		{pat: SortS + "nr", exp: []int{1, 3, 2, 0}},
		{pat: "?[hello]", exp: []int{0}},
		{pat: "![hello]", exp: []int{1, 2, 3}},
		{pat: "![~zzz]", exp: []int{0, 1, 3}},
		{pat: "?[=Sleepy, World]", exp: []int{3}},
	}

	for _, tcase := range cases {
		t.Run(tcase.pat, func(t *testing.T) {
			els := getElements()
			fil, pro := Search(tcase.pat, "")
			newEls := make([]*Element, 0, len(els))
			for _, el := range els {
				if fil(el) {
					newEls = append(newEls, el)
				}
			}

			if len(newEls) == 0 {
				return
			}

			newEls = pro(newEls)

			if len(newEls) != len(tcase.exp) {
				t.Errorf("Expected %d, got %d", len(tcase.exp), len(newEls))
				return
			}
			for i, el := range newEls {
				nam := el.Name
				ex := els[tcase.exp[i]].Name
				if nam != ex {
					t.Errorf("Expected %s, got %s", els[tcase.exp[i]].Name, el.Name)
				}
			}
		})
	}
}

type expectTar struct {
	pat string
	tar string
}

func TestTargeted(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	cases := []expectTar{
		{pat: "t", tar: traversal},
		{pat: "r", tar: traversal},
		{pat: "f", tar: files},
		{pat: "d", tar: dirs},
		{pat: "f" + SortS + "rn", tar: files},
		{pat: "d" + SortS + "rn", tar: dirs},
		{pat: "f" + SortS + "rn[hello]", tar: files},
		{pat: "d" + SortS + "rn[hello]", tar: dirs},
	}

	for _, tcase := range cases {
		t.Run(tcase.pat, func(t *testing.T) {
			opts := &Options{}
			fpp, ok := TargetedSearch(opts, tcase.pat)
			if fpp == nil || !ok {
				return
			}

			if fpp.Tar != tcase.tar {
				t.Errorf("Expected %s, got %s", tcase.tar, fpp.Tar)
			}
		})
	}

	l := expectTar{pat: "t[1:12]", tar: traversal}
	t.Run(l.pat, func(t *testing.T) {
		opts := &Options{}
		fpp, ok := TargetedSearch(opts, l.pat)
		if fpp == nil || !ok {
			return
		}

		if fpp.Tar != l.tar {
			t.Errorf("Expected %s, got %s", l.tar, fpp.Tar)
		}

		if opts.FromDepth != 1 || opts.ToDepth != 12 {
			t.Errorf("Expected 1, 12, got %d, %d", opts.FromDepth, opts.ToDepth)
		}
	})

	l = expectTar{pat: "t[1:]", tar: traversal}
	t.Run(l.pat, func(t *testing.T) {
		opts := &Options{}
		fpp, ok := TargetedSearch(opts, l.pat)
		if fpp == nil || !ok {
			return
		}

		if fpp.Tar != l.tar {
			t.Errorf("Expected %s, got %s", l.tar, fpp.Tar)
		}

		if opts.FromDepth != 1 || opts.ToDepth != math.MaxInt {
			t.Errorf("Expected 1, 0, got %d, %d", opts.FromDepth, opts.ToDepth)
		}
	})

	t.Run("dirs only", func(t *testing.T) {
		patterns := []string{"!f", "?d"}
		for _, pat := range patterns {
			opts := &Options{}
			fpp, ok := TargetedSearch(opts, pat)
			if fpp != nil || !ok {
				t.Errorf("Expected nil, got %v", fpp)
				return
			}

			if opts.DirOnly != true {
				t.Errorf("Expected true, got %t", opts.DirOnly)
			}
		}
	})

	t.Run("files only", func(t *testing.T) {
		patterns := []string{"!d", "?f"}
		for _, pat := range patterns {
			opts := &Options{}
			fpp, ok := TargetedSearch(opts, pat)
			if fpp != nil || !ok {
				t.Errorf("Expected nil, got %v", fpp)
				return
			}

			if opts.FileOnly != true {
				t.Errorf("Expected true, got %t", opts.FileOnly)
			}
		}
	})

	t.Run("negative result", func(t *testing.T) {
		patterns := []string{
			"c:/hello/world/",
			"hello",
			"e?",
			"e!",
			"e",
		}
		for _, pat := range patterns {
			opts := &Options{}
			fpp, ok := TargetedSearch(opts, pat)
			if ok {
				t.Errorf("Expected nil, got %v", fpp)
				return
			}
		}
	})
}
