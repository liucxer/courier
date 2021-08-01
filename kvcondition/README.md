# KV Condition

[![GoDoc Widget](https://godoc.org/github.com/liucxer/courier/kvcondition?status.svg)](https://godoc.org/github.com/liucxer/courier/kvcondition)
[![Build Status](https://travis-ci.org/go-courier/kvcondition.svg?branch=master)](https://travis-ci.org/go-courier/kvcondition)
[![codecov](https://codecov.io/gh/go-courier/kvcondition/branch/master/graph/badge.svg)](https://codecov.io/gh/go-courier/kvcondition)
[![Go Report Card](https://goreportcard.com/badge/github.com/liucxer/courier/kvcondition)](https://goreportcard.com/report/github.com/liucxer/courier/kvcondition)


Simple query language to describe condition for label

* equal rules:
  * value equal `key="value"` 
  * value not equal `key!="value"` 
  * value ends with `key$="value"` 
  * value starts with  `key^="value"` 
  * value contains `key*="value"` 
  * has key `key` 
* compose rules: 
  * and `&` 
  * or `|` 
  * hoisting operator precedence  `()` 

