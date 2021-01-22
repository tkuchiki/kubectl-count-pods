# kubectl-count-pods
A kubectl plugin to count the number of pods per status

## Installation

- Download from https://github.com/tkuchiki/kubectl-count-pods/releases
- `unzip kubectl-count-pods_os_arch.zip && mv kubectl-count-pods /usr/local/bin/`

## Usage

```
$ kubectl count pods -n NAMESPACE
+-----------+-------+
|  STATUS   | COUNT |
+-----------+-------+
| Pending   |     1 |
| Succeeded |     2 |
| Running   |     4 |
+-----------+-------+
|   TOTAL   |   7   |
+-----------+-------+
```
