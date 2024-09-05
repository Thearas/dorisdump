package parser

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

var (
	// The properties which value contains identifier.
	propertiesWithValueIds = lo.SliceToMap([]string{
		"bloom_filter_columns",
		"function_column.sequence_col",
	}, func(s string) (string, struct{}) { return s, struct{}{} })
)

func NewListener(hideSqlComment bool, modifyIdentifier func(id string, ignoreBuiltin bool) string) DorisParserListener {
	return &listener{hideSqlComment: hideSqlComment, modifyIdentifier: modifyIdentifier}
}

func NewErrListener(sqlId string) antlr.ErrorListener {
	return &errListener{ConsoleErrorListener: antlr.NewConsoleErrorListener(), sqlId: sqlId}
}

func NewParser(sqlId string, sqls string, listeners ...antlr.ParseTreeListener) *DorisParser {
	input := antlr.NewInputStream(sqls)
	lexer := NewDorisLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := NewDorisParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(NewErrListener(sqlId))

	for _, listener := range listeners {
		p.AddParseListener(listener)
	}

	return p
}

type errListener struct {
	*antlr.ConsoleErrorListener
	sqlId string
}

func (l *errListener) SyntaxError(_ antlr.Recognizer, _ any, line, column int, msg string, _ antlr.RecognitionException) {
	logrus.Errorf("sql %s parse error at line %d:%d %s\n", l.sqlId, line, column, msg)
}

type listener struct {
	*BaseDorisParserListener
	hideSqlComment   bool
	modifyIdentifier func(id string, ignoreBuiltin bool) string
}

func (l *listener) modifySymbolText(node antlr.TerminalNode, ignoreBuiltin bool) {
	symbol := node.GetSymbol()

	id := strings.Trim(symbol.GetText(), "`")
	symbol.SetText(l.modifyIdentifier(id, ignoreBuiltin))
}

func (l *listener) ExitUnquotedIdentifier(ctx *UnquotedIdentifierContext) {
	ignoreBuiltin := true
	child := ctx.GetChild(0)
	_, ok := child.(*NonReservedContext)
	if ok {
		// ignoreBuiltin = true
		// child = nonReserved.GetChild(0)
		return
	}
	l.modifySymbolText(child.(*antlr.TerminalNodeImpl), ignoreBuiltin)
}

func (l *listener) ExitQuotedIdentifier(ctx *QuotedIdentifierContext) {
	child := ctx.GetChild(0)
	l.modifySymbolText(child.(*antlr.TerminalNodeImpl), false)
}

func (l *listener) ExitPropertyItem(ctx *PropertyItemContext) {
	// e.g. "bloom_filter_columns" = "col1,col2"
	key := strings.Trim(ctx.GetKey().GetText(), `'"`)
	if _, ok := propertiesWithValueIds[key]; !ok {
		return
	}

	pvalue := ctx.PropertyValue()
	if pvalue.Constant() != nil {
		constant := pvalue.Constant()
		rawText := constant.GetText()
		quote := rawText[0]

		ids := strings.Split(rawText[1:len(rawText)-1], ",")
		for i, id := range ids {
			ids[i] = l.modifyIdentifier(strings.Trim(id, "`"), false)
		}

		symbol := constant.GetChild(0).(*antlr.TerminalNodeImpl).GetSymbol()
		symbol.SetText(fmt.Sprintf("%c%s%c", quote, strings.Join(ids, ","), quote))
	}
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
