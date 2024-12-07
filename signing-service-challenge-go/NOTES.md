# Solution Notes

### How to run
    I have added a makefile for building and runing. also for simplicity i added the ability to add some mock data.

```
    make build
     
    #run tests
    make test 
    
    # build and run    
    make run
    
    #in a new terminal
    run-mock 
  ```

#### Endpoints

- Signing-Device
  - Get All
    <br>
    Get devices through pagination starting from page 1, both are required <br/>
    sample:
    ``` curl --location 'http://localhost:8080/api/v0/signature-devices?pageNr=1&pageSize=4'```
  
  - Get By ID
    <br>
    Get devices by id <br/>
    sample: ```curl --location 'http://localhost:8080/api/v0/signature-device/3 ```
      <br/>
  
  - Create
    <br>
    Create a device
     <br>sample: ```curl --location 'http://localhost:8080/api/v0/signature-device/3 ```
 
- Signing-Creation
  - Get All
    <br>Get signatures through pagination starting from page 1, all parameters  are required <br/>
    sample: ```curl --location 'http://localhost:8080/api/v0/signing-creations?deviceId=4&pageNr=1&pageSize=10'```  
    <br/>
  
  - Sign
     <br>Sing Data endpoint. Both data and device_id are required.
      <br> sample:
    ```
    curl --location 'http://localhost:8080/api/v0/signing-creation' \
    --header 'Content-Type: application/json' \
    --data '{
    "device_id":"4",
    "data":"test4"
    }' 
    ```
    <br/>
    <br/>
    


### Think of how to expose these operations through a RESTful HTTP-based API.
    I have used the standart library here(since it was already there), and even though you can't achieve full restfull appplications
    with the that i think it should be ok.  I would use either gin or gorilla mux or even goa to write full restfull applications.

### In addition, list / retrieval operations for the resources generated in the previous operations should be made available to the customers.
    Yes i have create paginated response for both of the resources as showing below:

```
    curl --location 'http://localhost:8080/api/v0/signature-devices?pageNr=1&pageSize=4'
    curl --location 'http://localhost:8080/api/v0/signing-creations?deviceId=3&pageNr=2&pageSize=3'
```

### The system will be used by many concurrent clients accessing the same resources.
    I haven't implemented the sync in the inmemory storage since this is supposed to be a simple solution,
    but i have implemented locking mechanism in service layer.

### The signature_counter has to be strictly monotonically increasing and ideally without any gaps.
    Yes i belive i have achieved this by having a lock from getting counter to storing it this counter.

### The system currently only supports RSA and ECDSA as signature algorithms. Try to design the signing mechanism in a way that allows easy extension to other algorithms without changing the core domain logic.
    Yes as long as the new algorithms are added to the enum,have a marshaller(AlgorythmMarshaller) and implement Signer interface we should be ok.
    for less dependency i have implemented factory type  design pattern.

### For now it is enough to store signature devices in memory. Efficiency is not a priority for this. In the future we might want to scale out. As you design your storage logic, keep in mind that we may later want to switch to a relational database.
    As long as we have repositories that implement the required SignatureDeviceRepository and SignatureDeviceRepository respectively

### QA / Testing
    I have written some tests for the services which contain the bussines logic,
    but certenly we should have more coverage in a real application.


