# Pocketbase as a framework template.

## What is this?

This is a starter template that uses pocketbase as a framework using the [extend with go](https://pocketbase.io/docs/go-overview/) functionality. Thanks Pocketbase! Using a setup like this gives you something very close to a rails or django experience with the nice minimal code that Go provides. I really am liking it.

This allows us to create a very minimal lightweight website that should be very performant. We use pico css for a little bit of styling and htmx with some vanilla js sprinkled in to give us a little interactivity and make it feel like an app.
The demo comes with functionality to 'save a location' and 'delete a location'. It also has login, logout, register, and password reset(assuming you setup SMTP) capabilities.

## Getting Started
Clone the repo and run

```
go run . serve
```

Then login to the admin panel and create a new collection called favorite_locations with location text field and 'user_id' relation field call. This is optional, you can also just try out the auth stuff if you don't want to do this.

Example favorite_locations collection:

<img width="997" alt="image" src="https://github.com/user-attachments/assets/cbc45ad4-fafd-430e-a375-28cafa7f2236" />
