//  deploy  -  Deploy utility.  Deploy a set of files across a set of servers.
//
//  E.B.Smith  -  November, 2014


package main


import (
	"fmt"
	"bufio"
	"unicode"
	"unicode/utf8"
	"errors"
//	"strconv"
	)


//	Parse an identifier -- 

func IsValidIdentifierStartCharacter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
	}


func IsValidIdentifierCharacter(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_';
	}


func ScanIdentifier(data []byte, atEOF bool) (advance int, token []byte, err error) {

	// Skip leading spaces.
	var (
		r rune 
		width int = 0
		start int = 0
		)
	for width = 0; start < len(data); start += width {
		r, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(r) {
			break
			}
		}
	if atEOF && len(data) == 0 {
		return 0, nil, nil
		}

	// Make sure we're pointer at a valid identifier character.
	r, width = utf8.DecodeRune(data[start:])
	if !IsValidIdentifierStartCharacter(r) {
		oldAdvance := advance
		advance, token, _ = bufio.ScanWords(data, atEOF)
		advance += oldAdvance
		err = errors.New( fmt.Sprintf("Identifier expected. Scanning '%s'.", token) )
		return advance, token, err
		}

	// Scan while in identifier characters.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if !IsValidIdentifierCharacter(r) {
			return i + width, data[start:i], nil
			}
		}

	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
		}

	// Request more data.
	return 0, nil, nil
	}


func ScanManifest(data []byte, atEOF bool) (advance int, token []byte, err error) {

    advance, token, err = ScanIdentifier(data, atEOF)
    log(DULogDebug, "Advance: %d Token: %s Err: %v", advance, token, err)
	// if err == nil && token != nil {
	// 	_, err = strconv.ParseInt(string(token), 10, 32)
	// 	}
	return
    }











