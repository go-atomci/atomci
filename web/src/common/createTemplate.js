export default {
  methods: {
    doCancelCreate() {
      this.$refs.ruleForm.resetFields();
      this.dialogFormVisible = false;
    },
  },
};
