# Example implementation of the python fact client

### Deployment
First you need to create your deployment package. To generate a zip file containing all needed dependencies follow the [official aws guide](https://docs.aws.amazon.com/lambda/latest/dg/python-package.html#python-package-dependencies).
In your aws console create a new lambda function with the python 3.6 runtime and upload the generated zip.

### Console Log collection
Set the handler to  `console_fact_test.handler`
and run the test case. You should now be seeing the fact generated logs in your aws console.

### TCP Log collection
Setup an EC2 instance within a VPC and the default security group. Log into your instance via SSH and download the newest fact release or build it yourself using the go compiler.
Run ``./fact tcp [port] [parallel connections] [options]``. Make sure the port you chose is being forwarded in your security group.

In your Lambda function chnage the handler to  `tcp_fact_test.handler` and add `fact_tcp_address` and `fact_tcp_port` to the environment of the function representing your running TCP log collector on the EC2 instance.

Now you should be able to test the lambda function in the admin dashboard.