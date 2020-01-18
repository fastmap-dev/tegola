# HttpCache

httpcache configuration:

```toml
[cache]
type="http"
url="http://192.168.10.11/api/tile"
max_zoom=19

## Properties
The filecache config supports the following properties:

- `url` (string): [Required] http location.
- `max_zoom` (int): [Optional] the max zoom the cache should cache to. After this zoom, Set() calls will return before doing work.