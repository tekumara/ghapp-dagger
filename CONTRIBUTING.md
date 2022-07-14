# Contributing

## Prereqs

Install the following:

- docker: `brew install --cask docker`
- [dagger](https://dagger.io/): `brew install dagger/tap/dagger`
- [ngrok](https://ngrok.com/): `brew install ngrok/ngrok/ngrok`

## Getting started

1. Create a GitHub App via Settings -> Developer Settings -> New GitHub App:
   - Homepage URL: enter a dummy value
   - Webhook URL: enter a dummy value, this will be updated later
   - Permissions:
     - Checks: read & write
     - Contents: read-only
     - Metadata: read-only
   - Events:
     - ☑️ Check run
     - ☑️ Check suite
1. From your app's page, generate a private key and save it to disk.
1. Create _.envrc_ based on [.envrc.example](.envrc.example) with your Github App's settings (visible from the app's page).

## Development

1. Start ngrok to get a `https://*.ngrok.io` public URL that forwards to your laptop:

   ```
   ngrok http 8000
   ```

1. Update the Github App Webhook URL with your ngrok URL.

1. Run the server and continuously rebuild when files change:

   ```
   make run
   ```
