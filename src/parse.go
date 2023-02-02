package main

import (
	"fmt"
	"log"
	"strings"
)

var shorthandToAttr = map[rune]string{
	'l': "Length",
	'f': "Filename",
	'D': "Directory",
	'a': "Artist",
	'A': "AlbumArtist",
	't': "Title",
	'b': "Album",
	'y': "Date",
	//'n': "track number (01/12 -> 01)",
	//'N': "full track info (01/12 -> 01/12)",
	'g': "Genre",
	'c': "Composer",
	'p': "Performer",
	'd': "Disc",
	'C': "Comment",
	'P': "Priority",
}

func parseShorthand(ch rune, tr Track) string {
	attr, _ := tr.Attr(shorthandToAttr[ch])
	return attr
}

type Text interface {
	Format(Track, bool) string
}

func NewText(s MetadataFormat) (Text, error) {
	segments, err := NestedBrackets(s)
	if err != nil {
		return nil, err
	}
	return &text{segments: segments}, nil
}

type text struct {
	segments []conditionalText
}

func (t *text) Format(tr Track, colorized bool) string {
	// TODO: rename vars / vals
	vars := map[rune]string{}
	for r, s := range shorthandToAttr {
		vars[r], _ = tr.Attr(s)
	}

	var builder strings.Builder
	for _, segment := range t.segments {
		if !segment.eval(vars) {
			continue
		}

		vals := make([]interface{}, len(segment.vars))
		for i := range segment.vars {
			vals[i] = vars[segment.vars[i]]
		}
		builder.WriteString(fmt.Sprintf(segment.fmt, vals...))
	}
	if !colorized {
		return stripColors(builder.String())
	}
	return builder.String()
}

func stripColors(str string) string {
	var builder strings.Builder
	for i := 0; i < len(str); i++ {
		switch str[i] {
		case '$': // skip $ and next char
			i++
		default:
			builder.WriteByte(str[i])
		}
	}
	return builder.String()
}

type conditionalText struct {
	fmt  string
	vars []rune

	co condition
}

func (ct *conditionalText) eval(vars map[rune]string) bool {
	return ct.co.eval(vars)
}

type operator int

const (
	AND operator = iota
	OR
	NAND
	NOR
)

type condition struct {
	op operator
	cs []condition
	rs []rune
}

func (co condition) eval(vars map[rune]string) bool {
	switch co.op {
	case AND:
		for _, r := range co.rs {
			if v := vars[r]; v == "" {
				return false
			}
		}
		for _, c := range co.cs {
			if ok := c.eval(vars); !ok {
				return false
			}
		}
		return true
	case OR:
		for _, r := range co.rs {
			if v := vars[r]; v != "" {
				return true
			}
		}
		for _, c := range co.cs {
			if ok := c.eval(vars); ok {
				return true
			}
		}
		return true
	case NAND:
		return !condition{op: AND, rs: co.rs, cs: co.cs}.eval(vars)
	case NOR:
		return !condition{op: OR, rs: co.rs, cs: co.cs}.eval(vars)
	default:
		log.Panicln("condition missing operator")
		return true
	}
}

func and(rs []rune, tr Track) bool {
	for _, r := range rs {
		if parseShorthand(r, tr) == "" {
			return false
		}
	}
	return true
}

// TODO: add error handling
// TODO: split up, if possible
func NestedBrackets(s MetadataFormat) ([]conditionalText, error) {
	and := [][]rune{{}}
	andidx := []int{0}

	nand := [][][]rune{{{}}}
	nandidx := 0

	var alt bool

	var fmtbuilder strings.Builder
	var vars []rune

	type segment struct {
		fmt  string
		vars []rune

		and  []int
		nand [][]rune
	}
	var segments []segment
	writeSegment := func() {
		if len(fmtbuilder.String()) == 0 {
			return
		}

		segments = append(segments, segment{
			fmt:  fmtbuilder.String(),
			vars: append([]rune{}, vars...),

			and:  append([]int{}, andidx[1:]...),
			nand: [][]rune{}})
		for j := 0; j < len(nand[nandidx])-1; j++ {
			segments[len(segments)-1].nand = append(segments[len(segments)-1].nand, nand[nandidx][j])
		}

		fmtbuilder.Reset()
		vars = []rune{}
	}

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{':
			writeSegment()

			and = append(and, []rune{})
			andidx = append(andidx, len(and)-1)

			nandidx++
			if alt {
				if len(nand) > nandidx+1 {
					nand = nand[:nandidx+1]
				}
				nand[nandidx] = append(nand[nandidx], []rune{})
				alt = false
			} else {
				if len(nand) > nandidx {
					nand = nand[:nandidx]
				}
				nand = append(nand, [][]rune{{}})
			}
		case '}':
			writeSegment()

			andidx = andidx[:len(andidx)-1]
			nandidx--
		case '|':
			if i > 0 && s[i-1] == '}' {
				writeSegment()
				alt = true
			} else {
				fmtbuilder.WriteByte(s[i])
			}
		case '%':
			fmtbuilder.WriteString("%s")
			i++
			vars = append(vars, rune(s[i]))

			and[andidx[len(andidx)-1]] = append(and[andidx[len(andidx)-1]], rune(s[i]))
			nand[nandidx][len(nand[nandidx])-1] = append(nand[nandidx][len(nand[nandidx])-1], rune(s[i]))
		default:
			fmtbuilder.WriteByte(s[i])
		}
	}
	writeSegment()

	var res []conditionalText
	for _, s := range segments {
		var cs []condition

		andrs := []rune{}
		for _, idx := range s.and {
			andrs = append(andrs, and[idx]...)
		}
		if len(andrs) > 0 {
			cs = append(cs, condition{op: AND, rs: andrs})
		}

		for _, rs := range s.nand {
			if len(rs) > 0 {
				cs = append(cs, condition{op: NAND, rs: rs})
			}
		}

		ct := conditionalText{
			fmt:  s.fmt,
			vars: s.vars,
			co: condition{
				op: AND,
				cs: cs,
			},
		}
		res = append(res, ct)
	}
	return res, nil
}
