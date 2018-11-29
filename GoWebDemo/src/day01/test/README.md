[toc]
# 1. 前言
exiaConnector是一个project和scala的混合工程，该工程主要实现从ifusion读取（点，边）数据，经SNA算法计算结果，再将结果保存。工程主要包含四个子模块：
+ exia-algorithm：SNA算法子模块
+ exia-common：公共类，方法子模块
+ exia-example：SNA算法的调用和实现的一些示例子模块
+ spark-connector：spark从ifusion加载数据子模块

# 2. 编译
```
在根目录下运行：mvn clean install
```

# 3. 输出的algorithm jar包
```
exiaConnector/exia-algorithm/target/exia-algorithm-SNA.jar
```

# 4. 算法说明

## 4.1 路径搜索算法
### 4.1.1 前言
- 路径搜索算法:根据起点,终点,给定跳数,边的label.在给定的跳数范围内,找到起点到终点的满足label条件的所有路径.
- 实现思路:根据起点向外做扩展查询,找到目的点,则保存该路径,否则继续扩散查询.

### 4.1.2 spark-submit提交示例
```
spark-submit \
--class com.isinonet.exia.algorithm.pathSearch.PathSearch \
--master yarn \
--deploy-mode client \
--driver-memory 2g \
--driver-cores 2 \
--num-executors 3 \
--executor-cores 2 \
--executor-memory 2g \
/opt/exia/algorithms/exia-algorithm-SNA.jar \
--cluster=master:2181,idtp2:2181,idtp3:2181 \
--instance=instance \
--account=root \
--secret=123456 \
--es_locations=idtp2:9300,idtp3:9300,idtp4:9300,idtp5:9300 \
--es_index=ifusion \
--es_cluster=elasticsearch \
--visibility=visallo \
--prefix=ifusion \
--source=689f432a0d684bf089f0b3cb9e3b5a39 \
--target=8a7b6485156b44628fad524c524f9311 \
--hops=3 \
--output=hdfs://master:9000/findpath/result.json
```

**参数详解:**
- cluster：zk的集群
- instance：accumulo的instance名称
- account：accumulo的用户
- secret：Accumulo的用户密码
- es_locations：es的客户端连接location
- es_index：es的索引
- es_cluster：es的cluster名称
- visibility：visibility可见性
- prefix：table前缀
- source：搜索的起点
- target：搜索的终点
- hops=3：搜索的最大跳数
- output：执行结果状态保存文件

### 4.1.3 数据示例
**执行状态结果示例**
```
{
    "error_code":0,
    "message":"ok",
    "data":{"dataPath":"hdfs://master:9000/findpath/data.json"}
}
```

**数据结果示例**
```
{
    "startTime":"2018-10-26 13:49:55",
    "endTime":"2018-10-26 13:50:07",
    "labels":["http://www.isinonet.com/call#Message","http://www.isinonet.com/call#Call"],
    "hops":3,
    "pathCount":2,
    "graph":
        {"vertices":
            [
                {"_1":"http://www.isinonet.com/call#CallPhone","_2":["cfce488c22dd4882b37aeecd5398b21b","df83817dcdd5420588893fa478948579","5357b11666e84e96bcb3296e33d33f6f","2a8473d47157402ca30248af6f89ecba"]}
            ],
        "edges":
            [
                {"_1":"df83817dcdd5420588893fa478948579","_2":"cfce488c22dd4882b37aeecd5398b21b","_3":"http://www.isinonet.com/call#Message"},
                {"_1":"5357b11666e84e96bcb3296e33d33f6f","_2":"cfce488c22dd4882b37aeecd5398b21b","_3":"http://www.isinonet.com/call#Message"},
                {"_1":"df83817dcdd5420588893fa478948579","_2":"2a8473d47157402ca30248af6f89ecba","_3":"http://www.isinonet.com/call#Message"},
                {"_1":"2a8473d47157402ca30248af6f89ecba","_2":"5357b11666e84e96bcb3296e33d33f6f","_3":"http://www.isinonet.com/call#Call"}
            ]
        },
    "path":
        [
            {"_1":0,"_2":["df83817dcdd5420588893fa478948579","cfce488c22dd4882b37aeecd5398b21b"]},
            {"_1":1,"_2":["df83817dcdd5420588893fa478948579","2a8473d47157402ca30248af6f89ecba","5357b11666e84e96bcb3296e33d33f6f","cfce488c22dd4882b37aeecd5398b21b"]}
        ]
}
```

## 4.2 pagerank算法
### 4.2.1 前言
- 在一个有向图中，一个点有越多的入链，那么这个点就越重要，PR值就越高。入链就是指别的点指向该点的链接。
- 在这个实现中，主要根据一个图谱id，读出该图谱中的所有点和边，转为spark中对应的graphx，算法计算PR值，在对每一个点生成新的id，写入新的图谱。

