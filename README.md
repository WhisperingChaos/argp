# argp
Implements the partitioning performed by a Linux command shell of its command line input into tokens. However,
it decouples reading input by accepting any source that implements the io.Reader interface. Furthermore,
it can be called anytime while executing a process and it decouples the tokenized output from [os.Args](https://golang.org/pkg/os/#pkg-variables), so any array variable can accept the processed tokens.

After tokenizing input further processing is required to characterize each token as either an option (flag), an option value, or argument.  The go [flag package](https://golang.org/pkg/flag/) offers this functionality.

### Install
```go get github.com/WhisperingChaos/argp```

### Motivation

Enables development of a uniform console language that's consumable both when starting a process and during its execution.  This could be valuable, for example, to record and playback an interactive console conversation between an end user and the console process.  Therefore, instead of creating a different configuration file syntax, the console language would be used to configure the console.
