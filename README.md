# VATZ

> This is a pilot for initial commit in development branch. 
>
> Please, follow the instructions in gitbook to run [pilot version of Vatz](https://app.gitbook.com/o/-MiyxU38etxprZxihBZS/s/-Mj3CwiN6vyRfTZC-Ljw/general/contents/pilot)

{% hint style="info" %}

This branch will be used only on purpose of development, or to create a release candidate

{%  endhint %}


## Configuration for private repo access


Ref1: https://www.digitalocean.com/community/tutorials/how-to-use-a-private-go-module-in-your-own-project
Basic guide.

Ref2: https://stackoverflow.com/questions/69682030/how-to-go-get-private-repos-using-ssh-key-auth-with-password
We have to change the way go get invokes ssh, by disabling batch mode.

```
$ export GOPRIVATE=github.com/hqueue/vatz-secret
$ env GIT_SSH_COMMAND="ssh -o ControlMaster=no -o BatchMode=no" go get github.com/hqueue/vatz-secret
$ make
```

TODO: We have to automate above steps.
