# hex-decoding

When querying vehicle data using the OBD protocol, the response from obd query returns the value as a hex string. To access specific values, such as vehicle speed in kilometers per hour, odometer reading in kilometers, accelerator position in percentage, etc, we can look at the formula field provided by the vehicle-signal-decoding API endpoint. 

Given a formula, we can get the offset, length, factor parts (scale factor and offset adjustment), range parts(min and max values), and units
