## cfn - list up or delete CloudFormation stacks
> [!IMPORTANT]
> Not implemented yet. This is specifiction document.

The cfn command provides the following features:
- [x] List stacks
- [ ] Delete stacks
- [ ] Add tags to stacks
- [x] Interactive mode

### How to install
```shell
go install github.com/nao1215/rainbow/cmd/cfn@latest
```

### How to use
The cfn command allows you to specify a profile as an option, but it is more user-friendly to use the `AWS_PROFILE` environment variable.

### List stacks
```shell
cfn ls
```

### Delete stacks
```shell
cfn rm ${STACK_NAME}
```

### Add tags to stacks
```shell
cfn tag ${STACK_NAME} ${TAG_KEY}=${TAG_VALUE}
```

### Interactive mode
```shell
cfn
```

![cfn_tui](./cfn_tui.gif)
