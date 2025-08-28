<script setup lang="ts">

import {Button} from "@/components/ui/button/index.js";
import DialogGeneral from "@/components/DialogGeneral.vue";
import Login from "@/components/Login.vue";
import {cn} from "@/lib/utils.js";
import {useScreenSheetStore} from "@/stores/useScreenSheetStore.js";

const {isDialogOpen, setOpenDialog} = useScreenSheetStore()

const props = defineProps({
  className: {
    required: false,
  }
})
</script>

<template class="flex">
  <DialogGeneral description="Sign in to your account" title="Sign In"
                 :is-open="isDialogOpen('LoginDialog')">
    <template #trigger>
      <Button @click="setOpenDialog(true, {name: 'LoginDialog'})" size="lg"
              :class="cn(props.className)">Sign In
      </Button>
    </template>
    <template #content>
      <Login></Login>
    </template>
  </DialogGeneral>
</template>

<style scoped>

</style>