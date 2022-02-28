# Welcome to the Softcorp Block User Service

This service is one of the building blocks in the Softcorp Social Infrastructure.
You can run it by yourself by simply spinng up the docker container or set up a secure and scalable environment in the Softcorp Cloud.

## Run the softcorp-user-service program locally

In order to run the softcorp-user-service program locally, create an ```app/.env``` file containing the following
variables:

| Name                    | Type   | Explanation                                                                           |
 |-------------------------|---------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| MONGO_DB_NAME           | String | The name of the database where you want to store data from the Connect User Service.  |
| MONGO_USER_COLLECTION   | String | The name of the collection where you want to store data from the Connect User Service. |
| MONGO_URI               | String | The uri we can connect to your MongoDB through.                                       |
| GRPC_PORT       | Int    | The port where that you want the service to use.                                      |

