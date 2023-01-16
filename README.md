# git-switch

## Install
```bash
go install github.com/YukiAbyss/git-switch@latest
export PATH=$GOPATH/bin:$PATH
```

## Only supports ssh key
[Generating a new SSH key](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent)
```bash
ssh-keygen -t -C "your_email@example.com"
```

## Example
```bash
# new git a user
git-switch -n

# delete git a user
git-switch -d

# switch git a user
git-switch

# print shell cmd execution information
git-switch -o
```

## Check installed
```bash
$ ls -al $GOPATH/bin | grep git-switch
-rwxr-xr-x   1 yy  staff   4402130  1  5 12:10 git-switch

$ ls -al ~/.gitswitchv1.json 
-rw-r--r--  1 yy  staff  232  1 16 12:34 /Users/Yuki/.gitswitchv1.json

$ git-switch
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select a switch git user: 
  ▸ yuki         yy__yyyy@126.com       /Users/Yuki/.ssh/yuki_id_rsa
    nodereal     yuki.w@nodereal.io     /Users/Yuki/.ssh/nodereal_id_rsa
```

### After setting git-switch, you can view the configuration file 「~/.gitswitch.json」
```bash
$ cat ~/.gitswitch.json 
[
	{
		"name": "yuki",
		"email": "yy__yyyy@126.com",
		"ssh_key_file_path": "/Users/Yuki/.ssh/yuki_id_rsa"
	},
	{
		"name": "nodereal",
		"email": "yuki.w@nodereal.io",
		"ssh_key_file_path": "/Users/Yuki/.ssh/nodereal_id_rsa"
	}
]
```

### Some operations during switch

When git-switch [selects users](main.go#L139), the following 2 steps will be performed 
1. Delete ssh key
2. Add/Overwrite git global config
```bash
# clear ssh key and add ssh key
ssh-add -D
ssh-add ~/.ssh/id_rsa

# set git config
git config --global user.name {your name}
git config --global user.email {your email}
```



