# Git- Push
#### `git-push` is a cli tool for absolute beginner of github and git to push repo on github by no efforts.

![git-push gif](https://storage.googleapis.com/git-push/gif/git-push.gif)

- If you want to build git-push right away ,you need a working [Go environment](https://golang.org/doc/install). 
```
$ go get github.com/urvil38/git-push
$ go install
$ go build
```

# Downloads
### Linux and Macos:
1. For Linux and macos user you download git-push with `curl`:

    #### Linux:

    ```
    $ curl -LO https://storage.googleapis.com/git-push/v0.2/linux/git-push 
    ```

    #### MacOs:

    ```
    $ curl -LO https://storage.googleapis.com/git-push/v0.2/darwin/git-push
    ```

2. Make binary executable :
    ```
    $ chmod +x ./git-push
    ```

3. move binary file to bin directory :

    Typically these commands must be run as root or through `sudo`.
    ```
    $ mv git-push /usr/local/bin
    ```

### Windows:

You can download this binary from this url using your browser:

Go to this URL for download git-push:

```
https://storage.googleapis.com/git-push/v0.2/windows/git-push.exe
```

## Prerequisite:

1. You need to set `HOME` environment varible where git-push store it's configuration:

    - For linux and macos:

        Typically these environment variable is already set most of the time.You can confirm this using following command:
        ```
        $ echo $HOME
        ```
        If it's return nothing then set HOME environment varible by following command:
        ```
        $ export $HOME=/path/to/home/where/git-push/can/store/credentials
        ```

    - For windows:

        you must set the HOME environment variable to your chosen path(I suggest `c:\git-push`)
            
        Under Windows, you may set environment variables through the "Environment Variables" 
        button on the "Advanced" tab of the "System" control panel. Some versions of Windows 
        provide this control panel through the "Advanced System Settings" option inside 
        the "System" control panel.

2. You also need `git` install on your machine.If you don't have download from [here](https://git-scm.com/downloads)    