### 4.2.2 spark-submit提交示例
```
spark-submit \
--class com.isinonet.exia.algorithm.pagerank.PageRank \
--master yarn \
--deploy-mode client \
--driver-memory 2g \
--driver-cores 2 \
--num-executors 3 \
--executor-cores 2 \
--executor-memory 2g \
/opt/exia/algorithms/exia-algorithm-SNA.jar \
--cluster=master:2181,idtp2:2181,idtp3:2181 \
--instance=instance \
--account=root \
--secret=123456 \
--es_locations=idtp2:9300,idtp3:9300,idtp4:9300,idtp5:9300 \
--es_index=ifusion \
--es_cluster=elasticsearch \
--visibility=visallo \
--prefix=ifusion \
--es.nodes=idtp2,idtp3,idtp4,idtp5 \
--es.port=9200 \
--source=aa4ac8b972f24ad49c3d361dc34fad94 \
--target=6220dbe3e151411795f00cb9b583578e \
--weight_name=http://www.isinonet.com/insider#weight \
--output=hdfs://master:9000/pagerank/result.json
```

**参数详解:**
- cluster：zk的集群
- instance：accumulo的instance名称
- account：accumulo的用户
- secret：Accumulo的用户密码
- es_locations：es的客户端连接location
- es_index：es的索引
- es_cluster：es的cluster名称
- visibility：visibility可见性
- prefix：table前缀
- source：搜索的起点
- target：搜索的终点
- hops=3：搜索的最大终点
- es.nodes：es节点的nodes
- es.port=9200：es外部连接端口
- weight_name：保存PR值的key
- output：执行结果状态保存文件


## 4.3 ConnectedComponent算法
### 4.3.1 前言
- Connected Components即连通域，适用于无向图。一个连通域，指的是这里面的任意两个顶点，都能找到一个path连通。
- 在这个实现中，主要根据一个图谱id，和新的图谱id，将计算后带有cc值的新图保存到新的图谱中。

### 4.3.2 spark-submit提交示例
```
spark-submit \
--class com.isinonet.exia.algorithm.connectedComponent.ConnectedComponents \
--master yarn \
--deploy-mode client \
--driver-memory 2g \
--driver-cores 2 \
--num-executors 3 \
--executor-cores 2 \
--executor-memory 2g \
/opt/exia/algorithms/exia-algorithm-SNA.jar \
--cluster=master:2181,idtp2:2181,idtp3:2181 \
--instance=instance \
--account=root \
--secret=123456 \
--es_locations=idtp2:9300,idtp3:9300,idtp4:9300,idtp5:9300 \
--es_index=ifusion \
--es_cluster=elasticsearch \
--visibility=visallo \
--prefix=ifusion \
--es.nodes=idtp2,idtp3,idtp4,idtp5 \
--es.port=9200 \
--source=aa4ac8b972f24ad49c3d361dc34fad94 \
--target=6220dbe3e151411795f00cb9b583578e \
--cc_name=http://www.isinonet.com/insider#ccName \
--output=hdfs://master:9000/cc/result.json
```

**参数详解:**
- cluster：zk的集群
- instance：accumulo的instance名称
- account：accumulo的用户
- secret：Accumulo的用户密码
- es_locations：es的客户端连接location
- es_index：es的索引
- es_cluster：es的cluster名称
- visibility：visibility可见性
- prefix：table前缀
- source：搜索的起点
- target：搜索的终点
- hops=3：搜索的最大终点
- es.nodes：es节点的nodes
- es.port=9200：es外部连接端口
- cc_name：保存cc值的key
- output：执行结果状态保存文件


## 4.4 label Propagation算法
### 4.4.1 前言
- 对于网络中的每一个节点，在初始阶段，LabelPropagation算法对每一个节点一个唯一的标签，在每一个迭代的过程中，每一个节点根据与其相连的节点所属的标签改变自己的标签，更改的原则是选择与其相连的节点中所属标签最多的社区标签为自己的社区标签，这便是标签传播的含义。随着社区标签的不断传播，最终紧密连接的节点将有共同的标签。
- 在这个实现中，主要根据一个图谱id，和新的图谱id，将计算后带有lp值的新图保存到新的图谱中。

### 4.4.2 spark-submit提交示例
```
spark-submit \
--class com.isinonet.exia.algorithm.labelPropagation.LabelProp \
--master yarn \
--deploy-mode client \
--driver-memory 2g \
--driver-cores 2 \
--num-executors 3 \
--executor-cores 2 \
--executor-memory 2g \
/opt/exia/algorithms/exia-algorithm-SNA.jar \
--cluster=master:2181,idtp2:2181,idtp3:2181 \
--instance=instance \
--account=root \
--secret=123456 \
--es_locations=idtp2:9300,idtp3:9300,idtp4:9300,idtp5:9300 \
--es_index=ifusion \
--es_cluster=elasticsearch \
--visibility=visallo \
--prefix=ifusion \
--es.nodes=idtp2,idtp3,idtp4,idtp5 \
--es.port=9200 \
--source=aa4ac8b972f24ad49c3d361dc34fad94 \
--target=6220dbe3e151411795f00cb9b583578e \
--lp_name=http://www.isinonet.com/insider#lpName \
--iternum=10 \
--output=hdfs://master:9000/lp/result.json
```

**参数详解:**
- cluster：zk的集群
- instance：accumulo的instance名称
- account：accumulo的用户
- secret：Accumulo的用户密码
- es_locations：es的客户端连接location
- es_index：es的索引
- es_cluster：es的cluster名称
- visibility：visibility可见性
- prefix：table前缀
- source：搜索的起点
- target：搜索的终点
- hops=3：搜索的最大终点
- es.nodes：es节点的nodes
- es.port=9200：es外部连接端口
- lp_name：保存lp值的key
- iternum:算法迭代次数
- output：执行结果状态保存文件


