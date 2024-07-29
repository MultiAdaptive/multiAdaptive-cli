# multiAdaptive-cli

## Software Dependencies

| Dependency                                                    | Version  | Version Check Command |
| ------------------------------------------------------------- | -------- | --------------------- |
| [git](https://git-scm.com/)                                   | `^2`     | `git --version`       |
| [go](https://go.dev/)                                         | `^1.21`  | `go version`          |
| [make](https://linux.die.net/man/1/make)                      | `^3`     | `make --version`      |

## Build

  ```bash
    make build
  ```

## View All Broadcast Nodes

1. Run multiAdaptive-cli

```bash
./build/multiAdaptive-cli -privateKey="<your privateKey>" -advanced
```

2. Select View Broadcast Node Information

## View All Storage Nodes

1. Run multiAdaptive-cli

```bash
./build/multiAdaptive-cli -privateKey="<your privateKey>" -advanced
```

2. Select View Storage Node Information

## Register NodeGroup
1. Run multiAdaptive-cli

```bash
./build/multiAdaptive-cli -privateKey="<your privateKey>" -advanced
```

2. Select Register NodeGroup.
3. Enter the list of broadcast node addresses, separated by commas.
4. Enter the minimum number of signatures.
5. Obtain the nodeGroupKey (used for advanced testing).

## Register NameSpace
1. Run multiAdaptive-cli

```bash
./build/multiAdaptive-cli -privateKey="<your privateKey>" -advanced
```

2. Select Register NameSpace.
3. Enter the list of storage node addresses, separated by commas.
4. Obtain the nameSpaceKey (used for advanced testing).

## Participate in Testing
By default, this will participate in sending DA data for testing, sending data every five minutes. You can compare the information on the chain for verification.  
1. Run multiAdaptive-cli

```bash
./build/multiAdaptive-cli -privateKey="<your privateKey>" -advanced
```

2. Select General test

## Participate in Advanced Testing
1. Run multiAdaptive-cli

```bash
./build/multiAdaptive-cli -privateKey="<your privateKey>" -advanced
```

2. Select Advanced test.
3. Enter the nodeGroupKey that meets the policy.
4. Enter the nameSpaceKey you created (can be left empty).
5. Enter the data sending interval.

## Register as a Broadcast Node
1. Run multiAdaptive-cli

```bash
./build/multiAdaptive-cli -privateKey="<your privateKey>" -advanced
```

2. Select Register as a broadcast node.
3. Enter the URL where the node provides its services.
4. Enter the name of the node.
5. Enter the region where the node is located.
6. Enter the amount of staked tokens.
7. Enter the maximum storage size provided.

## Register as a Storage Node
1. Run multiAdaptive-cli

```bash
./build/multiAdaptive-cli -privateKey="<your privateKey>" -advanced
```

2. Select Register as a storage node.
3. Enter the URL where the node provides its services.
4. Enter the name of the node.
5. Enter the region where the node is located.
6. Enter the amount of staked tokens.
7. Enter the maximum storage size provided.

