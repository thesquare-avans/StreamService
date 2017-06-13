Fragmented Signed Data (FSD) protocol
=====================================

The Fragmented Signed Data (FSD) protocol can be used for transferring data from client to server. The data fragment is **always** signed using a signature to ensure data integrity (see below). No handshake is specified, the client talks _directly_ to the server. If a handshake is required, you can use another protocol for this or design your own.

Client connects to the server using TCP or another transport protocol and starts sending the fragmented data.

`Client ---------> Server`

The server _never_ sends data to the client and never response on data from the client, it only _accepts_ data.

## Fragment Header Format

`1 byte` is equal to `8 bits`.

|Order|Name        |Description                         |Type                                   |Length                 |
|-----|------------|------------------------------------|---------------------------------------|-----------------------|
|1    |Data length |The length of `Data` part in bytes. |Unsigned int (little endian)           |32 bits (4 bytes)      |
|2    |Signature   |Signature to verify the `Data` part.|RSA PKCS1 v1.5 signed SHA-256 hash of `Data` part.|2048 bits (256 bytes)|
|3    |Data        |The data fragment itself.           |Byte stream                            | Equal to `Data length`|

The total length of the fragment and fragment header can be calculated like this: `Data length + length of 'Data length' part + length of 'Signature' part = Data length + 4 + 256` in bytes.
