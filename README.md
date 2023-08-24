<h1 align="center">
CIGO
</h1>

<h4 align="center">A compiled CLI tool to help manage a <a href="https://en.wikipedia.org/wiki/Monorepo">monorepo</a> and automate CI/CD operations.</h4>

<p align="center">
  <a href="https://gitlab.cas.mcmaster.ca/alsaboaa/monorepo/-/graphs/main/charts">
    <img src="https://gitlab.cas.mcmaster.ca/alsaboaa/monorepo/badges/main/coverage.svg"
         alt="Coverage">
  </a>
  <a href="https://gitlab.cas.mcmaster.ca/alsaboaa/monorepo/-/pipelines">
    <img src="https://gitlab.cas.mcmaster.ca/alsaboaa/monorepo/badges/main/pipeline.svg"
         alt="Pipeline">
  </a>
</p>

## Directories

The documentation is stored under [docs/](docs), no PDFs, and the program code is under [src](src).

## Key Features

1. [x] Run project targets
2. [x] Manage project dependencies
3. [x] Manage target execution dependencies
4. [x] Detect updated projects
5. [x] Search projects
6. [x] List registered projects
7. [x] Cache target execution output
8. [x] Generate json schema for the JSON files
9. [x] A wizard to add a project to the workspace

## How To Use

The tool made here is a command line tool targeted at Unix Machines. It is written in Go, so you need to have Go installed on your machine to run the program, You can download it from [here](https://golang.org/dl/). We are using golang version 1.19.5.

The program can be compiled by first navigating to where the source code is located, [/src/](src) and using the command `go build`, which should produce a `cigo` binary.

You can call `./cigo -h` to get the help output.

`./cigo -h` should output the following:
```bash
Usage:
  cigo [OPTIONS] <command>

Application Options:
  -d, --dry      Print commands to run in order without running
                 anything.
  -V, --version  Print version

Help Options:
  -h, --help     Show this help message

Available commands:
  add-project
  create-schema
  get-changed
  list
  run
  search
```

Feel free to add the program to your path.

⚠️ **Warning:** This program is targeted to Unix machines. Windows is not supported. If you are using Windows, you can use [WSL](https://learn.microsoft.com/en-us/windows/wsl/install) to run the program.

### Sample Usage

This repo is already setup with the files that we need to run the program. In the root directory, we have our `workspace.json`
file which contains information about our workspace. We also have sample projects in the `/apps` directory, each setup with their own `project.json` files.

**Compiling the program:**

First navigate to wherever the source code is located.

`cd /src/`

Use the command `go build` to generate the `cigo` binary. A new file should appear in your current directory.

`go build`

You can then use the compiled tool to manage your monorepo.

**Creating schemas to test your workspace and project files against:**


Let's create JSON schema's that help us validate our workspace and project files. To create these schemas we
call the `create-schema` command, passing it the type of schema we want to create (`workspace' or `project').

`./cigo create-schema workspace`

**Searching for specific information in the projects of the workspace:**

You can search all the projects with key:value pairs by using the following command:

`./cigo search mainLanguage:cpp`

**Listing all the projects in the workspace:**

We have a list command that lists all the projects in the workspace:

`./cigo list`

**Getting changes in the workspace:**

This command will get all the projects that have been changed since the last commit. You can specify the branch or commit hash and the head commit to compare. If you don't specify -h , it will use the current HEAD commit.

`./cigo get-changed -b main -h HEAD`

**Running a projects Target:**

To run a target, you need to specify the project name and the target name. The target is what will be run, as it contains the commands that will be executed.

To specify a project, you need to use the flag `-p` or `--project` and to specify a target, you need to use the flag `-t` or `--target`.

`./cigo run -p proj_a -t build`

The command above would run the `build` target in the `proj_a` project.

**Adding a project to the workspace:**

To add a project to the workspace, you can use the `add-project` command. It will prompt you with questions regarding the project you want to add, and then will update the workspace and create the `project.json` file for you.

`./cigo add-project`


**Extra Information:**

If you need more information about the commands, you can use the help flag
`./cigo -h`

If you need more information about a specific command, you can use the help flag with the command name, for example to get help about the `run` command, you can use the following command:

`./cigo run -h`


## Credits
### Contributors

* Ahmed Ammar Al-Sabounchi (alsaboaa@mcmaster.ca)
* Ali Ahmed Khan (khana238@mcmaster.ca)
* Omar Alkersh (alkersho@mcmaster.ca)
* Tanveer Shakeel (shakeelt@mcmaster.ca)
* Usman Asad (asadu@mcmaster.ca)

## Related
There are other products that have the same purpose, but we couldn't find one that does exactly what we wanted to do. The relevant projects are:

1. [https://nx.dev/](NX)
2. [https://bazel.build/](bazel)
3. [https://rushjs.io/](Rushjs)
