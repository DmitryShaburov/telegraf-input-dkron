# telegraf-input-dkron
Telegraf input (exec) plugin for Dkron

Written @ https://www.sravni.ru

## Metrics
Provides four metrics for each job in Dkron:
* state: 0 for failed, 1 for successful last execution
* last_duration: duration of last successful execution
* success_count: number of successful executions
* error_count: number of failed executions

Output format for metrics is influx

## Usage
`telegraf-input-dkron http://dkron.domain.local:8080`

## Example telegraf config
```
[[inputs.exec]]
    commands = ["/opt/telegraf-input-dkron http://dkron.domain.local:8080"]
    timeout = "10s"
    data_format = "influx"
```
