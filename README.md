# compact
Compact binary encoding to minimize storage size and maximize read speed

Datatypes are reduced for datatypes
0 - 100 // 1 byte -- scores, type byte
0.0 - 1.0000  // 2 bytes (percent that goes to 1.000 @ 4 sig dig), type FP1
0.00 - 655.00 // 2 byte, type F16
lat, lon // 4 bytes, type float32

# Presorted
The file is sorted by latitude, longitude ascending
