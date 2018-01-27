# Git- Push
### **git-push** is a cli tool for absolute beginner of github and git to push repo on github by no efforts.

- If you want to build git-push right away ,you need a working [Go environment](https://golang.org/doc/install). 
```
$ go get github.com/urvil38/git-push
$ go install
$ go build
```

# Downloads

1. For Linux and macos user you download git-push with `curl`:

### For Linux:

```
$ curl -LO https://storage.googleapis.com/git-push/v0.1/linux/git-push 
```

### For MacOs:

```
$ curl -LO https://storage.googleapis.com/git-push/v0.1/darwin/git-push
```

2. Make binary executable :
```
$ chmod +x ./git-push
```

3. move binary file to bin directory ( For linux and macos users ) :

>Typically these commands must be run as root or through `sudo`.
```
$ mv git-cli /usr/local/bin
```

## For windows:

You can download this binary from this url using your browser:

> Go to this URL for download git-cli:

```
https://storage.googleapis.com/git-push/v0.1/windwos/git-push.exe
```

## Prerequisite:

You need to set HOME environment varible where git-cli store it's configuration:

- For linux and macos:

>Typically these environment variable is already set most of the time.You can confirm this using following command:
```
$ echo $HOME
```
>If it's return nothing then set HOME environment varible by following command:
```
$ export $HOME=/path/to/home/where/git-push/can/store/credentials
```

- For windows:

you must set the HOME environment variable to your chosen path(I suggest `c:\git-push`)
	
Under Windows, you may set environment variables through the "Environment Variables" 
button on the "Advanced" tab of the "System" control panel. Some versions of Windows 
provide this control panel through the "Advanced System Settings" option inside 
the "System" control panel.