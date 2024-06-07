# Applicator

A simple programming language based off of assembly built to make pngs (and maybe gifs later!).

## Installation

```bash
git clone https://github.com/candiruRocks/applicator.git
cd applicator/
```

## Usage

```bash
go build main.go
./main [aplr file] [result png file]
./main helloWorld.aplr out.png
```

## Programming Syntax and Instructions

All files begin with the length and width of the final png. e.g.
```
200 500
```
for a 200 x 500 png.

Applicator does not use indentation or functions. Instead, like in ASM, you specify **label**s that
are jumped to either from a **goto** or a jump instruction. For more information on this style, check
out NASM.

There are only 3 types in Applicator, integers, colors, and strings. A color is 3 8 bit unsigned integers
with an inferred 255 on the alpha channel according to Go's [color package](https://pkg.go.dev/image/color#RGBA).

The user may define variables of each type. There are already some predefined variables that may be modified.
They are length, width, white, black, red, green, and blue.

When you define a variable, depending on the type, the value of the variable is returned into the arguments
of the current instruction. e.g.
```
bg white
```
is equal to
```
bg 255 255 255
```

### Instructions

<details>
<summary>bg</summary>
Makes the background the specified color
```
bg [color]
```
</details>

<details>
<summary>point</summary>
Makes a point (x, y) into the specified color
```
point [x int] [y int] [color]
```
</details>

<details>
<summary>box</summary>
Makes a box from point (x1, y1) to (x2, y2) with specified color
```
box [x1 int] [y1 int] [x2 int] [y2 int] [color]
```
</details>

<details>
<summary>var</summary>
Sets or creates a variable equal to specified int
```
var [variableName] [int]
```
</details>

<details>
<summary>cvar</summary>
Sets or creates a variable equal to specified color
```
cvar [variableName] [color]
```
</details>

<details>
<summary>svar</summary>
Sets or creates a variable equal to specified string
```
svar [variableName] [string]
```
</details>

<details>
<summary>color</summary>
Sets or creates a variable equal to the r, g, or b value of given color
```
color [variable int] [r|g|b] [color]
```
</details>

<details>
<summary>add</summary>
Adds specified int to variable
```
add [variable int] [int]
```
</details>

<details>
<summary>sub</summary>
Subtracts specified int from variable
```
sub [variable int] [int]
```
</details>

<details>
<summary>mult</summary>
Sets variable to variable multiplied by int
```
mult [variable int] [int]
```
</details>

<details>
<summary>neg</summary>
Sets variable to opposite sign
```
neg [variable int]
```
</details>

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
