# GoUnity

GoUnity is a Go project that provides a client for managing Dell EMC Unity storage.


## Current State

**Under developing. Please contribute.**


## License

[Apache License version 2](LICENSE)


## Support Operations

- Query/Create LUNs
- Query Storage Pools
- Query Hosts
- Attach LUNs to Hosts


## Installation

```bash
go get github.com/murray-liang/gounity
```


## Tutorial

### Create a connection to Unity Systems

```go
unity, err := gounity.NewUnity("UnityMgmtIP", "username",
    "password", true)
if err != nil {
    panic(err)
}
```

### Query storage pools

```go
// List all the pools
pools, err := unity.GetPools()

// Get the pool by ID
pool, err := unity.GetPoolByID("Pool_1")
```

### Create LUNs
```go
// Create a 3GB LUN named `lunName` on `pool`
lun, err := pool.CreateLUN("lunName", 3)
```

### Query Hosts
```go
host, err := unity.GetHostByID("Host_1")
```

### Attach LUNs to Hosts
```go
hluNum, err := host.Attach(lun)
```

### More examples
*_test.go files of this package contains lots of detailed examples.


## Issues

If you have any questions or find any issues, please post [Github Issues](https://github.com/murray-LIANG/gounity/issues).
