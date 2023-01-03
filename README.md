# Dictionary - Exam 2023 (Silas Bue Handgaard)

## How to run

### Servers

The system will have 2 servers/replication managers. Open 2 terminals and run:

```shl
go run server/server.go 0
go run server/server.go 1
```

each in a seperate terminal. the first command will create the primary RM and
the second command will create the backup RM

### Clients

To run a client, in a new terminal run:

```shl
go run client/client.go 2
```

where 2 can be replaced with another number(except 0 or 1) if another port is desired.

#### Commands

### Add

To add something to the dictionary, enter:

```shl
add(<word>, <definition>)
```

where `<word>` is the word for which you want to add a definition and
`<definition>` is the definition of that word.

### Read

To read a definition from the dictionary, enter:

```shl
read(<word>)
```

where `<word>` is the word for which you want the definition.
