A Git implementation written in Go. Inspired and aided by [gitlet.js](http://gitlet.maryrosecook.com/) and [pygit](http://benhoyt.com/writings/pygit/). Undertaken in an effort to learn more about the inner workings of Git.

```
An implementation of Git written in Go.

Usage:
  mhgit [command]

Available Commands:
  add          Add file contents to the index
  cat-file     Provide content or type and size information for repository objects.
  commit       Record changes to the repository
  hash-object  Compute object ID and optionally creates a blob from a file.
  help         Help about any command
  init         Create an empty Git repository or reinitialize an existing one.
  ls-files     Show information about files in the index and the working tree
  rm           Remove files from the working tree and from the index
  status       Show the working tree status
  update-index Register file contents in the working tree to the index.
  write-tree   Create a tree object from the current index

Flags:
      --config string   config file (default is $HOME/.mhgit.yaml)
  -h, --help            help for mhgit

Use "mhgit [command] --help" for more information about a command.
```