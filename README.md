# gossip

The fourth installment in my lineage of IRC bots with varying intelligence.
1. [Swiggityspeare](https://github.com/raidancampbell/swiggityspeare): a Java bot with exec wrappers around [Karpathy's char-rnn](https://github.com/karpathy/char-rnn) LSTM neural network.  It was incredibly fragile.
1. [Stupidspeare](https://github.com/raidancampbell/stupidspeare): a Python reimplementation with none of the neural network intelligence.
1. [Sequelspeare](https://github.com/raidancampbell/sequelspeare): Stupidspeare, but with a Tensorflow implementation of `char-rnn`.
1. gossip: a Golang reimplementation.  No neural networks again, this incarnation focuses on implementing  an IRC bot without an IRC framework library.  [sorcix/irc.v2](https://github.com/sorcix/irc/tree/v2) is used for protocol parsing.


### Features

The usual array of IRC bot features
- karma, maintained through the generated `gossip.db` sqlite database
- `!ping` for liveness checks
- `!part` and `!die` support for leaving, and `!toggle` support to enable/disable features. authorized via configuration file
- HTML title text extraction
- a handful of other minor features, listed in the `gossip` package

### Environment

- golang, the code was written with 1.15, but probably works back to 1.13
- relevant network details filled out in `config.yaml`
- [optional] a splunk server with an HTTP Event Collector configured to receive logs
