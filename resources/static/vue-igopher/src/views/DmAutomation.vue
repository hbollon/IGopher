<template>
  <LateralNav/>
  <div class="d-flex flex-column" id="content-wrapper">
    <div id="content">
      <NavBar/>
      <DmAutomationPanel/>
    </div>
    <Footer/>
  </div><a class="border rounded scroll-to-top" href="#page-top"><i class="fas fa-angle-up"></i></a>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import LateralNav from '@/components/LateralNav.vue';
import NavBar from '@/components/NavBar.vue';
import Footer from '@/components/Footer.vue';
import DmAutomationPanel from '@/components/DmAutomationPanel.vue';

@Options({
  components: {
    LateralNav,
    NavBar,
    Footer,
    DmAutomationPanel,
  },
})
export default class DmAutomation extends Vue {}
</script>
