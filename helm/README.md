 ## Running the Helm Chart
 Before running the Helm Chart:
 1. Create a secret named connect-user-secret with the following properties:
       ```
    kubectl create secret generic connect-user-secret \
           --from-literal=MONGO_DB_USER= your mongo DB user \
           --from-literal=MONGO_DB_HOST= your mongo DB host \
           --from-literal=MONGO_DB_USER_PASSWORD= your mongo DB password  \
           --from-literal=CONNECT_USER_STATIC_KEY= 32 char long key \
           --from-literal=CONNECT_USER_LOCAL_KEY= 32 char long key
    ```
    If you rather want to specify the connection string yourself, you can create the above secret as follows:
       ```
        kubectl create secret generic connect-user-secret \
               --from-literal=MONGO_URI= your mongo client uri  \
               --from-literal=CONNECT_USER_STATIC_KEY= 32 char long key \
               --from-literal=CONNECT_USER_LOCAL_KEY= 32 char long key
        ```
 2. Then install the helm chart by running ``` make helm-install  ``` within the root folder. 
 
 3. In order to delete the Helm chart, run     ```make helm-delete    ``` within the root folder.