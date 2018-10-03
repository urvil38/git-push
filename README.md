# Git- Push
#### `git-push` is a cli tool help to initialize and push git repository to github or gitlab or bitbucket by no efforts.

[![asciicast](https://asciinema.org/a/rbv3Js0OYOdsqyoCNqJUotzdo.png)](https://asciinema.org/a/rbv3Js0OYOdsqyoCNqJUotzdo)

- If you want to build git-push right away ,you need a working [Go environment](https://golang.org/doc/install) and [dep](https://github.com/golang/dep).
```
$ git clone https://github.com/urvil38/git-push.git
$ dep ensure
$ go install
```

# Downloads
### Linux and Macos:
1. For Linux and macos user you download git-push with `curl`:

    #### Linux:

    ```
    $ curl -LO https://storage.googleapis.com/git-push-bin/v1.0/linux/git-push 
    ```

    #### MacOs:

    ```
    $ curl -LO https://storage.googleapis.com/git-push-bin/v1.0/darwin/git-push
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
https://storage.googleapis.com/git-push-bin/v1.0/windows/git-push.exe
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
        $ export HOME=/path/to/home/where/git-push/can/store/credentials
        ```

    - For windows:

        you must set the HOME environment variable to your chosen path(I suggest `c:\git-push`)

        You can do it by two way in windows:

        -  Using `Command Prompt` you can set this environment variable by following command:
            
            ```
            set HOME="c:\git-push"
            ```    
        -  Under Windows, you may set environment variables through the "Environment Variables" 
            button on the "Advanced" tab of the "System" control panel. Some versions of Windows 
            provide this control panel through the "Advanced System Settings" option inside 
            the "System" control panel.


2. You also need `git` install on your machine.If you don't have download from [here](https://git-scm.com/downloads).  
