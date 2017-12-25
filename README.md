# mackerel-plugin-passenger

## Install

```
% mkr plugin install nabewata07/mackerel-plugin-passenger
```

## Setting

### If you can see passenger status by simply executing `passenger-status`

```
[plugin.metrics.passenger]
command = "/path/to/mackerel-plugin-passenger"
```

### If you use bundler

```
[plugin.metrics.passenger]
command = "/path/to/mackerel-plugin-passenger -work-dir '/path/to/application_root' -bundle-path '/path/to/command/bundle'"
```

### If you use bundler and need to specify `passenger-status` path

```
[plugin.metrics.passenger]
command = "/path/to/mackerel-plugin-passenger -work-dir '/path/to/application_root' -bundle-path '/path/to/command/bundle' -status-pash '/path/to/command/passenger-status'"
```

