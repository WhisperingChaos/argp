/*
Package argp reimplements the partitioning performed by a Linux command shell of its command line input into tokens.  However, it decouples reading input by accepting any source that implements the io.Reader interface.  Furthermore, it can be called anytime while executing a process and it decouples the tokenized output from os.Args (https://golang.org/pkg/os/#pkg-variables), so any array variable can accept the processed tokens.

After tokenizing input further processing is required to characterize each token as either an option (flag), an option value, or argument.  The go flag package (https://golang.org/pkg/flag/) offers this functionality.

Motivation

Enables development of a uniform console language that's consumable both when starting a process and during its execution.  This could be valuable, for example, to record and playback an interactive console conversation between an end user and the console process.  Therefore, instead of creating a different configuration file syntax, the console language would be used to configure the console.
*/
package argp

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

/*
Parse converts input into tokens.  There are two token types: unencapsulated and encapsulated.  Unencapsulated tokens are composed from consecutive non-whitespace characters while encapsulated tokens are character strings delimited by either single or double quotes which may contain whitespace. Although the scanner presents these two token types to the parser, the parser will aggregate these tokens into a single one unless they're separated by whitespace.  Each individual resulting parser token is appended, as a new element, of  the output args slice.

Besides tokenizing the input, the parser will remove the delimiting quotes around encapsulated tokens.  Additionally, when delimited by double quotes, the parser will remove the escape character before either an embedded double quote or another escape character.  No such processing is performed for encapsulated tokens delimited by single quotes.  Finally, for unencapsulated tokens, all instances of the escape character are replaced by the character that follows the escape.

The above semantics must mirror the token processing normally performed by the Linux command shell before it executes the specified command.  The only known deviation being args[0]. It contains an argument, not the name of the command specified in the terminal session.  Perhaps args[0] should reflect the calling function name instead - to preserve uniformity?

*/
func Parse(
	argLine io.Reader, // all input provided by io.Reader is considered a single line.
) (
	args []string, // tokenized output of entire io.Reader input.  Note args[0] is a "valid" argument - not a command name.
	err error,
) {
	if argLine == nil {
		return nil, fmt.Errorf("nil io.Reader passed to parser")
	}
	escDbl, _ := regexp.Compile(`(\\\\)|(\\")`)
	escAll, _ := regexp.Compile(`\\.`)

	var argCurr string
	s := cliConfig(argLine)
	for s.Scan() {
		// Note the scanner relies on the parser semantics below to improve its
		// performance. If these semantics change, the scanner may also
		// require coding changes.
		switch tokenIDextract(s.Bytes()) {
		case tArgument:
			argCurr += escapeSubstitute(tokenExtract(s.Bytes()), escAll)
		case tArgumentEncap:
			argCurr += tokenExpose(tokenExtract(s.Bytes()), escDbl)
		case tWhiteSpace:
			if argCurr != "" {
				args = append(args, argCurr)
			}
			argCurr = ""
		}
	}
	if s.Err() != nil {
		return nil, s.Err()
	}
	if argCurr != "" {
		args = append(args, argCurr)
	}
	return
}
func tokenIDextract(token []byte) (id byte) {

	return token[0]
}
func tokenExtract(token []byte) (tok []byte) {

	return token[1:]
}
func tokenExpose(tokenQuoted []byte, escape *regexp.Regexp) (tok string) {

	quote := tokenQuoted[0]
	tokenNew := tokenQuoted[1:]
	tokenNew = tokenNew[:len(tokenNew)-1]
	if quote == byte('"') {
		tokenNew = escape.ReplaceAllFunc(tokenNew, escapeReplace)
	}
	return string(tokenNew)
}
func escapeSubstitute(token []byte, escAll *regexp.Regexp) (tok string) {

	return string(escAll.ReplaceAllFunc(token, escapeReplace))
}
func escapeReplace(in []byte) []byte {

	return []byte{in[1]}
}

type tokenID byte

const (
	// start with lower case t so scope is private to package.
	tArgument = iota
	tArgumentEncap
	tWhiteSpace
)

func cliConfig(args io.Reader) (s *bufio.Scanner) {

	s = bufio.NewScanner(args)
	s.Split(scanConfig())
	return
}
func scanConfig() func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// The static regexps were tested below using the "go playground".
	// Therefore, compilation errors aren't checked for at runtime.  Testing
	// is available:  https://play.golang.org/p/xdW8_VDNXcT

	whiteSpaceComplete, _ := regexp.Compile(`^[[:space:]]+`)

	tokenNotEncapComplete, _ := regexp.Compile(`^([^'"\\[:space:]]|(\\.))+`)

	tokenEncapComplete, _ := regexp.Compile(`(^"(([^\\"])|(\\.))*")|(^'[^']*')`)
	tokenEncapPartial, _ := regexp.Compile(`(^"(([^\\"])|(\\.))*$)|(^'[^']*$)`)

	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for {
			if part := whiteSpaceComplete.Find(data); part != nil {
				// Note because the parser will do "nothing" when
				// more than one whitespace token is generated in the
				// situation where forcing a buffer read would only
				// produce a single whitespace, it's therefore, unnecessary
				// to force a buffer read.
				return tokenGen(tWhiteSpace, part)
			}
			if part := tokenNotEncapComplete.Find(data); part != nil {
				// Note parser will simply concatenate two consecutive
				// non-encapsulated tokens, therefore, logic to
				// force a buffer read, to potentially generate only
				// a single token isn't needed.
				// Note this regexp may partition a token on an escape
				// character at the very end of the current buffer. This
				// situation should not be a problem as this escape character
				// will occupy the first byte of the buffer likely followed
				// by any remaining text.
				return tokenGen(tArgument, part)
			}
			if part := tokenEncapComplete.Find(data); part != nil {
				return tokenGen(tArgumentEncap, part)
			}
			if !atEOF {
				if part := tokenEncapPartial.Find(data); part != nil {
					// a partial encapsulated token isn't "complete" as its
					// enclosing companion character hasn't been processed
					// before reaching the end of the current buffer.
					return 0, nil, nil
				}
			}
			if atEOF && len(data) < 1 {
				return 0, nil, nil
			}
			return len(data), data, fmt.Errorf("unable to tokenize: '%v'", data)
		}
	}
}
func tokenGen(id tokenID, token []byte) (advance int, tok []byte, err error) {
	tok = append(tok, byte(id))
	tok = append(tok, token...)
	return len(token), tok, err
}
