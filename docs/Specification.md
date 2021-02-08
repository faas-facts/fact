# Specification

Fact is using a protobuf file format to transfer trace data between the collector and functions.
Fact supports multiple methods to collect traces form a function and persist them to disk. 
This specification describes the data contained in a Trace as well as the fact-client-library life-cycle.

## Traces
Find the protobuf specification for the trace exchange format at  [fact/trace.proto](../fact/trace.proto)

### Fields
 
| Name | Format | Description |
| ---- |  ------ | ----------- |
| ID   | uuid  based on based on RFC 4122 and DCE 1.1: Authentication and Security Services | a unique id for a function trace | 
| ChildOf |  same as ID | Null or a uuid of a parent function |
| Timestamp | (Time)[https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/timestamp.proto] | the time this trace was generated | 
| ContainerID | string | a unique function identifier that remains the same for a warm function | 
| HostID | string | a unique string identifying a function host, can be empty for FaaS platforms that can't be fingerprinted|
| BootTime | Time | the time a function was first started. This time states the same for a warm function | 
| Cost | float | accumulated cost of a function execution at this time | 
| RequestStartTime | Time | time the request started (sending of first byte) |
| StartTime | Time | the time a function invocation started | 
| Status | HTTP Code | the result of a function invocation, encoded as a HTTP Code | 
| EndTime | Time | the time a function invocation finished | 
| RequestEndTime | Time | the time the last byte of the request is received |
| CodeVersion | string| unique version of the deployed artefact (should increase only for code change)|
| ConfigVersion | string | unique version of the deployed configuration (should change if part of the direct configuration is changed) |
| Platform | `[A-Z]^{2,4}` | Platform identifier string, e.g., AWS for AWS Lambda Functions. |
| Region | string | Cloud region name the function is running at. Naming depends on the cloud provider | 
| Runtime | `(kernel) (kernel version) (Runtime) (Runtime Version)`| the system fingerprint string. Should contain the kernel, runtime (Python), runtime version (2.7) |
| Memory | `[1-9]^{5}` | Amount of allocated function memory. | 
| ExecutionLatency| Duration in ns | current execution duration for this trace. | 
| RequestResponseLatency | Duration in ns | duration from sending the request until receiving the last byte of the request |
| ExecutionDelay | Duration in ns | duration from send the last byte of an request to the start of the execution |
| TransportDelay | Duration in ns | duration from returning the response to receiving the first byte of response |
| Env | map | collection of environment variables available to the function | 
| Tags | map | user defined tags for this trace | 
| Logs | map | map of timestamps and user define log messages | 
| Args | array | array of user defined string data | 

## Fact-Client-Library Life-Cycle

*Library Pseudo Code:*

```
    import fact

    //executes once when the function runtime is started
    fact.boot(Collecter-Configuration, Tags)
    
    //FaaS Handler function called for an incomming event
    function handle(event,context) {
        fact.start()
        ... user code1 ...
        //optinal functions to log importnat evets druing function execution
        fact.update(context,Message,Tag)
        ... user code ...
        fact.done(context, Message, args...)
        return ...
    }
```

Fact Libraries use 4 distinct life-cycle methods:
 - boot: collect basic environment information of this function and opens the collector transport
 - start: logs the start of each invocation, important for `StartTime`
 - update: optional method that allows the developer to add log messages and tags during execution
 - done: logs the end of an invocation, important to calculate overall cost and execution time, this method can also collect tags, logs and user-defined data.
 
Each Fact library implements these presented methods. Fact libraries should be FaaS platform independent and use fingerprinting or config methods to adapt to platform specific apis. 

