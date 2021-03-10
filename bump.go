package main

import (
  "fmt"
  "os"
  "flag"
  "io"
)

const (
  filename = "VERSION"
  format = "%d.%d.%d%s\n"
)

func main() {

  var force = flag.Bool("f", false, "force create " + filename)
  var section = flag.String("b", "patch", "which section to bump: major or M, minor or m, patch or p")
  var commit = flag.String("c", "", "commit hash")
  var setman = flag.String("s", "", "set manually (overrides -b -c): M[.m[.p[-hash]]]")
  var major, minor, patch int
  var sections = map[string]*int{"major": &major, "M": &major, "minor": &minor, "m": &minor, "patch": &patch, "p": &patch}
  var err error
  var f *os.File
  var hash string

  flag.Parse()

  bumpme, valid := sections[*section]
  if ! valid {
    fmt.Fprintf(os.Stderr, "invalid bump: %s\n", *section)
    os.Exit(1)
  }

  if *force {
    f, err = os.Create(filename)
  } else {
    f, err = os.OpenFile(filename, os.O_RDWR, 0644)
  }
  if err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err.Error())
    os.Exit(1)
  }

  defer f.Close()

  if *setman != "" {
    fmt.Sscanf(*setman, format, &major, &minor, &patch, &hash)
  } else {
    fmt.Fscanf(f, format, &major, &minor, &patch, &hash)
    *bumpme++

    if *commit != "" {
      hash = "-" + *commit
    }
  }

  _, err = f.Seek(0,io.SeekStart)
  if err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err.Error())
    os.Exit(1)
  }

  fmt.Printf(format, major, minor, patch, hash)
  fmt.Fprintf(f, format, major, minor, patch, hash)
}
