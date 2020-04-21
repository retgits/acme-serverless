# ACME Serverless Fitness Store on AWS

Serverless is more than just Function-as-a-Service. There are many different forms to build and deploy serverless applications and one of them is using functions, messaging, and data storage on [AWS](https://aws.amazon.com). This tutorial walks you through setting up the services of the ACME Serverless Fitness Shop.

## Prerequisites

This tutorial leverages [AWS Lambda](https://aws.amazon.com/lambda/), [Amazon Simple Queue Service](https://aws.amazon.com/sqs/), and [Amazon DynamoDB](https://aws.amazon.com/dynamodb/). Sign up [here](https://portal.aws.amazon.com/billing/signup) for an AWS account. The deployments are done using [Pulumi](https://app.pulumi.com/signup).

To enable error reporting, you'll need a [Sentry.io account](https://sentry.io) and to get performance stats, you'll need a [Wavefront account](https://www.wavefront.com/sign-up/) and the [Wavefront API key](https://docs.wavefront.com/wavefront_api.html) for your account.

This tutorial assumes that you've used the [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) before and have configured it accordingly. If you haven't, you can run:

```bash
$ aws configure
AWS Access Key ID [None]: <YOUR_ACCESS_KEY_ID>
AWS Secret Access Key [None]: <YOUR_SECRET_ACCESS_KEY>
Default region name [None]:
Default output format [None]:
```

or manually create the `~/.aws/credentials` file and populate it with the expected settings:

```text
[default]
aws_access_key_id = <YOUR_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_SECRET_ACCESS_KEY>
```

## Step 1: Create DynamoDB table

> Amazon DynamoDB is a key-value and document database that delivers single-digit millisecond performance at any scale. It's a fully managed, multiregion, multimaster, durable database with built-in security, backup and restore, and in-memory caching for internet-scale applications.

With serverless applications you don't want to have synchronous invocations whenever possible, so you can leverage Amazon SQS to decouple the message producers from the message consumers.

```bash
## Clone the project
git clone https://github.com/retgits/acme-serverless

## Get the dependencies for the app (these dependencies are also used for other services in the ACME Serverless Fitness Shop)
go get ./...

## Change to the dynamodb pulumi directory
cd datastore/dynamodb/pulumi

## Create a new Pulumi stack
pulumi stack init <your pulumi org>/acmeserverless-dynamodb/dev
```

Pulumi is configured using a file called `Pulumi.dev.yaml`. You can generate the file using the command below:

```bash
AUTHOR=`whoami`

cat > Pulumi.dev.yaml <<EOF
config:
  aws:region: us-west-2
  awsconfig:dynamodb:
    billingmode: PAY_PER_REQUEST
    writecapacity: 5
    readcapacity: 5
  awsconfig:tags:
    author: $AUTHOR
    feature: acmeserverless
    team: vcs
    version: 0.2.0
EOF
```

To create the Pulumi stack, and create the Amazon DynamoDB table, run `pulumi up`. You can validate the DynamoDB table has been created successfully, by running:

```bash
$ pulumi stack
Current stack is dev:
    Owner: retgits
    Last updated: x days ago
    Pulumi version: v1.12.1
Current stack resources (3):
    TYPE                             NAME
    pulumi:pulumi:Stack              acmeserverless-dynamodb-dev
    ├─ aws:dynamodb/table:Table      dev-acmeserverless-dynamodb
    └─ pulumi:providers:aws          default

...
```

To keep track of the resources in Pulumi, you can add tags to your stack:

```bash
pulumi stack tag set app:name acmeserverless
pulumi stack tag set app:feature dynamodb
pulumi stack tag set app:domain infra
```

## Step 2: Seeding DynamoDB

To seed the DynamoDB table with random data, you can use the Go app in the [seed](./seed) directory. The app has two required flags:

* `region`: The region to send requests to
* `table`: The Amazon DynamoDB table to use

As an example, using the default settings, you can run

```bash
cd ../seed
TABLE=`pulumi stack output Table::Name`
go run main.go -region=us-west-2 -table=$TABLE
```

## Step 3: Create the SQS queues

> Amazon Simple Queue Service (SQS) is a fully managed message queuing service that enables you to decouple and scale microservices, distributed systems, and serverless applications.

```bash
## Change to the sqs pulumi directory
cd ../../../messaging/sqs/pulumi

## Create a new Pulumi stack
pulumi stack init <your pulumi org>/acmeserverless-sqs/dev
```

Just like with DynamoDB, you'll need a `Pulumi.dev.yaml`. You can generate the file using the command below:

```bash
AUTHOR=`whoami`

cat > Pulumi.dev.yaml <<EOF
config:
  aws:region: us-west-2
  awsconfig:tags:
    author: $AUTHOR
    feature: acmeserverless
    team: vcs
    version: 0.2.0
EOF
```

To create the Pulumi stack, and create the Amazon SQS queues, run `pulumi up`. You can validate the queues have been created successfully, by running:

```bash
$ pulumi stack
Current stack is dev:
    Owner: retgits
    Last updated: x days ago
    Pulumi version: v1.12.1
Current stack resources (8):
    TYPE                         NAME
    pulumi:pulumi:Stack          acmeserverless-sqs-dev
    ├─ aws:sqs/queue:Queue       dev-acmeserverless-sqs-payment-error
    ├─ aws:sqs/queue:Queue       dev-acmeserverless-sqs-shipment-error
    ├─ aws:sqs/queue:Queue       dev-acmeserverless-sqs-payment-request
    ├─ aws:sqs/queue:Queue       dev-acmeserverless-sqs-payment-response
    ├─ aws:sqs/queue:Queue       dev-acmeserverless-sqs-shipment-response
    ├─ aws:sqs/queue:Queue       dev-acmeserverless-sqs-shipment-request
    └─ pulumi:providers:aws      default

...
```

To keep track of the resources in Pulumi, you can add tags to your stack:

```bash
pulumi stack tag set app:name acmeserverless
pulumi stack tag set app:feature sqs
pulumi stack tag set app:domain infra
```

## Step 4: Build and deploy the services

### Step 4.1: Build the Cart service

> A cart service, because what is a shop without a cart to put stuff in?

The Cart service is to keep track of carts and items in the different carts.

```bash
## Clone the project
git clone https://github.com/retgits/acme-serverless-cart

## Change to the cart pulumi directory
cd acme-serverless-cart/pulumi

## Create the Pulumi.dev.yaml file
## These values are used for all of the other services, so using "export" makes sure you don't have to do it every time
export AWS_ACCOUNT_ID=<value>
export SENTRY_DSN=<value>
export WAVEFRONT_TOKEN=<value>
export WAVEFRONT_URL=<value>
AUTHOR=`whoami`

cat > Pulumi.dev.yaml <<EOF
config:
  aws:region: us-west-2
  awsconfig:generic:
    accountid: "$AWS_ACCOUNT_ID"
    sentrydsn: $SENTRY_DSN
    wavefronturl: $WAVEFRONT_URL
    wavefronttoken: $WAVEFRONT_TOKEN
  awsconfig:tags:
    author: $AUTHOR
    feature: acmeserverless
    team: vcs
    version: 0.2.0
EOF

## Create a pulumi stack
pulumi stack init <your pulumi org>/acmeserverless-cart/dev

## Run the Pulumi up command
pulumi up

## Validate the results
$ pulumi stack
Current stack is dev:
    Owner: retgits
    Last updated: x weeks ago
    Pulumi version: v1.14.0
Current stack resources (52):
    TYPE                                                      NAME
    pulumi:pulumi:Stack                                       acmeserverless-cart-dev
    ├─ aws:iam/role:Role                                      ACMEServerlessCartRole-lambda-cart-additem
    ├─ aws:iam/role:Role                                      ACMEServerlessCartRole-lambda-cart-all
    ├─ aws:iam/role:Role                                      ACMEServerlessCartRole-lambda-cart-clear
    ├─ aws:iam/role:Role                                      ACMEServerlessCartRole-lambda-cart-itemmodify
    ├─ aws:iam/role:Role                                      ACMEServerlessCartRole-lambda-cart-itemtotal
    ├─ aws:iam/role:Role                                      ACMEServerlessCartRole-lambda-cart-modify
    ├─ aws:iam/role:Role                                      ACMEServerlessCartRole-lambda-cart-total
    ├─ aws:iam/role:Role                                      ACMEServerlessCartRole-lambda-cart-user
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessCartPolicy-lambda-cart-additem
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-cart-additem
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessCartPolicy-lambda-cart-all
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-cart-all
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-cart-clear
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessCartPolicy-lambda-cart-clear
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessCartPolicy-lambda-cart-itemmodify
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-cart-itemmodify
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-cart-itemtotal
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessCartPolicy-lambda-cart-itemtotal
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-cart-modify
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessCartPolicy-lambda-cart-modify
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessCartPolicy-lambda-cart-total
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-cart-total
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-cart-user
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessCartPolicy-lambda-cart-user
    ├─ aws:apigateway/restApi:RestApi                         CartService
    ├─ aws:lambda/function:Function                           dev-lambda-cart-all
    ├─ aws:lambda/permission:Permission                       AllCartsAPIPermission
    ├─ aws:apigateway/integration:Integration                 AllCartsAPIIntegration
    ├─ aws:lambda/function:Function                           dev-lambda-cart-clear
    ├─ aws:lambda/permission:Permission                       ClearCartAPIPermission
    ├─ aws:apigateway/integration:Integration                 ClearCartAPIIntegration
    ├─ aws:lambda/function:Function                           dev-lambda-cart-itemmodify
    ├─ aws:apigateway/integration:Integration                 ItemModifyAPIIntegration
    ├─ aws:lambda/permission:Permission                       ItemModifyAPIPermission
    ├─ aws:lambda/function:Function                           dev-lambda-cart-itemtotal
    ├─ aws:lambda/permission:Permission                       ItemTotalAPIPermission
    ├─ aws:apigateway/integration:Integration                 ItemTotalAPIIntegration
    ├─ aws:lambda/function:Function                           dev-lambda-cart-total
    ├─ aws:lambda/permission:Permission                       CartTotalAPIPermission
    ├─ aws:apigateway/integration:Integration                 CartTotalAPIIntegration
    ├─ aws:lambda/function:Function                           dev-lambda-cart-modify
    ├─ aws:lambda/permission:Permission                       CartModifyAPIPermission
    ├─ aws:apigateway/integration:Integration                 CartModifyAPIIntegration
    ├─ aws:lambda/function:Function                           dev-lambda-cart-user
    ├─ aws:lambda/permission:Permission                       CartUserTotalAPIPermission
    ├─ aws:apigateway/integration:Integration                 CartUserTotalAPIIntegration
    ├─ aws:lambda/function:Function                           dev-lambda-cart-additem
    ├─ aws:lambda/permission:Permission                       AddItemAPIPermission
    ├─ aws:apigateway/integration:Integration                 AddItemAPIIntegration
    ├─ aws:apigateway/deployment:Deployment                   prod
    └─ pulumi:providers:aws                                   default

...

## Tag the resources
pulumi stack tag set app:name acmeserverless
pulumi stack tag set app:feature acmeserverless-cart
pulumi stack tag set app:domain cart

## Change back to the parent directory
cd ..
```

### Step 4.2: Build the Catalog service

> A catalog service, because what is a shop without a catalog to show off our awesome red pants?

The Catalog service is to register and serve the catalog of items sold by the shop.

```bash
## Clone the project
git clone https://github.com/retgits/acme-serverless-catalog

## Change to the catalog pulumi directory
cd acme-serverless-catalog/pulumi

## Create the Pulumi.dev.yaml file
cat > Pulumi.dev.yaml <<EOF
config:
  aws:region: us-west-2
  awsconfig:generic:
    accountid: "$AWS_ACCOUNT_ID"
    sentrydsn: $SENTRY_DSN
    wavefronturl: $WAVEFRONT_URL
    wavefronttoken: $WAVEFRONT_TOKEN
  awsconfig:tags:
    author: $AUTHOR
    feature: acmeserverless
    team: vcs
    version: 0.2.0
EOF

## Create a pulumi stack
pulumi stack init <your pulumi org>/acmeserverless-catalog/dev

## Run the Pulumi up command
pulumi up

## Validate the results
$ pulumi stack
Current stack is dev:
    Owner: retgits
    Last updated: x weeks ago
    Pulumi version: v1.14.0
Current stack resources (22):
    TYPE                                                      NAME
    pulumi:pulumi:Stack                                       acmeserverless-catalog-dev
    ├─ aws:iam/role:Role                                      ACMEServerlessCatalogRole-lambda-catalog-get
    ├─ aws:iam/role:Role                                      ACMEServerlessCatalogRole-lambda-catalog-all
    ├─ aws:iam/role:Role                                      ACMEServerlessCatalogRole-lambda-catalog-newproduct
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-catalog-get
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessCatalogPolicy-lambda-catalog-get
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessCatalogPolicy-lambda-catalog-all
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-catalog-all
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessCatalogPolicy-lambda-catalog-newproduct
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-catalog-newproduct
    ├─ aws:apigateway/restApi:RestApi                         CatalogService
    ├─ aws:lambda/function:Function                           dev-lambda-catalog-get
    ├─ aws:lambda/permission:Permission                       GetCatalogsAPIPermission
    ├─ aws:apigateway/integration:Integration                 GetCatalogsAPIIntegration
    ├─ aws:lambda/function:Function                           dev-lambda-catalog-all
    ├─ aws:lambda/permission:Permission                       AllCatalogsAPIPermission
    ├─ aws:apigateway/integration:Integration                 AllCatalogsAPIIntegration
    ├─ aws:lambda/function:Function                           dev-lambda-catalog-newproduct
    ├─ aws:lambda/permission:Permission                       NewCatalogAPIPermission
    ├─ aws:apigateway/integration:Integration                 NewCatalogAPIIntegration
    ├─ aws:apigateway/deployment:Deployment                   prod
    └─ pulumi:providers:aws                                   default

...

## Tag the resources
pulumi stack tag set app:name acmeserverless
pulumi stack tag set app:feature acmeserverless-catalog
pulumi stack tag set app:domain catalog

## Change back to the parent directory
cd ..
```

### Step 4.3: Build the Order service

> An order service, because what is a shop without actual orders to be shipped?

The Order service is to interact with the catalog, front-end, and make calls to the order services.

```bash
## Clone the project
git clone https://github.com/retgits/acme-serverless-order

## Change to the order pulumi directory
cd acme-serverless-order/pulumi

## Create the Pulumi.dev.yaml file
cat > Pulumi.dev.yaml <<EOF
config:
  aws:region: us-west-2
  awsconfig:generic:
    accountid: "$AWS_ACCOUNT_ID"
    sentrydsn: $SENTRY_DSN
    wavefronturl: $WAVEFRONT_URL
    wavefronttoken: $WAVEFRONT_TOKEN
  awsconfig:tags:
    author: $AUTHOR
    feature: acmeserverless
    team: vcs
    version: 0.2.0
EOF

## Create a pulumi stack
pulumi stack init <your pulumi org>/acmeserverless-order/dev

## Run the Pulumi up command
pulumi up

## Validate the results
$ pulumi stack
Current stack is dev:
    Owner: retgits
    Last updated: x weeks ago
    Pulumi version: v1.14.0
Current stack resources (32):
    TYPE                                                      NAME
    pulumi:pulumi:Stack                                       acmeserverless-order-dev
    ├─ aws:iam/role:Role                                      ACMEServerlessOrderRole-lambda-order-all
    ├─ aws:iam/role:Role                                      ACMEServerlessOrderRole-lambda-order-users
    ├─ aws:iam/role:Role                                      ACMEServerlessOrderRole-lambda-order-sqs-add
    ├─ aws:iam/role:Role                                      ACMEServerlessOrderRole-lambda-order-sqs-update
    ├─ aws:iam/role:Role                                      ACMEServerlessOrderRole-lambda-order-sqs-ship
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessOrderPolicy-lambda-order-all
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-order-all
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessOrderPolicy-lambda-order-users
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-order-users
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessOrderPolicy-lambda-order-sqs-add
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-order-sqs-add
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-order-sqs-update
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessOrderSQSPolicy-lambda-order-sqs-update
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessPaymentSQSPolicy-lambda-order-sqs-ship
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-order-sqs-ship
    ├─ aws:apigateway/restApi:RestApi                         OrderService
    ├─ aws:lambda/function:Function                           dev-lambda-order-all
    ├─ aws:lambda/permission:Permission                       OrderAllAPIPermission
    ├─ aws:apigateway/integration:Integration                 OrderAllAPIIntegration
    ├─ aws:lambda/function:Function                           dev-lambda-order-users
    ├─ aws:lambda/permission:Permission                       UserOrdersAPIPermission
    ├─ aws:apigateway/integration:Integration                 UserOrdersAPIIntegration
    ├─ aws:lambda/function:Function                           dev-lambda-order-sqs-add
    ├─ aws:lambda/permission:Permission                       OrderAddAPIPermission
    ├─ aws:apigateway/integration:Integration                 OrderAddAPIIntegration
    ├─ aws:apigateway/deployment:Deployment                   prod
    ├─ aws:lambda/function:Function                           dev-lambda-order-sqs-update
    ├─ aws:lambda/eventSourceMapping:EventSourceMapping       dev-lambda-order-sqs-update
    ├─ aws:lambda/function:Function                           dev-lambda-order-sqs-ship
    ├─ aws:lambda/eventSourceMapping:EventSourceMapping       dev-lambda-order-sqs-ship
    └─ pulumi:providers:aws                                   default

...

## Tag the resources
pulumi stack tag set app:name acmeserverless
pulumi stack tag set app:feature acmeserverless-order
pulumi stack tag set app:domain order

## Change back to the parent directory
cd ..
```

### Step 4.4: Build the Payment service

> A payment service, because nothing in life is really free...

The Payment service is to validate credit card payments. Currently the only validation performed is whether the card is acceptable.

```bash
## Clone the project
git clone https://github.com/retgits/acme-serverless-payment

## Change to the payment pulumi directory
cd acme-serverless-payment/pulumi

## Create the Pulumi.dev.yaml file
cat > Pulumi.dev.yaml <<EOF
config:
  aws:region: us-west-2
  awsconfig:generic:
    accountid: "$AWS_ACCOUNT_ID"
    sentrydsn: $SENTRY_DSN
    wavefronturl: $WAVEFRONT_URL
    wavefronttoken: $WAVEFRONT_TOKEN
  awsconfig:tags:
    author: $AUTHOR
    feature: acmeserverless
    team: vcs
    version: 0.2.0
EOF

## Create a pulumi stack
pulumi stack init <your pulumi org>/acmeserverless-payment/dev

## Run the Pulumi up command
pulumi up

## Validate the results
$ pulumi stack
Current stack is dev:
    Owner: retgits
    Last updated: x weeks ago
    Pulumi version: v1.14.0
Current stack resources (7):
    TYPE                                                      NAME
    pulumi:pulumi:Stack                                       acmeserverless-payment-dev
    ├─ aws:iam/role:Role                                      ACMEServerlessPaymentRole
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessPaymentSQSPolicy
    ├─ aws:lambda/function:Function                           dev-lambda-payment
    ├─ aws:lambda/eventSourceMapping:EventSourceMapping       dev-lambda-payment
    └─ pulumi:providers:aws                                   default

...

## Tag the resources
pulumi stack tag set app:name acmeserverless
pulumi stack tag set app:feature acmeserverless-payment
pulumi stack tag set app:domain payment

## Change back to the parent directory
cd ..
```

### Step 4.5: Build the Shipment service

> A shipping service, because what is a shop without a way to ship your purchases?

The Shipping service is, as the name implies, to ship products using a wide variety of shipping suppliers.

```bash
## Clone the project
git clone https://github.com/retgits/acme-serverless-shipment

## Change to the shipment pulumi directory
cd acme-serverless-shipment/pulumi

## Create the Pulumi.dev.yaml file
cat > Pulumi.dev.yaml <<EOF
config:
  aws:region: us-west-2
  awsconfig:generic:
    accountid: "$AWS_ACCOUNT_ID"
    sentrydsn: $SENTRY_DSN
    wavefronturl: $WAVEFRONT_URL
    wavefronttoken: $WAVEFRONT_TOKEN
  awsconfig:tags:
    author: $AUTHOR
    feature: acmeserverless
    team: vcs
    version: 0.2.0
EOF

## Create a pulumi stack
pulumi stack init <your pulumi org>/acmeserverless-shipment/dev

## Run the Pulumi up command
pulumi up

## Validate the results
$ pulumi stack
Current stack is dev:
    Owner: retgits
    Last updated: x weeks ago
    Pulumi version: v1.14.0
Current stack resources (7):
    TYPE                                                      NAME
    pulumi:pulumi:Stack                                       acmeserverless-shipment-dev
    ├─ aws:iam/role:Role                                      ACMEServerlessShipmentRole
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessShipmentSQSPolicy
    ├─ aws:lambda/function:Function                           dev-lambda-shipment
    ├─ aws:lambda/eventSourceMapping:EventSourceMapping       dev-lambda-shipment
    └─ pulumi:providers:aws                                   default

...

## Tag the resources
pulumi stack tag set app:name acmeserverless
pulumi stack tag set app:feature acmeserverless-shipment
pulumi stack tag set app:domain shipment

## Change back to the parent directory
cd ..
```

### Step 4.6: Build the User service

> A user service, because what is a shop without users to buy our awesome red pants?

The User service is to register and authenticate users using JWT tokens.

```bash
## Clone the project
git clone https://github.com/retgits/acme-serverless-user

## Change to the user pulumi directory
cd acme-serverless-user/pulumi

## Create the Pulumi.dev.yaml file
cat > Pulumi.dev.yaml <<EOF
config:
  aws:region: us-west-2
  awsconfig:generic:
    accountid: "$AWS_ACCOUNT_ID"
    sentrydsn: $SENTRY_DSN
    wavefronturl: $WAVEFRONT_URL
    wavefronttoken: $WAVEFRONT_TOKEN
  awsconfig:tags:
    author: $AUTHOR
    feature: acmeserverless
    team: vcs
    version: 0.2.0
EOF

## Create a pulumi stack
pulumi stack init <your pulumi org>/acmeserverless-user/dev

## Run the Pulumi up command
pulumi up

## Validate the results
$ pulumi stack
Current stack is dev:
    Owner: retgits
    Last updated: x weeks ago
    Pulumi version: v1.14.0
Current stack resources (40):
    TYPE                                                      NAME
    pulumi:pulumi:Stack                                       acmeserverless-user-dev
    ├─ aws:iam/role:Role                                      ACMEServerlessUserRole-lambda-user-refreshtoken
    ├─ aws:iam/role:Role                                      ACMEServerlessUserRole-lambda-user-get
    ├─ aws:iam/role:Role                                      ACMEServerlessUserRole-lambda-user-verifytoken
    ├─ aws:iam/role:Role                                      ACMEServerlessUserRole-lambda-user-all
    ├─ aws:iam/role:Role                                      ACMEServerlessUserRole-lambda-user-login
    ├─ aws:iam/role:Role                                      ACMEServerlessUserRole-lambda-user-register
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-user-refreshtoken
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessUserPolicy-lambda-user-refreshtoken
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-user-get
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessUserPolicy-lambda-user-get
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessUserPolicy-lambda-user-verifytoken
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-user-verifytoken
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-user-all
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessUserPolicy-lambda-user-all
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-user-login
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessUserPolicy-lambda-user-login
    ├─ aws:iam/rolePolicyAttachment:RolePolicyAttachment      AWSLambdaBasicExecutionRole-lambda-user-register
    ├─ aws:iam/rolePolicy:RolePolicy                          ACMEServerlessUserPolicy-lambda-user-register
    ├─ aws:apigateway/restApi:RestApi                         UserService
    ├─ aws:lambda/function:Function                           dev-lambda-user-refreshtoken
    ├─ aws:lambda/permission:Permission                       RefreshTokenAPIAPIPermission
    ├─ aws:apigateway/integration:Integration                 RefreshTokenAPIAPIIntegration
    ├─ aws:lambda/function:Function                           dev-lambda-user-verifytoken
    ├─ aws:lambda/permission:Permission                       VerifyTokenAPIPermission
    ├─ aws:apigateway/integration:Integration                 VerifyTokenAPIIntegration
    ├─ aws:lambda/function:Function                           dev-lambda-user-get
    ├─ aws:lambda/permission:Permission                       GetUserAPIPermission
    ├─ aws:apigateway/integration:Integration                 GetUserAPIIntegration
    ├─ aws:lambda/function:Function                           dev-lambda-user-register
    ├─ aws:lambda/permission:Permission                       RegisterUserAPIPermission
    ├─ aws:apigateway/integration:Integration                 RegisterUserAPIIntegration
    ├─ aws:lambda/function:Function                           dev-lambda-user-all
    ├─ aws:lambda/permission:Permission                       AllUsersAPIPermission
    ├─ aws:apigateway/integration:Integration                 AllUsersAPIIntegration
    ├─ aws:lambda/function:Function                           dev-lambda-user-login
    ├─ aws:apigateway/integration:Integration                 LoginUserAPIIntegration
    ├─ aws:lambda/permission:Permission                       LoginUserAPIPermission
    ├─ aws:apigateway/deployment:Deployment                   prod
    └─ pulumi:providers:aws                                   default

...

## Tag the resources
pulumi stack tag set app:name acmeserverless
pulumi stack tag set app:feature acmeserverless-user
pulumi stack tag set app:domain user

## Change back to the parent directory
cd ..
```

### Step 4.7: Build the Point-of-Sales (POS) service

> A point-of-sales app to sell our products in brick-and-mortar stores!

The Point-of-Sales is to serve as the front-end for Point-of-Sales locations.

```bash
## Clone the project
git clone https://github.com/retgits/acme-serverless-pos

## Update lines 3 through 5 in acme-serverless-pos/src/script.js with
# Catalog (line 3)
pulumi stack output Gateway::URL --stack retgits/acmeserverless-catalog/dev
# Users (line 4)
pulumi stack output Gateway::URL --stack retgits/acmeserverless-user/dev
# Order (line 4)
pulumi stack output Gateway::URL --stack retgits/acmeserverless-order/dev

## Change to the pos pulumi directory
cd acme-serverless-pos/pulumi

## Create the Pulumi.dev.yaml file
cat > Pulumi.dev.yaml <<EOF
config:
  aws:region: us-west-2
  awsconfig:s3:
    bucket: acme-serverless-$AUTHOR-pos
  awsconfig:tags:
    author: $AUTHOR
    feature: acmeserverless
    team: vcs
    version: 0.2.0
EOF

## Create a pulumi stack
pulumi stack init <your pulumi org>/acmeserverless-pos/dev

## Run the Pulumi up command
pulumi up

## Validate the results
$ pulumi stack
Current stack is dev:
    Owner: retgits
    Last updated: x month ago
    Pulumi version: v1.12.1
Current stack resources (8):
    TYPE                                     NAME
    pulumi:pulumi:Stack                      acmeserverless-pos-dev
    ├─ aws:s3/bucket:Bucket                  acme-serverless-pos
    ├─ aws:s3/bucketPolicy:BucketPolicy      bucketPolicy
    ├─ aws:s3/bucketObject:BucketObject      index.html
    ├─ aws:s3/bucketObject:BucketObject      favicon.png
    ├─ aws:s3/bucketObject:BucketObject      script.js
    ├─ aws:s3/bucketObject:BucketObject      style.css
    └─ pulumi:providers:aws                  default

...

## Tag the resources
pulumi stack tag set app:name acmeserverless
pulumi stack tag set app:feature acmeserverless-pos
pulumi stack tag set app:domain pos

## After the deployment you'll be able to see the Point-of-Sales app on
echo https://`pulumi stack output "POS URL"`
```

## Step 5: Test the services

After the deployments are complete, you can use the below cURL commands to test the functions.

### Step 5.1: Test the Cart service

```bash
URL=`pulumi stack output Gateway::URL --stack retgits/acmeserverless-cart/dev`
```

```bash
## Get all carts
curl --request GET \
  --url $URL/cart/all

## Add item to cart
curl --request POST \
  --url $URL/cart/item/add/499fd0be-63c1-4a24-bcdf-46694671b77f \
  --header 'content-type: application/json' \
  --data '{
    "description": "fitband for any age ",
    "itemid": "sdfsdfsfs",
    "name": "fitband",
    "price": 4.5,
    "quantity": 1
}'
```

### Step 5.2: Test the Catalog service

```bash
URL=`pulumi stack output Gateway::URL --stack retgits/acmeserverless-catalog/dev`
```

```bash
## Get product details
curl --request GET \
  --url $URL/products/050b7bdc-e993-4884-bb60-18323f9278dd

## Add new product
curl --request POST \
  --url $URL/product \
  --header 'content-type: application/json' \
  --data '{
    "name": "Tracker",
    "shortDescription": "Limited Edition Tracker",
    "description": "Limited edition Tracker with longer description",
    "imageurl1": "/static/images/tracker_square.jpg",
    "imageurl2": "/static/images/tracker_thumb2.jpg",
    "imageurl3": "/static/images/tracker_thumb3.jpg",
    "price": 149.99,
    "tags": [
        "tracker"
    ]
}'
```

### Step 5.3: Test the Order service

```bash
URL=`pulumi stack output Gateway::URL --stack retgits/acmeserverless-order/dev`
```

```bash
## Submit new order
curl --request POST \
  --url $URL/order/add/bbuttner0 \
  --header 'content-type: application/json' \
  --data '{
  "_id": "fa21f430-310a-4077-8280-2f091804e280",
  "status": "pending payment",
  "userid": "bbuttner0",
  "firstname": "Basia",
  "lastname": "Buttner",
  "address": {
    "street": "Jenna",
    "city": "Fresno",
    "zip": "93704",
    "state": "CA",
    "country": "United States"
  },
  "email": "bbuttner0@someplace.com",
  "delivery": "usps",
  "card": {
    "Type": "Visa",
    "Number": "4222222222222",
    "ExpiryMonth": 7,
    "ExpiryYear": 2021,
    "CVV": "672"
  },
  "cart": [
    {
      "id": "09c0a567-812e-49c1-9f63-e0f619db6466",
      "description": "arcu adipiscing molestie hendrerit at vulputate vitae nisl aenean lectus pellentesque eget nunc",
      "quantity": 5,
      "price": 6.3
    },
    {
      "id": "3a91c3b5-6273-4c86-82b9-896044ebde36",
      "description": "aliquet pulvinar sed nisl nunc rhoncus dui vel sem sed sagittis nam congue risus semper porta volutpat",
      "quantity": 5,
      "price": 7.6
    }
  ],
  "total": "1402"
}'

## Get all orders
curl --request GET \
  --url $URL/order/all
```

### Step 5.4: Test the User service

```bash
URL=`pulumi stack output Gateway::URL --stack retgits/acmeserverless-user/dev`
```

```bash
## Add a new user
curl --request POST \
  --url $URL/register \
  --header 'content-type: application/json' \
  --data '{
    "username":"peterp",
    "password":"vmware1!",
    "firstname":"amazing",
    "lastname":"spiderman",
    "email":"peterp@acmefitness.com"
}'

## Get the JWT tokens
curl --request POST \
  --url $URL/login \
  --header 'content-type: application/json' \
  --data '{
	"username": "peterp",
	"password": "vmware1!"
}'
```

## Step 6: Clean up

To remove all the Lambda functions, SQS queues, and the DynamoDB table, you can run:

```bash
pulumi stack rm --stack <your pulumi org>/acmeserverless-dynamodb/dev
pulumi stack rm --stack <your pulumi org>/acmeserverless-sqs/dev
pulumi stack rm --stack <your pulumi org>/acmeserverless-cart/dev
pulumi stack rm --stack <your pulumi org>/acmeserverless-catalog/dev
pulumi stack rm --stack <your pulumi org>/acmeserverless-order/dev
pulumi stack rm --stack <your pulumi org>/acmeserverless-payment/dev
pulumi stack rm --stack <your pulumi org>/acmeserverless-shipment/dev
pulumi stack rm --stack <your pulumi org>/acmeserverless-user/dev
pulumi stack rm --stack <your pulumi org>/acmeserverless-pos/dev
```