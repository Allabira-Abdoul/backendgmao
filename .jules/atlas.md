# Atlas Journal

## 2023-10-25 - [SRP/DIP]
Learning: The application layer (`MaintenanceService`) was tightly coupled to infrastructure-level auditing details. It contained logic for dual-publishing to an HTTP API and RabbitMQ, managing goroutines, and context extraction directly within the business logic struct.
Action: Extracted this infrastructure coordination into a secondary port `AuditLogger` and a `CompositeLogger` adapter. This aligns the application with SRP and DIP by making the service depend on a simple abstraction, encapsulating the complex asynchronous infrastructure logic entirely within the adapter.
