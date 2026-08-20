# Bug

A notification whose stored payload is `null` is loaded with a nil map. Writing a field to the returned payload panics.

# Trigger

Run the two targeted payload decoding tests from `collection.json` against the bug snapshot.

# Error

`panic: assignment to entry in nil map`
