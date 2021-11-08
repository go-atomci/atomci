<style>
.el-select .el-input {
  min-width: 132px;
}
</style>

<template>
  <el-input :placeholder="$t('bm.infrast.eSearchContent')" v-model.trim="filterTxt">
    <el-select slot="prepend" :placeholder="$t('bm.add.select')" v-model="searchSelect" @change="searchSelectChange">
      <el-option v-for="(item,index) in searchList" :key="index" :label="item.txt" :value="item.key"></el-option>
    </el-select>
  </el-input>
</template>
<script>
export default {
  props: ['searchList'],
  data() {
    return {
      filterTxt: '',
      searchSelect: '',
    };
  },
  created() {
    this.searchSelect = this.$props.searchList[0].key;
  },
  watch: {
    filterTxt(curVal) {
      this.$emit('changeFilterTxt', curVal, this.searchSelect);
    },
  },
  methods: {
    // 选择下拉列表框改变时
    searchSelectChange() {
      this.filterTxt = '';
      // 改变搜索字段时，原先搜索的值置空
      this.$props.filterTxt = '';
    },
  },
};
</script>
