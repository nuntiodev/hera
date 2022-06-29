## Welcome to Hera 👋

![example workflow](https://github.com/nuntiodev/hera/actions/workflows/build.yaml/badge.svg)

Building Nuntio Blocks: secure, efficient & open-source API blocks that run on any Linux-based environment and scales
massively. Nuntio Cloud is the commercial offering of the above technologies and help companies go from idea to
production faster, without losing control over their data or services. Our goal with Nuntio Cloud is to provide social
API-blocks that are secure, easy to set up, scales worldwide and that you can move from Nuntio Cloud to your private
data center whenever you want.

## Hera

Hera is a fully fledged user-management system written in **Golang**, designed to be secure & stateless and to run in Kubernetes. It handles
advanced data encryption models and features such as: normal and verfified authentication flows, creating, fetching and
updating user information with custom metatadata, sending and validating emails and much more!
Reach out to info@nuntio.io if you wanna give it a try, or sign-up for for a porfile
in [Nuntio Cloud](https://cloud.nuntio.io) if you wanna try our managed solution.

## Environment

| Name                                  | Type            | Description                                                                                                                                                                                                                                                                                                              | Default                 | Required |
|---------------------------------------|-----------------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------|----------|
| ENCRYPTION_KEYS                       | []String        | An array of encryption keys used to encrypt/decrypt data at rest. If you provide more than one key, the keys will automatically be used to create a master key. If you think your data has been compromised, add another encryption key to the list, and the service will automatically encrypt users under the new key. | []                      | No       |
| GRPC_PORT                             | int             | The port you wish to start the gRPC server on.                                                                                                                                                                                                                                                                           | 9000                    | No       |
| HTTP_PORT                             | int             | The port you wish to start the http server on.                                                                                                                                                                                                                                                                           | 9001                    | No       |
| ENABLE_GRPC_SERVER                    | bool            | This specifies weather or not you want to run the gRPC server.                                                                                                                                                                                                                                                           | true                    | No       |
| ENABLE_HTTP_SERVER                    | bool            | This specifies weather or not you want to run the http server.                                                                                                                                                                                                                                                           | false                   | No       |
| ACCESS_TOKEN_EXPIRY                   | Duration        | The expiration time for the access tokens. Should be a valid string duration such as: "30s", "30m" or "30h".                                                                                                                                                                                                             | 30m                     | No       |
| REFRESH_TOKEN_EXPIRY                  | Duration        | The expiration time for the refresh tokens. Should be a valid string duration such as: "30s", "30m" or "30h".                                                                                                                                                                                                            | 30d                     | No       |
| PUBLIC_KEY                            | RSA Public Key  | A public key used to validate access and refresh tokens.                                                                                                                                                                                                                                                                 | Auto-generate if empty. | No       |
| PRIVATE_KEY                           | RSA Private Key | A private key used to sign access and refresh tokens.                                                                                                                                                                                                                                                                    | Auto-generate if empty. | No       |
| MONGO_URI                             | String          | A URI for your Mongodb database.                                                                                                                                                                                                                                                                                         | ""                      | Yes      |
| INITIALIZE_SECRETS                    | Bool            | If set to true, the service will automatically create encryption secrsts and RSA public/private keys.                                                                                                                                                                                                                    | false                   | No       |
| INITIALIZE_ENGINE                     | String          | If set to "kubernetes", the service will create secrets used when the system starts up again. If set to "memory" the system will create new secrets when starting up again - do not use this option in production. This is only relevant if you initialize encryption secrtets.                                          | memory                  | No       |
| NEW_ENCRYPTION_KEY                    | String          | If provided, the system will automatically add a new encryption key to the system and encrypt users under that new key.                                                                                                                                                                                                  | ""                      | No       |
| MAX_EMAIL_VERIFICATION_AGE            | Duration        | Define how long the email verification code is valid.                                                                                                                                                                                                                                                                    | 5m                      | No       |

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Enjoy Hera? Show your gratitude by buying our team a cup of coffee!

We spend many hours building Hera and our other open-source projects - if you want to support our work, buy our
founder (and team) a cup of coffee below:

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/sinbadio)

