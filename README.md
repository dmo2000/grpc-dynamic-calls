# How to call GRPC methods dynamically in GO

This repository shows code examples of calling gprc dynamically in go as explained at 


## Usage

1. Install proto compiler (protoc) - https://protobuf.dev/installation/
2. Install dependencies executing `make deps`
3. Run `make run-server` to execute server, 
4. Run `make run-server-reflection-call` to execute client calls using server reflection

## Sequence diagram for calling GRPC method using server reflection

```mermaid
sequenceDiagram

participant Client
box rgb(128,128,128,0.25) GRPC Server
  participant Reflection as Reflection Service
  participant Greeter as Greeter Service
end

Client->>Reflection: Call ServerReflectionInfo and request Greeter service proto file  
Reflection--)Client: Return Greeter service proto file
Client->>Client: Extract method descriptor of SayHello from Greeter service proto file
Client->>Greeter: Call SayHello using method descriptor
Greeter--)Client: Return SayHello response
```

## Sequence diagram for calling GRPC method using proto file

```mermaid
sequenceDiagram

participant Client
box rgb(128,128,128,0.25) GRPC Server
  participant Greeter as Greeter Service
end

Client->>Client: Load proto file
Client->>Client: Extract method descriptor of SayHello from Greeter service proto file
Client->>Greeter: Call SayHello using method descriptor
Greeter--)Client: Return SayHello response
```


