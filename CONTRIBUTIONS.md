# Style guide
Usually we would do `snake_case` everywhere as Lord intended. Unfortunately Go wants us to use capital letters to mark  
public variables/functions. Silly.  
```go
PascalCase // for public functions and variables
WHATEVER_THIS_IS // for public consts
snake_case // for anything else
```
# Two idiots one terminal
Row and Col start from 1 -> so that was a small shot in the knee

Inputable components must have static width
Other components can have static and dynamic width
If new data comes from the server, we will rerender the component or window (most likely) so dimensions will be recalculated
Same goes to height
