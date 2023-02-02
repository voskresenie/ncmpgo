package main

import (
	"testing"
)

func (a condition) Equals(b condition) bool {
	if a.op != b.op {
		return false
	}
	if !eqSet(a.rs, b.rs) {
		return false
	}
	if len(a.cs) != len(b.cs) {
		return false
	}
	// TODO: make this not dependent on order...
	for i := range a.cs {
		if !a.cs[i].Equals(b.cs[i]) {
			return false
		}
	}
	return true
}

func (a conditionalText) Equals(b conditionalText) bool {
	if a.fmt != b.fmt {
		return false
	}
	if !eqSlice(a.vars, b.vars) {
		return false
	}
	if !a.co.Equals(b.co) {
		return false
	}
	return true
}

func TestFormat(t *testing.T) {
	testcases := []struct {
		fmt   MetadataFormat
		attrs map[string]string
		want  string
	}{
		{
			fmt:   "%a",
			attrs: map[string]string{"Artist": "TestArtist"},
			want:  "TestArtist",
		},
		{
			fmt:   "%a{ - %b}",
			attrs: map[string]string{"Artist": "TestArtist", "Album": ""},
			want:  "TestArtist",
		},
		{
			fmt:   "%a{ - %b}",
			attrs: map[string]string{"Artist": "TestArtist", "Album": "TestAlbum"},
			want:  "TestArtist - TestAlbum",
		},
		{
			fmt:   "{%a}|{%b}|{%c}",
			attrs: map[string]string{"Artist": "TestArtist", "Album": "TestAlbum", "Composer": "TestComposer"},
			want:  "TestArtist",
		},
		{
			fmt:   "{%a}|{%b}|{%c}",
			attrs: map[string]string{"Artist": "", "Album": "TestAlbum", "Composer": "TestComposer"},
			want:  "TestAlbum",
		},
		{
			fmt:   "{%a}|{%b}|{%c}",
			attrs: map[string]string{"Artist": "", "Album": "", "Composer": "TestComposer"},
			want:  "TestComposer",
		},
		{
			fmt:   "{%a - %b}|{%c - %d}",
			attrs: map[string]string{"Artist": "", "Album": "TestAlbum", "Composer": "TestComposer", "Disc": ""},
			want:  "",
		},
		{
			fmt:   "{%a - %b}|{%c - %d}",
			attrs: map[string]string{"Artist": "TestArtist", "Album": "", "Composer": "TestComposer", "Disc": "1"},
			want:  "TestComposer - 1",
		},
		{
			fmt:   "{%a - %b}|{%c - %d}",
			attrs: map[string]string{"Artist": "TestArtist", "Album": "TestAlbum", "Composer": "TestComposer", "Disc": "1"},
			want:  "TestArtist - TestAlbum",
		},
	}

	for _, testcase := range testcases {
		txt, _ := NewText(testcase.fmt)
		got := txt.Format(NewTrack(testcase.attrs), true)
		if got != testcase.want {
			t.Errorf("Format(\"%s\", true): incorrect result\n\twant %s\n\t got %s",
				testcase.fmt, testcase.want, got)
		}
	}
}

