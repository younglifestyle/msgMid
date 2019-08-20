## 理解GC

[Kafka学习笔记](https://github.com/lzhpo/Queue/blob/master/Kafka%E5%AD%A6%E4%B9%A0%E7%AC%94%E8%AE%B0.md)
Group的作用就在于让多个组织可以独立消费同一个topic。

Consumer Group （CG）：这是kafka用来实现一个topic消息的广播（发给所有的consumer）和单播（发给任意一个consumer）的手段。一个topic可以有多个CG。topic的消息会复制（不是真的复制，是概念上的）到所有的CG，但每个partion只会把消息发给该CG中的一个consumer。如果需要实现广播，只要每个consumer有一个独立的CG就可以了。要实现单播只要所有的consumer在同一个CG。用CG还可以将consumer进行自由的分组而不需要多次发送消息到不同的topic；

[kafka里面的group存在的意义和作用是什么呢？](https://www.zhihu.com/question/263587973)
比如你和老憨，从属于这家公司的两个team，业务需求都不一样，但是都需要这个topic，所以你们需要独立的消费，就是说同一条message给了那个组也得给我这个组。

但是在同一个组里，为了并行消费，可以设置多个机器（或者就是多个进程），每个机器的消费是不独立的，就是说给了你的话就不用再给我了，因为我们俩是一个组里面的同一个业务计算。