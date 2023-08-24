# Documents
This section contains all of the documents pertaining to the project. 

## Error's Compiling

### Glossary
If your glossary isn't showing up, it could be because your `{fileName}.gls` isn't being created by the `\makeglossaries` command. This could happen because of the order that you are trying to compile the document in.

#### Solution
The recipe for compiling that worked for me is `makeglossaries -> pdflatex`. You really just need to run `makeglossaries` before using your latex pdf compiler. In VSCode click the TeX icon on the explorer thing, and under `Build LaTeX project select the recepie you added.  _I spent way too much time on this lol._

If you don't want to think about it, just copy paste this into `latex-workshop` settings in `VSCode` (if you dont want to copy the entire section, just add the `makeglossaries` tool to your latex tools section and the `makeglossaries -> pdflatex` recipe to your latex recipes):
- ```
     "latex-workshop.latex.recipes":[
        {
            "name": "latexmk",
            "tools": [
                "latexmk"
            ]
        },
        {
            "name": "makeglossaries -> pdflatex",
            "tools": [
                "makeglossaries",
                "pdflatex"
            ]
        }
    ],
    "latex-workshop.latex.tools":[
        {
            "name": "latexmk",
            "command": "latexmk",
            "args": [
                "-shell-escape",
                "-synctex=1",
                "-interaction=nonstopmode",
                "-file-line-error",
                "-pdf",
                "-outdir=%OUTDIR%",
                "%DOC%"
            ]
        },
        {
            "name": "pdflatex",
            "command": "pdflatex",
            "args": [
                "-synctex=1",
                "-interaction=nonstopmode",
                "-file-line-error",
                "%DOC%"
            ]
        },
        {
            "name": "makeglossaries",
            "command": "makeglossaries",
            "args": [
              "%DOCFILE%"
            ]
          }
    ],```
