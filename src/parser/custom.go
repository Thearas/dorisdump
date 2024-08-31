package parser

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

func NewListener(hideSqlComment bool, modifyIdentifier func(id string, ignoreBuiltin bool) string) DorisParserListener {
	return &listener{hideSqlComment: hideSqlComment, modifyIdentifier: modifyIdentifier}
}

func NewParser(sqls string, listeners ...antlr.ParseTreeListener) *DorisParser {
	input := antlr.NewInputStream(sqls)
	lexer := NewDorisLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := NewDorisParser(stream)

	for _, listener := range listeners {
		p.AddParseListener(listener)
	}

	return p
}

type listener struct {
	*BaseDorisParserListener
	hideSqlComment   bool
	modifyIdentifier func(id string, ignoreBuiltin bool) string
}

func (l *listener) ExitUnquotedIdentifier(ctx *UnquotedIdentifierContext) {
	ignoreBuiltin := true
	child := ctx.GetChild(0)
	nonReserved, ok := child.(*NonReservedContext)
	if ok {
		ignoreBuiltin = false
		child = nonReserved.GetChild(0)
	}
	symbol := child.(*antlr.TerminalNodeImpl).GetSymbol()

	id := symbol.GetText()
	symbol.SetText(l.modifyIdentifier(id, ignoreBuiltin))
}

func (l *listener) ExitQuotedIdentifier(ctx *QuotedIdentifierContext) {
	child := ctx.GetChild(0)
	symbol := child.(*antlr.TerminalNodeImpl).GetSymbol()

	id := strings.Trim(symbol.GetText(), "`")
	symbol.SetText(l.modifyIdentifier(id, false))
}

// SimpleColumnDefContext & ColumnDefContext
func (l *listener) ExitSimpleColumnDef(ctx *SimpleColumnDefContext) {
	if l.hideSqlComment {
		hideComment(ctx, ctx.GetComment())
	}
}

func (l *listener) ExitColumnDef(ctx *ColumnDefContext) {
	if l.hideSqlComment {
		hideComment(ctx, ctx.GetComment())
	}
}

func hideComment(ctx antlr.ParserRuleContext, comment antlr.Token) {
	if comment == nil {
		return
	}
	text := comment.GetText()
	if len(text) <= len(`''`) {
		// empty comment
		return
	}

	newText := fmt.Sprintf(`'%s'`, strings.Repeat("*", len(text)))
	c := ctx.GetChild(ctx.GetChildCount() - 1)
	c.(*antlr.TerminalNodeImpl).GetSymbol().SetText(newText)
}

func (p *DorisParser) ToSQL() string {
	// parser and modify
	ms := p.MultiStatements()

	// get modified sql
	interval := antlr.NewInterval(ms.GetStart().GetTokenIndex(), ms.GetStop().GetTokenIndex())
	s := p.GetTokenStream().GetTextFromInterval(interval)

	return s
}
