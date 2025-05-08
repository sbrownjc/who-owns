# who-owns

`who-owns` is a dead simple script that:

- git clones `TheJumpCloud/inventory-mapping` to a temporary folder using SSH auth
- checks the repos and k8s namespaces for a matching pattern
- discards the temporary folder

## Installation

```sh
go install github.com/sbrownjc/who-owns@latest
```

## Usage

```sh
who-owns <repo/k8s_namespace>
```

### Example

```sh
$ who-owns policies-linux
Inventory last updated on Tue, 06 May 2025 09:59:55 EDT

Found repo TheJumpCloud/jumpcloud-policies-linux owned by Thorium
```
