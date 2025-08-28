<script setup lang="ts">

import {Button} from "@/components/ui/button/index.js";
import DialogGeneral from "@/components/DialogGeneral.vue";
import GeneralStepperForm from "@/components/GeneralStepperForm.vue";
import {useScreenSheetStore} from "@/stores/useScreenSheetStore.js";
import {cn} from "@/lib/utils.js";

const {isDialogOpen, setOpenDialog} = useScreenSheetStore()

const props = defineProps({
  className: {
    required: false,
  }
})
</script>

<template class="flex">
  <DialogGeneral :is-open="isDialogOpen('RegisterDialog')" title="Sign up"
                 description="Create your account">
    <template #trigger>
      <Button @click="setOpenDialog(true, {name: 'RegisterDialog'})" size="lg" variant="outline"
              :class="cn(props.className)">Sign Up
      </Button>
    </template>
    <template #content>
      <GeneralStepperForm></GeneralStepperForm>
    </template>
  </DialogGeneral>
</template>

<style scoped>

</style>