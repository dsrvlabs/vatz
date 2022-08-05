## To get started with VATZ 
Installation&Run guide)

> This instruction is mainly based on Linux:Ubuntu

### 1. VATZ 
1. Clone vatz (public) repository with:
    ```
    $ git clone git@github.com:dsrvlabs/vatz.git
    ```
2. Compile VATZ
    ```
    $ cd vatz
    $ make
    ```
    you will see binary named `vatz`
3. Initialize VATZ
    ```
    $ ./vatz init
    ```
    You will see config file that named default.yaml once you initialize VATZ.  

### 2. VATZ Plugins
1. Clone vatz-plugin official repository <br> (**This repository can be changed.**)
    ```
    $ git clone git@github.com:dsrvlabs/vatz-plugin-common.git
    ```
2. Compile plugins
    ```
    $ cd plugins/cosmos-sdk-blocksync
    $ make
    ```
   You can see binary named `cosmos-sdk-blocksync`

---

You've created two binaries vatz and vatz-plugin to run.  

## Usage
- Modify default.yaml (VATZ)
    - plugin_name
        - refer to `vatz-plugin-common/plugins/cosmos-sdk-blocksync/main.go` of official repository
            - Search `pluginName`in vat-plugin  and update the default.yaml with `pluginName`
    - Port
        - Check your machine used port number. If you set it to a port number that is already in use, an error will occur.
    - discord_secret
        - If you want to receive discord messages, register webhook URL.

- Execute each binary as below <br> <img width="1510" alt="image" src="https://user-images.githubusercontent.com/106724973/180962626-c4de859c-526d-4038-87b8-badebed44136.png">
- Alert notification will be sent if block height is not increased <br> <img width="223" alt="image" src="https://user-images.githubusercontent.com/106724973/180962941-4db55b66-ceb6-4cd4-bd53-96d21bdf5596.png">
