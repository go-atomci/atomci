import backend from '../api/backend';

export default {
  data() {
    return {
      noDataTxt: this.$t('bm.add.noResult'),
    };
  },
  computed: {
    filteredList() {
      const list = this.curList || [];
      const filterVal = this.filterTxt;
      return filterVal.length > 0
        ? list.filter((element) => {
          return `${element[this.searchSelect]}`.indexOf(filterVal) > -1;
        })
        : list;
    },
    dataList() {
      // 默认一页显示十条  当前为第一页   后续的值由分页控件控制
      const start = (this.$refs.page && this.$refs.page.pageStart) || 1;
      const end = (this.$refs.page && this.$refs.page.pageLast) || 10;
      return this.filteredList.slice(start - 1, end);
    },
  },
  methods: {
    changeFilterTxt(val, searchSelect) {
      this.filterTxt = val;
      this.searchSelect = searchSelect;
    },
  },
};
