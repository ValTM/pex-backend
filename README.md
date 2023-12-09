# Counter app backend

## Description

Example backend for a counter task.
The counters can increment, decrement, and reset.
Each counter is identified by unique UUID, so they are separated for each client.
Counters can not be decreased to less than 0.

Technologies used:
- go-chi
- go-chi render

## Next steps

- counter storage
- better logging
- read server address and other settings from env/file
- more error handling
