# Petit Chat 🐈‍⬛

The program enables peer-to-peer chat which doesn't require a centralized server
to relay messages between clients. It is implemented based on
[libp2p](https://github.com/libp2p/go-libp2p) for network layers and written in
[Go](https://golang.org). Basic ideas start from "Nobody wants his/her private
messages get relayed and stored by untrustful third-party service providers even
though they advertise the services' concrete security protocol".  In other
words, any kind of data transmitted between clients remain on their side, not on
the side of middle men who could take off masks and reveal real face of 'Big
Brother'.

## Road Map
- [x] Build Hierarchy
- [x] Make packets flow anyway
- [x] Encapsulate messages to include metadata
- [x] Advance MsgBox
- [ ] Empower MsgBox
  - [x] Encrypt msg
  - [x] Sign & Verify msg
  - [x] Enhance msgTypeHello
  - [x] Add msgBox authentication
  - [x] Refactor msg type, struct, interface, etc.
  - [x] Enhance msgBox internal metadata
  - [x] Make msgBox only relay msgs
  - [ ] Advance msgBox store
- [ ] Conduct sub-tasks
  - [ ] Prettify CLI
  - [ ] Increase test code coverage up to 80%
  - [ ] Inspect open source licenses
  - [ ] Update go modules to up to date
  - [ ] Implement logger
- [ ] Release `v0.0.1`