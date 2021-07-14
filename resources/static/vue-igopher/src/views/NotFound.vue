<template>
  <div class="container-fluid">
    <div class="text-center mt-5">
        <div class="error mx-auto" data-text="404">
            <p class="m-0">404</p>
        </div>
        <p class="text-dark mb-5 lead">Page Not Found</p>
        <p class="text-black-50 mb-0">It looks like you found a glitch in the matrix...</p><router-link class="nav-link" to="/">← Back to Dashboard</router-link>
    </div>
  </div>
</template>

<script>
import { Vue, Options } from "vue-class-component";

@Options({
  title: "404 Not Found",
})
export default class NotFound extends Vue {

}
</script>

<style lang="scss">

</style>