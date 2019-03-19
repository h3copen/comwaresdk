`comwaresdk`
----------

This module is used for grpc connection with H3C devices. 

### grpc_service
---
grpc_service.proto — Defines the public RPC methods. In sdksessoin, we just used two methods： 

    service GrpcService {   // gRPC methods

        rpc Login (LoginRequest) returns (LoginReply) {}    // Login method

        rpc Logout (LogoutRequest) returns (LogoutReply) {}     // Logout method

    }

### sdksession
---
NewClient used to create a NewGrpcServiceClient and get a grpcSession token.  
Example usage:  

    var(
        address      string = "192.168.18.102"  
        port         uint = 50051  
        username     string = "2"  
        password     string = "123456"  
        grpcSession  *sdk.GrpcSession
    )
    
    grpcSession, err = sdk.NewClient(address, port, username , password)
    if err != nil {
    log.Println("Failed to open session.")
    }
    defer grpcSession.Close()

    ctx_with_token, cancel := sdk.CtxWithToken(grpcSession.Token, time.Second)
    defer cancel()