func TestNestedBrackets(t *testing.T) {
	testcases := []struct {
		name string
		fmt  MetadataFormat
		want []conditionalText
		err  error
	}{
		{
			fmt: "%a",
			want: []conditionalText{{
				fmt:  "%s",
				vars: []rune{'a'},
			}},
		},
		{
			fmt: "%a{%b}",
			want: []conditionalText{{
				fmt:  "%s",
				vars: []rune{'a'},
			}, {
				fmt:  "%s",
				vars: []rune{'b'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'b'},
					}},
				},
			}},
		},
		{
			fmt: "{%a}%b",
			want: []conditionalText{{
				fmt:  "%s",
				vars: []rune{'a'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'b'},
			}},
		},
		{
			fmt: "{%a}{%b}",
			want: []conditionalText{{
				fmt:  "%s",
				vars: []rune{'a'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'b'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'b'},
					}},
				},
			}},
		},
		{
			fmt: "{%a}|{%b}",
			want: []conditionalText{{
				fmt:  "%s",
				vars: []rune{'a'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'b'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'b'},
					}, {
						op: NAND,
						rs: []rune{'a'},
					}},
				},
			}},
		},
		{
			fmt: "{%a}|{%b}|{%c}",
			want: []conditionalText{{
				fmt:  "%s",
				vars: []rune{'a'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'b'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'b'},
					}, {
						op: NAND,
						rs: []rune{'a'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'c'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'c'},
					}, {
						op: NAND,
						rs: []rune{'a'},
					}, {
						op: NAND,
						rs: []rune{'b'},
					}},
				},
			}},
		},
		{
			fmt: "{%a%b}|{%c%d}",
			want: []conditionalText{{
				fmt:  "%s%s",
				vars: []rune{'a', 'b'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a', 'b'},
					}},
				},
			}, {
				fmt:  "%s%s",
				vars: []rune{'c', 'd'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'c', 'd'},
					}, {
						op: NAND,
						rs: []rune{'a', 'b'},
					}},
				},
			}},
		},
		{
			fmt: "{%a{%b{%c}|{%d}}|{%e}}|{%f}",
			want: []conditionalText{{
				fmt:  "%s",
				vars: []rune{'a'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'b'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a', 'b'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'c'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a', 'b', 'c'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'d'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a', 'b', 'd'},
					}, {
						op: NAND,
						rs: []rune{'c'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'e'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a', 'e'},
					}, {
						op: NAND,
						rs: []rune{'b'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'f'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'f'},
					}, {
						op: NAND,
						rs: []rune{'a'},
					}},
				},
			}},
		},
		{
			fmt: "{{%a}|{%b}}|{{%c}|{%d}}",
			want: []conditionalText{{
				fmt:  "%s",
				vars: []rune{'a'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'b'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'b'},
					}, {
						op: NAND,
						rs: []rune{'a'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'c'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'c'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'d'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'d'},
					}, {
						op: NAND,
						rs: []rune{'c'},
					}},
				},
			}},
		},
		{
			fmt: "{%a{%b{%c}}%d}{%e%f}",
			want: []conditionalText{{
				fmt:  "%s",
				vars: []rune{'a'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a', 'd'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'b'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a', 'b', 'd'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'c'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a', 'b', 'c', 'd'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'d'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a', 'd'},
					}},
				},
			}, {
				fmt:  "%s%s",
				vars: []rune{'e', 'f'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'e', 'f'},
					}},
				},
			}},
		},
		{
			fmt: "{%a{%b}|{%c}}|{%d}|{%e}",
			want: []conditionalText{{
				fmt:  "%s",
				vars: []rune{'a'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'b'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a', 'b'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'c'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'a', 'c'},
					}, {
						op: NAND,
						rs: []rune{'b'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'d'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'d'},
					}, {
						op: NAND,
						rs: []rune{'a'},
					}},
				},
			}, {
				fmt:  "%s",
				vars: []rune{'e'},
				co: condition{
					op: AND,
					cs: []condition{{
						op: AND,
						rs: []rune{'e'},
					}, {
						op: NAND,
						rs: []rune{'a'},
					}, {
						op: NAND,
						rs: []rune{'d'},
					}},
				},
			}},
		},
	}

	for i, testcase := range testcases {
		got, err := NestedBrackets(testcase.fmt)
		if len(got) != len(testcase.want) {
			t.Errorf("\nNestedBrackets(\"%s\"): test %d `%s`: incorrect result length\n\twant %d\n\t got %d",
				testcase.fmt, i, testcase.name, len(testcase.want), len(got))
		} else {
			for j := range got {
				if !got[j].Equals(testcase.want[j]) {
					t.Errorf("\nNestedBrackets(\"%s\"): test %d `%s`: incorrect result at index %d\n\twant %+v\n\t got %+v",
						testcase.fmt, i, testcase.name, j, testcase.want, got)
				}
			}
		}
		if err != testcase.err {
			t.Errorf("\nParse(\"%s\", tr): test %d `%s`: incorrect error\n\twant %s\n\t got %s",
				testcase.fmt, i, testcase.name, testcase.err, err)
		}
	}
}
