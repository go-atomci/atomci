<style scoped>
.commonDialog .footer {
  margin-bottom: 0;
}

.footer .demonstration {
  line-height: 42px;
  float: right;
  margin-right: 8px;
}

.el-pagination {
  float: right;
  margin-top: 5px;
}
</style>

<template>
  <div class="footer clearfix" v-if="list.length>0">
    <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange" :current-page="currentPage" :page-sizes="sizeRange" :page-size="pageSize" :layout="showItems" :total="this.total">
    </el-pagination>
    <span class="demonstration">{{$t('bm.other.total', {total: this.total})}}</span>
  </div>
</template>
<script>
import { mapGetters } from 'vuex';

export default {
  props: ['list'],
  data() {
    return {
      sizeRange: [10, 20, 50, 100],
      pageSize: 10,
      currentPage: 1,
      showItems: 'sizes, prev, pager, next, jumper',
      total: 0,
    };
  },
  methods: {
    handleSizeChange(val) {
      this.pageSize = val;
      this.currentPage = 1;
      this.$emit('getlist');
    },
    handleCurrentChange(val) {
      this.currentPage = val;
      this.$emit('getlist', false);
    },
  },
  computed: {
    pageStart() {
      return this.pageSize * (this.currentPage - 1) + 1;
    },
    pageLast() {
      const last = this.pageSize * this.currentPage;
      return last;
    },
  },
};
</script>
