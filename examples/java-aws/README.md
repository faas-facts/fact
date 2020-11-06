 ## Example implementation of the java fact client
 ### Setup
 First add the fact dependency to the pom 
 ``` 
      <dependency>
        <groupId>io.github.fact</groupId>
        <artifactId>client</artifactId>
        <version>0.1.1</version>
      </dependency> 
```
Make sure Maven is configured to use Github packages if you haven't follow [this guide](https://docs.github.com/en/free-pro-team@latest/packages/using-github-packages-with-your-projects-ecosystem/configuring-apache-maven-for-use-with-github-packages).
Use `mvn package` to generate the jar for your lambda.

### Deployment
In your aws console create a new lambda function with the java 11 runtime. Upload the generated jar named `*-with-dependencies.jar` 
Set the handler to `fact.ConsoleFactTest::handleRequest` now you should be good to go.

### TCP Log collection
Setup an EC2 instance within a VPC and the default security group. Log into your instance via SSH and download the newest fact release or build it yourself using the go compiler.
Run ``./fact tcp [port] [parallel connections] [options]``. Make sure the port you chose is being forwarded in your security group.

In your Lambda function chnage the handler to  `fact.TcpFactTest::handleRequest` and add `fact_tcp_address` and `fact_tcp_port` to the environment of the function representing your running TCP log collector on the EC2 instance.

Now you should be able to test the lambda function in the admin dashboard.
