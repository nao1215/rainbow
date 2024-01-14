## Static Web Site Distribution With CloudFront and S3
### Overview
The simplest way to deploy the static website is to store the content in Amazon S3 (Simple Storage Service) and distribute it using CloudFront (Content Delivery Network).
  
This infrastructure configuration looks like the diagram below.　　
![./s3_cloudfront.png](./s3_cloudfront.png)

The configuration is characterized by its simplicity and the following features:

1. Cost-effectiveness
2. Responsive performance through effective utilization of caching (Cache Distribution pattern)
   
However, there are constraints. For example, if there is a functionality to rewrite a Relational Database on the client side, it cannot be accommodated with the infrastructure configuration depicted in the diagram.


#### Not allowed to access S3 directly
As a premise, you can host a static website using S3. In this context, a static website refers to content on individual web pages being static, although client-side scripts may be included.
  
In other words, S3 content can be publicly accessible, allowing direct access to S3. However, enabling public access to S3 poses security risks and the potential for information leakage. In general, public access to S3 should be disabled. For instance, there is a risk of personal information being stolen from S3 by third parties, or the possibility of delivering compromised JavaScript containing malicious code.
  
To prevent such scenarios, it is essential to appropriately configure the S3 bucket policy.

#### Access Log
[WIP]

#### Chache
[WIP]