package main

import (
  "os"
  "strings"
  "math"
  "strconv"
  "fmt"
  "image"
  "image/draw"
  "image/color"
  "image/png"
)

type DataType int

const (
  dataNum DataType = iota
  dataColor
  dataString
)

var vars map[string]int = make(map[string]int)
var colorVars map[string]color.RGBA = make(map[string]color.RGBA)
var stringVars map[string]string = make(map[string]string)
var cmpEqual bool = true
var cmpGreater bool = true

func main() {
  data, err := os.ReadFile(os.Args[1])
  check(err)

  unCheckedinput := strings.Split(string(data), "\n") 
  imgLength, err := strconv.Atoi(strings.Fields(unCheckedinput[0])[0]); check(err)
  imgWidth, err := strconv.Atoi(strings.Fields(unCheckedinput[0])[1]); check(err)

  populateVars(imgLength, imgWidth)
  var stack []int
  var jumpLine int
  unCheckedinput = unCheckedinput[1:]

  var input []string
  for _, line := range unCheckedinput {
    trimmedLine := strings.TrimSpace(line)
    if trimmedLine != "" && string(trimmedLine[0]) != ";" {
        input = append(input, trimmedLine)
      }
    } 
  
  
  img := image.NewRGBA(image.Rect(0, 0, imgLength, imgWidth))
  
  for lineNum := 0; lineNum < len(input); lineNum++ {
    tokens := strings.Split(input[lineNum], " ")
    switch tokens[0] {
      case "bg":
        apDraw(img, img.Bounds(), tokenColor(tokens, 1))
      case "point":
        img.Set(numHandler(tokens[1]), numHandler(tokens[2]), tokenColor(tokens, 3))
      case "box":
        apDraw(img, tokenRect(tokens, 1), tokenColor(tokens, 5))
      case "var":
        vars[tokens[1]] = numHandler(tokens[2])
      case "cvar":
        colorVars[tokens[1]] = tokenColor(tokens, 2)
      case "svar":
        if len(tokens) == 3 {
          stringVars[tokens[1]] = stringHandler(tokens[2])
        } else {
          strLength := strings.Index(input[lineNum], tokens[1])
          stringVars[tokens[1]] = input[lineNum][strLength + len(tokens[1]) + 1:]
        }
      case "color":
        // color redNum r red
        colorName := tokens[3]
        switch tokens[2] {
          case "r":
            vars[tokens[1]] = int(colorVars[colorName].R)
          case "g":
            vars[tokens[1]] = int(colorVars[colorName].G)
          case "b":
            vars[tokens[1]] = int(colorVars[colorName].B)
        }
      case "add":
        vars[tokens[1]] = vars[tokens[1]] + numHandler(tokens[2])
      case "sub":
        vars[tokens[1]] = vars[tokens[1]] - numHandler(tokens[2])

      case "mult":
        vars[tokens[1]] = vars[tokens[1]] * numHandler(tokens[2])
      case "neg":
        vars[tokens[1]] = -vars[tokens[1]]
      case "abs":
        vars[tokens[1]] = int(math.Abs(float64(vars[tokens[1]])))
      case "push":
        stack = append(stack, numHandler(tokens[1]))
      case "pop":
        vars[tokens[1]] = stack[len(stack) - 1]
        stack = stack[:len(stack)-1]
      case "cmp":
        cmpEqual = numHandler(tokens[1]) == numHandler(tokens[2])
        cmpGreater = numHandler(tokens[1]) > numHandler(tokens[2])
      case "goto":
        jumpLine = lineNum
        lineNum = gotoLine(input, tokens[1])
      case "je":
        if cmpEqual {
          jumpLine = lineNum
          lineNum = gotoLine(input, tokens[1])
        }
      case "jne":
        if !cmpEqual {
          jumpLine = lineNum
          lineNum = gotoLine(input, tokens[1])
        }
      case "jg":
        if cmpGreater && !cmpEqual {
          jumpLine = lineNum
          lineNum = gotoLine(input, tokens[1])
        }
      case "jl":
        if !cmpGreater && !cmpEqual {
          jumpLine = lineNum
          lineNum = gotoLine(input, tokens[1])
        }
      case "jge":
        if cmpGreater || cmpEqual {
          jumpLine = lineNum
          lineNum = gotoLine(input, tokens[1])
        }
      case "jle":
        if !cmpGreater || cmpEqual {
          jumpLine = lineNum
          lineNum = gotoLine(input, tokens[1])
        }
      case "return":
        lineNum = jumpLine
      case "print":
        if len(tokens) == 2 {
          switch detectType(tokens[1]) {
            case dataNum:
              fmt.Println(numHandler(tokens[1]))
            case dataColor:
              printColor := colorVars[tokens[1]]
              r := printColor.R 
              g := printColor.G 
              b := printColor.B 
              fmt.Printf("%d %d %d\n", r, g, b)
            case dataString:
              fmt.Println(stringHandler(tokens[1]))
          }
        } else {
          fmt.Println(input[lineNum][6:])
        }
      case "exit":
        defer os.Exit(numHandler(tokens[1]))
        break
    }
  }
  outFile, err := os.Create(os.Args[2])
  check(err)
  defer outFile.Close()

  err = png.Encode(outFile, img)
  check(err)
  fmt.Print()
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func apDraw(img *image.RGBA, bounds image.Rectangle, RGBcolor color.RGBA) {
  draw.Draw(img, bounds, &image.Uniform{RGBcolor}, image.ZP, draw.Src)
}

func tokenColor(tokens []string, offset int) color.RGBA {
  if val, ok := colorVars[tokens[offset]]; ok {
    return val
  }
  r := uint8(numHandler(tokens[offset]))
  g := uint8(numHandler(tokens[offset + 1]))
  b := uint8(numHandler(tokens[offset + 2]))
  return color.RGBA{r, g, b, 255}
}

func tokenRect(tokens []string, offset int) image.Rectangle {
  point1 := numHandler(tokens[offset])
  point2 := numHandler(tokens[offset + 1])
  point3 := numHandler(tokens[offset + 2])
  point4 := numHandler(tokens[offset + 3])

  return image.Rect(point1, point2, point3, point4)
}

func numHandler(numString string) int {
  if val, ok := vars[numString]; ok {
    return val
  }
  val, err := strconv.Atoi(numString); check(err)
  return val
}

func stringHandler(str string) string {
  if val, ok := stringVars[str]; ok {
    return val
  }
  return str
}

func detectType(unkownType string) DataType {
  if _, ok := vars[unkownType]; ok {
    return dataNum
  }
  if _, ok := colorVars[unkownType]; ok {
    return dataColor
  }
  return dataString
}

func populateVars(length, width int) {
  vars["length"] = length
  vars["width"] = width

  colorVars["white"] = makeColor(255, 255, 255)
  colorVars["black"] = makeColor(0, 0, 0)
  colorVars["red"]   = makeColor(255, 0, 0)
  colorVars["green"] = makeColor(0, 255, 0)
  colorVars["blue"]  = makeColor(0, 0, 255)
}

func makeColor(r, g, b uint8) color.RGBA {
  return color.RGBA{r, g, b, 255}
}

func gotoLine(input []string, gotoString string) int{
  for gotoLineNum := range input {
    currLine := strings.Fields(input[gotoLineNum])
      if currLine[0] == "label" && currLine[1] == gotoString {
        return gotoLineNum 
      }
  }
  fmt.Printf("Could not find label %s\n", gotoString)
  os.Exit(1)
  return len(input) - 1
}
