# Easy Git Modules Manager

## Readable slow install (parallel install off)

<img src="demo-slow.gif" width="700px" />


## Fast install (default, parallel install on)

<img src="demo-fast.gif" width="700px" />


## How it works

1. The list of git repositories to clone is read from the config file (`package.json` by default, can be changed) from `dependencies` and `devDependencies` fields.
2. All git repositories are cloned into the specified folder (`src/modules` by default, can be changed)
3. Before cloning a repo, the local folder with this repo is deleted if it exists. 
4. If a repository has unsaved changes, it is skipped so you won't lose any changes


## How to use

> If you're on mac, you can just use already compiled `mod` file from this repo.
> If you're not on mac:
> 1. Clone this repo
> 2. Install go
> 3. Run `go build -o mod` inside the cloned repo to generate a compiled file for your operating system (instead of `mod` you can use any filename)

1. Copy the `mod` file with the app to the place where you want to manage multiple git repositories.
2. Copy `go.env` file from this repository to the same place. It contains the environment variables required for app to run.
3. You can change any environment variables you want in `go.env` + the app also reads a file `go.env.local`
4. `!Important!` If you're cloning git repositories with ssh - set the path to your ssh key in the env config file.

## Env config example

```.env
CONFIG_FILE="package.json"
MODULES_DIR="src/modules"
SSH_KEY_PATH="/Users/your_username/.ssh/id_rsa"
SSH_KEY_PASSWORD=""
```

`CONFIG_FILE` - the name of the JSON file from which the list of git repositories to clone is read. Inside, the fields `dependencies` and `devDependencies` are used (for compatibility with standart package.json format)<br>
`MODULES_DIR` - relative path to the directory to which git repositories will be cloned. Don't forget to add it to `.gitignore`<br>
`SSH_KEY_PATH` - absolute(!) path to your local private ssh key<br>
`SSH_KEY_PASSWORD` - your ssh key password (if there's none - leave empty)<br>

## Available commands

```bash
    # Help. List all available commands
    ./mod -h
    # Run the app
    ./mod 
    # Run the app without cloning repositories in parallel. Much slower, but better log readability.
    ./mod -parallel-install=false
    # Forcibly delete the whole directory with your git repositories before running the app.
    ./mod -safe-install=false 
    # See the list of git repositories that have unsaved local changes. No repositories are cloned when using this flag
    ./mod -show-changed-modules=true 
```

## Config file example

1. Repository `test-module-js` with the name `test`. Normal ssh link
2. Repository `storefront-query-builder`. Normal https link
3. Repository `scrape-search-ai`. Https link with the branch `dev`
4. Repository `Magnific-Popup` with the name `magnific-popup`. Ssh link with the tag `1.8.0` (the app uses `\d(\..*)+` regex for determining that it's a tag). 
5. Repository `Course_Bash-Programming` with the name `Course_Bash`. Ssh link with a specific commit hash

```json
{
   "devDependencies": {
      "test": "git@github.com:SergeyDarn/test-module-js.git",
      "storefront-query-builder": "https://github.com/DivanteLtd/storefront-query-builder.git",
   },
   "dependencies": {
      "scrape-search-ai": "https://github.com/SergeyDarn/scrape-search-ai#dev",
      "magnific-popup": "git@github.com:SergeyDarn/Magnific-Popup.git#1.8.0",
      "Course_Bash": "git@github.com:SergeyDarn/Course_Bash-Programming.git#3a7a19020151b45a29896c9142723efe5b11a061",
   }
}
```
