## Install from sources (Ubuntu)
- Clone vatz (public) repository with:
```
$ git cone git@github.com:dsrvlabs/vatz.git
```
Then compile:
```
$ cd vatz
$ make
```
you can see binary named `vatz`

- Clone vatz-plugin official repository with: (**This repository can be changed.**)
```
$ git clone git@github.com:dsrvlabs/vatz-plugin-common.git
```
Then compile:
```
$ cd plugins/cosmos-sdk-blocksync
$ make
```
You can see binary named `cosmos-sdk-blocksync`
Until now, you have made two binaries vatz and vatz-plugin.

## Usage
- Modify default.yaml (vatz)
    - plugin_name
        - refer to `vatz-plugin-common/plugins/cosmos-sdk-blocksync/main.go` of official repository
            - Search `pluginName`in vat-plugin  and update the default.yaml with `pluginName`
    - Port
        - Check your machine used port number. If you set it to a port number that is already in use, an error will occur.
    - discord_secret
        - If you want to receive discord messages, register webhook URL.

- Execute two binaries each other like below: <br> <img width="1510" alt="image" src="https://user-images.githubusercontent.com/106724973/180962626-c4de859c-526d-4038-87b8-badebed44136.png">
- If block height is not increased, <br> <img width="223" alt="image" src="https://user-images.githubusercontent.com/106724973/180962941-4db55b66-ceb6-4cd4-bd53-96d21bdf5596.png">
