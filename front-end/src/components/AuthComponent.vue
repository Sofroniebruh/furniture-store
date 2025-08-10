<script setup>

import GeneralStepperForm from "@/components/GeneralStepperForm.vue";
import DialogGeneral from "@/components/DialogGeneral.vue";
import Login from "@/components/Login.vue";
import {Button} from "@/components/ui/button/index.js";
import {useScreenSheetStore} from "@/stores/useScreenSheetStore.js";

const {isDialogOpen, setOpenDialog} = useScreenSheetStore()
</script>

<template>
  <div class="flex flex-col justify-center items-center gap-12">
    <h1 class="font-semibold text-4xl">FUMI</h1>
    <div class="flex justify-around items-center h-full gap-3">
      <DialogGeneral :is-open="isDialogOpen('RegisterDialog')" title="Sign up"
                     description="Create your account">
        <template #trigger>
          <Button @click="setOpenDialog(true, {name: 'RegisterDialog'})" size="lg" variant="outline"
                  class="text-lg cursor-pointer">Sign Up
          </Button>
        </template>
        <template #content>
          <GeneralStepperForm></GeneralStepperForm>
        </template>
      </DialogGeneral>
      <DialogGeneral description="Sign in to your account" title="Sign In"
                     :is-open="isDialogOpen('LoginDialog')">
        <template #trigger>
          <Button @click="setOpenDialog(true, {name: 'LoginDialog'})" size="lg"
                  class="text-lg cursor-pointer bg-[#c9a275] hover:bg-[#dbb384]">Sign In
          </Button>
        </template>
        <template #content>
          <Login></Login>
        </template>
      </DialogGeneral>
    </div>
  </div>
</template>