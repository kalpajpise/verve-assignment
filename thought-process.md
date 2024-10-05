# High-Throughput Go REST Service: Implementation Approach and Design Considerations

## Overview
This document outlines the approach and design considerations for implementing a high-throughput REST service in Go, capable of processing at least 10,000 requests per second. The service includes features such as request deduplication, periodic processing, and integration with external systems.

## Key Components

**HTTP Server:** Using gorilla/mux for routing and handling HTTP requests.

**Redis Store:** Unique id's are maintained by using a distributed redis store, particulary redis `SET`, which only allows unique values.

**Periodic Processing:** Using a goroutine with a ticker to run a job every clock minute to reset the unique id's count and publish it to the streaming service, AWS kinesis in this case.

## Design Considerations
- Performance

    ***Goroutines:*** Utilizing Go's concurrency model to handle multiple requests simultaneously and schedule parallel job for distributed streaming.

   ***Efficient Data Structures:*** Using a set to natively store unique id's. This will efficiently handle the case where different app instances will receive the same id but it will be counted only once.

<br />

- Scalability

    *Stateless Design:* The core service is stateless, allowing for easy horizontal scaling.

    *Load Balancer Compatibility:* The deduplication logic works behind a load balancer, by using redis as the store.

<br />

- Reliability

    *Error Handling:* Robust error handling and logging is implemented throughout the application.

    *Graceful Shutdown:* Implemented graceful shutdown to ensure all requests are processed before the service stops.


<br />

- Extensibility

    **Modular Design:** Structured the code in a way that makes it easy to add new features or modify existing ones.

    *Configuration:* Using environment variables for easy application config with proper validation.