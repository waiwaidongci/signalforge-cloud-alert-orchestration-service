# Bug

Canceled request contexts are ignored by the notification dispatcher and by the email, log, and webhook delivery paths.

# Trigger

Run the three targeted canceled-context tests from `collection.json` against the bug snapshot.

# Error

`canceled channel errors = [<nil> <nil>]`
