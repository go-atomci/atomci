<style>
.no-result {
  text-align: center;
  margin: 20px 0;
}

.no-result .title {
  font-size: 16px;
  color: #ccc;
  font-weight: 300;
  line-height: 50px;
}

.no-result .no-result-icon {
  background-image: url('../../assets/noresult.png');
  background-repeat: no-repeat;
  background-position: center;
  background-size: 200px 100%;
  width: 100%;
  height: 150px;
}
</style>

<template>
  <div class="no-result">
    <div class="title">
      {{ title || this.$t('bm.serviceM.noResultsFound')}}
    </div>
    <div class="no-result-icon"></div>
  </div>
</template>
<script>
export default {
  props: ['title'],
};
</script>
