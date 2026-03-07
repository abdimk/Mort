<div>
<a><img src="https://i.ibb.co/Q76fcW3s/Screenshot-2025-12-21-120851.png" alt="Screenshot 2025 12 21 120851" border="0"></a>
    <h1 align="center"> Loham</h1>
<p align="center">Real-Time Distributed Cache</p>
<p align="center"><em>This is for experimental propose </em></p>
</div>


---
## About
<p> Mort is a real-time distributed cache built for high-performance systems. Written in Go, it provides ultra-low latency data access, automatic distribution across nodes, and efficient memory management. Mort is designed to handle large-scale workloads with features like intelligent caching, horizontal scalability, and high availability, making it ideal for modern applications, microservices, and data-intensive platforms that need fast and reliable data access. </p>


> [!NOTE]
> This is Expermental and the core idea of the project is to impliment all the feauture that redis had and improve the thing that
> Implement the core features Redis has (must-have)
> Add features Redis lacks or does poorly (differentiation)
> Leverage Go strengths (concurrency, networking)

Core Features
- TinyLFU eviction
- Hot key replication
- Built-in rate limiting
- Cache stampede protection
- Multi-region replication
- WASM compute
- Native observability
- Auto-scaling cluster

