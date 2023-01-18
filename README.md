# monitoring-agent-check-nt-replacement
Replacement for check_nt built ontop of monitoring-agent's /os_specific API

Essentially we can convert calls to 
```
check_nt -H HOST -p 12489 -v COUNTER -l -l "\\Memory\\Available Bytes","Available Bytes","Bytes"
``` 

into something like:

```
monitoring-agent-check-nt-replacement -host HOST -username username -password password -counter "\\Memory\\Available Bytes" -label "Available Bytes" -unit "Bytes" 
```
