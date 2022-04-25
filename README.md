## Welcome to the Nuntio User Block 👋

![example workflow](https://github.com/nuntiodev/nuntio-user-block/actions/workflows/build.yaml/badge.svg)

Building Nuntio Blocks: secure, efficient & open-source API blocks that run on any Linux-based environment and scales
massively. Nuntio Cloud is the commercial offering of the above technologies and help companies go from idea to
production faster, without losing control over their data or services. Our goal with Nuntio Cloud is to provide social
API-blocks that are secure, easy to set up, scales worldwide and that you can move from Nuntio Cloud to your private
data center whenever you want.

## Nuntio User Block
The Nuntio User Block is a fully fledged user-management system, designed to be secure & stateless and to run in Kubernetes. It handles advanced data encryption models and features such as: normal and verfified authentication flows, creating, fetching and updating user information with custom metatadata and much more! Reach out to info@softcorp.io if you wanna give it a try.

## Environment

| Name                  | Type            | Description                                                                                                                                                                                                                                                                                                             | Default | Required |
|-----------------------|-----------------|:------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------| -------- |
| ENCRYPTION_KEYS       | []String        | An array of encryption keys used to encrypt/decrypt data at rest. If you provide more than one key, the keys will automatically be used to create a master key. If you think your data has been compromised, add another encryption key to the list, and the service will automatically encrypt users under the new key. | []     | No       |
| GRPC_PORT             | int             | The port you wish to start the gRPC server on.                                                                                                                                                                                                                                                                          | 9000   | No       |
| ACCESS_TOKEN_EXPIRY   | Duration        | The expiration time for the access tokens. Should be a valid string duration such as: "30s", "30m" or "30h".                                                                                                                                                                                                            | 30m    | No       |
| REFRESH_TOKEN_EXPIRY  | Duration        | The expiration time for the refresh tokens. Should be a valid string duration such as: "30s", "30m" or "30h".                                                                                                                                                                                                           | 30d    | No       |
| PUBLIC_KEY            | RSA Public Key  | A public key used to validate access and refresh tokens.                                                                                                                                                                                                                                                                | Empty - required. | Yes      |
| PRIVATE_KEY           | RSA Private Key | A private key used to sign access and refresh tokens.                                                                                                                                                                                                                                                                   | Empty - required. | Yes      |
| VALIDATE_PASSWORD     | Bool            | Use this if you want the service to validate all password automatically. Else passwords will only be validated if explicitly stated.                                                                                                                                                                                    | false  | No       |
| MONGO_URI             | String          | A URI for your Mongodb database. Use this if you want to provide the URL yourself, else use user and password authentication.                                                                                                                                                                                           | ""     | No       |
| MONGO_DB_USER         | String          | The username for your Mongodb user.                                                                                                                                                                                                                                                                                     | ""     | No       |
| MONGO_DB_USER_PASSWORD | String          | The password for your Mongodb user.                                                                                                                                                                                                                                                                                     | ""     | No       |
| MONGO_DB_HOST         | String          | The hostname for your Mongodb user.                                                                                                                                                                                                                                                                                     | ""     | No       |
| INITIALIZE_SECRETS    | Bool            | If set to true, the service will automatically create encryption secrsts and RSA public/private keys.                                                                                                                                                                                                                   | false  | No       |
| INITIALIZE_ENGINE     | String          | If set to "kubernetes", the service will create secrets used when the system starts up again. If set to "memory" the system will create new secrets when starting up again - do not use this option in production. This is only relevant if you initialize encryption secrtets.                                         | memory | No       |
| NEW_ENCRYPTION_KEY    | String          | If provided, the system will automatically add a new encryption key to the system and encrypt users under that new key.                                                                                                                                                                                                 | ""     | No       |
| ACTIVE_MEASUREMENT_EXPIRES_AT | Duration | If frontend is sending data to the server about user active engagement, this can be used to save the data for a specific amount of time. | 3 days | No |

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Buy me a cup of coffee!

I spend many hours building these API blocks - if you want to support my work, buy me a cup of coffee below:

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/sinbadio)

