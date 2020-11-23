# direnv-gc
Extension to [direnv](https://direnv.net/) that keep tracks of your direnvs and allows you to clean up unused ones.

## Installing
MacOS and Linux binaries are available on the [releases](https://github.com/jonas-p/direnv-gc/releases) page.

You can also build and install it using Go:
```bash
go get -u github.com/jonas-p/direnv-gc
```

After installing and making sure it's available in your PATH, add `eval "$(direnv-gc hook)"` to your direnvrc
file (`~/.direnvrc` or `~/.config/direnv/direnvrc`). This will add a hook to direnv and update the database
everytime direnv loads a new environment.

## Usage
Run `direnv-gc` to clean up unused (10 days) environments.

```bash
$ direnv-gc
Removing /home/user/dev/project-a/.direnv (2f5175bc49c993c455a0371cff31797c41e5d350f2b549478367a9dd86941a31)
Removing /home/user/dev/project-b/.direnv (b5f140cd54a79c81e8b5083bd4245efe5d9dc7ff3df34cfd980ebb76e222b982)
Cleaned up 2 environments, saving a total of 10mb.
```

You can use the flag `--days` to specify how long an environment should have been inactive (not loaded) for before
removing it.

For additional flags see `direnv-gc --help`.


