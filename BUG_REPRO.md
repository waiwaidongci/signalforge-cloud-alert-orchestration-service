# Bug

Scheduler recovery and silence-expiry workers dereference a nil runner, dependency container, database, or logger during startup.

# Trigger

Run the three targeted scheduler nil-path tests from `collection.json` against the bug snapshot.

# Error

`recovery worker panicked with missing dependencies`

`silence worker panicked with missing dependencies`

The logger case terminates with a nil-pointer panic in `TestSchedulerLoggerHandlesNil`.
