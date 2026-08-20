# Bug

The silence time-state boundary differs across the domain, application clock, and SQL filtering layers, so ended or future silences can appear active.

# Trigger

Run the three targeted silence boundary tests from `collection.json` against the bug snapshot.

# Error

`silence remained active at its end boundary`

The application test also reports that the repository received the system clock instead of the service clock, and the persistence test returns both the future and active rules.
