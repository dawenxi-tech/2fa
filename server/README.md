# Server

The Server is a simpel golang Server that runs NATS Jetstream for storage and HTMX for Web.

NATS Jetstream embedded into Server, so that any records that a client changes are stored encrypted on the server and synced with any other clients.

Auth can use NATS Auth. NATS is nice in that it has a auth and authz systme baked in.

Web GUI using HTMX, so that you can login using passkeys and then see whats going on with your enrolled clients, and view the records.

Any other ideas, that you have please add !!!

## Stack

Just a list of bits of the tech stack.

https://github.com/maxpert/marmot/blob/master/stream/embedded_nats.go shows how to embed NATS Jetstream Server in a golang Server. 
- Marmot is a Database synchronisation system that sync dataabses across different data cneters btw.

https://htmx.org is my current favourite was to build Web. Its gloriously simple.

https://htmx.org/server-examples/#go has some golang packages that make life easy.

## Hosting

Fly is free up to a certain point and light.

https://fly.io/docs/languages-and-frameworks/golang/

Typically you can deploy to 3 data centers and the Client will connect to the nearest one.
If once of the data centers goes down, everything keeps working and when it comes back up, that instacne will just catchup off the NATS Jetstream. 


## Build

See Makefile.

## Release

We can release the server to github and have fly action auto-update the Server running on fly.io.


## Client builds

The Web site can host all Client and Server downloads page for Desktops and Mobiles ( andorid side loading ).

At the moment these are tagged releases on github. But we can also have the NATS Server get told of a new release and NATS can store the Releases in it NATS KV store.

## Client Updates

Because each client is connected to NATS, when a client update occurs, we can push an update to the clients as an event, and then the cleont can download it. This is why NATS is realyl nice. there is no polling and its real time.

