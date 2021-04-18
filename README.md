# simplefsm

基于 [raft](https://github.com/bitcapybara/raft) 的极简状态机，模拟设备的启 `RUNNING` 停 `STOPPED` 状态，并通过对外暴露的Http接口下发 `START`  和 `STOP` 命令进行控制。

### 接口

* `GET /state` ：查询设备状态
* `POST /applyCommand` ：下发命令。传递url参数 `command=START` 或 `command=STOP` 
* `POST /sleep` ：模拟客户端下线，调用此接口后，客户端不会再响应除 `POST /awake` 外的任何请求
* `POST /awake` ：模拟客户端重新上线，开始接收请求
* `POST /changeConfig` ：成员变更请求。
* `POST /addLearner` ：添加 Learner 节点。
* `POST /transferLeadership` ：领导权转移请求

