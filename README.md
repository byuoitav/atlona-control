# atlona-control
Microservice to RESTfully control Atlona video switchers

This microservice will immediately execute one connection and then will add a delay of 1 second for all subsiquent commands until all connections have completed.  Atlona switchers require 500ms before each command and this code includes 500ms of delay when reading the responses from the device.  

The following API guides were used and switchers have been tested:

https://atlona.com/pdf/AT-UHD-SW-52ED_API.pdf

https://atlona.com/pdf/AT-OME-PS62_API.pdf

https://atlona.com/pdf/AT-GAIN-60_API.pdf

The main package is located in the cmd folder
