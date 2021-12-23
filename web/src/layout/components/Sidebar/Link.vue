<template>
  <component :is="type" v-bind="linkProps(to)">
    <slot />
  </component>
</template>

<script>
import { isExternal } from '@/utils/validate'
import { mapGetters } from 'vuex'

export default {
  props: {
    to: {
      type: String,
      required: true
    }
  },
  computed: {
    isExternal() {
      return isExternal(this.to)
    },
    type() {
      if (this.isExternal) {
        return 'a'
      }
      return 'router-link'
    },
    ...mapGetters({
      projectIDgetter: 'projectID',
    }),
    projectID() {
      let projectID = this.projectIDgetter
      if ( projectID === undefined ) {
        projectID =  this.$route.params.projectID
      }
      return projectID
    }
  },
  methods: {
    linkProps(to) {
      if (this.isExternal) {
        return {
          href: to,
          target: '_blank',
          rel: 'noopener'
        }
      }
      if (to.startsWith('/project/')) {
        to = to.replace(':projectID', this.projectID)
      }
      return {
        to: to
      }
    }
  }
}
</script>
