package ui

import "strings"

func helpText() string {
	return strings.TrimSpace(`Navigation
  <j/k>       move selection
  <enter>     shell into selected running instance
  <r>         refresh
  </>         search instances

Actions
  <e>         edit instance config
  <s>         stop running instance
  <S>         start stopped instance
  <d>         delete stopped instance
  <l>         view logs
  <c>         view console logs

Context
  <R>         switch remote
  <p>         switch project

General
  <?>         help
  <esc/q>     back
  <ctrl+c>    quit
  <q>         quit from instance table`)
}
