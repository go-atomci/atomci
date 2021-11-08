<style>
.commonDialog .footer {
  margin-bottom: 0;
}

.footer {
  float: right;
}
.footer .demonstration {
  line-height: 42px;
}

.el-pagination {
  float: right;
  margin-top: 5px;
}
</style>

<template>
  <div class="footer" v-if="list.length>0">
    <span class="demonstration">{{$t('bm.other.total', {total: list.length})}}</span>
    <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange" :current-page="currentPage" :page-sizes="sizeRange" :page-size="pageSize" :layout="showItems" :total="list.length">
    </el-pagination>
  </div>
</template>
<script>

export default {
  // hasSizes 是否显示切换页码
  // beforeChangeCurrent 在切换页码之前
  // handleCurrent 切换完页面之后
  props: ['list', 'handleSize', 'beforeChangeCurrent', 'handleCurrent', 'hasSizes'],
  data() {
    return {
      sizeRange: [10, 20, 50, 100],
      pageSize: 10,
      currentPage: 1,
      showItems: 'sizes, prev, pager, next, jumper',
    };
  },
  methods: {
    handleSizeChange(val) {
      this.pageSize = val;
      this.$props.handleSize && this.$props.handleSize();
    },
    handleCurrentChange(val) {
      this.$props.beforeChangeCurrent && this.$props.beforeChangeCurrent();
      this.currentPage = val;
      if (this.$props.handleCurrent) {
        this.$nextTick(() => {
          this.$props.handleCurrent();
        });
      }
    },
  },
  mounted() {
    if (this.$props.hasSizes === false) {
      this.showItems = 'prev, pager, next, jumper';
    }
  },
  computed: {
    sum() {
      return this.$props.list.length;
    },
    pageStart() {
      return (this.pageSize) * (this.currentPage - 1) + 1;
    },
    pageLast() {
      let last = (this.pageSize) * this.currentPage;
      last = Math.min(last, this.sum);
      return last;
    },
  },
};
</script>
