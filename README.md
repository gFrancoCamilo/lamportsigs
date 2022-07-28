# A Go Implementation of Lamport Signatures

A simple Go implementation of the [Lamport signature scheme](https://en.wikipedia.org/wiki/Lamport_signature). The implementation generates a key-pair $(pk,sk)$, where $pk$ is a public-key and $sk$ is the secret-key. Then, given a message $m$ the implementation generates a signature $\sigma$, so the user can send $(m,\sigma)$. Finally, the implementation given $(pk,m,\sigma)$ verifies if the signature $\sigma$ is valid or not. 

## References 
    Lamport, Leslie. "Constructing digital signatures from a one way function." (1979).