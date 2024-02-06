## Lambda batch with EventBridge (CloudWatch Events)
### Overview
Scheduled batch execution is one way to perform processing asynchronously in batches. There are two main advantages of scheduled batch execution:  

1. Automation: Scheduled batch execution eliminates the need for manual task execution by running tasks at regular intervals. This enables automation of workflows, saving time and effort.
2. Consistency: Using scheduled batch execution ensures tasks are performed at predetermined frequencies, maintaining consistent processing. This ensures data updates or processing occurs at the appropriate times.

### Architecture
The pattern using Lambda and EventBridge is a very simple design. EventBridge triggers Lambda at specified times. However, due to the timeout constraints of Lambda, this infrastructure pattern is not suitable for long-running processes.

The architecture of this infrastructure configuration is as follows:

![lambda-batch](./lambda-batch-with-event-bridge.svg)

### How to deploy
```shell
$ make deploy
```
The dependencies are [Golang](https://go.dev/doc/install) and [AWS Serverless Application Model (SAM) CLI](https://github.com/aws/aws-sam-cli). Please make sure to install them beforehand.
