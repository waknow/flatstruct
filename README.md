# flatstruct

convert tree like structure into two dimensional table rows

# todo

- [ ] flat and unflat
- [ ] set primary object
- [ ] set muti-primary objects
- [ ] detect loop

# thoughs

- convert struct to tree, and print into console
- only one node can be marked as primary
- records need contain primary node path and object type, used to unflat struct

# problem

- if tree rebuild and rebuild back, the first tree and the last will be mirror each other. maybe need a flag to mark witch operation is rebuild back