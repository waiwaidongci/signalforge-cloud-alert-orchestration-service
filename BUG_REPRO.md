# Bug

An oversized request body loses its `http.MaxBytesError` chain while JSON decoding and is mapped to a generic client or server error instead of payload-too-large.

# Trigger

Run the three targeted request-body limit tests from `collection.json` against the bug snapshot.

# Error

`expected MaxBytesError in chain, got 请求体不是合法 JSON: http: request body too large`

The response mapping tests return status 500 or 400 instead of 413.
