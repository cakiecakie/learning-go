# Golang 内存分配

## 基本方针

1. 每次向操作系统申请大块内存，减少系统调用。
2. 将申请到的内存分割成特定大小的小块，构成链表。
3. 为对象分配内存时，只从链表中取出特定大小的内存块。
4. 回收对象内存时，把内存块放回原链表。
5. 如果闲置内存过多时，归还部分给操作系统。

## 内存分配器

### 内存块

内存分配器管理两种形式的内存块

1. span： 有多个地址连续的page组成的大块内存。 span 大小不固定，会发生裁剪和合并的操作。
2. object：将span 按特定大小分成多个小块，一个小块只放一个对象。object 按照 8 byte 的倍数分成多种。（分配器会尝试把多个小对象放在一个 object 里）

### 内存分配器的组件（thread cache malloc 架构）

1. cache：每个运行期的工作线程都会绑定一个 cache，用来无锁的分配 object块。
2. central：为所有 cache 提供切分好的后备 span。
3. heap：管理闲置的span，负责从操作系统申请内存。

### 内存分配流程

1. 计算分配对象对应的 size class。
2. 根据 size class 和要分配内存的对象是否包含指针（noscan）得到 span class, 根据 span class 在 cache.alloc 数组里找对应的 span。
3. 在 span.freelist 中取出 object。
4. 如果 span.freelist 为空，就去 central 获取新的相同 span class 的 span。
5. 如果 central.nonempty 为空，就去 heap.free/freelarge 然后切分成该 size class 的 object 链表。
6. 如果 heap 没有大小合适的 span，就去操作系统申请新内存块。

###  内存释放流程

1. 将标记为可回收的 object 块还给所属的 span.freelist。
2. 把 span 放回 central，可以被任意 cache 重新去用。
3. 如果 span 中的所有 object 都被回收，则还给 heap，可以被重新切分成另一种 size class 使用。
4. 定期扫描 heap 里的闲置 span，释放回操作系统。

* 大对象（32 KB）的分配和释放直接从 heap.freelarge 中分配和释放。大对象直接分配一个 span，对应的 size class 为0。

### tcmalloc 的好处

1. cache 可以实现具体线程的高性能分配。
2. 从属于某个cache 的span 如果使用率不高，就会造成浪费，回收可以把该 span 交还给 central，方便别的cache 使用。
3. 把span 交还给 heap 后，可以将其重新切分成新的 size，保证不同 size 的object 需求平衡。

### 部分源码

```go
type mheap struct {  // 全局唯一，不在 arena区
	lock mutex  // 保护mspan 的分配和归还
	allspans []*mspan
	arenas [1 << arenaL1Bits]*[1 << arenaL2Bits]*heapArena  // 一个二维数组，记录所有 arena 的 metadata
	curArena struct {  // 当前在使用的 arena 区
		base, end uintptr
	}
	central [numSpanClasses]struct {  // 用 span class 做索引的 central 数组
             	mcentral mcentral
		pad      [cpu.CacheLinePadSize - unsafe.Sizeof(mcentral{})%cpu.CacheLinePadSize]byte  // 用来对齐的，保证所有的 mcentral 的 lock 都只用来管理自己的 span。（目前没看懂）
        }
}

type heapArena struct {  // arena 的 metadata，自身保存在 arena 之外
	bitmap [heapArenaBitmapBytes]byte  // 两个bit 对应 arena 区一个地址空间（8B）。分别表示是不是指针，是否已扫描，被GC 使用。
	spans [pagesPerArena]*mspan  // 保存 page 与丛属 mspan 的对应关系。对于已经分配出去的span，page对应具体的mspan，目前是free 的span，只有开始和结束两个 page对应这个span，中间的page不保证（猜测与span的合并分裂有关）。从来没被分配过的page指向 nil。
}

type mspan struct {  // span 的 metadata，arena 之外。
	next *mspan
	prev *mspan
	startAddr uintptr  // span 的第一个字节的地址
	npages uintptr  // 共有多少个 page
	nelems uintptr  // 共有多少个 object
	freeindex uintptr  // 下一个object 的地址
	allocBits  *gcBits  // span 中 object 使用情况
	gcmarkBits *gcBits  // GC 标记的情况，GC 结束后就是新的 allocBits
	allocCache uint64  // freeindex 之后allocBits 的情况，注意 freeindex 之前一定被分配出去，之后的可能已经被分配。
}

type mcentral struct {
	lock      mutex  // 用于分配和归还 span 的锁
	spanclass spanClass  // 这个central 里 span的类型（存指针非指针，size class 多大）
	nonempty  mSpanList  // 还存在空余 object 的 span 列表
	empty     mSpanList  // 没有空余 object 或者已经被 mcache 取走的 span
	nmalloc uint64  // 这个 central 上被分配出去的 object 的个数，GC 使用
}

type mcache struct {
	tiny             uintptr  // 当前用于 tiny 对象分配的object 的地址
 	tinyoffset       uintptr  // 在这个 obejct 里分配对象的 offset
	local_tinyallocs uintptr  // tiny 对象被分配的次数
	alloc [numSpanClasses]*mspan  // 目前有的 mspan
}
```
