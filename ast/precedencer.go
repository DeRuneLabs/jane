// Copyright (c) 2024 - DeRuneLabs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package ast

type precedencer struct {
	pairs [][]any
}

func (p *precedencer) set(level uint, expr any) {
	for i, pair := range p.pairs {
		pair_level := pair[0].(uint)
		if level > pair_level {
			first := p.pairs[:i]
			appended := append([][]any{{level, expr}}, p.pairs[i:]...)
			p.pairs = append(first, appended...)
			return
		}
	}
	p.pairs = append(p.pairs, []any{level, expr})
}

func (p *precedencer) get_lower() any {
	for i := len(p.pairs) - 1; i >= 0; i-- {
		data := p.pairs[i][1]
		if data != nil {
			return data
		}
	}
	return nil
}
