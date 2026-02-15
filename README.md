# txtllm
![txtllm](https://github.com/user-attachments/assets/0d3c1266-56e2-4f84-9da7-8ffa638416a1)



DNS server that answers `TXT` and `A` queries using an LLM.
### Is this useful in any way or form?
Maybe. If you want a quick answer to something without opening a browser or installing a cli tool, you can use `dig` to query this server. Or if you are on a restrictive network (Guest Wifi/Corporate firewall) that prohibits access to HTTPS/HTTP (ports 80/443) but allows port 53 for DNS resolution, this could be the only way to talk to an LLM.
### Can it search the web?
Yes. Use any openrouter model with the web plugin and it should incorporate web search. Note that web searches incur extra costs.
### Why?
Why not?
