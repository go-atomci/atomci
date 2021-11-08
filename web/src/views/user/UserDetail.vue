<template>
  <div class="page-content buManage">
    <div class="portlet-body">
       <el-tabs v-model="activeName" @tab-click="handleClick">
         <el-tab-pane :label="$t('bm.add.beRole')" name="groupRole">
           <user-detail ref="groupRole" :isUser="true" isType="binding"></user-detail>
         </el-tab-pane>
         <el-tab-pane label="个人约束" name="groupPermission">
           <user-constraints ref="groupPermission" :isUser="true"></user-constraints>
         </el-tab-pane>
       </el-tabs>
    </div>
  </div>
</template>
<script>
import UserDetail from '@/components/view/UserDetail';
import UserConstraints from '@/components/view/UserConstraints';

export default {
  data() {
    return {
      activeName: 'groupRole',
    };
  },
  components: {
    UserDetail,
    UserConstraints,
  },
  mounted() {
    this.$refs[this.activeName].getList();
  },
  methods: {
    handleClick(tab) {
      this.activeName = tab.name;
      this.$nextTick(() => {
        if (!this.$refs[tab.name].curList || this.$refs[tab.name].curList.length === 0) {
          this.$refs[tab.name].getList();
        }
      });
    },
  },
};
</script>
