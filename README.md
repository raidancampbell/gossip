# gossip

The fourth installment in my lineage of IRC bots with varying intelligence.
1. [Swiggityspeare](https://github.com/raidancampbell/swiggityspeare): a Java bot with exec wrappers around [Karpathy's char-rnn](https://github.com/karpathy/char-rnn) LSTM neural network.  It was incredibly fragile.
1. [Stupidspeare](https://github.com/raidancampbell/stupidspeare): a Python reimplementation with none of the neural network intelligence.
1. [Sequelspeare](https://github.com/raidancampbell/sequelspeare): Stupidspeare, but with a Tensorflow implementation of `char-rnn`.

### Environment

- golang, the code was written with 1.15, but probably works back to 1.13
- relevant network details filled out in `config.yaml`
- [optional]: a splunk server listed to receive logs
